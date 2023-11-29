package client_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/slinky/service/client"
	"github.com/skip-mev/slinky/service/types/mocks"

	"github.com/skip-mev/slinky/service"
	servicetypes "github.com/skip-mev/slinky/service/types"
	"github.com/skip-mev/slinky/x/oracle/types"
)

const (
	clientTimeout = 1 * time.Second
	delay         = 5 * time.Second
)

type ClientTestSuite struct {
	suite.Suite

	// mock oracle
	mockOracle *mocks.Oracle
	// LocalClient
	client *client.LocalClient
	ctx    context.Context
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (s *ClientTestSuite) SetupTest() {
	// create mock oracle
	s.mockOracle = mocks.NewOracle(s.T())
	// create local client
	s.client = client.NewLocalClient(s.mockOracle, clientTimeout)

	s.ctx = context.Background()
	// expect oracle to start
	s.mockOracle.On("Start", s.ctx).Return(nil)
	s.mockOracle.On("IsRunning").Return(false).Once()
	// start the LocalClient
	s.Require().NoError(s.client.Start(s.ctx))
}

func (s *ClientTestSuite) TearDownTest() {
	s.mockOracle.On("Stop").Return(nil)
	// stop the LocalClient
	s.Require().NoError(s.client.Stop(s.ctx))
}

// test that the client times out if the oracle times out
func (s *ClientTestSuite) TestTimeout() {
	// expect oracle to timeout after a long delay
	s.mockOracle.On("IsRunning").Return(true).Once()
	s.mockOracle.On("GetPrices").Return(nil, context.Canceled).After(delay)
	// query prices
	_, err := s.client.Prices(context.Background(), &service.QueryPricesRequest{})
	s.Require().Error(err)
	s.Require().Equal(context.Canceled, err)
}

// test that the client returns if the context is cancelled before request is returned
func (s *ClientTestSuite) TestClientReturnsWhenContextCancelled() {
	// expect oracle to timeout after a long delay
	s.mockOracle.On("IsRunning").Return(true).Once()

	// may be called in go-routine, may not be, but shld not return before context cancellation
	s.mockOracle.On("GetPrices").Return(nil, nil).Maybe()
	s.mockOracle.On("GetLastSyncTime").Return(time.Now()).Maybe().After(clientTimeout / 2)
	// query prices
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := s.client.Prices(ctx, &service.QueryPricesRequest{})
	s.Require().Error(err)
	s.Require().Equal(context.Canceled, err)
}

// test that the client returns if the oracle is not running
func (s *ClientTestSuite) TestClientReturnsWhenOracleNotRunning() {
	// expect oracle to timeout after a long delay
	s.mockOracle.On("IsRunning").Return(false).Once()
	// query prices
	_, err := s.client.Prices(context.Background(), &service.QueryPricesRequest{})
	s.Require().Error(err)
	s.Require().Equal(servicetypes.ErrorOracleNotRunning, err)
}

// test that the client returns valid prices
func (s *ClientTestSuite) TestClientReturnsValidPrices() {
	// set the mock oracle to return price-data
	s.mockOracle.On("IsRunning").Return(true)
	cp1 := types.CurrencyPair{
		Base:  "BTC",
		Quote: "USD",
	}

	cp2 := types.CurrencyPair{
		Base:  "ETH",
		Quote: "USD",
	}

	s.mockOracle.On("GetPrices").Return(map[types.CurrencyPair]*big.Int{
		cp1: big.NewInt(100),
		cp2: big.NewInt(200),
	})

	ts := time.Now()
	s.mockOracle.On("GetLastSyncTime").Return(ts)

	// call from client
	resp, err := s.client.Prices(context.Background(), &service.QueryPricesRequest{})
	s.Require().NoError(err)

	// check response
	s.Require().Equal(resp.Prices[cp1.ToString()], big.NewInt(100).String())
	s.Require().Equal(resp.Prices[cp2.ToString()], big.NewInt(200).String())
}
