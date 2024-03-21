package dydx_test

import (
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

const dYdXResponseValid = `
{
	"market_params": [
	  {
		"id": 0,
		"pair": "BTC-USD",
		"exponent": -5,
		"min_exchanges": 3,
		"min_price_change_ppm": 1000,
		"exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"BTCUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"BTCUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"BTC-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"btcusdt\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XXBTZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"BTC_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}"
	  },
	  {
		"id": 1,
		"pair": "ETH-USD",
		"exponent": -6,
		"min_exchanges": 3,
		"min_price_change_ppm": 1000,
		"exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"ETHUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"ETHUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ETH-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"ethusdt\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XETHZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"ETH-USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"ETH_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ETH-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}"
	  },
	  {
		"id": 1000000,
		"pair": "USDT-USD",
		"exponent": -9,
		"min_exchanges": 3,
		"min_price_change_ppm": 1000,
		"exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"USDCUSDT\",\"invert\":true},{\"exchangeName\":\"Bybit\",\"ticker\":\"USDCUSDT\",\"invert\":true},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"USDT-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"ethusdt\",\"adjustByMarket\":\"ETH-USD\",\"invert\":true},{\"exchangeName\":\"Kraken\",\"ticker\":\"USDTZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"BTC-USD\",\"invert\":true},{\"exchangeName\":\"Okx\",\"ticker\":\"USDC-USDT\",\"invert\":true}]}"
	  }
	],
	"pagination": {
	  "next_key": null,
	  "total": "58"
	}
}
`

const dYdXResponseInvalid = `
{
	"market_params": [
	  {
		"id": 0,
		"pair": "BTC-USD",
		"exponent": -5,
		"min_exchanges": 3,
		"min_price_change_ppm": 1000,
		"exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"BTCUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"BTCUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"BTC-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"btcusdt\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XXBTZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"BTC_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"BTC-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}"
	  },
	  {
		"id": 1,
		"pair": "ETH-USD",
		"exponent": -6,
		"min_exchanges": 3,
		"min_price_change_ppm": 1000,
		"exchange_config_json": "{\"exchanges\":[{\"exchangeName\":\"Binance\",\"ticker\":\"ETHUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Bybit\",\"ticker\":\"ETHUSDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"CoinbasePro\",\"ticker\":\"ETH-USD\"},{\"exchangeName\":\"Huobi\",\"ticker\":\"ethusdt\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Kraken\",\"ticker\":\"XETHZUSD\"},{\"exchangeName\":\"Kucoin\",\"ticker\":\"ETH-USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Mexc\",\"ticker\":\"ETH_USDT\",\"adjustByMarket\":\"USDT-USD\"},{\"exchangeName\":\"Okx\",\"ticker\":\"ETH-USDT\",\"adjustByMarket\":\"USDT-USD\"}]}"
	  },
	],
	"pagination": {
	  "next_key": null,
	  "total": "58"
	}
}
`

var convertedResponse = mmtypes.GetMarketMapResponse{
	MarketMap: mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			"BTC/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     slinkytypes.NewCurrencyPair("BTC", "USD"),
					Decimals:         5,
					MinProviderCount: 3,
				},
				Paths: mmtypes.Paths{
					Paths: []mmtypes.Path{
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "binance_api",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "bybit_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "coinbase_api",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "coinbase_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "huobi_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "kraken_api",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "kucoin_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "mexc_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "okx_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
					},
				},
				Providers: mmtypes.Providers{
					Providers: []mmtypes.ProviderConfig{
						{
							Name:           "binance_api",
							OffChainTicker: "BTCUSDT",
						},
						{
							Name:           "bybit_ws",
							OffChainTicker: "BTCUSDT",
						},
						{
							Name:           "coinbase_api",
							OffChainTicker: "BTC-USD",
						},
						{
							Name:           "coinbase_ws",
							OffChainTicker: "BTC-USD",
						},
						{
							Name:           "huobi_ws",
							OffChainTicker: "btcusdt",
						},
						{
							Name:           "kraken_api",
							OffChainTicker: "XXBTZUSD",
						},
						{
							Name:           "kucoin_ws",
							OffChainTicker: "BTC-USDT",
						},
						{
							Name:           "mexc_ws",
							OffChainTicker: "BTCUSDT",
						},
						{
							Name:           "okx_ws",
							OffChainTicker: "BTC-USDT",
						},
					},
				},
			},

			"ETH/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     slinkytypes.NewCurrencyPair("ETH", "USD"),
					Decimals:         6,
					MinProviderCount: 3,
				},
				Providers: mmtypes.Providers{
					Providers: []mmtypes.ProviderConfig{
						{
							Name:           "binance_api",
							OffChainTicker: "ETHUSDT",
						},
						{
							Name:           "bybit_ws",
							OffChainTicker: "ETHUSDT",
						},
						{
							Name:           "coinbase_api",
							OffChainTicker: "ETH-USD",
						},
						{
							Name:           "coinbase_ws",
							OffChainTicker: "ETH-USD",
						},
						{
							Name:           "huobi_ws",
							OffChainTicker: "ethusdt",
						},
						{
							Name:           "kraken_api",
							OffChainTicker: "XETHZUSD",
						},
						{
							Name:           "kucoin_ws",
							OffChainTicker: "ETH-USDT",
						},
						{
							Name:           "mexc_ws",
							OffChainTicker: "ETHUSDT",
						},
						{
							Name:           "okx_ws",
							OffChainTicker: "ETH-USDT",
						},
					},
				},
				Paths: mmtypes.Paths{
					Paths: []mmtypes.Path{
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "binance_api",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "bybit_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "coinbase_api",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "coinbase_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "huobi_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "kraken_api",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "kucoin_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "mexc_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "okx_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
					},
				},
			},
			"USDT/USD": {
				Ticker: mmtypes.Ticker{
					CurrencyPair:     slinkytypes.NewCurrencyPair("USDT", "USD"),
					Decimals:         9,
					MinProviderCount: 3,
				},
				Providers: mmtypes.Providers{
					Providers: []mmtypes.ProviderConfig{
						{
							Name:           "binance_api",
							OffChainTicker: "USDCUSDT",
						},
						{
							Name:           "bybit_ws",
							OffChainTicker: "USDCUSDT",
						},
						{
							Name:           "coinbase_api",
							OffChainTicker: "USDT-USD",
						},
						{
							Name:           "coinbase_ws",
							OffChainTicker: "USDT-USD",
						},
						{
							Name:           "kraken_api",
							OffChainTicker: "USDTZUSD",
						},
						{
							Name:           "okx_ws",
							OffChainTicker: "USDC-USDT",
						},
					},
				},
				Paths: mmtypes.Paths{
					Paths: []mmtypes.Path{
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "binance_api",
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       true,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "bybit_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       true,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "coinbase_api",
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "coinbase_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "huobi_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       true,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "kraken_api",
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "kucoin_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       true,
								},
								{
									Provider:     mmtypes.IndexPrice,
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "okx_ws",
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       true,
								},
							},
						},
					},
				},
			},
		},
		AggregationType: mmtypes.AggregationType_INDEX_PRICE_AGGREGATION,
	},
}
