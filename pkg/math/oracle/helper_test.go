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
	BTC_USD = mmtypes.Ticker{
		CurrencyPair:     constants.BITCOIN_USD.CurrencyPair,
		Decimals:         constants.BITCOIN_USD.Decimals,
		MinProviderCount: 3,
	}
	BTC_USDT = constants.BITCOIN_USDT

	ETH_USD = mmtypes.Ticker{
		CurrencyPair:     constants.ETHEREUM_USD.CurrencyPair,
		Decimals:         constants.ETHEREUM_USD.Decimals,
		MinProviderCount: 3,
	}
	ETH_USDT = constants.ETHEREUM_USDT

	USDT_USD = mmtypes.Ticker{
		CurrencyPair:     constants.USDT_USD.CurrencyPair,
		Decimals:         constants.USDT_USD.Decimals,
		MinProviderCount: 2,
	}
	USDC_USDT = constants.USDC_USDT

	PEPE_USD = mmtypes.Ticker{
		CurrencyPair:     constants.PEPE_USD.CurrencyPair,
		Decimals:         constants.PEPE_USD.Decimals,
		MinProviderCount: 1,
	}
	PEPE_USDT = constants.PEPE_USDT

	logger = zap.NewExample()

	// Marketmap is a test market map that contains a set of tickers, providers, and paths.
	// In particular all of the paths correspond to the desired "index prices" i.e. the
	// prices we actually want to resolve to.
	marketmap = mmtypes.MarketMap{
		Tickers: map[string]mmtypes.Ticker{
			BTC_USD.String():   BTC_USD,
			BTC_USDT.String():  BTC_USDT,
			USDT_USD.String():  USDT_USD,
			USDC_USDT.String(): USDC_USDT,
			ETH_USD.String():   ETH_USD,
			ETH_USDT.String():  ETH_USDT,
			PEPE_USDT.String(): PEPE_USDT,
			PEPE_USD.String():  PEPE_USD,
		},
		Paths: map[string]mmtypes.Paths{
			BTC_USD.String(): {
				Paths: []mmtypes.Path{
					{
						// COINBASE BTC/USD = BTC/USD
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USD.CurrencyPair,
								Invert:       false,
								Provider:     coinbase.Name,
							},
						},
					},
					{
						// COINBASE BTC/USDT * INDEX USDT/USD = BTC/USD
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USDT.CurrencyPair,
								Invert:       false,
								Provider:     coinbase.Name,
							},
							{
								CurrencyPair: USDT_USD.CurrencyPair,
								Invert:       false,
								Provider:     oracle.IndexProviderPrice,
							},
						},
					},
					{
						// BINANCE BTC/USDT * INDEX USDT/USD = BTC/USD
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USDT.CurrencyPair,
								Invert:       false,
								Provider:     binance.Name,
							},
							{
								CurrencyPair: USDT_USD.CurrencyPair,
								Invert:       false,
								Provider:     oracle.IndexProviderPrice,
							},
						},
					},
				},
			},
			ETH_USD.String(): {
				Paths: []mmtypes.Path{
					{
						// COINBASE ETH/USD = ETH/USD
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: ETH_USD.CurrencyPair,
								Invert:       false,
								Provider:     coinbase.Name,
							},
						},
					},
					{
						// COINBASE ETH/USDT * INDEX USDT/USD = ETH/USD
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: ETH_USDT.CurrencyPair,
								Invert:       false,
								Provider:     coinbase.Name,
							},
							{
								CurrencyPair: USDT_USD.CurrencyPair,
								Invert:       false,
								Provider:     oracle.IndexProviderPrice,
							},
						},
					},
					{
						// BINANCE ETH/USDT * INDEX USDT/USD = ETH/USD
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: ETH_USDT.CurrencyPair,
								Invert:       false,
								Provider:     binance.Name,
							},
							{
								CurrencyPair: USDT_USD.CurrencyPair,
								Invert:       false,
								Provider:     oracle.IndexProviderPrice,
							},
						},
					},
				},
			},
			USDT_USD.String(): {
				Paths: []mmtypes.Path{
					{
						// COINBASE USDT/USD = USDT/USD
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: USDT_USD.CurrencyPair,
								Invert:       false,
								Provider:     coinbase.Name,
							},
						},
					},
					{
						// COINBASE USDC/USDT ^ -1 = USDT/USD
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: USDC_USDT.CurrencyPair,
								Invert:       true,
								Provider:     coinbase.Name,
							},
						},
					},
					{
						// BINANCE USDT/USD = USDT/USD
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: USDT_USD.CurrencyPair,
								Invert:       false,
								Provider:     binance.Name,
							},
						},
					},

					{
						// Kucoin BTC/USDT ^-1 * INDEX BTC/USD = USDT/USD
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: BTC_USDT.CurrencyPair,
								Invert:       true,
								Provider:     kucoin.Name,
							},
							{
								CurrencyPair: BTC_USD.CurrencyPair,
								Invert:       false,
								Provider:     oracle.IndexProviderPrice,
							},
						},
					},
				},
			},
			PEPE_USD.String(): {
				Paths: []mmtypes.Path{
					{
						// BINANCE PEPE/USDT * INDEX USDT/USD = PEPE/USD
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: PEPE_USDT.CurrencyPair,
								Invert:       false,
								Provider:     binance.Name,
							},
							{
								CurrencyPair: USDT_USD.CurrencyPair,
								Invert:       false,
								Provider:     oracle.IndexProviderPrice,
							},
						},
					},
				},
			},
		},
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
	// Convert the price to a big float so we can perform the multiplication.
	floatPrice := big.NewFloat(price)

	// Scale the price and convert it to a big.Int.
	one := oracle.ScaledOne(decimals)
	scaledPrice := new(big.Float).Mul(floatPrice, new(big.Float).SetInt(one))
	intPrice, _ := scaledPrice.Int(nil)
	return intPrice
}
