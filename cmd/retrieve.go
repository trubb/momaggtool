package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// retrieveFundData fetches monthly data with a lookback of 14 months
func retrieveFundData(f Fund) (AzaResponse, error) {
	urlWithLookBack := func(lb int) string {
		to := time.Now().Format("2006-01-02")
		from := time.Now().AddDate(0, -lb, 0).Format("2006-01-02")
		return fmt.Sprintf(
			"https://www.avanza.se/_api/price-chart/stock/%d?from=%s&to=%s&resolution=month",
			f.AzaID,
			from,
			to,
		)
	}

	var response AzaResponse
	var errs []error

	n := 6
	lookback := 14
	for i := range n {
		// reduce lookback by 2 months each iteration [0..5] down to 14-(5*2)=4 months
		url := urlWithLookBack(lookback - i*2)

		response = AzaResponse{}

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return AzaResponse{}, err
		}
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			slog.Error("failure to retrieve", "fund", f.AzaID, "error", err)
			errs = append(errs, fmt.Errorf("failure to retrieve fund '%d', err: %w", f.AzaID, err))
			continue
		}
		if resp.StatusCode != http.StatusOK {
			slog.Error("non-ok response status", "status", resp.Status, "fund", f.AzaID, "url", url, "iteration", i)
			errs = append(errs, fmt.Errorf("non-ok response status '%s', fund '%d'", resp.Status, f.AzaID))
			continue
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			slog.Error("failure to decode", "fund", f.AzaID, "error", err)
			errs = append(errs, err)
			continue
		} else {
			break
		}
	}

	if len(errs) == n { // if all attempts failed we return the errors in a blob
		return AzaResponse{}, errors.Join(errs...)
	}
	// else we clearly succeded
	return response, nil
}

func retrieveFundsData(funds []Fund) ([]AzaResponse, error) {
	responseData := make([]AzaResponse, 0, len(funds))
	var errs []error

	for _, fund := range funds {
		azaResp, err := retrieveFundData(fund)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if azaResp.From != "" {
			responseData = append(responseData, azaResp)
		}
	}

	return responseData, errors.Join(errs...)
}
