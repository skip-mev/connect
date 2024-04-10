package uniswapv3_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/apis/defi/uniswapv3"
	"github.com/skip-mev/slinky/providers/apis/defi/uniswapv3/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

func TestFetch(t *testing.T) {
	testCases := []struct {
		name     string
		tickers  []types.ProviderTicker
		client   func() uniswapv3.EVMClient
		expected types.PriceResponse
	}{
		{
			name:    "no tickers",
			tickers: []types.ProviderTicker{},
			client: func() uniswapv3.EVMClient {
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
			client: func() uniswapv3.EVMClient {
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
			client: func() uniswapv3.EVMClient {
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
			client: func() uniswapv3.EVMClient {
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
			client: func() uniswapv3.EVMClient {
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
			client: func() uniswapv3.EVMClient {
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
