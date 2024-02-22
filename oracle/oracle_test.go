package oracle_test

import (
	"context"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/testutils"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	providerCfg1 = config.ProviderConfig{
		Name: "api1",
		API: config.APIConfig{
			Interval:   500 * time.Millisecond,
			Timeout:    250 * time.Millisecond,
			MaxQueries: 10,
			Enabled:    true,
			Name:       "api1",
			URL:        "http://test.com",
		},
	}
	providerCfg2 = config.ProviderConfig{
		Name: "websocket1",
		WebSocket: config.WebSocketConfig{
			MaxBufferSize:                 10,
			Enabled:                       true,
			ReconnectionTimeout:           250 * time.Millisecond,
			WSS:                           "wss://test.com",
			Name:                          "websocket1",
			ReadBufferSize:                config.DefaultReadBufferSize,
			WriteBufferSize:               config.DefaultWriteBufferSize,
			HandshakeTimeout:              config.DefaultHandshakeTimeout,
			EnableCompression:             config.DefaultEnableCompression,
			ReadTimeout:                   config.DefaultReadTimeout,
			WriteTimeout:                  config.DefaultWriteTimeout,
			PingInterval:                  config.DefaultPingInterval,
			MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
		},
	}
)

type OracleTestSuite struct {
	suite.Suite
	random *rand.Rand

	logger *zap.Logger

	// Oracle config
	currencyPairs []mmtypes.Ticker
}

func TestOracleSuite(t *testing.T) {
	suite.Run(t, new(OracleTestSuite))
}

func (s *OracleTestSuite) SetupTest() {
	s.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	s.logger = zap.NewExample()

	s.currencyPairs = []mmtypes.Ticker{
		constants.BITCOIN_USD,
		constants.ETHEREUM_USD,
		constants.ATOM_USD,
	}
}

func (s *OracleTestSuite) TestStopWithContextCancel() {
	testCases := []struct {
		name    string
		factory types.PriceProviderFactory
	}{
		{
			name: "no providers",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProvider, error) {
				return nil, nil
			},
		},
		{
			name: "1 provider",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProvider, error) {
				provider := testutils.CreateAPIProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				// Create the provider factory.
				providers := []types.PriceProvider{provider}
				return providers, nil
			},
		},
		{
			name: "multiple providers",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProvider, error) {
				provider1 := testutils.CreateAPIProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				provider2 := testutils.CreateWebSocketProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					time.Second,
					s.currencyPairs,
					providerCfg2,
					s.logger,
					nil,
				)

				// Create the provider factory.
				providers := []types.PriceProvider{provider1, provider2}
				return providers, nil
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cfg := config.OracleConfig{
				UpdateInterval: 1 * time.Second,
			}

			providers, err := tc.factory(cfg)
			s.Require().NoError(err)

			oracle, err := oracle.New(
				oracle.WithLogger(s.logger),
				oracle.WithProviders(providers),
				oracle.WithUpdateInterval(cfg.UpdateInterval),
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
		factory  types.PriceProviderFactory
		duration time.Duration
	}{
		{
			name: "no providers",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProvider, error) {
				return nil, nil
			},
			duration: 1 * time.Second,
		},
		{
			name: "1 provider",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProvider, error) {
				provider := testutils.CreateAPIProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				// Create the provider factory.
				providers := []types.PriceProvider{provider}
				return providers, nil
			},
			duration: 1 * time.Second,
		},
		{
			name: "multiple providers",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProvider, error) {
				provider1 := testutils.CreateAPIProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				provider2 := testutils.CreateWebSocketProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					time.Second,
					s.currencyPairs,
					providerCfg2,
					s.logger,
					nil,
				)

				// Create the provider factory.
				providers := []types.PriceProvider{provider1, provider2}
				return providers, nil
			},
			duration: 1 * time.Second,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cfg := config.OracleConfig{
				UpdateInterval: 1 * time.Second,
			}

			providers, err := tc.factory(cfg)
			s.Require().NoError(err)

			oracle, err := oracle.New(
				oracle.WithUpdateInterval(cfg.UpdateInterval),
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
		factory  types.PriceProviderFactory
		duration time.Duration
	}{
		{
			name: "1 provider",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProvider, error) {
				provider := testutils.CreateAPIProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				// Create the provider factory.
				providers := []types.PriceProvider{provider}
				return providers, nil
			},
			duration: 1 * time.Second,
		},
		{
			name: "multiple providers",
			factory: func(
				config.OracleConfig,
			) ([]types.PriceProvider, error) {
				provider1 := testutils.CreateAPIProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					s.logger,
					providerCfg1,
					s.currencyPairs,
					nil,
				)

				provider2 := testutils.CreateWebSocketProviderWithGetResponses[mmtypes.Ticker, *big.Int](
					s.T(),
					time.Second,
					s.currencyPairs,
					providerCfg2,
					s.logger,
					nil,
				)

				// Create the provider factory.
				providers := []types.PriceProvider{provider1, provider2}
				return providers, nil
			},
			duration: 1 * time.Second,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			cfg := config.OracleConfig{
				UpdateInterval: 1 * time.Second,
			}

			providers, err := tc.factory(cfg)
			s.Require().NoError(err)

			oracle, err := oracle.New(
				oracle.WithUpdateInterval(cfg.UpdateInterval),
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

func checkFn(o oracle.Oracle) func() bool {
	return func() bool {
		return !o.IsRunning()
	}
}
