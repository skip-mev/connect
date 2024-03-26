package oracle_test

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
	providermocks "github.com/skip-mev/slinky/providers/types/mocks"
	mmtypes "github.com/skip-mev/slinky/x/mm2/types"
)

func (s *OracleTestSuite) TestProviders() {
	testCases := []struct {
		name           string
		factory        types.PriceProviderFactoryI
		expectedPrices types.TickerPrices
	}{
		{
			name: "1 provider with no prices",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProviderI, error) {
				provider := testutils.CreateAPIProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
					200*time.Millisecond,
				)

				providers := []types.PriceProviderI{provider}
				return providers, nil
			},
			expectedPrices: types.TickerPrices{},
		},
		{
			name: "1 provider with prices",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProviderI, error) {
				resolved := types.ResolvedPrices{
					s.currencyPairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response := providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, nil)
				responses := []providertypes.GetResponse[mmtypes.Ticker, *big.Int]{response}
				provider := testutils.CreateAPIProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					responses,
					200*time.Millisecond,
				)

				providers := []types.PriceProviderI{provider}
				return providers, nil
			},
			expectedPrices: types.TickerPrices{
				s.currencyPairs[0]: big.NewInt(100),
			},
		},
		{
			name: "multiple providers with prices",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProviderI, error) {
				resolved := types.ResolvedPrices{
					s.currencyPairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response := providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, nil)
				responses := []providertypes.GetResponse[mmtypes.Ticker, *big.Int]{response}
				provider := testutils.CreateAPIProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					responses,
					200*time.Millisecond,
				)

				resolved2 := types.ResolvedPrices{
					s.currencyPairs[0]: {
						Value:     big.NewInt(200),
						Timestamp: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response2 := providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved2, nil)
				responses2 := []providertypes.GetResponse[mmtypes.Ticker, *big.Int]{response2}
				provider2 := testutils.CreateWebSocketProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					time.Second*2,
					s.currencyPairs,
					providerCfg2,
					s.logger,
					responses2,
				)

				providers := []types.PriceProviderI{provider, provider2}
				return providers, nil
			},
			expectedPrices: types.TickerPrices{
				s.currencyPairs[0]: big.NewInt(150),
			},
		},
		{
			name: "multiple providers with 1 provider erroring on start",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProviderI, error) {
				resolved := types.ResolvedPrices{
					s.currencyPairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response := providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, nil)
				responses := []providertypes.GetResponse[mmtypes.Ticker, *big.Int]{response}
				provider := testutils.CreateAPIProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					responses,
					200*time.Millisecond,
				)

				providers := []types.PriceProviderI{provider, s.noStartProvider("provider2")}
				return providers, nil
			},
			expectedPrices: types.TickerPrices{
				s.currencyPairs[0]: big.NewInt(100),
			},
		},
		{
			name: "1 provider with stale prices",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProviderI, error) {
				resolved := types.ResolvedPrices{
					s.currencyPairs[0]: {
						Value:     big.NewInt(100),
						Timestamp: time.Date(1738, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response := providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, nil)
				responses := []providertypes.GetResponse[mmtypes.Ticker, *big.Int]{response}
				provider := testutils.CreateAPIProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					responses,
					200*time.Millisecond,
				)

				providers := []types.PriceProviderI{provider}
				return providers, nil
			},
			expectedPrices: types.TickerPrices{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cfg := config.OracleConfig{
				UpdateInterval: 1 * time.Second,
			}

			providers, err := tc.factory(cfg)
			s.Require().NoError(err)

			ctx, cancel := context.WithTimeout(context.Background(), 4*cfg.UpdateInterval)
			defer cancel()

			for _, provider := range providers {
				go func() {
					provider.Start(ctx)
				}()
			}

			testOracle, err := oracle.New(
				oracle.WithUpdateInterval(cfg.UpdateInterval),
				oracle.WithLogger(s.logger),
				oracle.WithProviders(providers),
			)
			s.Require().NoError(err)

			go func() {
				testOracle.Start(ctx)
			}()

			// Wait for the oracle to start and update.
			time.Sleep(2 * cfg.UpdateInterval)

			// Get the prices.
			prices := testOracle.GetPrices()
			s.Require().Equal(tc.expectedPrices, prices)

			// Stop the oracle.
			testOracle.Stop()

			time.Sleep(5 * cfg.UpdateInterval)

			// Ensure that the oracle is not running.
			checkFn := func() bool {
				return !testOracle.IsRunning()
			}
			s.Eventually(checkFn, 5*time.Second, 100*time.Millisecond)

			// Ensure that the providers are not running.
			for _, provider := range providers {
				providerCheckFn := func() bool {
					return !provider.IsRunning()
				}
				s.Eventually(providerCheckFn, 5*time.Second, 100*time.Millisecond)
			}
		})
	}
}

func (s *OracleTestSuite) noStartProvider(name string) types.PriceProviderI {
	provider := providermocks.NewProvider[mmtypes.Ticker, *big.Int](s.T())

	provider.On("Name").Return(name).Maybe()
	provider.On("Start", mock.Anything).Return(fmt.Errorf("no rizz start")).Maybe()
	provider.On("GetData").Return(nil).Maybe()
	provider.On("Type").Return(providertypes.API).Maybe()
	provider.On("IsRunning").Return(false).Maybe()

	return provider
}
