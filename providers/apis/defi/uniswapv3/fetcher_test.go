package uniswapv3_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/apis/defi/uniswapv3"
	"github.com/skip-mev/slinky/providers/apis/defi/uniswapv3/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func TestFetch(t *testing.T) {
	testCases := []struct {
		name     string
		tickers  []mmtypes.Ticker
		client   func() uniswapv3.EVMClient
		expected types.PriceResponse
	}{
		{
			name:    "no tickers",
			tickers: []mmtypes.Ticker{},
			client: func() uniswapv3.EVMClient {
				c := mocks.NewEVMClient(t)
				c.On("BatchCallContext", context.Background(), []rpc.BatchElem{}).Return(nil)
				return c
			},
			expected: types.PriceResponse{
				Resolved:   map[mmtypes.Ticker]providertypes.ResolvedResult[*big.Int]{},
				UnResolved: map[mmtypes.Ticker]providertypes.UnresolvedResult{},
			},
		},
		{
			name: "fails to retrieve pool for an  empty ticker",
			tickers: []mmtypes.Ticker{
				{},
			},
			client: func() uniswapv3.EVMClient {
				return mocks.NewEVMClient(t)
			},
			expected: types.PriceResponse{
				Resolved: map[mmtypes.Ticker]providertypes.ResolvedResult[*big.Int]{},
				UnResolved: map[mmtypes.Ticker]providertypes.UnresolvedResult{
					{}: {},
				},
			},
		},
		{
			name: "fails to make a batch call",
			tickers: []mmtypes.Ticker{
				weth_usdc_ticker,
			},
			client: func() uniswapv3.EVMClient {
				return createEVMClientWithResponse(t, fmt.Errorf("failed to make a batch call"), nil, nil)
			},
			expected: types.PriceResponse{
				Resolved: map[mmtypes.Ticker]providertypes.ResolvedResult[*big.Int]{},
				UnResolved: map[mmtypes.Ticker]providertypes.UnresolvedResult{
					weth_usdc_ticker: {},
				},
			},
		},
		{
			name: "batch request has an error for a single ticker",
			tickers: []mmtypes.Ticker{
				weth_usdc_ticker,
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
				Resolved: map[mmtypes.Ticker]providertypes.ResolvedResult[*big.Int]{},
				UnResolved: map[mmtypes.Ticker]providertypes.UnresolvedResult{
					weth_usdc_ticker: {},
				},
			},
		},
		{
			name: "batch request returns a result that cannot be parsed",
			tickers: []mmtypes.Ticker{
				weth_usdc_ticker,
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
				Resolved: map[mmtypes.Ticker]providertypes.ResolvedResult[*big.Int]{},
				UnResolved: map[mmtypes.Ticker]providertypes.UnresolvedResult{
					weth_usdc_ticker: {},
				},
			},
		},
		{
			name: "weth/usdc mainnet result",
			tickers: []mmtypes.Ticker{
				weth_usdc_ticker,
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
				Resolved: map[mmtypes.Ticker]providertypes.ResolvedResult[*big.Int]{
					weth_usdc_ticker: {
						Value: func() *big.Int {
							v, ok := new(big.Float).SetString("3.313131879703878971626114658316303e+21")
							require.True(t, ok)
							i, _ := v.Int(nil)
							return i
						}(),
					},
				},
				UnResolved: map[mmtypes.Ticker]providertypes.UnresolvedResult{},
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
				math.VerifyPrice(t, result.Value, response.Resolved[ticker].Value, acceptableDelta)
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
		ticker := mmtypes.Ticker{}
		_, err := fetcher.GetPool(ticker)
		require.Error(t, err)
	})

	t.Run("ticker does not have valid metadata", func(t *testing.T) {
		expected := uniswapv3.PoolConfig{
			Address: "0x1234",
		}
		ticker := mmtypes.Ticker{
			Metadata_JSON: expected.MustToJSON(),
		}
		_, err := fetcher.GetPool(ticker)
		require.Error(t, err)
	})

	t.Run("ticker is not json formatted", func(t *testing.T) {
		ticker := mmtypes.Ticker{
			Metadata_JSON: "not json, just a string",
		}
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
		ticker := mmtypes.Ticker{
			Metadata_JSON: expected.MustToJSON(),
		}
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
