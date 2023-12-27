package oracle_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	metric_mocks "github.com/skip-mev/slinky/oracle/metrics/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
	provider_mocks "github.com/skip-mev/slinky/providers/types/mocks"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

type OracleMetricsTestSuite struct {
	suite.Suite

	// mocked providers
	mockProvider1 *provider_mocks.Provider[oracletypes.CurrencyPair, *big.Int]
	mockProvider2 *provider_mocks.Provider[oracletypes.CurrencyPair, *big.Int]

	// mock metrics
	mockMetrics *metric_mocks.Metrics

	o *oracle.Oracle
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
	s.mockProvider1 = provider_mocks.NewProvider[oracletypes.CurrencyPair, *big.Int](s.T())
	s.mockProvider1.On("Name").Return("provider1").Maybe()

	s.mockProvider2 = provider_mocks.NewProvider[oracletypes.CurrencyPair, *big.Int](s.T())
	s.mockProvider2.On("Name").Return("provider2").Maybe()

	// mock metrics
	s.mockMetrics = metric_mocks.NewMetrics(s.T())

	oracleConfig := config.OracleConfig{
		InProcess:      true,
		RemoteAddress:  "",
		UpdateInterval: oracleTicker,
	}
	factory := func(*zap.Logger, config.OracleConfig) ([]providertypes.Provider[oracletypes.CurrencyPair, *big.Int], error) {
		return []providertypes.Provider[oracletypes.CurrencyPair, *big.Int]{s.mockProvider1, s.mockProvider2}, nil
	}

	var err error
	s.o, err = oracle.New(
		zap.NewNop(),
		oracleConfig,
		factory,
		aggregator.ComputeMedian(),
		s.mockMetrics,
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

// test Tick metrics are updated correctly
func (s *OracleMetricsTestSuite) TestTickMetric() {
	// expect tick to be called
	s.mockMetrics.On("AddTick").Return()

	s.mockProvider1.On("Name").Return("provider1")
	s.mockProvider1.On("Start", mock.Anything).Return(nil)
	s.mockProvider1.On("LastUpdate").Return(time.Now().Add(time.Hour))
	s.mockProvider1.On("GetData", mock.Anything).Return(nil, nil)

	s.mockProvider2.On("Name").Return("provider2")
	s.mockProvider2.On("Start", mock.Anything).Return(nil)
	s.mockProvider2.On("LastUpdate").Return(time.Now().Add(time.Hour))
	s.mockProvider2.On("GetData", mock.Anything).Return(nil, nil)

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
