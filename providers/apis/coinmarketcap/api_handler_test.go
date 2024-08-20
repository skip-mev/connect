package coinmarketcap_test

import (
	"math/big"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/apis/coinmarketcap"
	"github.com/skip-mev/connect/v2/providers/base/testutils"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

var (
	btcusd = types.DefaultProviderTicker{
		OffChainTicker: "1",
	}
	ethusd = types.DefaultProviderTicker{
		OffChainTicker: "2",
	}
)

func TestCreateURL(t *testing.T) {
	cases := []struct {
		name        string
		tickers     []types.ProviderTicker
		url         string
		expectedErr bool
	}{
		{
			name:        "no tickers",
			tickers:     []types.ProviderTicker{},
			url:         "",
			expectedErr: true,
		},
		{
			name: "single valid currency pair",
			tickers: []types.ProviderTicker{
				btcusd,
			},
			url:         "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest?id=1",
			expectedErr: false,
		},
		{
			name: "multiple valid currency pairs",
			tickers: []types.ProviderTicker{
				btcusd,
				ethusd,
			},
			url:         "https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest?id=1,2",
			expectedErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coinmarketcap.NewAPIHandler(coinmarketcap.DefaultAPIConfig)
			require.NoError(t, err)

			url, err := h.CreateURL(tc.tickers)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.url, url)
		})
	}
}

func TestParseResponse(t *testing.T) {
	cases := []struct {
		name     string
		tickers  []types.ProviderTicker
		response *http.Response
		expected types.PriceResponse
	}{
		{
			name: "single valid ticker",
			tickers: []types.ProviderTicker{
				btcusd,
			},
			response: testutils.CreateResponseFromJSON(`
{
	"data":{
		"1":{
			"id":1,
			"name":"Bitcoin",
			"symbol":"BTC",
			"slug":"bitcoin",
			"is_active":1,
			"is_fiat":0,
			"circulating_supply":17199862,
			"total_supply":17199862,
			"max_supply":21000000,
			"date_added":"2013-04-28T00:00:00.000Z",
			"num_market_pairs":331,
			"cmc_rank":1,
			"last_updated":"2018-08-09T21:56:28.000Z",
			"tags":[
				"mineable"
			],
			"platform":null,
			"self_reported_circulating_supply":null,
			"self_reported_market_cap":null,
			"quote":{
				"USD":{
				"price":6602.60701122,
				"volume_24h":4314444687.5194,
				"volume_change_24h":-0.152774,
				"percent_change_1h":0.988615,
				"percent_change_24h":4.37185,
				"percent_change_7d":-12.1352,
				"percent_change_30d":-12.1352,
				"market_cap":852164659250.2758,
				"market_cap_dominance":51,
				"fully_diluted_market_cap":952835089431.14,
				"last_updated":"2018-08-09T21:56:28.000Z"
				}
			}
		}
	},
	"status":{
		"timestamp":"2024-06-26T08:00:37.500Z",
		"error_code":0,
		"error_message":"",
		"elapsed":10,
		"credit_count":1,
		"notice":""
	}
}
			`),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{
					btcusd: {
						Value: big.NewFloat(6602.60701122),
					},
				},
				types.UnResolvedPrices{},
			),
		},
		{
			name: "response with failure status",
			tickers: []types.ProviderTicker{
				btcusd,
			},
			response: testutils.CreateResponseFromJSON(`
{
	"status":{
		"timestamp":"2024-06-26T08:00:37.500Z",
		"error_code":1,
		"error_message":"this does not mog",
		"elapsed":10,
		"credit_count":1,
		"notice":""
	}
}
			`),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{},
				},
			),
		},
		{
			name: "bad response",
			tickers: []types.ProviderTicker{
				btcusd,
			},
			response: testutils.CreateResponseFromJSON("this also does not mog"),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{},
				},
			),
		},
		{
			name: "no quote for ticker",
			tickers: []types.ProviderTicker{
				btcusd,
			},
			response: testutils.CreateResponseFromJSON(`
{
	"data":{
		"1":{
			"id":1,
			"name":"Bitcoin",
			"symbol":"BTC",
			"slug":"bitcoin",
			"is_active":1,
			"is_fiat":0,
			"circulating_supply":17199862,
			"total_supply":17199862,
			"max_supply":21000000,
			"date_added":"2013-04-28T00:00:00.000Z",
			"num_market_pairs":331,
			"cmc_rank":1,
			"last_updated":"2018-08-09T21:56:28.000Z",
			"tags":[
				"mineable"
			],
			"platform":null,
			"self_reported_circulating_supply":null,
			"self_reported_market_cap":null,
			"quote":{}
		}
	},
	"status":{
		"timestamp":"2024-06-26T08:00:37.500Z",
		"error_code":0,
		"error_message":"",
		"elapsed":10,
		"credit_count":1,
		"notice":""
	}
}
			`),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{},
				},
			),
		},
		{
			name: "no USD quote for ticker",
			tickers: []types.ProviderTicker{
				btcusd,
			},
			response: testutils.CreateResponseFromJSON(`
{
	"data":{
		"1":{
			"id":1,
			"name":"Bitcoin",
			"symbol":"BTC",
			"slug":"bitcoin",
			"is_active":1,
			"is_fiat":0,
			"circulating_supply":17199862,
			"total_supply":17199862,
			"max_supply":21000000,
			"date_added":"2013-04-28T00:00:00.000Z",
			"num_market_pairs":331,
			"cmc_rank":1,
			"last_updated":"2018-08-09T21:56:28.000Z",
			"tags":[
				"mineable"
			],
			"platform":null,
			"self_reported_circulating_supply":null,
			"self_reported_market_cap":null,
			"quote":{
				"MOG":{
					"price":6602.60701122,
					"volume_24h":4314444687.5194,
					"volume_change_24h":-0.152774,
					"percent_change_1h":0.988615,
					"percent_change_24h":4.37185,
					"percent_change_7d":-12.1352,
					"percent_change_30d":-12.1352,
					"market_cap":852164659250.2758,
					"market_cap_dominance":51,
					"fully_diluted_market_cap":952835089431.14,
					"last_updated":"2018-08-09T21:56:28.000Z"
				}
			}
		}
	},
	"status":{
		"timestamp":"2024-06-26T08:00:37.500Z",
		"error_code":0,
		"error_message":"",
		"elapsed":10,
		"credit_count":1,
		"notice":""
	}
}
			`),
			expected: types.NewPriceResponse(
				types.ResolvedPrices{},
				types.UnResolvedPrices{
					btcusd: providertypes.UnresolvedResult{},
				},
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h, err := coinmarketcap.NewAPIHandler(coinmarketcap.DefaultAPIConfig)
			require.NoError(t, err)

			// Update the cache since it is assumed that createURL is executed before ParseResponse.
			_, err = h.CreateURL(tc.tickers)
			require.NoError(t, err)

			now := time.Now()
			resp := h.ParseResponse(tc.tickers, tc.response)

			require.Len(t, resp.Resolved, len(tc.expected.Resolved))
			require.Len(t, resp.UnResolved, len(tc.expected.UnResolved))

			for cp, result := range tc.expected.Resolved {
				require.Contains(t, resp.Resolved, cp)
				r := resp.Resolved[cp]
				require.Equal(t, result.Value.SetPrec(18), r.Value.SetPrec(18))
				require.True(t, r.Timestamp.After(now))
			}

			for cp := range tc.expected.UnResolved {
				require.Contains(t, resp.UnResolved, cp)
				require.Error(t, resp.UnResolved[cp])
			}
		})
	}
}
