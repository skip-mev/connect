package oracle_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	client "github.com/skip-mev/slinky/service/clients/oracle"
	clientmock "github.com/skip-mev/slinky/service/clients/oracle/mocks"
	"github.com/skip-mev/slinky/service/metrics"
	metrics_mock "github.com/skip-mev/slinky/service/metrics/mocks"
	"github.com/skip-mev/slinky/service/servers/oracle/types"
)

type MetricsClientTestSuite struct {
	suite.Suite
	m          *metrics_mock.Metrics    // mocked metrics
	mockClient *clientmock.OracleClient // mocked client
	client     *client.MetricsClient    // metrics client
}

func TestMetricsClientTestSuite(t *testing.T) {
	suite.Run(t, new(MetricsClientTestSuite))
}

func (s *MetricsClientTestSuite) SetupSubTest() {
	s.m = metrics_mock.NewMetrics(s.T())
	s.mockClient = clientmock.NewOracleClient(s.T())
	s.client = client.NewMetricsClient(log.NewNopLogger(), s.mockClient, s.m)
}

// test that responses are updated correctly
func (s *MetricsClientTestSuite) TestResponses() {
	s.Run("test that correct responses are reported correctly", func() {
		// expect a normal response
		ctx := context.Background()
		req := &types.QueryPricesRequest{}
		res := &types.QueryPricesResponse{}
		s.mockClient.On("Prices", ctx, req).Return(res, nil).Once()
		// expect a normal response
		s.m.On("AddOracleResponse", metrics.StatusSuccess).Return().Once()
		s.m.On("ObserveOracleResponseLatency", mock.Anything).Return().Once()

		// call the client
		actualRes, actualErr := s.client.Prices(ctx, req)
		s.Require().NoError(actualErr)
		s.Require().Equal(res, actualRes)
		s.m.AssertExpectations(s.T())
		s.mockClient.AssertExpectations(s.T())
	})

	s.Run("test that error responses are reported correctly", func() {
		// expect an error response
		ctx := context.Background()
		req := &types.QueryPricesRequest{}
		s.mockClient.On("Prices", ctx, req).Return(nil, fmt.Errorf("error"))
		// expect an error response
		s.m.On("AddOracleResponse", metrics.StatusFailure).Return()
		s.m.On("ObserveOracleResponseLatency", mock.Anything).Return()

		// call the client
		_, actualErr := s.client.Prices(ctx, req)
		s.Require().Error(actualErr)
		s.m.AssertExpectations(s.T())
		s.mockClient.AssertExpectations(s.T())
	})
}

// test that histogram observations are updated correctly
func (s *MetricsClientTestSuite) TestResponseLatency() {
	s.Run("test that response latency is reported correctly, no latency", func() {
		// expect a normal response
		ctx := context.Background()
		req := &types.QueryPricesRequest{}
		res := &types.QueryPricesResponse{}
		s.mockClient.On("Prices", ctx, req).Return(res, nil).Once()
		// expect a normal response
		s.m.On("AddOracleResponse", metrics.StatusSuccess).Return().Once()
		s.m.On("ObserveOracleResponseLatency", mock.Anything).Return().Once().Run(func(args mock.Arguments) {
			// expect to be within +/- 10ms of 0ms
			s.InDelta(0, args.Get(0), float64(10*time.Millisecond))
		})

		// call the client
		actualRes, actualErr := s.client.Prices(ctx, req)
		s.Require().NoError(actualErr)
		s.Require().Equal(res, actualRes)
		s.m.AssertExpectations(s.T())
		s.mockClient.AssertExpectations(s.T())
	})

	s.Run("test that response latency is reported correctly, with latency", func() {
		// expect a normal response
		ctx := context.Background()
		req := &types.QueryPricesRequest{}
		res := &types.QueryPricesResponse{}
		s.mockClient.On("Prices", ctx, req).Return(res, nil).Once().After(100 * time.Millisecond)
		// expect a normal response
		s.m.On("AddOracleResponse", metrics.StatusSuccess).Return().Once()
		s.m.On("ObserveOracleResponseLatency", mock.Anything).Return().Once().Run(func(args mock.Arguments) {
			// expect to be within +/- 10ms of 100ms
			s.InDelta(100*time.Millisecond, args.Get(0), float64(10*time.Millisecond))
		})

		// call the client
		actualRes, actualErr := s.client.Prices(ctx, req)
		s.Require().NoError(actualErr)
		s.Require().Equal(res, actualRes)
		s.m.AssertExpectations(s.T())
		s.mockClient.AssertExpectations(s.T())
	})
}
