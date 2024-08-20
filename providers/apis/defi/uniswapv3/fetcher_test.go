package uniswapv3_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/apis/defi/ethmulticlient"
	"github.com/skip-mev/connect/v2/providers/apis/defi/ethmulticlient/mocks"
	"github.com/skip-mev/connect/v2/providers/apis/defi/uniswapv3"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

func TestFetch(t *testing.T) {
	testCases := []struct {
		name     string
		tickers  []types.ProviderTicker
		client   func() ethmulticlient.EVMClient
		expected types.PriceResponse
	}{
		{
			name:    "no tickers",
			tickers: []types.ProviderTicker{},
			client: func() ethmulticlient.EVMClient {
				c := mocks.NewEVMClient(t)
				c.On("BatchCallContext", context.Background(), []rpc.BatchElem{}).Return(nil)
				return c
			},
			expected: types.PriceResponse{
				Resolved:   map[types.ProviderTicker]providertypes.ResolvedResult[*big.Float]{},
				UnResolved: map[types.ProviderTicker]providertypes.UnresolvedResult{},
			},
		},
		{
			name: "fails to retrieve pool for an empty ticker",
			tickers: []types.ProviderTicker{
				types.NewProviderTicker("WETH/USDC", ""),
			},
			client: func() ethmulticlient.EVMClient {
				return mocks.NewEVMClient(t)
			},
			expected: types.PriceResponse{
				Resolved: map[types.ProviderTicker]providertypes.ResolvedResult[*big.Float]{},
				UnResolved: map[types.ProviderTicker]providertypes.UnresolvedResult{
					types.NewProviderTicker("WETH/USDC", ""): {},
				},
			},
		},
		{
			name: "fails to make a batch call",
			tickers: []types.ProviderTicker{
				wethusdcTicker,
			},
			client: func() ethmulticlient.EVMClient {
				return createEVMClientWithResponse(t, fmt.Errorf("failed to make a batch call"), nil, nil)
			},
			expected: types.PriceResponse{
				Resolved: map[types.ProviderTicker]providertypes.ResolvedResult[*big.Float]{},
				UnResolved: map[types.ProviderTicker]providertypes.UnresolvedResult{
					wethusdcTicker: {},
				},
			},
		},
		{
			name: "batch request has an error for a single ticker",
			tickers: []types.ProviderTicker{
				wethusdcTicker,
			},
			client: func() ethmulticlient.EVMClient {
				batchErrors := []error{
					fmt.Errorf("request for ticker did not return a result"),
				}
				responses := []string{
					"",
				}
				return createEVMClientWithResponse(t, nil, responses, batchErrors)
			},
			expected: types.PriceResponse{
				Resolved: map[types.ProviderTicker]providertypes.ResolvedResult[*big.Float]{},
				UnResolved: map[types.ProviderTicker]providertypes.UnresolvedResult{
					wethusdcTicker: {},
				},
			},
		},
		{
			name: "batch request returns a result that cannot be parsed",
			tickers: []types.ProviderTicker{
				wethusdcTicker,
			},
			client: func() ethmulticlient.EVMClient {
				batchErrors := []error{
					nil,
				}
				responses := []string{
					"not a valid result",
				}
				return createEVMClientWithResponse(t, nil, responses, batchErrors)
			},
			expected: types.PriceResponse{
				Resolved: map[types.ProviderTicker]providertypes.ResolvedResult[*big.Float]{},
				UnResolved: map[types.ProviderTicker]providertypes.UnresolvedResult{
					wethusdcTicker: {},
				},
			},
		},
		{
			name: "weth/usdc mainnet result",
			tickers: []types.ProviderTicker{
				wethusdcTicker,
			},
			client: func() ethmulticlient.EVMClient {
				batchErrors := []error{
					nil,
				}
				responses := []string{
					"0x00000000000000000000000000000000000043dd3b966e761000000000000000000000000000000000000000000000000000000000000000000000000002fabf000000000000000000000000000000000000000000000000000000000000057900000000000000000000000000000000000000000000000000000000000005a000000000000000000000000000000000000000000000000000000000000005a000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
				}
				return createEVMClientWithResponse(t, nil, responses, batchErrors)
			},
			expected: types.PriceResponse{
				Resolved: map[types.ProviderTicker]providertypes.ResolvedResult[*big.Float]{
					wethusdcTicker: {
						Value: big.NewFloat(3313.131879703878971626114658316303),
					},
				},
				UnResolved: map[types.ProviderTicker]providertypes.UnresolvedResult{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fetcher := createPriceFetcherWithClient(t, tc.client())

			response := fetcher.Fetch(context.Background(), tc.tickers)
			require.Equal(t, len(tc.expected.Resolved), len(response.Resolved))
			require.Equal(t, len(tc.expected.UnResolved), len(response.UnResolved))

			for ticker, result := range tc.expected.Resolved {
				require.Contains(t, response.Resolved, ticker)
				require.Equal(t, result.Value.SetPrec(40), response.Resolved[ticker].Value.SetPrec(40))
			}

			for ticker := range tc.expected.UnResolved {
				require.Contains(t, response.UnResolved, ticker)
			}
		})
	}
}

func TestGetPool(t *testing.T) {
	fetcher := createPriceFetcher(t)

	t.Run("ticker is empty", func(t *testing.T) {
		ticker := types.NewProviderTicker("", "")
		_, err := fetcher.GetPool(ticker)
		require.Error(t, err)
	})

	t.Run("ticker does not have valid metadata", func(t *testing.T) {
		expected := uniswapv3.PoolConfig{
			Address: "0x1234",
		}
		ticker := types.NewProviderTicker("WETH/USDC", expected.MustToJSON())
		_, err := fetcher.GetPool(ticker)
		require.Error(t, err)
	})

	t.Run("ticker is not json formatted", func(t *testing.T) {
		ticker := types.NewProviderTicker("WETH/USDC", "not json, something else")
		_, err := fetcher.GetPool(ticker)
		require.Error(t, err)
	})

	t.Run("ticker has valid metadata", func(t *testing.T) {
		expected := uniswapv3.PoolConfig{
			Address:       "0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8", // WETH/USDC
			BaseDecimals:  18,
			QuoteDecimals: 6,
			Invert:        true,
		}
		ticker := types.NewProviderTicker("WETH/USDC", expected.MustToJSON())
		pool, err := fetcher.GetPool(ticker)
		require.NoError(t, err)
		require.Equal(t, expected, pool)
	})
}

func TestParseSqrtPriceX96(t *testing.T) {
	fetcher := createPriceFetcher(t)

	t.Run("result does not map to a string pointer", func(t *testing.T) {
		_, err := fetcher.ParseSqrtPriceX96(42)
		require.Error(t, err)
	})

	t.Run("result is nil", func(t *testing.T) {
		_, err := fetcher.ParseSqrtPriceX96(nil)
		require.Error(t, err)
	})

	t.Run("result is a nil string pointer", func(t *testing.T) {
		_, err := fetcher.ParseSqrtPriceX96((*string)(nil))
		require.Error(t, err)
	})

	t.Run("result is a empty string pointer", func(t *testing.T) {
		_, err := fetcher.ParseSqrtPriceX96(new(string))
		require.Error(t, err)
	})

	t.Run("result cannot be unpacked by the uniswap abi", func(t *testing.T) {
		result := new(string)
		*result = "0x1234"
		_, err := fetcher.ParseSqrtPriceX96(result)
		require.Error(t, err)
	})

	t.Run("mainnet result for BTC/USDT", func(t *testing.T) {
		result := new(string)
		*result = "0x000000000000000000000000000000000000001a105ec774c5175b820a157dac000000000000000000000000000000000000000000000000000000000000febe000000000000000000000000000000000000000000000000000000000000000d0000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"
		sqrtPriceX96, err := fetcher.ParseSqrtPriceX96(result)
		require.NoError(t, err)

		expectedResult, ok := new(big.Int).SetString("2064998566460012397847876304300", 10)
		require.True(t, ok)
		require.Equal(t, expectedResult, sqrtPriceX96)
	})

	t.Run("mainnet result for MOG/ETH", func(t *testing.T) {
		result := new(string)
		*result = "0x00000000000000000000000000000000000000000000d4df8e2f67f1e094d4dafffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc8f1c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"
		sqrtPriceX96, err := fetcher.ParseSqrtPriceX96(result)
		require.NoError(t, err)

		expectedResult, ok := new(big.Int).SetString("1005265563818767859111130", 10)
		require.True(t, ok)
		require.Equal(t, expectedResult, sqrtPriceX96)
	})
}

func TestNewPriceFetcher(t *testing.T) {
	ctx := context.TODO()

	testcases := []struct {
		name   string
		logger *zap.Logger
		api    config.APIConfig
		err    error
	}{
		{
			name:   "no logger errors",
			logger: nil,
			err:    fmt.Errorf("logger cannot be nil"),
		},
		{
			name:   "invalid api config errors",
			logger: logger,
			api: config.APIConfig{
				Enabled: true,
			},
			err: fmt.Errorf("invalid api config: "),
		},
		{
			name:   "invalid provider name errors",
			logger: logger,
			api: config.APIConfig{
				Name: "uniswapv3_api-foobar",
			},
			err: fmt.Errorf("invalid api config name uniswapv3_api-foobar"),
		},
		{
			name:   "invalid provider name errors",
			logger: logger,
			api: config.APIConfig{
				Name: "uniswapv3_api-ethereum",
			},
			err: fmt.Errorf("api config for uniswapv3_api-ethereum is not enabled"),
		},
		{
			name:   "no url or endpoints errors",
			logger: logger,
			api: config.APIConfig{
				Enabled:          true,
				Timeout:          1,
				ReconnectTimeout: 1,
				Interval:         1,
				MaxQueries:       1,
				Name:             "uniswapv3_api-ethereum",
			},
			err: fmt.Errorf("invalid api config"),
		},
		{
			name:   "multiclient failure errors",
			logger: logger,
			api: config.APIConfig{
				Enabled:          true,
				Timeout:          1,
				ReconnectTimeout: 1,
				Interval:         1,
				MaxQueries:       1,
				Endpoints: []config.Endpoint{
					{URL: "foobar", Authentication: config.Authentication{APIKey: "foobar", APIKeyHeader: "foobar"}},
				},
				Name: "uniswapv3_api-ethereum",
			},
			err: fmt.Errorf("failed to dial go ethereum client"),
		},
		{
			name:   "multiclient success",
			logger: logger,
			api: config.APIConfig{
				Enabled:          true,
				Timeout:          1,
				ReconnectTimeout: 1,
				Interval:         1,
				MaxQueries:       1,
				Endpoints: []config.Endpoint{
					{URL: "http://localhost:0", Authentication: config.Authentication{APIKey: "foobar", APIKeyHeader: "foobar"}},
				},
				Name: "uniswapv3_api-ethereum",
			},
			err: nil,
		},
		{
			name:   "url success",
			logger: logger,
			api: config.APIConfig{
				Enabled:          true,
				Timeout:          1,
				ReconnectTimeout: 1,
				Interval:         1,
				MaxQueries:       1,
				Endpoints:        []config.Endpoint{{URL: "http://localhost:0"}},
				Name:             "uniswapv3_api-ethereum",
			},
			err: nil,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			pf, err := uniswapv3.NewPriceFetcher(
				ctx,
				tc.logger,
				metrics.NewNopAPIMetrics(),
				tc.api,
			)
			if tc.err != nil {
				require.ErrorContains(t, err, tc.err.Error())
			} else {
				require.NoError(t, err)
				require.NotNil(t, pf)
			}
		})
	}
}
