package uniswapv3_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	pkgtypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/apis/uniswapv3"
	"github.com/skip-mev/slinky/providers/apis/uniswapv3/mocks"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	// precision is the precision used for big.Float calculations. Specifically
	// this is used to ensure that float values are the same within a certain
	// precision.
	precision = 30

	// acceptableDelta is the acceptable difference between the expected and actual price.
	// In this case, we use a delta of 1e-8. This means we will accept any price that is
	// within 1e-8 of the expected price.
	acceptableDelta = 1e-8
)

var (
	logger, _ = zap.NewDevelopment()
	m         = metrics.NewNopAPIMetrics()

	// PoolConfigs used for testing
	weth_usdc_cfg = uniswapv3.PoolConfig{
		Address:       "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
		BaseDecimals:  18,
		QuoteDecimals: 6,
		Invert:        true,
	}
	eth_usdc_cfg = uniswapv3.PoolConfig{
		Address:       "0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640",
		BaseDecimals:  18,
		QuoteDecimals: 6,
		Invert:        true,
	}
	mog_eth_cfg = uniswapv3.PoolConfig{
		Address:       "0x7832310Cd0de39c4cE0A635F34d9a4B5b47fd434",
		BaseDecimals:  18,
		QuoteDecimals: 18,
		Invert:        false,
	}
	btc_usdt_cfg = uniswapv3.PoolConfig{
		Address:       "0x9Db9e0e53058C89e5B94e29621a205198648425B",
		BaseDecimals:  8,
		QuoteDecimals: 6,
		Invert:        false,
	}

	// Tickers used for testing
	weth_usdc_ticker = mmtypes.Ticker{
		CurrencyPair: pkgtypes.CurrencyPair{
			Base:  "WETH",
			Quote: "USDC",
		},
		Decimals:      18,
		Metadata_JSON: weth_usdc_cfg.ToJSON(),
	}
	eth_usdc_ticker = mmtypes.Ticker{
		CurrencyPair: pkgtypes.CurrencyPair{
			Base:  "ETH",
			Quote: "USDC",
		},
		Decimals:      18,
		Metadata_JSON: eth_usdc_cfg.ToJSON(),
	}
	mog_eth_ticker = mmtypes.Ticker{
		CurrencyPair: pkgtypes.CurrencyPair{
			Base:  "MOG",
			Quote: "ETH",
		},
		Decimals:      18,
		Metadata_JSON: mog_eth_cfg.ToJSON(),
	}
	btc_usdt_ticker = mmtypes.Ticker{
		CurrencyPair: pkgtypes.CurrencyPair{
			Base:  "BTC",
			Quote: "USDT",
		},
		Decimals:      8,
		Metadata_JSON: btc_usdt_cfg.ToJSON(),
	}
)

func createPriceFetcher(
	t *testing.T,
) *uniswapv3.UniswapV3PriceFetcher {
	t.Helper()

	client := mocks.NewEVMClient(t)
	fetcher, err := uniswapv3.NewUniswapV3PriceFetcher(
		logger,
		m,
		uniswapv3.DefaultAPIConfig,
		client,
	)
	require.NoError(t, err)

	return fetcher
}

func createPriceFetcherWithClient(
	t *testing.T,
	client uniswapv3.EVMClient,
) *uniswapv3.UniswapV3PriceFetcher {
	t.Helper()

	fetcher, err := uniswapv3.NewUniswapV3PriceFetcher(
		logger,
		m,
		uniswapv3.DefaultAPIConfig,
		client,
	)
	require.NoError(t, err)

	return fetcher
}

func createEVMClientWithResponse(
	t *testing.T,
	failedRequestErr error,
	responses []string,
	errs []error,
) uniswapv3.EVMClient {
	t.Helper()

	c := mocks.NewEVMClient(t)
	if failedRequestErr != nil {
		c.On("BatchCallContext", mock.Anything, mock.Anything).Return(failedRequestErr)
	} else {
		c.On("BatchCallContext", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			elems, ok := args.Get(1).([]rpc.BatchElem)
			require.True(t, ok)
			fmt.Println(elems)

			require.True(t, ok)
			require.Equal(t, len(elems), len(responses))
			require.Equal(t, len(elems), len(errs))

			for i, elem := range elems {
				elem.Result = &responses[i]
				elem.Error = errs[i]
				elems[i] = elem
			}

			fmt.Println(elems)
		})
	}

	return c
}

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
