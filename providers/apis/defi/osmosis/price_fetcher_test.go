package osmosis_test

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	oracleconfig "github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/providers/apis/defi/osmosis"
	"github.com/skip-mev/connect/v2/providers/apis/defi/osmosis/mocks"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
)

const (
	USDTTokenDenom        = "USDT"
	BTCTokenDenom         = "BTC"
	BTCUSDTPoolID  uint64 = 1
	ETHTokenDenom         = "ETH"
	ETHUSDTPoolID  uint64 = 10
	MOGTokenDenom         = "MOG"
	SOLTokenDenom         = "SOL"
	MOGSOLPoolID   uint64 = 11
)

func TestTickerMetadataValidateBasic(t *testing.T) {
	tcs := []struct {
		name string
		osmosis.TickerMetadata
		expFail bool
	}{
		{
			name: "invalid base token denom",
			TickerMetadata: osmosis.TickerMetadata{
				PoolID:          ETHUSDTPoolID,
				BaseTokenDenom:  "",
				QuoteTokenDenom: USDTTokenDenom,
			},
			expFail: true,
		},
		{
			name: "invalid quote token denom",
			TickerMetadata: osmosis.TickerMetadata{
				PoolID:          ETHUSDTPoolID,
				BaseTokenDenom:  ETHTokenDenom,
				QuoteTokenDenom: "",
			},
			expFail: true,
		},
		{
			name: "valid",
			TickerMetadata: osmosis.TickerMetadata{
				PoolID:          ETHUSDTPoolID,
				BaseTokenDenom:  ETHTokenDenom,
				QuoteTokenDenom: USDTTokenDenom,
			},
			expFail: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.TickerMetadata.ValidateBasic()
			if tc.expFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test Provider init.
func TestProviderInit(t *testing.T) {
	t.Run("config fails validate basic", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:    true,
			MaxQueries: 0,
		}

		_, err := osmosis.NewAPIPriceFetcher(
			zap.NewNop(),
			cfg,
			metrics.NewNopAPIMetrics(),
		)

		require.Error(t, err)
	})

	t.Run("config has invalid endpoints", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:    true,
			MaxQueries: 0,
			Endpoints: []oracleconfig.Endpoint{
				{
					URL: "", // invalid url
				},
				{
					URL: "https://osmosis.io",
				},
			},
		}

		_, err := osmosis.NewAPIPriceFetcher(
			zap.NewNop(),
			cfg,
			metrics.NewNopAPIMetrics(),
		)

		require.Error(t, err)
	})

	t.Run("incorrect provider name", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:          true,
			MaxQueries:       2,
			Interval:         1 * time.Second,
			Timeout:          2 * time.Second,
			ReconnectTimeout: 2 * time.Second,
			Endpoints:        []oracleconfig.Endpoint{{URL: "https://osmosis.io"}},
			Name:             osmosis.Name + "a",
		}

		_, err := osmosis.NewAPIPriceFetcher(
			zap.NewNop(),
			cfg,
			metrics.NewNopAPIMetrics(),
		)
		require.Error(t, err)
	})

	t.Run("api not enabled", func(t *testing.T) {
		cfg := oracleconfig.APIConfig{
			Enabled:          false,
			MaxQueries:       2,
			Interval:         1 * time.Second,
			Timeout:          2 * time.Second,
			ReconnectTimeout: 2 * time.Second,
			Name:             osmosis.Name,
		}

		_, err := osmosis.NewAPIPriceFetcher(
			zap.NewNop(),
			cfg,
			metrics.NewNopAPIMetrics(),
		)
		require.Error(t, err, "config is not enabled")
	})
}

// Test getting prices.
func TestProviderFetch(t *testing.T) {
	btcUSDTMetadata := osmosis.TickerMetadata{
		PoolID:          BTCUSDTPoolID,
		BaseTokenDenom:  BTCTokenDenom,
		QuoteTokenDenom: USDTTokenDenom,
	}
	ethUSDTMetadata := osmosis.TickerMetadata{
		PoolID:          ETHUSDTPoolID,
		BaseTokenDenom:  ETHTokenDenom,
		QuoteTokenDenom: USDTTokenDenom,
	}
	mogSOLMetadata := osmosis.TickerMetadata{
		PoolID:          MOGSOLPoolID,
		BaseTokenDenom:  MOGTokenDenom,
		QuoteTokenDenom: SOLTokenDenom,
	}

	var (
		expectedBTCUSDTPrice = "10"
		expectedETHUSDTPrice = "11"
		expectedMOGSOLPRICE  = "12"
	)

	tickers := []types.DefaultProviderTicker{
		{
			OffChainTicker: "BTC/USDC",
			JSON:           marshalDataToJSON(btcUSDTMetadata),
		},
		{
			OffChainTicker: "ETH/USDT",
			JSON:           marshalDataToJSON(ethUSDTMetadata),
		},
		{
			OffChainTicker: "MOG/SOL",
			JSON:           marshalDataToJSON(mogSOLMetadata),
		},
	}

	t.Run("single valid ticker", func(t *testing.T) {
		client := mocks.NewClient(t)
		pf, err := newPriceFetcher(client)
		require.NoError(t, err)

		ctx := context.Background()

		client.On("SpotPrice", mock.Anything, btcUSDTMetadata.PoolID, btcUSDTMetadata.BaseTokenDenom,
			btcUSDTMetadata.QuoteTokenDenom,
		).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{
				SpotPrice: expectedBTCUSDTPrice,
			},
		}, nil).Once()

		ts := defaultTickersToProviderTickers([]types.DefaultProviderTicker{tickers[0]})
		resp := pf.Fetch(ctx, ts)
		// expect a failed response
		require.Equal(t, 1, len(resp.Resolved))
		require.Equal(t, 0, len(resp.UnResolved))
	})

	t.Run("failing query", func(t *testing.T) {
		client := mocks.NewClient(t)
		pf, err := newPriceFetcher(client)
		require.NoError(t, err)

		ctx := context.Background()

		err = fmt.Errorf("error")

		client.On("SpotPrice", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{
				SpotPrice: "",
			},
		}, err).Times(3)

		ts := defaultTickersToProviderTickers(tickers)
		resp := pf.Fetch(ctx, ts)
		// expect a failed response
		require.Equal(t, 0, len(resp.Resolved))
		require.Equal(t, 3, len(resp.UnResolved))

		for _, result := range resp.UnResolved {
			require.True(t, strings.Contains(result.Error(), err.Error()))
		}
	})

	t.Run("unexpected ticker in query", func(t *testing.T) {
		client := mocks.NewClient(t)
		pf, err := newPriceFetcher(client)
		require.NoError(t, err)

		ctx := context.Background()

		mogtia := types.DefaultProviderTicker{
			OffChainTicker: "MOG/TIA",
			JSON:           "{}",
		}
		resp := pf.Fetch(ctx, []types.ProviderTicker{
			mogtia,
		})
		// expect a failed response
		require.Equal(t, len(resp.Resolved), 0)
		require.Equal(t, len(resp.UnResolved), 1)

		for _, result := range resp.UnResolved {
			t.Log(result.Error())
			require.True(t, strings.Contains(result.Error(), osmosis.NoOsmosisMetadataForTickerError("MOG/TIA").Error()))
		}
	})

	t.Run("multi-asset one failing", func(t *testing.T) {
		client := mocks.NewClient(t)
		pf, err := newPriceFetcher(client)
		require.NoError(t, err)

		ctx := context.Background()

		err = fmt.Errorf("error")

		client.On("SpotPrice", mock.Anything, btcUSDTMetadata.PoolID, btcUSDTMetadata.BaseTokenDenom,
			btcUSDTMetadata.QuoteTokenDenom,
		).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{SpotPrice: ""},
		}, err).Once()

		client.On("SpotPrice", mock.Anything, ethUSDTMetadata.PoolID, ethUSDTMetadata.BaseTokenDenom,
			ethUSDTMetadata.QuoteTokenDenom,
		).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{SpotPrice: expectedETHUSDTPrice},
		}, nil).Once()

		client.On("SpotPrice", mock.Anything, mogSOLMetadata.PoolID, mogSOLMetadata.BaseTokenDenom,
			mogSOLMetadata.QuoteTokenDenom,
		).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{SpotPrice: expectedMOGSOLPRICE},
		}, nil).Once()

		ts := defaultTickersToProviderTickers(tickers)
		resp := pf.Fetch(ctx, ts)

		// expect a failed response
		require.Equal(t, 2, len(resp.Resolved))
		require.Equal(t, 1, len(resp.UnResolved))

		for _, result := range resp.UnResolved {
			require.True(t, strings.Contains(result.Error(), err.Error()))
		}
	})

	t.Run("multi-asset success", func(t *testing.T) {
		client := mocks.NewClient(t)
		pf, err := newPriceFetcher(client)
		require.NoError(t, err)

		ctx := context.Background()

		client.On("SpotPrice", mock.Anything, btcUSDTMetadata.PoolID, btcUSDTMetadata.BaseTokenDenom,
			btcUSDTMetadata.QuoteTokenDenom,
		).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{SpotPrice: expectedBTCUSDTPrice},
		}, nil).Once()

		client.On("SpotPrice", mock.Anything, ethUSDTMetadata.PoolID, ethUSDTMetadata.BaseTokenDenom,
			ethUSDTMetadata.QuoteTokenDenom,
		).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{SpotPrice: expectedETHUSDTPrice},
		}, nil).Once()

		client.On("SpotPrice", mock.Anything, mogSOLMetadata.PoolID, mogSOLMetadata.BaseTokenDenom,
			mogSOLMetadata.QuoteTokenDenom,
		).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{SpotPrice: expectedMOGSOLPRICE},
		}, nil).Once()

		ts := defaultTickersToProviderTickers(tickers)
		resp := pf.Fetch(ctx, ts)

		// expect a failed response
		require.Equal(t, 3, len(resp.Resolved))
		require.Equal(t, 0, len(resp.UnResolved))
	})
}

func marshalDataToJSON(obj interface{}) string {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func newPriceFetcher(client *mocks.Client) (*osmosis.APIPriceFetcher, error) {
	cfg := oracleconfig.APIConfig{
		Enabled:          true,
		MaxQueries:       2,
		Interval:         1 * time.Second,
		Timeout:          2 * time.Second,
		ReconnectTimeout: 2 * time.Second,
		Name:             osmosis.Name,
		Endpoints:        []oracleconfig.Endpoint{{URL: "https://osmosis.io"}},
	}

	return osmosis.NewAPIPriceFetcherWithClient(
		zap.NewExample(),
		cfg,
		metrics.NewNopAPIMetrics(),
		client,
	)
}

func defaultTickersToProviderTickers(tickers []types.DefaultProviderTicker) []types.ProviderTicker {
	providerTickers := make([]types.ProviderTicker, len(tickers))
	for i, ticker := range tickers {
		providerTickers[i] = ticker
	}
	return providerTickers
}
