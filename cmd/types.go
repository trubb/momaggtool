package main

import "context"

// TODO rename to something more neutral? could technically be a stock or derivative
type Fund struct {
	AzaID int    `toml:"aza_id"`
	Name  string `toml:"name"`
}

type Funds struct {
	Funds []Fund `toml:"funds"`
}

type Config struct {
	ctx           context.Context
	ctxCancel     context.CancelFunc
	fundsFilePath string
	Funds         []Fund
}

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
	Ohlc                 []AzaOhlc   `json:"ohlc"`
	Metadata             AzaMetadata `json:"metadata"`
	From                 string      `json:"from"`
	To                   string      `json:"to"`
	PreviousClosingPrice float64     `json:"previousClosingPrice"`
}
