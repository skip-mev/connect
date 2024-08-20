package coinmarketcap

import (
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

// NOTE: All documentation for this file can be located on the CoinMarketCap
// API documentation: https://coinmarketcap.com/api/documentation/v1/#operation/getV2CryptocurrencyQuotesLatest.
// This API does not require a subscription to use (i.e. No API key is required), but will
// get rate limited if too many requests are made.

const (
	// Name is the name of the CoinMarketCap provider.
	Name = "coinmarketcap_api"

	// URL is the base URL of the CoinMarketCap API that is used with API keys.
	URL = "https://pro-api.coinmarketcap.com"

	// URL is the base URL of the Coinbase API. This includes the base and quote
	// currency pairs that need to be inserted into the URL.
	Endpoint = "%s/v2/cryptocurrency/quotes/latest?id=%s"

	// APIKeyHeader is the header that is used to pass the API key to the CoinMarketCap API.
	APIKeyHeader = "X-CMC_PRO_API_KEY" //nolint

	// DefaultQuoteDenom is the default denomination provided by the CoinMarketCap API.
	DefaultQuoteDenom = "USD"
)

// DefaultAPIConfig is the default configuration for the CoinMarketCap API.
var DefaultAPIConfig = config.APIConfig{
	Name:             Name,
	Atomic:           true,
	Enabled:          true,
	Timeout:          3000 * time.Millisecond,
	Interval:         2000 * time.Millisecond,
	ReconnectTimeout: 2000 * time.Millisecond,
	MaxQueries:       1,
	Endpoints: []config.Endpoint{
		{
			URL: URL,
		},
	},
}

// CoinMarketCapResponse is the response from the CoinMarketCap API.
//
// ex.
//
//	{
//		"data":{
//		   "1":{
//			  "id":1,
//			  "name":"Bitcoin",
//			  "symbol":"BTC",
//			  "slug":"bitcoin",
//			  "is_active":1,
//			  "is_fiat":0,
//			  "circulating_supply":17199862,
//			  "total_supply":17199862,
//			  "max_supply":21000000,
//			  "date_added":"2013-04-28T00:00:00.000Z",
//			  "num_market_pairs":331,
//			  "cmc_rank":1,
//			  "last_updated":"2018-08-09T21:56:28.000Z",
//			  "tags":[
//				 "mineable"
//			  ],
//			  "platform":null,
//			  "self_reported_circulating_supply":null,
//			  "self_reported_market_cap":null,
//			  "quote":{
//				 "USD":{
//					"price":6602.60701122,
//					"volume_24h":4314444687.5194,
//					"volume_change_24h":-0.152774,
//					"percent_change_1h":0.988615,
//					"percent_change_24h":4.37185,
//					"percent_change_7d":-12.1352,
//					"percent_change_30d":-12.1352,
//					"market_cap":852164659250.2758,
//					"market_cap_dominance":51,
//					"fully_diluted_market_cap":952835089431.14,
//					"last_updated":"2018-08-09T21:56:28.000Z"
//				 }
//			  }
//		   }
//		},
//		"status":{
//		   "timestamp":"2024-06-26T08:00:37.500Z",
//		   "error_code":0,
//		   "error_message":"",
//		   "elapsed":10,
//		   "credit_count":1,
//		   "notice":""
//		}
//	}
//
// ref: https://coinmarketcap.com/api/documentation/v1/#operation/getV2CryptocurrencyQuotesLatest
type CoinMarketCapResponse struct { //nolint
	Data   map[string]CoinMarketCapData `json:"data"`
	Status CoinMarketCapStatus          `json:"status"`
}

// CoinMarketCapData is the data from the CoinMarketCap API.
type CoinMarketCapData struct { //nolint
	ID     int64                         `json:"id"`
	Name   string                        `json:"name"`
	Symbol string                        `json:"symbol"`
	Slug   string                        `json:"slug"`
	Quote  map[string]CoinMarketCapQuote `json:"quote"`
}

// CoinMarketCapQuote is the quote from the CoinMarketCap API.
type CoinMarketCapQuote struct { //nolint
	Price float64 `json:"price"`
}

// CoinMarketCapStatus is the status from the CoinMarketCap API.
type CoinMarketCapStatus struct { //nolint
	ErrorCode    int64  `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}
