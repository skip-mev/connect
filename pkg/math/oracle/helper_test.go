package oracle_test

import (
	"go.uber.org/zap"

	pkgtypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/apis/binance"
	"github.com/skip-mev/connect/v2/providers/apis/coinbase"
	"github.com/skip-mev/connect/v2/providers/websockets/kucoin"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var (
	// Create some custom tickers for testing.
	BTC_USD = mmtypes.Ticker{
		CurrencyPair:     pkgtypes.NewCurrencyPair("BTC", "USD"),
		Decimals:         8,
		MinProviderCount: 3,
		Enabled:          true,
	}

	ETH_USD = mmtypes.Ticker{
		CurrencyPair:     pkgtypes.NewCurrencyPair("ETH", "USD"),
		Decimals:         11,
		MinProviderCount: 3,
		Enabled:          true,
	}

	USDT_USD = mmtypes.Ticker{
		CurrencyPair:     pkgtypes.NewCurrencyPair("USDT", "USD"),
		Decimals:         6,
		MinProviderCount: 2,
		Enabled:          true,
	}

	PEPE_USD = mmtypes.Ticker{
		CurrencyPair:     pkgtypes.NewCurrencyPair("PEPE", "USD"),
		Decimals:         18,
		MinProviderCount: 1,
		Enabled:          true,
	}

	logger = zap.NewExample()

	// Marketmap is a test market map that contains a set of tickers, providers, and paths.
	// In particular, all paths correspond to the desired "index prices" i.e. the
	// prices we actually want to resolve to.
	marketmap = mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			BTC_USD.String(): {
				Ticker: BTC_USD,
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           coinbase.Name,
						OffChainTicker: "BTC-USD",
					},
					{
						Name:           coinbase.Name,
						OffChainTicker: "BTC-USDT",
						NormalizeByPair: &pkgtypes.CurrencyPair{
							Base:  "USDT",
							Quote: "USD",
						},
					},
					{
						Name:           binance.Name,
						OffChainTicker: "BTCUSDT",
						NormalizeByPair: &pkgtypes.CurrencyPair{
							Base:  "USDT",
							Quote: "USD",
						},
					},
				},
			},
			ETH_USD.String(): {
				Ticker: ETH_USD,
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           coinbase.Name,
						OffChainTicker: "ETH-USD",
					},
					{
						Name:           coinbase.Name,
						OffChainTicker: "ETH-USDT",
						NormalizeByPair: &pkgtypes.CurrencyPair{
							Base:  "USDT",
							Quote: "USD",
						},
					},
					{
						Name:           binance.Name,
						OffChainTicker: "ETHUSDT",
						NormalizeByPair: &pkgtypes.CurrencyPair{
							Base:  "USDT",
							Quote: "USD",
						},
					},
				},
			},
			USDT_USD.String(): {
				Ticker: USDT_USD,
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						Name:           coinbase.Name,
						OffChainTicker: "USDT-USD",
					},
					{
						Name:           coinbase.Name,
						OffChainTicker: "USDC-USDT",
						Invert:         true,
					},
					{
						Name:           binance.Name,
						OffChainTicker: "USDTUSD",
					},
					{
						Name:           kucoin.Name,
						OffChainTicker: "BTC-USDT",
						Invert:         true,
						NormalizeByPair: &pkgtypes.CurrencyPair{
							Base:  "BTC",
							Quote: "USD",
						},
					},
				},
			},
			PEPE_USD.String(): {
				Ticker: PEPE_USD,
				ProviderConfigs: []mmtypes.ProviderConfig{
					{
						OffChainTicker: "PEPEUSDT",
						Name:           binance.Name,
						NormalizeByPair: &pkgtypes.CurrencyPair{
							Base:  "USDT",
							Quote: "USD",
						},
					},
				},
			},
		},
	}
)
