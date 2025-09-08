package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

type UnrateResponse struct {
	Status       string   `json:"status"`
	ResponseTime int      `json:"responseTime"`
	Message      []string `json:"message"`
	Results      struct {
		Series []struct {
			SeriesID string `json:"seriesID"`
			Data     []struct {
				Year       string                   `json:"year"`
				Period     string                   `json:"period"`
				PeriodName string                   `json:"periodName"`
				Latest     string                   `json:"latest,omitempty"`
				Value      string                   `json:"value"`
				Footnotes  []map[string]interface{} `json:"footnotes"`
			} `json:"data"`
		} `json:"series"`
	} `json:"Results"`
}

func handleUnrate(ctx context.Context) (float64, float64, bool, error) {
	res, err := fetchUnrate(ctx)
	if err != nil {
		return 0, 0, false, err
	}

	if len(res.Results.Series) == 0 || len(res.Results.Series[0].Data) == 0 {
		return 0, 0, false, fmt.Errorf("no data in unrate response")
	}

	// latest data point
	ldp := res.Results.Series[0].Data[0]
	slog.Debug("latest unrate data", "data", fmt.Sprintf("year: %s, period: %s, value: %s", ldp.Year, ldp.Period, ldp.Value))
	unrate, err := strconv.ParseFloat(ldp.Value, 64)
	if err != nil {
		return 0, 0, false, fmt.Errorf("failed to parse latest unrate value: %w", err)
	}

	// find the 12-month moving average
	if len(res.Results.Series[0].Data) < 12 {
		return unrate, 0, false, fmt.Errorf("not enough data points to calculate 12-month moving average")
	}
	var sum float64
	for i := range 12 {
		v, err := strconv.ParseFloat(res.Results.Series[0].Data[i].Value, 64)
		if err != nil {
			slog.Error("failed to parse unrate value", "error", err)
		}
		sum += v
	}

	unrateMa12 := sum / 12

	triggerCondition := unrate > unrateMa12

	return unrate, unrateMa12, triggerCondition, nil
}

func fetchUnrate(ctx context.Context) (UnrateResponse, error) {
	// latest value
	const unrate = "LNS14000000"
	unrateURL := fmt.Sprintf("https://api.bls.gov/publicAPI/v2/timeseries/data/%s", unrate)

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, unrateURL, nil)
	if err != nil {
		return UnrateResponse{}, err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("failure to retrieve unrate", "error", err)
		return UnrateResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		slog.Error("non-ok response status", "status", resp.Status)
		return UnrateResponse{}, fmt.Errorf("non-ok response status '%s'", resp.Status)
	}
	defer resp.Body.Close()

	var response UnrateResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		slog.Error("failure to decode unrate response", "error", err)
	}

	return response, nil
}
