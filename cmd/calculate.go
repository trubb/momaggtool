package main

import "errors"

// TODO calculate performance, SMA3 etc from retrieved fund data

// TODO need to check that the data is enough to work with
// > 3 months of data
// sorted correctly
// no missing data
// retry if data is missing etc (perhaps in retrieve)

func sma(periodcount int, fundData AzaResponse) (any, error) {
	// TODO implement
	return nil, nil
}

func calculatePerformance(fundsData []AzaResponse) ([]any, error) {
	// TODO implement
	var errs []error
	for _, fd := range fundsData {
		sma3, err := sma(3, fd)
		if err != nil {
			errs = append(errs, err)
		}
		_ = sma3
	}

	return nil, errors.Join(errs...)
}
