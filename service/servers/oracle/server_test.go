package oracle_test

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/mocks"
	"github.com/skip-mev/connect/v2/oracle/types"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	client "github.com/skip-mev/connect/v2/service/clients/oracle"
	"github.com/skip-mev/connect/v2/service/metrics"
	server "github.com/skip-mev/connect/v2/service/servers/oracle"
	stypes "github.com/skip-mev/connect/v2/service/servers/oracle/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

const (
	localhost     = "localhost"
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
	port       string
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupTest() {
	// mock logger
	logger := zap.NewExample()

	s.mockOracle = mocks.NewOracle(s.T())
	s.srv = server.NewOracleServer(s.mockOracle, logger)

	// listen on a random port and extract that port number
	ln, err := net.Listen("tcp", localhost+":0")
	s.Require().NoError(err)
	_, s.port, err = net.SplitHostPort(ln.Addr().String())
	s.Require().NoError(err)

	s.client, err = client.NewClient(
		log.NewTestLogger(s.T()),
		localhost+":"+s.port,
		timeout,
		metrics.NewNopMetrics(),
		client.WithBlockingDial(), // block on dialing the server
	)

	s.Require().NoError(err)

	s.httpClient = http.DefaultClient

	// create context
	s.ctx, s.cancel = context.WithCancel(context.Background())

	// start server + client w/ context
	go s.srv.StartServerWithListener(s.ctx, ln)

	dialCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Require().NoError(s.client.Start(dialCtx))

	// Health check
	for i := 0; ; i++ {
		_, err := s.httpClient.Get(fmt.Sprintf("http://%s:%s/slinky/oracle/v1/version", localhost, s.port))
		if err == nil {
			break
		}
		if i == 10 {
			s.T().Fatal("failed to connect to server")
		}
		time.Sleep(1 * time.Second)
	}
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
	s.mockOracle.EXPECT().IsRunning().Return(false)

	// call from client
	_, err := s.client.Prices(context.Background(), &stypes.QueryPricesRequest{})

	// expect oracle not running error
	s.Require().Equal(err.Error(), grpcErrPrefix+server.ErrOracleNotRunning.Error())
}

func (s *ServerTestSuite) TestOracleServerTimeout() {
	// set the mock oracle to delay GetPrices response (delay for absurd time)
	s.mockOracle.EXPECT().IsRunning().Return(true)
	s.mockOracle.On("GetPrices").Return(nil).After(delay)

	// call from client
	_, err := s.client.Prices(context.Background(), &stypes.QueryPricesRequest{})

	// expect deadline exceeded error
	s.Require().Error(err)
}

func (s *ServerTestSuite) TestOracleServerPrices() {
	// set the mock oracle to return price-data
	s.mockOracle.EXPECT().IsRunning().Return(true)
	cp1 := mmtypes.Ticker{
		CurrencyPair: connecttypes.CurrencyPair{
			Base:  "BTC",
			Quote: "USD",
		},
		Decimals: 8,
	}

	cp2 := mmtypes.Ticker{
		CurrencyPair: connecttypes.CurrencyPair{
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
	httpResp, err := s.httpClient.Get(fmt.Sprintf("http://%s:%s/connect/oracle/v2/prices", localhost, s.port))
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
				CurrencyPair:     connecttypes.CurrencyPair{Base: "ETH", Quote: "USD"},
				Decimals:         420,
				MinProviderCount: 79,
				Enabled:          true,
				Metadata_JSON:    "",
			},
			ProviderConfigs: []mmtypes.ProviderConfig{
				{
					Name:           "FOO",
					OffChainTicker: "BAR",
					NormalizeByPair: &connecttypes.CurrencyPair{
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
