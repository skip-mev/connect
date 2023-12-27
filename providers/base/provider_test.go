package base_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	"github.com/skip-mev/slinky/providers/base/mocks"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	logger = zap.NewNop()
	cfg    = config.ProviderConfig{
		Name:     "test",
		Path:     "test",
		Timeout:  time.Millisecond * 50,
		Interval: time.Millisecond * 100,
	}
	pairs = []oracletypes.CurrencyPair{
		{
			Base:  "BTC",
			Quote: "USD",
		},
		{
			Base:  "ETH",
			Quote: "USD",
		},
	}
)

func TestStart(t *testing.T) {
	t.Parallel()

	t.Run("closes on cancel", func(t *testing.T) {
		handler := mocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](t)

		provider, err := base.NewProvider(logger, cfg, handler)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err = provider.Start(ctx)
		require.Equal(t, context.Canceled, err)
	})

	t.Run("closes with deadline", func(t *testing.T) {
		handler := mocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](t)
		handler.On("Get", mock.Anything).Return(nil, nil).Maybe()

		provider, err := base.NewProvider(logger, cfg, handler)
		require.NoError(t, err)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Interval*2)
		defer cancel()

		err = provider.Start(ctx)
		require.Equal(t, context.DeadlineExceeded, err)
	})
}

func TestGetData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		provider       func() *base.BaseProvider[oracletypes.CurrencyPair, *big.Int]
		expectedPrices map[oracletypes.CurrencyPair]*big.Int
		expectedUpdate bool
	}{
		{
			"no price",
			func() *base.BaseProvider[oracletypes.CurrencyPair, *big.Int] {
				handler := mocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				handler.On("Get", mock.Anything).Return(
					map[oracletypes.CurrencyPair]*big.Int{},
					nil,
				).Maybe()

				provider, err := base.NewProvider(logger, cfg, handler)
				require.NoError(t, err)

				return provider
			},
			map[oracletypes.CurrencyPair]*big.Int{},
			false,
		},
		{
			"1 price",
			func() *base.BaseProvider[oracletypes.CurrencyPair, *big.Int] {
				handler := mocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				handler.On("Get", mock.Anything).Return(
					map[oracletypes.CurrencyPair]*big.Int{
						pairs[0]: big.NewInt(100),
					},
					nil,
				).Maybe()

				provider, err := base.NewProvider(logger, cfg, handler)
				require.NoError(t, err)

				return provider
			},
			map[oracletypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
			},
			true,
		},
		{
			"multiple prices",
			func() *base.BaseProvider[oracletypes.CurrencyPair, *big.Int] {
				handler := mocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				handler.On("Get", mock.Anything).Return(
					map[oracletypes.CurrencyPair]*big.Int{
						pairs[0]: big.NewInt(100),
						pairs[1]: big.NewInt(200),
					},
					nil,
				).Maybe()

				provider, err := base.NewProvider(logger, cfg, handler)
				require.NoError(t, err)

				return provider
			},
			map[oracletypes.CurrencyPair]*big.Int{
				pairs[0]: big.NewInt(100),
				pairs[1]: big.NewInt(200),
			},
			true,
		},
		{
			"continues updating even with error",
			func() *base.BaseProvider[oracletypes.CurrencyPair, *big.Int] {
				handler := mocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				handler.On("Get", mock.Anything).Return(
					map[oracletypes.CurrencyPair]*big.Int{},
					fmt.Errorf("big oopsie"),
				).Maybe()

				provider, err := base.NewProvider(logger, cfg, handler)
				require.NoError(t, err)

				return provider
			},
			map[oracletypes.CurrencyPair]*big.Int{},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			provider := tc.provider()
			now := time.Now()

			// Start the provider with a timeout of 2x the interval. This should allow it to
			// update at least once.
			experimentTime := cfg.Interval * 2
			ctx, cancel := context.WithTimeout(context.Background(), experimentTime)
			defer cancel()

			go func() {
				err := provider.Start(ctx)
				require.Equal(t, context.DeadlineExceeded, err)
			}()

			// Sleep to allow the goroutine to close.
			time.Sleep(experimentTime * 2)

			prices := provider.GetData()
			require.Equal(t, tc.expectedPrices, prices)

			latestUpdate := provider.LastUpdate()
			if tc.expectedUpdate {
				require.True(t, latestUpdate.After(now))
			} else {
				require.False(t, latestUpdate.After(now))
			}
		})
	}
}

func TestNewProvider(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		cfg     config.ProviderConfig
		pairs   []oracletypes.CurrencyPair
		handler func() base.APIDataHandler[oracletypes.CurrencyPair, *big.Int]
		expErr  bool
	}{
		{
			"valid",
			cfg,
			pairs,
			func() base.APIDataHandler[oracletypes.CurrencyPair, *big.Int] {
				handler := mocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](t)
				return handler
			},
			false,
		},
		{
			"no handler",
			cfg,
			pairs,
			func() base.APIDataHandler[oracletypes.CurrencyPair, *big.Int] {
				return nil
			},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := tc.handler()
			provider, err := base.NewProvider(logger, tc.cfg, handler)
			if tc.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.cfg.Name, provider.Name())
			}
		})
	}
}
