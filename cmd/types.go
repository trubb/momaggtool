package main

import "context"

// TODO rename to something more neutral?
// could technically be a stock or derivative
type Fund struct {
	AzaID int    `toml:"aza_id"`
	Name  string `toml:"name"`
}

type Funds struct {
	Funds []Fund `toml:"funds"`
}

type FundInfo struct {
	Fund
	AzaResponse

	ThreeMonthPerformance float64
	Sma                   float64
	SmaPeriod             int
	SmaDistance           float64
}

type Service struct {
	ctx           context.Context
	ctxCancel     context.CancelFunc
	fundsFilePath string
	FundInfo      map[int]FundInfo
}

// read funds from file
// put them in a map keyed by aza_id
// but what should the map content type be?
//

/*
{
   "ohlc": [
      {
         "timestamp": 1735686000000,
         "open": 113.65379,
         "close": 116.93708,
         "low": 111.89387,
         "high": 117.11919,
         "totalVolumeTraded": 0
      },
      {
         "timestamp": 1738364400000,
         "open": 115.67039,
         "close": 110.95573,
         "low": 108.80958,
         "high": 116.10879,
         "totalVolumeTraded": 0
      },
      {
         "timestamp": 1740783600000,
         "open": 106.41978,
         "close": 99.5495,
         "low": 99.5495,
         "high": 106.41978,
         "totalVolumeTraded": 0
      }
   ],
   "metadata": {
      "resolution": {
         "chartResolution": "month",
         "availableResolutions": [
            "day",
            "week",
            "month"
         ]
      }
   },
   "from": "2024-12-10",
   "to": "2025-03-11",
   "previousClosingPrice": 116.655000
}
*/

// WARN seems to only contain data up to the last complete month
// This means that the SMA will be a bit off if you are too late in the month
// One way to solve it would be to fetch ~90 days of daily data and calculate the SMA from that
// The problem then will be to get a full 90 days worth of trading days
// This has knock-on effects e.g. in the fund data retrieval function
type AzaOhlc struct {
	Timestamp         int     `json:"timestamp"`
	Open              float64 `json:"open"`
	Close             float64 `json:"close"`
	Low               float64 `json:"low"`
	High              float64 `json:"high"`
	TotalVolumeTraded int     `json:"totalVolumeTraded"`
}

type AzaMetadata struct {
	AzaResolution `json:"resolution"`
}

type AzaResolution struct {
	ChartResolution      string   `json:"chartResolution"`
	AvailableResolutions []string `json:"availableResolutions"`
}

type AzaResponse struct {
	Ohlc     []AzaOhlc   `json:"ohlc"`
	Metadata AzaMetadata `json:"metadata"`
	From     string      `json:"from"`
	To       string      `json:"to"`
	// WARN this value seems to represent the NAV of the day immediately preceeding the "From" date
	PreviousClosingPrice float64 `json:"previousClosingPrice"`
}

type AzaPerformance struct {
	OneDay      float64 `json:"developmentOneDay"`
	OneWeek     float64 `json:"developmentOneWeek"`
	OneMonth    float64 `json:"developmentOneMonth"`
	ThreeMonths float64 `json:"developmentThreeMonths"`
	SixMonths   float64 `json:"developmentSixMonths"`
	OneYear     float64 `json:"developmentOneYear"`
	ThisYear    float64 `json:"developmentThisYear"`
	ThreeYears  float64 `json:"developmentThreeYears"`
	FiveYears   float64 `json:"developmentFiveYears"`
	TenYears    float64 `json:"developmentTenYears"`
}
