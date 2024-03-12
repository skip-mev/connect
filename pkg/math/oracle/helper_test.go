package oracle_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	"github.com/skip-mev/slinky/providers/apis/coinbase"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	// acceptableDelta is the acceptable difference between the expected and actual price.
	// In this case, we use a delta of 1e-8. This means we will accept any price that is
	// within 1e-8 of the expected price.
	acceptableDelta = 1e-8

	logger, _ = zap.NewDevelopment()
	marketmap = mmtypes.MarketMap{
		Tickers: map[string]mmtypes.Ticker{
			constants.BITCOIN_USD.String():   constants.BITCOIN_USD,
			constants.BITCOIN_USDT.String():  constants.BITCOIN_USDT,
			constants.USDT_USD.String():      constants.USDT_USD,
			constants.USDC_USDT.String():     constants.USDC_USDT,
			constants.ETHEREUM_USD.String():  constants.ETHEREUM_USD,
			constants.ETHEREUM_USDT.String(): constants.ETHEREUM_USDT,
		},
		Providers: map[string]mmtypes.Providers{
			constants.BITCOIN_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.BITCOIN_USD],
				},
			},
			constants.BITCOIN_USDT.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.BITCOIN_USDT],
				},
			},
			constants.USDT_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.USDT_USD],
				},
			},
			constants.USDC_USDT.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.USDC_USDT],
				},
			},
			constants.ETHEREUM_USD.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.ETHEREUM_USD],
				},
			},
			constants.ETHEREUM_USDT.String(): {
				Providers: []mmtypes.ProviderConfig{
					coinbase.DefaultMarketConfig[constants.ETHEREUM_USDT],
				},
			},
		},
		Paths: map[string]mmtypes.Paths{
			constants.BITCOIN_USD.String(): {
				Paths: []mmtypes.Path{
					{
						Operations: []mmtypes.Operation{
							{
								CurrencyPair: constants.BITCOIN_USD.CurrencyPair,
								Invert:       false,
								Provider:     coinbase.Name,
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
