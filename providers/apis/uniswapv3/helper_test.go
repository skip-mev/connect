package uniswapv3_test

import (
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

	// PoolConfigs used for testing.
	weth_usdc_cfg = uniswapv3.PoolConfig{ //nolint
		Address:       "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
		BaseDecimals:  18,
		QuoteDecimals: 6,
		Invert:        true,
	}

	// Tickers used for testing.
	weth_usdc_ticker = mmtypes.Ticker{ //nolint
		CurrencyPair: pkgtypes.CurrencyPair{
			Base:  "WETH",
			Quote: "USDC",
		},
		Decimals:      18,
		Metadata_JSON: weth_usdc_cfg.ToJSON(),
	}
)

func createPriceFetcher(
	t *testing.T,
) *uniswapv3.PriceFetcher {
	t.Helper()

	client := mocks.NewEVMClient(t)
	fetcher, err := uniswapv3.NewPriceFetcher(
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
) *uniswapv3.PriceFetcher {
	t.Helper()

	fetcher, err := uniswapv3.NewPriceFetcher(
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

			require.True(t, ok)
			require.Equal(t, len(elems), len(responses))
			require.Equal(t, len(elems), len(errs))

			for i, elem := range elems {
				elem.Result = &responses[i]
				elem.Error = errs[i]
				elems[i] = elem
			}
		})
	}

	return c
}
