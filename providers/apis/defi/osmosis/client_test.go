package osmosis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/apis/defi/osmosis"
	"github.com/skip-mev/connect/v2/providers/apis/defi/osmosis/mocks"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
)

// TestMultiClient tests the MultiClient.
func TestMultiClient(t *testing.T) {
	cfg := osmosis.DefaultAPIConfig
	cfg.Endpoints = []config.Endpoint{
		{
			URL: "http://localhost:8899",
		},
		{
			URL: "http://localhost:8899/",
			Authentication: config.Authentication{
				APIKey:       "test",
				APIKeyHeader: "X-API-Key",
			},
		},
		{
			URL: "http://localhost:8899/",
		},
	}

	client1 := mocks.NewClient(t)
	client2 := mocks.NewClient(t)
	client3 := mocks.NewClient(t)
	client, err := osmosis.NewMultiClient(
		zap.NewNop(),
		cfg,
		metrics.NewNopAPIMetrics(),
		[]osmosis.Client{client1, client2, client3},
	)
	require.NoError(t, err)

	t.Run("test MultiClient From endpoints", func(t *testing.T) {
		t.Run("invalid endpoint", func(t *testing.T) {
			tempCfg := cfg
			tempCfg.Endpoints = nil

			_, err := osmosis.NewMultiClientFromEndpoints(
				zap.NewNop(),
				tempCfg,
				metrics.NewNopAPIMetrics(),
			)
			require.Error(t, err)
		})

		t.Run("endpoints with / wo authentication", func(t *testing.T) {
			_, err := osmosis.NewMultiClientFromEndpoints(
				zap.NewNop(),
				cfg,
				metrics.NewNopAPIMetrics(),
			)
			require.NoError(t, err)
		})
	})

	// test adherence to the context
	t.Run("test failures an underlying client", func(t *testing.T) {
		var (
			poolID        uint64 = 1
			baseAsset            = "test1"
			quoteAsset           = "test2"
			expectedPrice        = "10"
		)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// mocks
		client1.On("SpotPrice", mock.Anything, poolID, baseAsset, quoteAsset).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{SpotPrice: expectedPrice},
		}, nil).Once()

		client2.On("SpotPrice", mock.Anything, poolID, baseAsset, quoteAsset).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{SpotPrice: expectedPrice},
		}, nil).Once()

		client3.On("SpotPrice", mock.Anything, poolID, baseAsset, quoteAsset).Return(osmosis.WrappedSpotPriceResponse{},
			fmt.Errorf("error")).Once()

		resp, err := client.SpotPrice(ctx, poolID, baseAsset, quoteAsset)
		require.NoError(t, err)

		require.Equal(t, expectedPrice, resp.SpotPrice)
	})

	// test correct aggregation of responses
	t.Run("test correct aggregation of responses", func(t *testing.T) {
		var (
			poolID        uint64 = 1
			baseAsset            = "test1"
			quoteAsset           = "test2"
			expectedPrice        = "11"
		)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		// mocks
		client1.On("SpotPrice", mock.Anything, poolID, baseAsset, quoteAsset).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{SpotPrice: expectedPrice},
		}, nil).Once()
		client2.On("SpotPrice", mock.Anything, poolID, baseAsset, quoteAsset).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{SpotPrice: expectedPrice},
		}, nil).Once()
		client3.On("SpotPrice", mock.Anything, poolID, baseAsset, quoteAsset).Return(osmosis.WrappedSpotPriceResponse{
			SpotPriceResponse: osmosis.SpotPriceResponse{SpotPrice: expectedPrice},
		}, nil).Once()

		resp, err := client.SpotPrice(ctx, poolID, baseAsset, quoteAsset)
		require.NoError(t, err)

		require.Equal(t, expectedPrice, resp.SpotPrice)
	})
}
