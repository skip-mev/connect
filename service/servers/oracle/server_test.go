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
	"github.com/skip-mev/slinky/oracle/types"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	client "github.com/skip-mev/slinky/service/clients/oracle"
	"github.com/skip-mev/slinky/service/metrics"
	server "github.com/skip-mev/slinky/service/servers/oracle"
	stypes "github.com/skip-mev/slinky/service/servers/oracle/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
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
	logger := zap.NewExample()

	s.mockOracle = mocks.NewOracle(s.T())
	s.srv = server.NewOracleServer(s.mockOracle, logger)

	var err error
	s.client, err = client.NewClient(
		log.NewTestLogger(s.T()),
		localhost+":"+port,
		timeout,
		metrics.NewNopMetrics(),
		client.WithBlockingDial(), // block on dialing the server
	)

	s.Require().NoError(err)

	s.httpClient = http.DefaultClient

	// create context
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// start server + client w/ context
	go s.srv.StartServer(s.ctx, localhost, port)

	dialCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Require().NoError(s.client.Start(dialCtx))
}

// teardown test suite.
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
	cp1 := mmtypes.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "BTC",
			Quote: "USD",
		},
		Decimals: 8,
	}

	cp2 := mmtypes.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "ETH",
			Quote: "USD",
		},
		Decimals: 8,
	}

	s.mockOracle.On("GetPrices").Return(types.Prices{
		cp1.String(): big.NewFloat(100.1),
		cp2.String(): big.NewFloat(200.1),
	})
	ts := time.Now()
	s.mockOracle.On("GetLastSyncTime").Return(ts)

	// call from grpc client
	resp, err := s.client.Prices(context.Background(), &stypes.QueryPricesRequest{})
	s.Require().NoError(err)

	// check response
	s.Require().Equal(resp.Prices[cp1.String()], big.NewInt(100).String())
	s.Require().Equal(resp.Prices[cp2.String()], big.NewInt(200).String())
	// check timestamp

	s.Require().Equal(resp.Timestamp, ts.UTC())

	// call from http client
	httpResp, err := s.httpClient.Get(fmt.Sprintf("http://%s:%s/slinky/oracle/v1/prices", localhost, port))
	s.Require().NoError(err)

	// check response
	s.Require().Equal(http.StatusOK, httpResp.StatusCode)
	respBz, err := io.ReadAll(httpResp.Body)
	s.Require().NoError(err)
	s.Require().Contains(string(respBz), fmt.Sprintf(`{"prices":{"%s":"100","%s":"200"},"timestamp":`, cp1.String(), cp2.String()))
}

func (s *ServerTestSuite) TestOracleMarketMap() {
	dummyMarketMap := mmtypes.MarketMap{Markets: map[string]mmtypes.Market{
		"foo": {
			Ticker: mmtypes.Ticker{
				CurrencyPair:     slinkytypes.CurrencyPair{Base: "ETH", Quote: "USD"},
				Decimals:         420,
				MinProviderCount: 79,
				Enabled:          true,
				Metadata_JSON:    "",
			},
			ProviderConfigs: []mmtypes.ProviderConfig{
				{
					Name:           "FOO",
					OffChainTicker: "BAR",
					NormalizeByPair: &slinkytypes.CurrencyPair{
						Base:  "FOO",
						Quote: "BAR",
					},
				},
			},
		},
	}}
	expectedJSON, err := dummyMarketMap.Marshal()
	_ = expectedJSON
	s.Require().NoError(err)
	s.mockOracle.On("GetMarketMap", mock.Anything).Return(dummyMarketMap).Once()

	res, err := s.client.MarketMap(context.Background(), &stypes.QueryMarketMapRequest{})
	s.Require().NoError(err)
	s.Require().Equal(*res.GetMarketMap(), dummyMarketMap)
}

// test that the oracle server closes when expected.
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
