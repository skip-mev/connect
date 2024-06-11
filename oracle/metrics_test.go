package oracle_test

import (
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

/*
import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	metricmocks "github.com/skip-mev/slinky/oracle/metrics/mocks"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base/testutils"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmclienttypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)
*/

var _ oracle.PriceAggregator = &noOpPriceAggregator{}

type noOpPriceAggregator struct{}

func (n noOpPriceAggregator) SetProviderPrices(provider string, prices types.Prices) {
}

func (n noOpPriceAggregator) UpdateMarketMap(m mmtypes.MarketMap) {
}

func (n noOpPriceAggregator) AggregatePrices() {
}

func (n noOpPriceAggregator) GetPrices() types.Prices {
	return types.Prices{}
}

func (n noOpPriceAggregator) Reset() {
}

/*
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

	handler, mmp := createTestMarketMapProvider(s.T(), []mmclienttypes.Chain{{ChainID: "foo"}})
	handler.On("CreateURL", mock.Anything).Return("", nil).Maybe()
	handler.On("ParseResponse", mock.Anything).Return(responses, nil).Maybe()

	// mock metrics
	s.mockMetrics = metricmocks.NewMetrics(s.T())
	var err error
	s.o, err = oracle.New(
		config.OracleConfig{UpdateInterval: 1 * time.Millisecond, MaxPriceAge: 60 * time.Second, Host: "foo", Port: "10"},
		noOpPriceAggregator{},
		oracle.WithMetrics(s.mockMetrics),
		oracle.WithPriceProviders(providers...),
		oracle.WithMarketMapProvider(mmp),
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
	s.mockMetrics.On("SetSlinkyBuildInfo").Return()
	s.mockMetrics.On("AddTick").Return()

	t1 := s.o.GetLastSyncTime()
	// wait for a tick on the oracle
	go func() {
		s.Require().NoError(s.o.Start(context.Background()))
	}()

	// wait for a tick
	for {
		t2 := s.o.GetLastSyncTime()
		if !t2.Equal(t1) {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// assert expectations
	s.mockMetrics.AssertExpectations(s.T())
	s.o.Stop()
}

*/
