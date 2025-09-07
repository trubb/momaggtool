package main

import (
	"errors"
	"fmt"
	"log/slog"
)

// TODO calculate performance, SMA3 etc from retrieved fund data

// TODO may not be needed if we can grab it?
func calculateSMA(funds map[int]FundInfo) error {
	const periodLength = 90 // days

	var errs []error

	for _, fd := range funds {
		slog.Debug("Calculating SMA", "id", fd.AzaID, "name", fd.Name, "data_points", len(fd.Ohlc))
		sma, err := simpleMovingAverage(periodLength, fd.Ohlc, fd.Name) // TODO from where do we select the period length?
		if err != nil {
			errs = append(errs, err)
		}

		lastClose := fd.Ohlc[len(fd.Ohlc)-1].Close

		smaDistance := ((lastClose - sma) / sma) * 100 // percentage

		fund, ok := funds[fd.AzaID]
		if ok {
			fund.Sma = sma
			fund.SmaDistance = smaDistance
			fund.SmaPeriod = periodLength

			funds[fd.AzaID] = fund

			slog.Debug(
				"calcPerf",
				"sma_period", fund.SmaPeriod,
				"sma", fmt.Sprintf("%f", fund.Sma),
				"distance", fmt.Sprintf("%f%%", fund.SmaDistance),
				"above?", lastClose > sma,
				"last close", lastClose,
				"id", fd.AzaID,
				"name", fd.Name,
			)
			// TODO above/below sma?
		}

		// TODO at this point we have the SMA and the last close price
		// this allows us to calculate the distance between the two
	}

	return errors.Join(errs...)
}

// TODO need to check that the data is enough to work with
// > 3 months of data
// sorted correctly
// no missing data
// retry if data is missing etc (perhaps in retrieve)

func simpleMovingAverage(periodLength int, data []AzaOhlc, fundName string) (float64, error) {
	if len(data) < periodLength {
		return 0, fmt.Errorf("not enough data points: %s, %d < %d", fundName, len(data), periodLength)
	}
	// grab, from the last <periodLength>, the closing prices
	// sum them up and divide by <periodLength>
	// return the sum
	var sum float64
	//	for i := 0; i < len(data) && i < 90; i++ { // TODO do we need a guard for i vs sma period count?
	for i := len(data); 0 < i && len(data)-periodLength < i; i-- {

		sum += data[len(data)-i].Close // will go 0<-len(data)
	}

	result := sum / float64(periodLength)

	return result, nil
}
