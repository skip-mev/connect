package oracle_test

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/providers/base"
	basemocks "github.com/skip-mev/slinky/providers/base/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
	providermocks "github.com/skip-mev/slinky/providers/types/mocks"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/mock"
)

func (s *OracleTestSuite) TestProviders() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testCases := []struct {
		name           string
		factory        providertypes.ProviderFactory[oracletypes.CurrencyPair, *big.Int]
		expectedPrices map[oracletypes.CurrencyPair]*big.Int
	}{
		{
			name: "no providers",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				return nil, nil
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
		{
			name: "1 provider with no prices",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				// Create the provider.
				handler := basemocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](s.T())
				handler.On("Get", mock.Anything).Return(nil, nil).Maybe()

				providerCfg := config.ProviderConfig{
					Interval: 250 * time.Millisecond,
					Name:     "provider1",
				}
				provider, err := base.NewProvider(
					s.logger,
					providerCfg,
					handler,
				)
				s.Require().NoError(err)

				return []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider}, nil
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
		{
			name: "1 provider with prices",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				// Create the provider.
				handler := basemocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](s.T())
				handler.On("Get", mock.Anything).Return(
					map[oracletypes.CurrencyPair]*big.Int{
						s.currencyPairs[0]: big.NewInt(100),
					},
					nil,
				).Maybe()

				providerCfg := config.ProviderConfig{
					Interval: 250 * time.Millisecond,
					Name:     "provider1",
				}
				provider, err := base.NewProvider(
					s.logger,
					providerCfg,
					handler,
				)
				s.Require().NoError(err)

				return []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider}, nil
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{
				s.currencyPairs[0]: big.NewInt(100),
			},
		},
		{
			name: "multiple providers with prices",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				// Create the provider.
				handler := basemocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](s.T())
				handler.On("Get", mock.Anything).Return(
					map[oracletypes.CurrencyPair]*big.Int{
						s.currencyPairs[0]: big.NewInt(100),
					},
					nil,
				).Maybe()

				providerCfg := config.ProviderConfig{
					Interval: 250 * time.Millisecond,
					Name:     "provider1",
				}
				provider1, err := base.NewProvider(
					s.logger,
					providerCfg,
					handler,
				)
				s.Require().NoError(err)

				// Create the provider.
				handler = basemocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](s.T())
				handler.On("Get", mock.Anything).Return(
					map[oracletypes.CurrencyPair]*big.Int{
						s.currencyPairs[0]: big.NewInt(200),
					},
					nil,
				).Maybe()

				providerCfg = config.ProviderConfig{
					Interval: 250 * time.Millisecond,
					Name:     "provider2",
				}
				provider2, err := base.NewProvider(
					s.logger,
					providerCfg,
					handler,
				)
				s.Require().NoError(err)

				return []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider1, provider2}, nil
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{
				s.currencyPairs[0]: big.NewInt(150),
			},
		},
		{
			name: "multiple providers with 1 provider erroring on start",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				// Create the provider.
				handler := basemocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](s.T())
				handler.On("Get", mock.Anything).Return(
					map[oracletypes.CurrencyPair]*big.Int{
						s.currencyPairs[0]: big.NewInt(100),
					},
					nil,
				).Maybe()

				providerCfg := config.ProviderConfig{
					Interval: 250 * time.Millisecond,
					Name:     "provider1",
				}
				provider1, err := base.NewProvider(
					s.logger,
					providerCfg,
					handler,
				)
				s.Require().NoError(err)

				return []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{
					provider1, s.noStartProvider("provider2"),
				}, nil
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{
				s.currencyPairs[0]: big.NewInt(100),
			},
		},
		{
			name: "1 provider with stale prices",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				// Create the provider.
				handler := basemocks.NewAPIDataHandler[oracletypes.CurrencyPair, *big.Int](s.T())
				handler.On("Get", mock.Anything).Return(
					map[oracletypes.CurrencyPair]*big.Int{
						s.currencyPairs[0]: big.NewInt(100),
					},
					nil,
				).Maybe()

				providerCfg := config.ProviderConfig{
					Interval: 10 * time.Second, // 10 seconds is greater than the max price age.
					Name:     "provider1",
				}
				provider, err := base.NewProvider(
					s.logger,
					providerCfg,
					handler,
				)
				s.Require().NoError(err)

				return []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider}, nil
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
		{
			name: "1 provider that panics on get prices",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				provider := s.panicProvider("provider1")
				return []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider}, nil
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cfg := config.OracleConfig{
				UpdateInterval: 1 * time.Second,
			}
			oracle, err := oracle.New(
				s.logger,
				cfg,
				tc.factory,
				s.aggregationFn,
				metrics.NewNopMetrics(),
			)
			s.Require().NoError(err)

			go func() {
				oracle.Start(ctx)
			}()

			// Wait for the oracle to start and update.
			time.Sleep(3 * cfg.UpdateInterval)

			// Get the prices.
			prices := oracle.GetPrices()
			s.Require().Equal(tc.expectedPrices, prices)

			// Stop the oracle.
			oracle.Stop()

			time.Sleep(2 * time.Second)

			// Ensure that the oracle is not running.
			checkFn := func() bool {
				return !oracle.IsRunning()
			}
			s.Eventually(checkFn, 5*time.Second, 100*time.Millisecond)
		})
	}
}

func (s *OracleTestSuite) noStartProvider(name string) providertypes.Provider[oracletypes.CurrencyPair, *big.Int] {
	provider := providermocks.NewProvider[oracletypes.CurrencyPair, *big.Int](s.T())

	provider.On("Name").Return(name).Maybe()
	provider.On("Start", mock.Anything).Return(fmt.Errorf("no rizz error")).Maybe()
	provider.On("LastUpdate").Return(time.Now()).Maybe()
	provider.On("GetData").Return(nil, nil).Maybe()

	return provider
}

func (s *OracleTestSuite) panicProvider(name string) providertypes.Provider[oracletypes.CurrencyPair, *big.Int] {
	provider := providermocks.NewProvider[oracletypes.CurrencyPair, *big.Int](s.T())

	provider.On("Name").Return(name).Maybe()
	provider.On("GetData").Panic("no rizz panic").Maybe()
	provider.On("Start", mock.Anything).Return(nil).Maybe()
	provider.On("LastUpdate").Return(time.Now().Add(1 * time.Hour)).Maybe()

	return provider
}
