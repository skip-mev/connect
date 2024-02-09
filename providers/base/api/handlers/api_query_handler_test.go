package handlers_test

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/oracle/config"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/providers/base/api/errors"
	"github.com/skip-mev/slinky/providers/base/api/handlers"
	"github.com/skip-mev/slinky/providers/base/api/handlers/mocks"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
	mockmetrics "github.com/skip-mev/slinky/providers/base/api/metrics/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

var (
	logger  = zap.NewExample()
	btcusd  = slinkytypes.NewCurrencyPair("BTC", "USD")
	ethusd  = slinkytypes.NewCurrencyPair("ETH", "USD")
	atomusd = slinkytypes.NewCurrencyPair("ATOM", "USD")

	constantURL = "http://fetchdata.org:8080"

	cfg = config.APIConfig{
		Enabled:    true,
		Timeout:    500 * time.Millisecond,
		Interval:   1 * time.Second,
		MaxQueries: 1,
		Atomic:     true,
		URL:        constantURL,
		Name:       "handler1",
	}

	nonAtomicCfg = config.APIConfig{
		Enabled:    true,
		Timeout:    500 * time.Millisecond,
		Interval:   1 * time.Second,
		MaxQueries: 1,
		Atomic:     false,
		URL:        constantURL,
		Name:       "handler1",
	}
)

func TestAPIQueryHandler(t *testing.T) {
	testCases := []struct {
		name           string
		requestHandler func() handlers.RequestHandler
		apiHandler     func() handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int]
		metrics        func() metrics.APIMetrics
		capacity       int
		ids            []slinkytypes.CurrencyPair
		atomic         bool
		responses      providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]
	}{
		{
			name: "no ids to query",
			requestHandler: func() handlers.RequestHandler {
				return mocks.NewRequestHandler(t)
			},
			apiHandler: func() handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[slinkytypes.CurrencyPair, *big.Int](t)
				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)
				return m
			},
			capacity: 0,
			ids:      []slinkytypes.CurrencyPair{},
			atomic:   false,
			responses: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved:   map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				UnResolved: map[slinkytypes.CurrencyPair]error{},
			},
		},
		{
			name: "single id to query with no errors and an atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Times(1)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] {
				expectedIDs := []slinkytypes.CurrencyPair{btcusd}

				h := mocks.NewAPIDataHandler[slinkytypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", expectedIDs).Return(constantURL, nil).Times(1)

				resolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				response := providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
					resolved,
					nil,
				)

				h.On("ParseResponse", expectedIDs, newValidResponse()).Return(response).Times(1)

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", mock.Anything).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), metrics.Success).Times(1)

				return m
			},
			capacity: 1,
			ids:      []slinkytypes.CurrencyPair{btcusd},
			atomic:   true,
			responses: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				},
				UnResolved: map[slinkytypes.CurrencyPair]error{},
			},
		},
		{
			name: "single id to query with no errors and a non-atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Times(1)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] {
				expectedIDs := []slinkytypes.CurrencyPair{btcusd}

				h := mocks.NewAPIDataHandler[slinkytypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", expectedIDs).Return(constantURL, nil).Times(1)

				resolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				response := providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
					resolved,
					nil,
				)

				h.On("ParseResponse", expectedIDs, newValidResponse()).Return(response).Times(1)

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", mock.Anything).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), metrics.Success).Times(1)

				return m
			},
			capacity: 1,
			ids:      []slinkytypes.CurrencyPair{btcusd},
			atomic:   false,
			responses: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				},
				UnResolved: map[slinkytypes.CurrencyPair]error{},
			},
		},
		{
			name: "single id to query with rate limit errors and atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newRateLimitResponse(), nil).Times(1)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[slinkytypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []slinkytypes.CurrencyPair{btcusd}).Return(constantURL, nil).Times(1)

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", mock.Anything).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), metrics.RateLimit).Times(1)

				return m
			},
			capacity: 1,
			ids:      []slinkytypes.CurrencyPair{btcusd},
			atomic:   true,
			responses: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				UnResolved: map[slinkytypes.CurrencyPair]error{
					btcusd: errors.ErrRateLimit,
				},
			},
		},
		{
			name: "single id to query with unexpected status code errors and atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newUnexpectedStatusCodeResponse(), nil).Times(1)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[slinkytypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []slinkytypes.CurrencyPair{btcusd}).Return(constantURL, nil).Times(1)

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", mock.Anything).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), metrics.UnexpectedStatusCode).Times(1)

				return m
			},
			capacity: 1,
			ids:      []slinkytypes.CurrencyPair{btcusd},
			atomic:   true,
			responses: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				UnResolved: map[slinkytypes.CurrencyPair]error{
					btcusd: errors.ErrUnexpectedStatusCodeWithCode(http.StatusInternalServerError),
				},
			},
		},
		{
			name: "single id to query with failure to make request and atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(nil, fmt.Errorf("client has no rizz")).Times(1)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[slinkytypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []slinkytypes.CurrencyPair{btcusd}).Return(constantURL, nil).Times(1)

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", mock.Anything).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), metrics.DoRequest).Times(1)

				return m
			},
			capacity: 1,
			ids:      []slinkytypes.CurrencyPair{btcusd},
			atomic:   true,
			responses: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{},
				UnResolved: map[slinkytypes.CurrencyPair]error{
					btcusd: errors.ErrDoRequestWithErr(fmt.Errorf("client has no rizz")),
				},
			},
		},
		{
			name: "multiple ids to query with no errors and atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Times(1)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] {
				expectedIDs := []slinkytypes.CurrencyPair{btcusd, ethusd, atomusd}

				h := mocks.NewAPIDataHandler[slinkytypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", expectedIDs).Return(constantURL, nil).Times(1)

				resolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
					ethusd: {
						Value: big.NewInt(200),
					},
					atomusd: {
						Value: big.NewInt(300),
					},
				}
				response := providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
					resolved,
					nil,
				)

				h.On("ParseResponse", expectedIDs, newValidResponse()).Return(response).Times(1)

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", mock.Anything).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), metrics.Success).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(ethusd)), metrics.Success).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(atomusd)), metrics.Success).Times(1)

				return m
			},
			capacity: 3,
			ids:      []slinkytypes.CurrencyPair{btcusd, ethusd, atomusd},
			atomic:   true,
			responses: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
					ethusd: {
						Value: big.NewInt(200),
					},
					atomusd: {
						Value: big.NewInt(300),
					},
				},
				UnResolved: map[slinkytypes.CurrencyPair]error{},
			},
		},
		{
			name: "multiple ids to query with no errors and non-atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Times(3)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[slinkytypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []slinkytypes.CurrencyPair{btcusd}).Return(constantURL, nil).Times(1)
				h.On("CreateURL", []slinkytypes.CurrencyPair{ethusd}).Return(constantURL, nil).Times(1)
				h.On("CreateURL", []slinkytypes.CurrencyPair{atomusd}).Return(constantURL, nil).Times(1)

				btcResolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				btcResponse := providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
					btcResolved,
					nil,
				)

				ethResolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					ethusd: {
						Value: big.NewInt(200),
					},
				}
				ethResponse := providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
					ethResolved,
					nil,
				)

				atomResolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					atomusd: {
						Value: big.NewInt(300),
					},
				}
				atomResponse := providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
					atomResolved,
					nil,
				)

				h.On("ParseResponse", []slinkytypes.CurrencyPair{btcusd}, newValidResponse()).Return(btcResponse).Times(1)
				h.On("ParseResponse", []slinkytypes.CurrencyPair{ethusd}, newValidResponse()).Return(ethResponse).Times(1)
				h.On("ParseResponse", []slinkytypes.CurrencyPair{atomusd}, newValidResponse()).Return(atomResponse).Times(1)

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", mock.Anything).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), metrics.Success).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(ethusd)), metrics.Success).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(atomusd)), metrics.Success).Times(1)

				return m
			},
			capacity: 3,
			ids:      []slinkytypes.CurrencyPair{btcusd, ethusd, atomusd},
			atomic:   false,
			responses: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
					ethusd: {
						Value: big.NewInt(200),
					},
					atomusd: {
						Value: big.NewInt(300),
					},
				},
				UnResolved: map[slinkytypes.CurrencyPair]error{},
			},
		},
		{
			name: "multiple ids to query with some errors and non-atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Times(2)
				h.On("Do", mock.Anything, constantURL+"eth").Return(newRateLimitResponse(), nil).Times(1)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[slinkytypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []slinkytypes.CurrencyPair{btcusd}).Return(constantURL, nil).Times(1)
				h.On("CreateURL", []slinkytypes.CurrencyPair{ethusd}).Return(constantURL+"eth", nil).Times(1)
				h.On("CreateURL", []slinkytypes.CurrencyPair{atomusd}).Return(constantURL, nil).Times(1)

				btcResolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				btcResponse := providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
					btcResolved,
					nil,
				)

				atomResolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					atomusd: {
						Value: big.NewInt(300),
					},
				}
				atomResponse := providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
					atomResolved,
					nil,
				)

				h.On("ParseResponse", []slinkytypes.CurrencyPair{btcusd}, newValidResponse()).Return(btcResponse).Times(1)
				h.On("ParseResponse", []slinkytypes.CurrencyPair{atomusd}, newValidResponse()).Return(atomResponse).Times(1)

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", mock.Anything).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), metrics.Success).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(ethusd)), metrics.RateLimit).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(atomusd)), metrics.Success).Times(1)

				return m
			},
			capacity: 3,
			ids:      []slinkytypes.CurrencyPair{btcusd, ethusd, atomusd},
			atomic:   false,
			responses: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
					atomusd: {
						Value: big.NewInt(300),
					},
				},
				UnResolved: map[slinkytypes.CurrencyPair]error{
					ethusd: errors.ErrRateLimit,
				},
			},
		},
		{
			name: "multiple ids to query with no errors and non-atomic handler and capacity on concurrent requests",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				// Delay the responses by 1 second to ensure that the requests are made sequentially.
				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Times(3).After(1 * time.Second)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[slinkytypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[slinkytypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []slinkytypes.CurrencyPair{btcusd}).Return(constantURL, nil).Times(1)
				h.On("CreateURL", []slinkytypes.CurrencyPair{ethusd}).Return(constantURL, nil).Times(1)
				h.On("CreateURL", []slinkytypes.CurrencyPair{atomusd}).Return(constantURL, nil).Times(1)

				btcResolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				btcResponse := providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
					btcResolved,
					nil,
				)

				ethResolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					ethusd: {
						Value: big.NewInt(200),
					},
				}
				ethResponse := providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
					ethResolved,
					nil,
				)

				atomResolved := map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					atomusd: {
						Value: big.NewInt(300),
					},
				}
				atomResponse := providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](
					atomResolved,
					nil,
				)

				h.On("ParseResponse", []slinkytypes.CurrencyPair{btcusd}, newValidResponse()).Return(btcResponse).Times(1)
				h.On("ParseResponse", []slinkytypes.CurrencyPair{ethusd}, newValidResponse()).Return(ethResponse).Times(1)
				h.On("ParseResponse", []slinkytypes.CurrencyPair{atomusd}, newValidResponse()).Return(atomResponse).Times(1)

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", mock.Anything).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), metrics.Success).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(ethusd)), metrics.Success).Times(1)
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(atomusd)), metrics.Success).Times(1)

				return m
			},
			capacity: 1,
			ids:      []slinkytypes.CurrencyPair{btcusd, ethusd, atomusd},
			atomic:   false,
			responses: providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int]{
				Resolved: map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
					ethusd: {
						Value: big.NewInt(200),
					},
					atomusd: {
						Value: big.NewInt(300),
					},
				},
				UnResolved: map[slinkytypes.CurrencyPair]error{},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var apiCfg config.APIConfig
			if tc.atomic {
				apiCfg = cfg
			} else {
				apiCfg = nonAtomicCfg
			}

			handler, err := handlers.NewAPIQueryHandler[slinkytypes.CurrencyPair, *big.Int](
				logger,
				apiCfg,
				tc.requestHandler(),
				tc.apiHandler(),
				tc.metrics(),
			)
			require.NoError(t, err)

			responseCh := make(chan providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int], tc.capacity)
			go func() {
				handler.Query(context.Background(), tc.ids, responseCh)
				close(responseCh)
			}()

			expectedResponses := tc.responses
			for resp := range responseCh {
				for id, result := range resp.Resolved {
					require.Equal(t, expectedResponses.Resolved[id], result)
					delete(expectedResponses.Resolved, id)
				}

				for id, err := range resp.UnResolved {
					require.Equal(t, expectedResponses.UnResolved[id], err)
					delete(expectedResponses.UnResolved, id)
				}
			}

			// Ensure all responses are account for.
			require.Empty(t, expectedResponses.Resolved)
			require.Empty(t, expectedResponses.UnResolved)
		})
	}
}

func newRateLimitResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Body:       nil,
	}
}

func newUnexpectedStatusCodeResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       nil,
	}
}

func newValidResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       nil,
	}
}
