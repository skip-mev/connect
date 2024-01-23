package oracle_test

import (
	"context"
	"fmt"
	"math/big"
	"time"

	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"

	"go.uber.org/zap"

	"github.com/stretchr/testify/mock"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
	providermocks "github.com/skip-mev/slinky/providers/types/mocks"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
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
				config.OracleMetricsConfig,
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
				config.OracleMetricsConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				provider := testutils.CreateAPIProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				providers := []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider}
				return providers, nil
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
		{
			name: "1 provider with prices",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
				config.OracleMetricsConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					s.currencyPairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)
				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{response}
				provider := testutils.CreateAPIProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					responses,
				)

				providers := []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider}
				return providers, nil
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
				config.OracleMetricsConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					s.currencyPairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)
				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{response}
				provider := testutils.CreateAPIProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					responses,
				)

				resolved2 := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					s.currencyPairs[0]: {
						Value:     big.NewInt(200),
						Timestamp: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response2 := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved2, nil)
				responses2 := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{response2}
				provider2 := testutils.CreateWebSocketProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					time.Second*2,
					providerCfg2,
					s.logger,
					responses2,
				)

				providers := []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider, provider2}
				return providers, nil
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
				config.OracleMetricsConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					s.currencyPairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)
				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{response}
				provider := testutils.CreateAPIProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					responses,
				)

				providers := []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider, s.noStartProvider("provider2")}
				return providers, nil
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
				config.OracleMetricsConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					s.currencyPairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: time.Date(1738, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)
				responses := []providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{response}
				provider := testutils.CreateAPIProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					responses,
				)

				providers := []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider}
				return providers, nil
			},
			expectedPrices: map[oracletypes.CurrencyPair]*big.Int{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cfg := config.OracleConfig{
				UpdateInterval: 1 * time.Second,
				InProcess:      true,
				ClientTimeout:  1 * time.Second,
			}
			metricsCfg := config.OracleMetricsConfig{
				Enabled: false,
			}

			providers, err := tc.factory(s.logger, cfg, metricsCfg)
			s.Require().NoError(err)

			testOracle, err := oracle.New(
				cfg,
				oracle.WithLogger(s.logger),
				oracle.WithProviders(providers),
			)
			s.Require().NoError(err)

			go func() {
				s.Require().NoError(testOracle.Start(ctx))
			}()

			// Wait for the oracle to start and update.
			time.Sleep(3 * cfg.UpdateInterval)

			// Get the prices.
			prices := testOracle.GetPrices()
			s.Require().Equal(tc.expectedPrices, prices)

			// Stop the oracle.
			testOracle.Stop()

			// Ensure that the oracle is not running.
			checkFn := func() bool {
				return !testOracle.IsRunning()
			}
			s.Eventually(checkFn, 5*time.Second, 100*time.Millisecond)
		})
	}
}

func (s *OracleTestSuite) noStartProvider(name string) providertypes.Provider[oracletypes.CurrencyPair, *big.Int] {
	provider := providermocks.NewProvider[oracletypes.CurrencyPair, *big.Int](s.T())

	provider.On("Name").Return(name).Maybe()
	provider.On("Start", mock.Anything).Return(fmt.Errorf("no rizz error")).Maybe()
	provider.On("GetData").Return(nil).Maybe()
	provider.On("Type").Return(providermetrics.API)

	return provider
}
