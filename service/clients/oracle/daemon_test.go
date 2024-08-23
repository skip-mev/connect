package oracle_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/service/clients/oracle"
	"github.com/skip-mev/connect/v2/service/clients/oracle/mocks"
	"github.com/skip-mev/connect/v2/service/servers/oracle/types"
)

func TestNewPriceDaemon(t *testing.T) {
	testCases := []struct {
		name   string
		logger log.Logger
		cfg    config.AppConfig
		client oracle.OracleClient
		err    bool
	}{
		{
			name:   "valid",
			logger: log.NewNopLogger(),
			cfg: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				ClientTimeout: time.Second,
				Interval:      time.Second,
				PriceTTL:      time.Second * 2,
			},
			client: &oracle.NoOpClient{},
			err:    false,
		},
		{
			name:   "nil logger",
			logger: nil,
			cfg: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				ClientTimeout: time.Second,
				Interval:      time.Second,
				PriceTTL:      time.Second * 2,
			},
			client: &oracle.NoOpClient{},
			err:    true,
		},
		{
			name:   "invalid config",
			logger: log.NewNopLogger(),
			cfg: config.AppConfig{
				Enabled: true,
			},
			client: &oracle.NoOpClient{},
			err:    true,
		},
		{
			name:   "nil client",
			logger: log.NewNopLogger(),
			cfg: config.AppConfig{
				Enabled:       true,
				OracleAddress: "localhost:8080",
				ClientTimeout: time.Second,
				Interval:      time.Second,
				PriceTTL:      time.Second * 2,
			},
			client: nil,
			err:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := oracle.NewPriceDaemon(tc.logger, tc.cfg, tc.client)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestPriceDaemon_Start(t *testing.T) {
	logger := log.NewTestLogger(t)
	cfg := config.AppConfig{
		Enabled:       true,
		OracleAddress: "localhost:8080",
		ClientTimeout: time.Second,
		Interval:      time.Millisecond * 100,
		PriceTTL:      time.Second,
	}

	t.Run("stops with context cancel", func(t *testing.T) {
		client := mocks.NewOracleClient(t)
		client.On("Start", mock.Anything).Return(nil).Once()
		client.On("Prices", mock.Anything, mock.Anything).Return(&types.QueryPricesResponse{}, nil).Maybe()
		client.On("Stop").Return(nil).Once()

		d, err := oracle.NewPriceDaemon(logger, cfg, client)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(time.Millisecond * 300)
			cancel()
		}()

		err = d.Start(ctx)
		require.Equal(t, err, context.Canceled)
	})

	t.Run("can correctly store the last result", func(t *testing.T) {
		prices := map[string]string{
			"btc/usd": "10000",
		}

		client := mocks.NewOracleClient(t)
		client.On("Start", mock.Anything).Return(nil).Once()
		client.On("Prices", mock.Anything, mock.Anything).Return(&types.QueryPricesResponse{
			Prices: prices,
		}, nil).Maybe()
		client.On("Stop").Return(nil).Once()

		d, err := oracle.NewPriceDaemon(logger, cfg, client)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(time.Millisecond * 300)
			cancel()
		}()

		err = d.Start(ctx)
		require.Equal(t, err, context.Canceled)

		time.Sleep(cfg.Interval * 2)
		resp, err := d.Prices(context.Background(), &types.QueryPricesRequest{})
		require.NoError(t, err)
		require.Equal(t, prices, resp.Prices)
	})

	t.Run("will return an error if the latest response is too old", func(t *testing.T) {
		prices := map[string]string{
			"btc/usd": "10000",
		}
		cfg := config.AppConfig{
			Enabled:       true,
			OracleAddress: "localhost:8080",
			ClientTimeout: time.Second,
			PriceTTL:      time.Millisecond * 250,
			Interval:      time.Millisecond * 100,
		}

		client := mocks.NewOracleClient(t)
		client.On("Start", mock.Anything).Return(nil).Once()
		client.On("Prices", mock.Anything, mock.Anything).Return(&types.QueryPricesResponse{Prices: prices}, nil).Once()
		client.On("Prices", mock.Anything, mock.Anything).Return(&types.QueryPricesResponse{}, fmt.Errorf("failed to make request")).Maybe()
		client.On("Stop").Return(nil).Once()

		d, err := oracle.NewPriceDaemon(logger, cfg, client)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(time.Millisecond * 150)
			cancel()
		}()

		err = d.Start(ctx)
		require.Equal(t, err, context.Canceled)

		resp, err := d.Prices(context.Background(), &types.QueryPricesRequest{})
		require.NoError(t, err)
		require.Equal(t, prices, resp.Prices)

		time.Sleep(cfg.PriceTTL * 2)
		resp, err = d.Prices(context.Background(), &types.QueryPricesRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("returns an error if it never started", func(t *testing.T) {
		client := mocks.NewOracleClient(t)
		d, err := oracle.NewPriceDaemon(logger, cfg, client)
		require.NoError(t, err)

		resp, err := d.Prices(context.Background(), &types.QueryPricesRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("stops after channel receive", func(t *testing.T) {
		client := mocks.NewOracleClient(t)
		client.On("Start", mock.Anything).Return(nil).Once()
		client.On("Prices", mock.Anything, mock.Anything).Return(&types.QueryPricesResponse{}, nil).Maybe()
		client.On("Stop").Return(nil).Once()

		d, err := oracle.NewPriceDaemon(logger, cfg, client)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			time.Sleep(time.Millisecond * 300)
			d.Stop()
		}()

		err = d.Start(ctx)
		require.NoError(t, err)
	})

	t.Run("client only returns errors", func(t *testing.T) {
		client := mocks.NewOracleClient(t)
		client.On("Start", mock.Anything).Return(nil).Once()
		client.On("Prices", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to make request")).Maybe()
		client.On("Stop").Return(nil).Once()

		d, err := oracle.NewPriceDaemon(logger, cfg, client)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(time.Millisecond * 300)
			cancel()
		}()

		err = d.Start(ctx)
		require.Error(t, err)

		resp, err := d.Prices(context.Background(), &types.QueryPricesRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
