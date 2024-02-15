package oracle_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/slinky/oracle"
	metricmocks "github.com/skip-mev/slinky/oracle/metrics/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
	providermocks "github.com/skip-mev/slinky/providers/types/mocks"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

type OracleMetricsTestSuite struct {
	suite.Suite

	// mocked providers
	mockProvider1 *providermocks.Provider[mmtypes.Ticker, *big.Int]
	mockProvider2 *providermocks.Provider[mmtypes.Ticker, *big.Int]

	// mock metrics
	mockMetrics *metricmocks.Metrics

	o oracle.Oracle
}

const (
	oracleTicker = 1 * time.Second
	provider1    = "provider1"
	provider2    = "provider2"
)

func TestOracleMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(OracleMetricsTestSuite))
}

func (s *OracleMetricsTestSuite) SetupTest() {
	// mock providers
	s.mockProvider1 = providermocks.NewProvider[mmtypes.Ticker, *big.Int](s.T())
	s.mockProvider1.On("Name").Return("provider1").Maybe()

	s.mockProvider2 = providermocks.NewProvider[mmtypes.Ticker, *big.Int](s.T())
	s.mockProvider2.On("Name").Return("provider2").Maybe()

	// mock metrics
	s.mockMetrics = metricmocks.NewMetrics(s.T())

	var err error
	s.o, err = oracle.New(
		oracle.WithUpdateInterval(oracleTicker),
		oracle.WithProviders(
			[]providertypes.Provider[mmtypes.Ticker, *big.Int]{
				s.mockProvider1,
				s.mockProvider2,
			},
		),
		oracle.WithMetrics(s.mockMetrics),
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

	s.mockProvider1.On("Name").Return("provider1")
	s.mockProvider1.On("Start", mock.Anything).Return(nil)
	s.mockProvider1.On("Type").Return(providertypes.API)
	s.mockProvider1.On("GetData").Return(nil)

	s.mockProvider2.On("Name").Return("provider2")
	s.mockProvider2.On("Start", mock.Anything).Return(nil)
	s.mockProvider2.On("Type").Return(providertypes.API)
	s.mockProvider2.On("GetData").Return(nil, nil)

	// wait for a tick on the oracle
	go func() {
		s.o.Start(context.Background())
	}()

	// wait for a tick
	time.Sleep(4 * oracleTicker)

	// assert expectations
	s.mockMetrics.AssertExpectations(s.T())
	s.o.Stop()
}
