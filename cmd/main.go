package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"
)

const fundsFileFlag = "fundsfile"

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

	config := Config{}
	config.ctx, config.ctxCancel = context.WithCancel(c.Context)

	err := config.init(c)
	if err != nil {
		return err
	}

	err = config.start()
	if err != nil {
		return err
	}

	slog.Debug("Exiting")

	return nil
}

func (conf *Config) init(c *cli.Context) error {

	slog.Debug("Initializing")

	if c.NumFlags() > 0 {
		if c.IsSet(fundsFileFlag) {
			conf.fundsFilePath = c.String(fundsFileFlag)
		} else {
			return fmt.Errorf("no flags set")
		}
	}

	funds, err := parseFunds(conf.fundsFilePath)
	if err != nil {
		return err
	}
	conf.Funds = funds

	return nil
}

func (conf *Config) start() error {
	slog.Debug("Starting")

	fundData, err := retrieveFundsData(conf.Funds)
	if err != nil {
		slog.Error("Failed to retrieve some or all data", "error", err)
		if len(fundData) == 0 {
			return err
		}
	}
	slog.Debug("Retrieved data", "data", fmt.Sprintf("%+v", fundData), "len", len(fundData))

	// TODO calculate performance, SMA3 etc from retrieved data
	perf, err := calculatePerformance(fundData)
	if err != nil {
		slog.Error("Failed to calculate performance", "error", err)
		return err
	}
	_ = perf

	// TODO output results sorted by performance

	return nil
}
