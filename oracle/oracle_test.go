package oracle_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
)

var (
	providerCfg1 = config.ProviderConfig{
		Name: "api1",
		API: config.APIConfig{
			Interval:         500 * time.Millisecond,
			Timeout:          250 * time.Millisecond,
			ReconnectTimeout: 250 * time.Millisecond,
			MaxQueries:       10,
			Enabled:          true,
			Name:             "api1",
			URL:              "http://test.com",
		},
		Type: "price_provider",
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
		Type: "price_provider",
	}
)

type OracleTestSuite struct {
	suite.Suite
	random *rand.Rand

	logger *zap.Logger

	// Oracle config
	currencyPairs []types.ProviderTicker
}

func TestOracleSuite(t *testing.T) {
	suite.Run(t, new(OracleTestSuite))
}

func (s *OracleTestSuite) SetupTest() {
	s.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	s.logger = zap.NewExample()

	s.currencyPairs = []types.ProviderTicker{
		types.NewProviderTicker("test", "BTC/USD", "{}"),
		types.NewProviderTicker("test", "ETH/USD", "{}"),
		types.NewProviderTicker("test", "ATOM/USD", "{}"),
	}
}
