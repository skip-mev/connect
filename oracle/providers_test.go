package oracle_test

import (
	"context"
	"math/big"
	"time"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	mathtestutils "github.com/skip-mev/slinky/pkg/math/testutils"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

func (s *OracleTestSuite) TestProviders() {
	testCases := []struct {
		name           string
		factory        func() []*types.PriceProvider
		expectedPrices types.Prices
	}{
		{
			name: "1 provider with no prices",
			factory: func() []*types.PriceProvider {
				provider := testutils.CreateAPIProviderWithGetResponses[types.ProviderTicker, *big.Float](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
					200*time.Millisecond,
				)

				providers := []*types.PriceProvider{provider}
				return providers
			},
			expectedPrices: types.Prices{},
		},
		{
			name: "1 provider with prices",
			factory: func() []*types.PriceProvider {
				resolved := types.ResolvedPrices{
					s.currencyPairs[0]: {
						Value:     big.NewFloat(100),
						Timestamp: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response := providertypes.NewGetResponse[types.ProviderTicker, *big.Float](resolved, nil)
				responses := []providertypes.GetResponse[types.ProviderTicker, *big.Float]{response}
				provider := testutils.CreateAPIProviderWithGetResponses[types.ProviderTicker, *big.Float](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					responses,
					200*time.Millisecond,
				)

				providers := []*types.PriceProvider{provider}
				return providers
			},
			expectedPrices: types.Prices{
				s.currencyPairs[0].String(): big.NewFloat(100),
			},
		},
		{
			name: "multiple providers with prices",
			factory: func() []*types.PriceProvider {
				resolved := types.ResolvedPrices{
					s.currencyPairs[0]: {
						Value:     big.NewFloat(100),
						Timestamp: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response := providertypes.NewGetResponse[types.ProviderTicker, *big.Float](resolved, nil)
				responses := []providertypes.GetResponse[types.ProviderTicker, *big.Float]{response}
				provider := testutils.CreateAPIProviderWithGetResponses[types.ProviderTicker, *big.Float](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					responses,
					200*time.Millisecond,
				)

				resolved2 := types.ResolvedPrices{
					s.currencyPairs[0]: {
						Value:     big.NewFloat(200),
						Timestamp: time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response2 := providertypes.NewGetResponse[types.ProviderTicker, *big.Float](resolved2, nil)
				responses2 := []providertypes.GetResponse[types.ProviderTicker, *big.Float]{response2}
				provider2 := testutils.CreateWebSocketProviderWithGetResponses[types.ProviderTicker, *big.Float](
					s.T(),
					time.Second*2,
					s.currencyPairs,
					providerCfg2,
					s.logger,
					responses2,
				)

				providers := []*types.PriceProvider{provider, provider2}
				return providers
			},
			expectedPrices: types.Prices{
				s.currencyPairs[0].String(): big.NewFloat(150),
			},
		},
		{
			name: "1 provider with stale prices",
			factory: func() []*types.PriceProvider {
				resolved := types.ResolvedPrices{
					s.currencyPairs[0]: {
						Value:     big.NewFloat(100),
						Timestamp: time.Date(1738, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}
				response := providertypes.NewGetResponse[types.ProviderTicker, *big.Float](resolved, nil)
				responses := []providertypes.GetResponse[types.ProviderTicker, *big.Float]{response}
				provider := testutils.CreateAPIProviderWithGetResponses[types.ProviderTicker, *big.Float](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					responses,
					200*time.Millisecond,
				)

				providers := []*types.PriceProvider{provider}
				return providers
			},
			expectedPrices: types.Prices{},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cfg := config.OracleConfig{
				UpdateInterval: 1 * time.Second,
			}

			providers := tc.factory()
			ctx, cancel := context.WithTimeout(context.Background(), 4*cfg.UpdateInterval)
			defer cancel()

			for _, provider := range providers {
				p := provider
				go func() {
					// context deadline exceeded
					s.Require().Error(p.Start(ctx))
				}()
			}

			testOracle, err := oracle.New(
				oracle.WithUpdateInterval(cfg.UpdateInterval),
				oracle.WithLogger(s.logger),
				oracle.WithProviders(providers),
				oracle.WithPriceAggregator(mathtestutils.NewMedianAggregator()),
			)
			s.Require().NoError(err)

			go func() {
				s.Require().NoError(testOracle.Start(ctx))
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
