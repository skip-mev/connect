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
									Provider:     "Binance",
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
									Provider:     "Bybit",
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
									Provider:     "CoinbasePro",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "Huobi",
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
									Provider:     "Kraken",
									CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "Kucoin",
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
									Provider:     "Mexc",
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
									Provider:     "Okx",
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
							Name:           "Binance",
							OffChainTicker: "BTCUSDT",
						},
						{
							Name:           "Bybit",
							OffChainTicker: "BTCUSDT",
						},
						{
							Name:           "CoinbasePro",
							OffChainTicker: "BTC-USD",
						},
						{
							Name:           "Huobi",
							OffChainTicker: "btcusdt",
						},
						{
							Name:           "Kraken",
							OffChainTicker: "XXBTZUSD",
						},
						{
							Name:           "Kucoin",
							OffChainTicker: "BTC-USDT",
						},
						{
							Name:           "Mexc",
							OffChainTicker: "BTC_USDT",
						},
						{
							Name:           "Okx",
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
							Name:           "Binance",
							OffChainTicker: "ETHUSDT",
						},
						{
							Name:           "Bybit",
							OffChainTicker: "ETHUSDT",
						},
						{
							Name:           "CoinbasePro",
							OffChainTicker: "ETH-USD",
						},
						{
							Name:           "Huobi",
							OffChainTicker: "ethusdt",
						},
						{
							Name:           "Kraken",
							OffChainTicker: "XETHZUSD",
						},
						{
							Name:           "Kucoin",
							OffChainTicker: "ETH-USDT",
						},
						{
							Name:           "Mexc",
							OffChainTicker: "ETH_USDT",
						},
						{
							Name:           "Okx",
							OffChainTicker: "ETH-USDT",
						},
					},
				},
				Paths: mmtypes.Paths{
					Paths: []mmtypes.Path{
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "Binance",
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
									Provider:     "Bybit",
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
									Provider:     "CoinbasePro",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "Huobi",
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
									Provider:     "Kraken",
									CurrencyPair: slinkytypes.NewCurrencyPair("ETH", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "Kucoin",
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
									Provider:     "Mexc",
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
									Provider:     "Okx",
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
							Name:           "Binance",
							OffChainTicker: "USDCUSDT",
						},
						{
							Name:           "Bybit",
							OffChainTicker: "USDCUSDT",
						},
						{
							Name:           "CoinbasePro",
							OffChainTicker: "USDT-USD",
						},
						{
							Name:           "Kraken",
							OffChainTicker: "USDTZUSD",
						},
						{
							Name:           "Okx",
							OffChainTicker: "USDC-USDT",
						},
					},
				},
				Paths: mmtypes.Paths{
					Paths: []mmtypes.Path{
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "Binance",
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       true,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "Bybit",
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       true,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "CoinbasePro",
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "Huobi",
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
									Provider:     "Kraken",
									CurrencyPair: slinkytypes.NewCurrencyPair("USDT", "USD"),
									Invert:       false,
								},
							},
						},
						{
							Operations: []mmtypes.Operation{
								{
									Provider:     "Kucoin",
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
									Provider:     "Okx",
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
