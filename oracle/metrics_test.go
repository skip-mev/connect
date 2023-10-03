package oracle_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/metrics"
	metric_mocks "github.com/skip-mev/slinky/oracle/metrics/mocks"
	"github.com/skip-mev/slinky/oracle/types"
	provider_mocks "github.com/skip-mev/slinky/oracle/types/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type OracleMetricsTestSuite struct {
	suite.Suite

	// mocked providers
	mockProvider1 *provider_mocks.Provider
	mockProvider2 *provider_mocks.Provider

	// mock metrics
	mockMetrics *metric_mocks.Metrics

	o *oracle.Oracle
}

const (
	oracleTicker = time.Second
	provider1    = "provider1"
	provider2    = "provider2"
)

func TestOracleMetricsTestSuite(t *testing.T) {
	suite.Run(t, new(OracleMetricsTestSuite))
}

func (s *OracleMetricsTestSuite) SetupTest() {
	// mock providers
	s.mockProvider1 = provider_mocks.NewProvider(s.T())
	s.mockProvider2 = provider_mocks.NewProvider(s.T())

	// mock metrics
	s.mockMetrics = metric_mocks.NewMetrics(s.T())

	s.o = oracle.New(
		log.NewNopLogger(),
		oracleTicker,
		[]types.Provider{s.mockProvider1, s.mockProvider2},
		types.ComputeMedian(),
		s.mockMetrics,
	)
}

// test Tick metrics are updated correctly
func (s *OracleMetricsTestSuite) TestTickMetric() {
	// expect tick to be called
	s.mockMetrics.On("AddTick").Return()
	s.mockMetrics.On("AddProviderResponse", provider1, mock.Anything).Return()
	s.mockMetrics.On("AddProviderResponse", provider2, mock.Anything).Return()
	s.mockMetrics.On("ObserveProviderResponseLatency", provider1, mock.Anything).Return()
	s.mockMetrics.On("ObserveProviderResponseLatency", provider2, mock.Anything).Return()

	s.mockProvider1.On("GetPrices", mock.Anything).Return(nil, nil)
	s.mockProvider2.On("GetPrices", mock.Anything).Return(nil, nil)
	s.mockProvider1.On("Name").Return("provider1")
	s.mockProvider2.On("Name").Return("provider2")

	// wait for a tick on the oracle
	go func() {
		s.o.Start(context.Background())
	}()

	// wait for a tick
	time.Sleep(2 * oracleTicker)

	// assert expectations
	s.mockMetrics.AssertExpectations(s.T())
	s.o.Stop()
}

// test ProviderResponseMetric from a provider the metrics are updated correctly
func (s *OracleMetricsTestSuite) TestProviderResponseMetric() {
	// expect tick to be called
	s.mockMetrics.On("AddTick").Return()
	s.mockMetrics.On("ObserveProviderResponseLatency", provider1, mock.Anything).Return()
	s.mockMetrics.On("ObserveProviderResponseLatency", provider2, mock.Anything).Return()
	s.mockMetrics.On("AddProviderResponse", provider1, metrics.StatusFailure).Return()
	s.mockMetrics.On("AddProviderResponse", provider2, metrics.StatusSuccess).Return()

	s.mockProvider1.On("GetPrices", mock.Anything).Return(nil, errors.New("provider1 error"))
	s.mockProvider2.On("GetPrices", mock.Anything).Return(nil, nil)
	s.mockProvider1.On("Name").Return("provider1")
	s.mockProvider2.On("Name").Return("provider2")

	// wait for a tick on the oracle
	go func() {
		s.o.Start(context.Background())
	}()

	// wait for a tick
	time.Sleep(2 * oracleTicker)

	// assert expectations
	s.mockMetrics.AssertExpectations(s.T())
	s.o.Stop()
}

// Test ProviderResponseTimeMetrics are updated correctly
func (s *OracleMetricsTestSuite) TestProviderResponseTimeMetric() {
	// expect tick to be called
	s.mockMetrics.On("AddTick").Return()
	s.mockMetrics.On("ObserveProviderResponseLatency", provider1, mock.Anything).Return().Run(func(args mock.Arguments) {
		// expect to be within +/- 100ms of 500ms
		assert.InDelta(s.T(), 100*time.Millisecond, args.Get(1), float64(20*time.Millisecond)) // delta may need to be tuned (this is arbitrary)
	})
	s.mockMetrics.On("ObserveProviderResponseLatency", provider2, mock.Anything).Return().Run(func(args mock.Arguments) {
		// expect to be within +/- 100ms of 1000ms
		assert.InDelta(s.T(), 150*time.Millisecond, args.Get(1), float64(20*time.Millisecond))
	})
	s.mockMetrics.On("AddProviderResponse", provider1, mock.Anything).Return()
	s.mockMetrics.On("AddProviderResponse", provider2, mock.Anything).Return()

	s.mockProvider1.On("GetPrices", mock.Anything).Return(nil, nil).After(100 * time.Millisecond)
	s.mockProvider2.On("GetPrices", mock.Anything).Return(nil, nil).After(150 * time.Millisecond)
	s.mockProvider1.On("Name").Return("provider1")
	s.mockProvider2.On("Name").Return("provider2")

	// wait for a tick on the oracle
	go func() {
		s.o.Start(context.Background())
	}()

	// wait for a tick
	time.Sleep(2 * oracleTicker)

	// assert expectations
	s.mockMetrics.AssertExpectations(s.T())
	s.o.Stop()
}
