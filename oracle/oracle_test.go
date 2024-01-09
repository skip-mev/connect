package oracle_test

import (
	"context"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	providerCfg1 = config.ProviderConfig{
		Interval:   500 * time.Millisecond,
		Timeout:    250 * time.Millisecond,
		Path:       "testpath",
		Name:       "provider1",
		MaxQueries: 10,
	}
	providerCfg2 = config.ProviderConfig{
		Interval:   500 * time.Millisecond,
		Timeout:    250 * time.Millisecond,
		Path:       "testpath2",
		Name:       "provider2",
		MaxQueries: 10,
	}
)

type OracleTestSuite struct {
	suite.Suite
	random *rand.Rand

	logger *zap.Logger

	// Oracle config
	currencyPairs []oracletypes.CurrencyPair
	aggregationFn aggregator.AggregateFn[string, map[oracletypes.CurrencyPair]*big.Int]
}

func TestOracleSuite(t *testing.T) {
	suite.Run(t, new(OracleTestSuite))
}

func (s *OracleTestSuite) SetupTest() {
	s.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	s.logger = zap.NewExample()

	s.currencyPairs = []oracletypes.CurrencyPair{
		oracletypes.NewCurrencyPair("BITCOIN", "USD"),
		oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
		oracletypes.NewCurrencyPair("COSMOS", "USD"),
	}
	s.aggregationFn = aggregator.ComputeMedian()
}

func (s *OracleTestSuite) TestStopWithContextCancel() {
	testCases := []struct {
		name    string
		factory providertypes.ProviderFactory[oracletypes.CurrencyPair, *big.Int]
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
		},
		{
			name: "1 provider",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
				config.OracleMetricsConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				provider := testutils.CreateProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				// Create the provider factory.
				providers := []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider}
				return providers, nil
			},
		},
		{
			name: "multiple providers",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
				config.OracleMetricsConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				provider1 := testutils.CreateProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				provider2 := testutils.CreateProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg2,
					s.currencyPairs,
					nil,
				)

				// Create the provider factory.
				providers := []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider1, provider2}
				return providers, nil
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cfg := config.OracleConfig{
				UpdateInterval: 1 * time.Second,
				ClientTimeout:  1 * time.Second,
				InProcess:      true,
			}
			metricsCfg := config.OracleMetricsConfig{
				Enabled: false,
			}

			providers, err := tc.factory(s.logger, cfg, metricsCfg)
			s.Require().NoError(err)

			oracle, err := oracle.New(
				cfg,
				oracle.WithLogger(s.logger),
				oracle.WithProviders(providers),
			)
			s.Require().NoError(err)

			ctx, cancel := context.WithCancel(context.Background())

			// Start the oracle. This should automatically stop.
			go func() {
				err = oracle.Start(ctx)
				s.Require().Equal(err, context.Canceled)
			}()

			// Wait for the experiment to run.
			time.Sleep(2 * time.Second)
			cancel()

			// Ensure that the oracle is not running.
			s.Eventually(checkFn(oracle), 3*time.Second, 100*time.Millisecond)
		})
	}
}

func (s *OracleTestSuite) TestStopWithContextDeadline() {
	testCases := []struct {
		name     string
		factory  providertypes.ProviderFactory[oracletypes.CurrencyPair, *big.Int]
		duration time.Duration
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
			duration: 1 * time.Second,
		},
		{
			name: "1 provider",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
				config.OracleMetricsConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				provider := testutils.CreateProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				// Create the provider factory.
				providers := []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider}
				return providers, nil
			},
			duration: 1 * time.Second,
		},
		{
			name: "multiple providers",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
				config.OracleMetricsConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				provider1 := testutils.CreateProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				provider2 := testutils.CreateProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg2,
					s.currencyPairs,
					nil,
				)

				// Create the provider factory.
				providers := []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider1, provider2}
				return providers, nil
			},
			duration: 1 * time.Second,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cfg := config.OracleConfig{
				UpdateInterval: 1 * time.Second,
				ClientTimeout:  1 * time.Second,
				InProcess:      true,
			}
			metricsCfg := config.OracleMetricsConfig{
				Enabled: false,
			}

			providers, err := tc.factory(s.logger, cfg, metricsCfg)
			s.Require().NoError(err)

			oracle, err := oracle.New(
				cfg,
				oracle.WithLogger(s.logger),
				oracle.WithProviders(providers),
			)
			s.Require().NoError(err)

			ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(tc.duration))
			defer cancel()

			// Start the oracle. This should automatically stop.
			go func() {
				err = oracle.Start(ctx)
				s.Require().Equal(err, context.DeadlineExceeded)
			}()

			// Ensure that the oracle is not running.
			s.Eventually(checkFn(oracle), 2*tc.duration, 100*time.Millisecond)
		})
	}
}

func (s *OracleTestSuite) TestStop() {
	testCases := []struct {
		name     string
		factory  providertypes.ProviderFactory[oracletypes.CurrencyPair, *big.Int]
		duration time.Duration
	}{
		{
			name: "1 provider",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
				config.OracleMetricsConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				provider := testutils.CreateProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				// Create the provider factory.
				providers := []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider}
				return providers, nil
			},
			duration: 1 * time.Second,
		},
		{
			name: "multiple providers",
			factory: func(
				*zap.Logger,
				config.OracleConfig,
				config.OracleMetricsConfig,
			) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
				provider1 := testutils.CreateProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				provider2 := testutils.CreateProviderWithGetResponses[oracletypes.CurrencyPair, *big.Int](
					s.T(),
					s.logger,
					providerCfg2,
					s.currencyPairs,
					nil,
				)

				// Create the provider factory.
				providers := []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{provider1, provider2}
				return providers, nil
			},
			duration: 1 * time.Second,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cfg := config.OracleConfig{
				UpdateInterval: 1 * time.Second,
				ClientTimeout:  1 * time.Second,
				InProcess:      true,
			}
			metricsCfg := config.OracleMetricsConfig{
				Enabled: false,
			}

			providers, err := tc.factory(s.logger, cfg, metricsCfg)
			s.Require().NoError(err)

			oracle, err := oracle.New(
				cfg,
				oracle.WithLogger(s.logger),
				oracle.WithProviders(providers),
			)
			s.Require().NoError(err)

			// Start the oracle. This should automatically stop.
			go func() {
				oracle.Start(context.Background())
			}()

			// Wait for the experiment to run.
			time.Sleep(tc.duration)

			// Ensure that the oracle is not running.
			oracle.Stop()

			// Ensure that the oracle is not running.
			s.Eventually(checkFn(oracle), 2*tc.duration, 100*time.Millisecond)
		})
	}
}

func checkFn(o *oracle.Oracle) func() bool {
	return func() bool {
		return !o.IsRunning()
	}
}
