package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	fundsFileFlag = "fundsfile"
	logLevelFlag  = "loglevel"
)

func main() {
	app := &cli.App{
		Name:   "momaggtool",
		Usage:  "A tool for dealing with trend following strategies",
		Action: run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     fundsFileFlag,
				Aliases:  []string{"ff"},
				Usage:    "Load funds from toml-formatted `FILE`",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "loglevel",
				Aliases: []string{"l"},
				Usage:   "Set the log `LEVEL` (debug, info, warn, error)",
				Value:   "info",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run(c *cli.Context) error {
	// TODO set log level from arg
	slog.SetDefault(slog.New(slogHandler("")))

	slog.Info("Running momaggtool")

	svc := Service{}
	svc.ctx, svc.ctxCancel = context.WithCancel(c.Context)
	svc.FundInfo = make(map[int]FundInfo)

	if c.NumFlags() > 0 {
		if c.IsSet(fundsFileFlag) {
			svc.fundsFilePath = c.String(fundsFileFlag)
		} else {
			return fmt.Errorf("no funds file specified")
		}

		if c.IsSet(logLevelFlag) {
			slog.SetDefault(
				slog.New(
					slogHandler(
						c.String(logLevelFlag),
					),
				),
			)
		}
	}

	err := svc.init(c)
	if err != nil {
		return err
	}

	err = svc.start()
	if err != nil {
		return err
	}

	slog.Debug("Exiting")

	return nil
}

func (svc *Service) init(c *cli.Context) error {

	slog.Debug("Initializing")

	funds, err := parseFunds(svc.fundsFilePath)
	if err != nil {
		return err
	}
	for _, fund := range funds {
		svc.FundInfo[fund.AzaID] = FundInfo{Fund: fund}
	}

	return nil
}

func (svc *Service) start() error {
	slog.Debug("Starting")

	err := retrieveFundsData(svc.ctx, svc.FundInfo)
	if err != nil {
		slog.Error("Failed to retrieve some or all data", "error", err)
	}

	// TODO Check if we have any data at all

	// TODO calculate performance, SMA3M (SMA90d) etc from retrieved data
	err = calculateSMA(svc.FundInfo)
	if err != nil {
		slog.Error("Failed to calculate performance", "error", err)
		return err // TODO perhaps don't shit the bed if only one fails?
	}

	for _, v := range svc.FundInfo {
		slog.Info(
			"pre-ordering",
			"Name", v.Name,
			"AzaID", v.AzaID,
			"ThreeMonthPerformance", v.ThreeMonthPerformance,
			"SmaPeriod", v.SmaPeriod,
			"Sma", v.Sma,
			"SmaDistance", fmt.Sprintf("%.2f%%", v.SmaDistance*100),
		)
	}

	err = display(svc.FundInfo)
	if err != nil {
		slog.Error("Failed to display results", "error", err)
		return err
	}

	// TODO Fetch UNRATE

	// TODO pick top 3 funds based on our criteria
	// reuse print() but the set that is finally picked may differ from the ordered list
	// due to SMA distance etc
	// which is triggered by BLS UNRATE statistics

	return nil
}
