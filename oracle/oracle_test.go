package oracle_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle"
	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/oracle/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
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
			Endpoints:        []config.Endpoint{{URL: "http://test.com"}},
		},
		Type: "price_provider",
	}
	providerCfg2 = config.ProviderConfig{
		Name: "websocket1",
		WebSocket: config.WebSocketConfig{
			MaxBufferSize:       10,
			Enabled:             true,
			ReconnectionTimeout: 250 * time.Millisecond,
			Endpoints: []config.Endpoint{
				{
					URL: "ws://localhost:8080",
				},
			},
			Name:                          "websocket1",
			ReadBufferSize:                config.DefaultReadBufferSize,
			WriteBufferSize:               config.DefaultWriteBufferSize,
			HandshakeTimeout:              config.DefaultHandshakeTimeout,
			EnableCompression:             config.DefaultEnableCompression,
			ReadTimeout:                   config.DefaultReadTimeout,
			WriteTimeout:                  config.DefaultWriteTimeout,
			PingInterval:                  config.DefaultPingInterval,
			MaxSubscriptionsPerConnection: config.DefaultMaxSubscriptionsPerConnection,
			MaxSubscriptionsPerBatch:      config.DefaultMaxSubscriptionsPerBatch,
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

	marketmap mmtypes.MarketMap
}

func TestOracleSuite(t *testing.T) {
	suite.Run(t, new(OracleTestSuite))
}

func (s *OracleTestSuite) SetupTest() {
	s.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	s.logger = zap.NewExample()

	s.currencyPairs = []types.ProviderTicker{
		types.NewProviderTicker("BTC/USD", "{}"),
		types.NewProviderTicker("ETH/USD", "{}"),
		types.NewProviderTicker("ATOM/USD", "{}"),
	}
	s.marketmap = mmtypes.MarketMap{Markets: map[string]mmtypes.Market{
		btcusdtCP.String(): {
			Ticker: mmtypes.Ticker{
				CurrencyPair:     btcusdtCP,
				MinProviderCount: 1,
				Decimals:         8,
				Enabled:          true,
			},
			ProviderConfigs: []mmtypes.ProviderConfig{
				{
					Name:           providerCfg1.Name,
					OffChainTicker: coinbasebtcusd.GetOffChainTicker(),
				},
				{
					Name:           providerCfg2.Name,
					OffChainTicker: coinbasebtcusd.GetOffChainTicker(),
				},
			},
		},
	}}
}

func (s *OracleTestSuite) TestGetMarketMap() {
	dummyMM := marketMap
	o, err := oracle.New(oracleCfg, noOpPriceAggregator{}, oracle.WithMarketMap(dummyMM))
	s.Require().NoError(err)

	gotMM := o.GetMarketMap()
	s.Require().Equal(dummyMM, gotMM)

	// test when no option provided, should just give empty mm
	o, err = oracle.New(oracleCfg, noOpPriceAggregator{})
	s.Require().NoError(err)
	gotMM = o.GetMarketMap()
	s.Require().Equal(gotMM, mmtypes.MarketMap{})
}

func (s *OracleTestSuite) TestErrorsWhenNoAggregator() {
	_, err := oracle.New(oracleCfg, nil)
	s.Require().ErrorContains(err, "aggregator is required")
}
