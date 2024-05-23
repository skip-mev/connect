package oracle_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle"
	metricmocks "github.com/skip-mev/slinky/oracle/metrics/mocks"
	"github.com/skip-mev/slinky/oracle/types"
	mathtestutils "github.com/skip-mev/slinky/pkg/math/testutils"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

type OracleMetricsTestSuite struct {
	suite.Suite

	// mock metrics
	mockMetrics *metricmocks.Metrics

	o oracle.Oracle
}

const (
	oracleTicker = 1 * time.Second
)

func TestOracleMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(OracleMetricsTestSuite))
}

func (s *OracleMetricsTestSuite) SetupTest() {
	ids := []types.ProviderTicker{
		types.NewProviderTicker("BTCUSD", "{}"),
		types.NewProviderTicker("ETHUSD", "{}"),
	}

	// mock providers
	resolved := types.ResolvedPrices{}
	response := providertypes.NewGetResponse[types.ProviderTicker, *big.Float](resolved, nil)
	responses := []providertypes.GetResponse[types.ProviderTicker, *big.Float]{response}
	provider := testutils.CreateAPIProviderWithGetResponses[types.ProviderTicker, *big.Float](
		s.T(),
		zap.NewNop(),
		providerCfg1,
		ids,
		responses,
		200*time.Millisecond,
	)

	resolved2 := types.ResolvedPrices{}
	response2 := providertypes.NewGetResponse[types.ProviderTicker, *big.Float](resolved2, nil)
	responses2 := []providertypes.GetResponse[types.ProviderTicker, *big.Float]{response2}
	provider2 := testutils.CreateWebSocketProviderWithGetResponses[types.ProviderTicker, *big.Float](
		s.T(),
		time.Second*2,
		ids,
		providerCfg2,
		zap.NewNop(),
		responses2,
	)

	providers := []*types.PriceProvider{provider, provider2}

	// mock metrics
	s.mockMetrics = metricmocks.NewMetrics(s.T())

	var err error
	s.o, err = oracle.New(
		oracle.WithUpdateInterval(oracleTicker),
		oracle.WithProviders(providers),
		oracle.WithMetrics(s.mockMetrics),
		oracle.WithPriceAggregator(mathtestutils.NewMedianAggregator()),
	)
	s.Require().NoError(err)
}

// TearDownTest is run after each test in the suite.
func (s *OracleMetricsTestSuite) TearDownTest(_ *testing.T) {
	checkFn := func() bool {
		return !s.o.IsRunning()
	}
	s.Eventually(checkFn, 5*time.Second, 100*time.Millisecond)
}

// test Tick metrics are updated correctly.
func (s *OracleMetricsTestSuite) TestTickMetric() {
	// expect tick to be called
	s.mockMetrics.On("AddTick").Return()
	s.mockMetrics.On("SetSlinkyBuildInfo").Return()

	// wait for a tick on the oracle
	go func() {
		s.Require().NoError(s.o.Start(context.Background()))
	}()

	// wait for a tick
	time.Sleep(4 * oracleTicker)

	// assert expectations
	s.mockMetrics.AssertExpectations(s.T())
	s.o.Stop()
}
