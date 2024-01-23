package oracle_test

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"

	"github.com/skip-mev/slinky/oracle/mocks"
	client "github.com/skip-mev/slinky/service/clients/oracle"
	"github.com/skip-mev/slinky/service/metrics"
	server "github.com/skip-mev/slinky/service/servers/oracle"
	stypes "github.com/skip-mev/slinky/service/servers/oracle/types"
	"github.com/skip-mev/slinky/x/oracle/types"
)

const (
	localhost     = "localhost"
	port          = "8080"
	timeout       = 1 * time.Second
	delay         = 20 * time.Second
	grpcErrPrefix = "rpc error: code = Unknown desc = "
)

type ServerTestSuite struct {
	suite.Suite

	srv        *server.OracleServer
	mockOracle *mocks.Oracle
	client     client.OracleClient
	httpClient *http.Client
	ctx        context.Context
	cancel     context.CancelFunc
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupTest() {
	// mock logger
	logger := zap.NewNop()

	s.mockOracle = mocks.NewOracle(s.T())
	s.srv = server.NewOracleServer(s.mockOracle, logger)

	var err error
	s.client, err = client.NewClient(
		log.NewTestLogger(s.T()),
		localhost+":"+port,
		timeout,
		metrics.NewNopMetrics(),
	)
	s.Require().NoError(err)

	s.httpClient = http.DefaultClient

	// create context
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// expect oracle to start
	s.mockOracle.On("Start", mock.Anything).Return(nil)

	// start server + client w/ context
	go s.srv.StartServer(s.ctx, localhost, port)
	s.Require().NoError(s.client.Start())
}

// teardown test suite
func (s *ServerTestSuite) TearDownTest() {
	// close server
	s.srv.Close()
	defer s.cancel()

	// wait for the server to finish
	select {
	case <-s.srv.Done():
	case <-time.After(2 * time.Second):
		s.T().Fatal("server failed to stop")
	}

	// close client
	s.Require().NoError(s.client.Stop())
}

func (s *ServerTestSuite) TestOracleServerNotRunning() {
	// set the mock oracle to not be running
	s.mockOracle.On("IsRunning").Return(false)

	// call from client
	_, err := s.client.Prices(context.Background(), &stypes.QueryPricesRequest{})

	// expect oracle not running error
	s.Require().Equal(err.Error(), grpcErrPrefix+server.ErrOracleNotRunning.Error())
}

func (s *ServerTestSuite) TestOracleServerTimeout() {
	// set the mock oracle to delay GetPrices response (delay for absurd time)
	s.mockOracle.On("IsRunning").Return(true)
	s.mockOracle.On("GetPrices").Return(nil).After(delay)

	// call from client
	_, err := s.client.Prices(context.Background(), &stypes.QueryPricesRequest{})

	// expect deadline exceeded error
	s.Require().Equal(err.Error(), status.FromContextError(context.DeadlineExceeded).Err().Error())
}

func (s *ServerTestSuite) TestOracleServerPrices() {
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

	// call from grpc client
	resp, err := s.client.Prices(context.Background(), &stypes.QueryPricesRequest{})
	s.Require().NoError(err)

	// check response
	s.Require().Equal(resp.Prices[cp1.ToString()], big.NewInt(100).String())
	s.Require().Equal(resp.Prices[cp2.ToString()], big.NewInt(200).String())
	// check timestamp

	s.Require().Equal(resp.Timestamp, ts.UTC())

	// call from http client
	httpResp, err := s.httpClient.Get(fmt.Sprintf("http://%s:%s/slinky/oracle/v1/prices", localhost, port))
	s.Require().NoError(err)

	// check response
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	respBz, err := io.ReadAll(httpResp.Body)
	s.Require().NoError(err)
	s.Require().Contains(string(respBz), fmt.Sprintf(`{"prices":{"%s":"100","%s":"200"},"timestamp":`, cp1.ToString(), cp2.ToString()))
}

// test that the oracle server closes when expected
func (s *ServerTestSuite) TestOracleServerClose() {
	// close the server, and check that no requests are received
	s.cancel()

	// wait for server to close
	select {
	case <-s.srv.Done():
	case <-time.After(1 * time.Second):
		s.T().Fatal("server failed to stop")
	}

	// expect requests to server to timeout
	_, err := s.client.Prices(context.Background(), &stypes.QueryPricesRequest{})

	// expect request to have failed (connection is closed)
	s.Require().NotNil(err)
}

func TestOracleFailureStopsServer(t *testing.T) {
	// create mock oracle
	mockOracle := mocks.NewOracle(t)
	mockOracle.On("Start", mock.Anything).Return(fmt.Errorf("failed to start oracle"))

	// create server
	srv := server.NewOracleServer(mockOracle, zap.NewNop())

	// start the server, and expect immediate closure
	go srv.StartServer(context.Background(), localhost, port)

	// wait for server to close
	select {
	case <-srv.Done():
	case <-time.After(1 * time.Second):
		t.Fatal("server failed to stop")
	}
}
