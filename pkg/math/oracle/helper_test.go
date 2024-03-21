package oracle_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	"github.com/skip-mev/slinky/providers/apis/binance"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	"github.com/skip-mev/slinky/providers/websockets/kucoin"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	// acceptableDelta is the acceptable difference between the expected and actual price.
	// In this case, we use a delta of 1e-8. This means we will accept any price that is
	// within 1e-8 of the expected price.
	acceptableDelta = 1e-8

	// Create some custom tickers for testing.
	BTC_USD = mmtypes.Market{
		Ticker: mmtypes.Ticker{
			CurrencyPair:     constants.BITCOIN_USD.CurrencyPair,
			Decimals:         constants.BITCOIN_USD.Decimals,
			MinProviderCount: 3,
		},
		Paths: mmtypes.Paths{
			Paths: []mmtypes.Path{
				{
					// COINBASE BTC/USD = BTC/USD
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
							Invert:       false,
							Provider:     coinbase.Name,
						},
					},
				},
				{
					// COINBASE BTC/USDT * INDEX USDT/USD = BTC/USD
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
							Invert:       false,
							Provider:     coinbase.Name,
						},
						{
							CurrencyPair: constants.USDT_USD.CurrencyPair,
							Invert:       false,
							Provider:     mmtypes.IndexPrice,
						},
					},
				},
				{
					// BINANCE BTC/USDT * INDEX USDT/USD = BTC/USD
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
							Invert:       false,
							Provider:     binance.Name,
						},
						{
							CurrencyPair: constants.USDT_USD.CurrencyPair,
							Invert:       false,
							Provider:     mmtypes.IndexPrice,
						},
					},
				},
			},
		},
		Providers: mmtypes.Providers{
			Providers: []mmtypes.ProviderConfig{
				coinbase.DefaultProviderConfig[constants.BITCOIN_USD],
			},
		},
	}

	BTC_USDT = mmtypes.Market{
		Ticker: constants.BITCOIN_USDT,
		Paths: mmtypes.Paths{
			Paths: []mmtypes.Path{
				{
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
							Invert:       false,
							Provider:     coinbase.Name,
						},
					},
				},
				{
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
							Invert:       false,
							Provider:     binance.Name,
						},
					},
				},
			},
		},
		Providers: mmtypes.Providers{
			Providers: []mmtypes.ProviderConfig{
				coinbase.DefaultProviderConfig[constants.BITCOIN_USDT],
				binance.DefaultNonUSProviderConfig[constants.BITCOIN_USDT],
				kucoin.DefaultProviderConfig[constants.BITCOIN_USDT],
			},
		},
	}

	ETH_USD = mmtypes.Market{
		Ticker: mmtypes.Ticker{
			CurrencyPair:     constants.ETHEREUM_USD.CurrencyPair,
			Decimals:         constants.ETHEREUM_USD.Decimals,
			MinProviderCount: 3,
		},
		Paths: mmtypes.Paths{
			Paths: []mmtypes.Path{
				{
					// COINBASE ETH/USD = ETH/USD
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.ETHEREUM_USD.CurrencyPair,
							Invert:       false,
							Provider:     coinbase.Name,
						},
					},
				},
				{
					// COINBASE ETH/USDT * INDEX USDT/USD = ETH/USD
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.ETHEREUM_USDT.CurrencyPair,
							Invert:       false,
							Provider:     coinbase.Name,
						},
						{
							CurrencyPair: constants.USDT_USD.CurrencyPair,
							Invert:       false,
							Provider:     mmtypes.IndexPrice,
						},
					},
				},
				{
					// BINANCE ETH/USDT * INDEX USDT/USD = ETH/USD
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.ETHEREUM_USDT.CurrencyPair,
							Invert:       false,
							Provider:     binance.Name,
						},
						{
							CurrencyPair: constants.USDT_USD.CurrencyPair,
							Invert:       false,
							Provider:     mmtypes.IndexPrice,
						},
					},
				},
			},
		},
		Providers: mmtypes.Providers{
			Providers: []mmtypes.ProviderConfig{
				coinbase.DefaultProviderConfig[constants.ETHEREUM_USD],
			},
		},
	}

	ETH_USDT = mmtypes.Market{
		Ticker: constants.ETHEREUM_USDT,
		Paths:  mmtypes.Paths{},
		Providers: mmtypes.Providers{
			Providers: []mmtypes.ProviderConfig{
				coinbase.DefaultProviderConfig[constants.ETHEREUM_USDT],
				binance.DefaultNonUSProviderConfig[constants.ETHEREUM_USDT],
			},
		},
	}

	USDT_USD = mmtypes.Market{
		Ticker: mmtypes.Ticker{
			CurrencyPair:     constants.USDT_USD.CurrencyPair,
			Decimals:         constants.USDT_USD.Decimals,
			MinProviderCount: 2,
		},
		Paths: mmtypes.Paths{
			Paths: []mmtypes.Path{
				{
					// COINBASE USDT/USD = USDT/USD
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.USDT_USD.CurrencyPair,
							Invert:       false,
							Provider:     coinbase.Name,
						},
					},
				},
				{
					// COINBASE USDC/USDT ^ -1 = USDT/USD
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.USDC_USDT.CurrencyPair,
							Invert:       true,
							Provider:     coinbase.Name,
						},
					},
				},
				{
					// BINANCE USDT/USD = USDT/USD
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.USDT_USD.CurrencyPair,
							Invert:       false,
							Provider:     binance.Name,
						},
					},
				},

				{
					// Kucoin BTC/USDT ^-1 * INDEX BTC/USD = USDT/USD
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.BITCOIN_USDT.CurrencyPair,
							Invert:       true,
							Provider:     coinbase.Name,
						},
						{
							CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
							Invert:       false,
							Provider:     mmtypes.IndexPrice,
						},
					},
				},
			},
		},
		Providers: mmtypes.Providers{
			Providers: []mmtypes.ProviderConfig{
				coinbase.DefaultProviderConfig[constants.USDT_USD],
				binance.DefaultNonUSProviderConfig[constants.USDT_USD],
			},
		},
	}

	USDC_USDT = mmtypes.Market{
		Ticker: constants.USDC_USDT,
		Paths:  mmtypes.Paths{},
		Providers: mmtypes.Providers{
			Providers: []mmtypes.ProviderConfig{
				coinbase.DefaultProviderConfig[constants.USDC_USDT],
			},
		},
	}

	PEPE_USD = mmtypes.Market{
		Ticker: mmtypes.Ticker{
			CurrencyPair:     constants.PEPE_USD.CurrencyPair,
			Decimals:         constants.PEPE_USD.Decimals,
			MinProviderCount: 1,
		},
		Paths:     mmtypes.Paths{},
		Providers: mmtypes.Providers{},
	}

	PEPE_USDT = mmtypes.Market{
		Ticker: constants.PEPE_USDT,
		Paths: mmtypes.Paths{
			Paths: []mmtypes.Path{
				{
					// BINANCE PEPE/USDT * INDEX USDT/USD = PEPE/USD
					Operations: []mmtypes.Operation{
						{
							CurrencyPair: constants.PEPE_USDT.CurrencyPair,
							Invert:       false,
							Provider:     binance.Name,
						},
						{
							CurrencyPair: constants.USDT_USD.CurrencyPair,
							Invert:       false,
							Provider:     mmtypes.IndexPrice,
						},
					},
				},
			},
		},
		Providers: mmtypes.Providers{
			Providers: []mmtypes.ProviderConfig{
				binance.DefaultNonUSProviderConfig[constants.PEPE_USDT],
			},
		},
	}

	logger = zap.NewExample()

	// Marketmap is a test market map that contains a set of tickers, providers, and paths.
	// In particular all paths correspond to the desired "index prices" i.e. the
	// prices we actually want to resolve to.
	marketmap = mmtypes.MarketMap{
		Markets: map[string]mmtypes.Market{
			BTC_USD.Ticker.String():   BTC_USD,
			BTC_USDT.Ticker.String():  BTC_USDT,
			USDT_USD.Ticker.String():  USDT_USD,
			USDC_USDT.Ticker.String(): USDC_USDT,
			ETH_USD.Ticker.String():   ETH_USD,
			ETH_USDT.Ticker.String():  ETH_USDT,
			PEPE_USDT.Ticker.String(): PEPE_USDT,
			PEPE_USD.Ticker.String():  PEPE_USD,
		},

		AggregationType: mmtypes.AggregationType_INDEX_PRICE_AGGREGATION,
	}
)

// verifyPrice verifies that the expected price matches the actual price within an acceptable delta.
func verifyPrice(t *testing.T, expected, actual *big.Int) {
	t.Helper()

	zero := big.NewInt(0)
	if expected.Cmp(zero) == 0 {
		require.Equal(t, zero, actual)
		return
	}

	var diff *big.Float
	if expected.Cmp(actual) > 0 {
		diff = new(big.Float).Sub(new(big.Float).SetInt(expected), new(big.Float).SetInt(actual))
	} else {
		diff = new(big.Float).Sub(new(big.Float).SetInt(actual), new(big.Float).SetInt(expected))
	}

	scaledDiff := new(big.Float).Quo(diff, new(big.Float).SetInt(expected))
	delta, _ := scaledDiff.Float64()
	t.Logf("expected price: %s; actual price: %s; diff %s", expected.String(), actual.String(), diff.String())
	t.Logf("acceptable delta: %.25f; actual delta: %.25f", acceptableDelta, delta)

	switch {
	case delta == 0:
		// If the difference between the expected and actual price is 0, the prices match.
		// No need for a delta comparison.
		return
	case delta <= acceptableDelta:
		// If the difference between the expected and actual price is within the acceptable delta,
		// the prices match.
		return
	default:
		// If the difference between the expected and actual price is greater than the acceptable delta,
		// the prices do not match.
		require.Fail(t, "expected price does not match the actual price; delta is too large")
	}
}

// createPrice creates a price with the given number of decimals.
func createPrice(price float64, decimals uint64) *big.Int {
	// Convert the price to a big float, so we can perform the multiplication.
	floatPrice := big.NewFloat(price)

	// Scale the price and convert it to a big.Int.
	one := oracle.ScaledOne(decimals)
	scaledPrice := new(big.Float).Mul(floatPrice, new(big.Float).SetInt(one))
	intPrice, _ := scaledPrice.Int(nil)
	return intPrice
}
