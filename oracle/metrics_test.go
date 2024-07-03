package oracle_test

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	metricmocks "github.com/skip-mev/slinky/oracle/metrics/mocks"
	"github.com/skip-mev/slinky/oracle/types"
	mathtestutils "github.com/skip-mev/slinky/pkg/math/testutils"
	"github.com/skip-mev/slinky/providers/base/testutils"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
)

func (s *OracleTestSuite) TestMetrics() {
	cfg := config.OracleConfig{
		UpdateInterval: 1 * time.Second,
		MaxPriceAge:    1 * time.Minute,
		Providers:      nil,
		Metrics:        oracleCfg.Metrics,
		Host:           oracleCfg.Host,
		Port:           oracleCfg.Port,
	}
	provider := testutils.CreateAPIProviderWithGetResponses[types.ProviderTicker, *big.Float](
		s.T(),
		s.logger,
		providerCfg1,
		s.currencyPairs,
		nil,
		200*time.Millisecond,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*cfg.UpdateInterval)
	defer cancel()

	metrics := metricmocks.NewMetrics(s.T())
	testOracle, err := oracle.New(
		cfg,
		mathtestutils.NewMedianAggregator(),
		oracle.WithLogger(s.logger),
		oracle.WithPriceProviders(provider),
		oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),
		oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory),
		oracle.WithMarketMap(s.marketmap),
		oracle.WithMetrics(metrics),
	)
	s.Require().NoError(err)

	go func() {
		err := testOracle.Start(ctx)
		if err != nil {
			if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
				s.T().Errorf("Start() should have returned context.Canceled error. Got: %v", err)
			}
		}
	}()

	metrics.On("SetSlinkyBuildInfo").Return()
	metrics.On("AddTick").Return()

	time.Sleep(2 * cfg.UpdateInterval)
	testOracle.Stop() // block on the oracle actually closing
	metrics.AssertExpectations(s.T())
}
