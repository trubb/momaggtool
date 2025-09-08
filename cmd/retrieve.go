package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// retrieveFundsData retrieves time series and performance data for all funds in the provided map
func retrieveFundsData(ctx context.Context, funds map[int]FundInfo) error {
	var errs []error
	for _, fd := range funds {
		c1, cCancel := context.WithTimeout(ctx, 15*time.Second)
		defer cCancel()

		azaResp, err := retrieveFundData(c1, fd.AzaID)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if azaResp.From != "" {
			slog.Debug(
				"retrieved data",
				"from", azaResp.From,
				"to", azaResp.To,
				"ohlc", len(azaResp.Ohlc),
				"fund", fd.AzaID,
				"name", fd.Name,
			)

			c2, c2Cancel := context.WithTimeout(ctx, 15*time.Second)
			defer c2Cancel()
			perf, err := retrieveFundsPerformance(c2, fd.AzaID)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			fd.AzaResponse = azaResp
			fd.ThreeMonthPerformance = perf.ThreeMonths
			funds[fd.AzaID] = fd
		}
	}

	return errors.Join(errs...)
}

// retrieveFundData fetches monthly data with a lookback of 14 months
//
// TODO change signature to require the URL as a string or a function returning a string to allow testing
// TODO could likely be simplified
// TODO could likely be merged with retrieveFundsPerformance
func retrieveFundData(ctx context.Context, azaID int) (AzaResponse, error) {

	const resolution = "day"

	urlWithLookBack := func(lb int) string {
		to := time.Now().Format("2006-01-02")
		from := time.Now().AddDate(0, 0, -lb).Format("2006-01-02")
		// TODO day resolution better???
		return fmt.Sprintf(
			"https://www.avanza.se/_api/price-chart/stock/%d?from=%s&to=%s&resolution=%s",
			azaID,
			from,
			to,
			resolution,
		)
	}

	var response AzaResponse
	var errs []error

	n := 6
	lookback := 150 // 130 days lookback should give us at least 90 valid days of trading
	for i := range n {

		// reduce lookback by 2 days each iteration [0..5] down to 130-(5*2)=120 days
		// if this fails we have encountered a fund which was recently created and need to deal with it manually
		url := urlWithLookBack(lookback - i*2)

		slog.Debug("retrieving data", "fund", azaID, "url", url, "iteration", i)

		client := &http.Client{}
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return AzaResponse{}, err
		}
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

		resp, err := client.Do(req)
		if err != nil {
			slog.Error("failure to retrieve", "fund", azaID, "error", err)
			errs = append(errs, fmt.Errorf("failure to retrieve fund '%d', err: %w", azaID, err))
			continue
		}
		if resp.StatusCode != http.StatusOK {
			slog.Error("non-ok response status", "status", resp.Status, "fund", azaID, "url", url, "iteration", i)
			errs = append(errs, fmt.Errorf("non-ok response status '%s', fund '%d'", resp.Status, azaID))
			continue
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			slog.Error("failure to decode", "fund", azaID, "error", err)
			errs = append(errs, err)

			time.Sleep(3 * time.Second) // be nice to the server

			continue
		} else {
			break
		}
	}

	slog.Debug("retrieved data", "fund", azaID, "data_points", len(response.Ohlc), "from", response.From, "to", response.To, "errors:", len(errs))

	return response, errors.Join(errs...)
}

// TODO will need handler for performance as well
/*
`GET https://www.avanza.se/_api/fund-reference/development/56127`
```json
{
"developmentOneDay":-0.00798,
"developmentOneWeek":0.03194,
"developmentOneMonth":0.21597,
"developmentThreeMonths":0.8695,
"developmentSixMonths":1.81212,
"developmentOneYear":4.39963,
"developmentThisYear":0.71543,
"developmentThreeYears":11.131816,
"developmentFiveYears":12.44839,
"developmentTenYears":15.304605
}```
`GET https://www.avanza.se/_api/fund-guide/chart/timeperiods/56127`...
```json
[
{"timePeriod":"one_month","change":0.0021597,"startDate":"2025-02-15"},
{"timePeriod":"three_months","change":0.008695000000000001,"startDate":"2024-12-15"},
{"timePeriod":"this_year","change":0.0071543,"startDate":"2025-01-02"},
{"timePeriod":"one_year","change":0.0439963,"startDate":"2024-03-15"},
{"timePeriod":"three_years","change":0.11131816,"startDate":"2022-03-15"},
{"timePeriod":"five_years","change":0.1244839,"startDate":"2020-03-15"},
{"timePeriod":"infinity","change":0.4049114151154967,"startDate":"2006-10-31"}
]```
*/

// retrieveFundsPerformance fetches performance data for a fund
//
//	TODO could likely be merged with retrieveFundData
func retrieveFundsPerformance(ctx context.Context, azaID int) (AzaPerformance, error) {

	url := fmt.Sprintf("https://www.avanza.se/_api/fund-reference/development/%d", azaID)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return AzaPerformance{}, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return AzaPerformance{}, fmt.Errorf("failure to retrieve fund '%d', err: %w", azaID, err)
	}
	if resp.StatusCode != http.StatusOK {
		return AzaPerformance{}, fmt.Errorf("non-ok response status '%s', fund '%d'", resp.Status, azaID)
	}
	defer resp.Body.Close()

	var response AzaPerformance
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return AzaPerformance{}, fmt.Errorf("failure to decode fund '%d', err: %w", azaID, err)
	}
	slog.Debug("retrieved performance data", "fund", azaID, "data", fmt.Sprintf("%+v", response))

	return response, nil
}
