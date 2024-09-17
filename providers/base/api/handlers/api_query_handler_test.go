package handlers_test

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/providers/base/api/errors"
	"github.com/skip-mev/connect/v2/providers/base/api/handlers"
	"github.com/skip-mev/connect/v2/providers/base/api/handlers/mocks"
	"github.com/skip-mev/connect/v2/providers/base/api/metrics"
	mockmetrics "github.com/skip-mev/connect/v2/providers/base/api/metrics/mocks"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var (
	logger  = zap.NewExample()
	btcusd  = connecttypes.NewCurrencyPair("BTC", "USD")
	ethusd  = connecttypes.NewCurrencyPair("ETH", "USD")
	atomusd = connecttypes.NewCurrencyPair("ATOM", "USD")

	constantURL = "http://fetchdata.org:8080"

	cfg = config.APIConfig{
		Enabled:          true,
		Timeout:          500 * time.Millisecond,
		Interval:         250 * time.Millisecond,
		ReconnectTimeout: 250 * time.Millisecond,
		MaxQueries:       1,
		Atomic:           true,
		Endpoints:        []config.Endpoint{{URL: constantURL}},
		Name:             "handler1",
	}

	nonAtomicCfg = config.APIConfig{
		Enabled:          true,
		Timeout:          500 * time.Millisecond,
		Interval:         250 * time.Millisecond,
		ReconnectTimeout: 250 * time.Millisecond,
		MaxQueries:       3,
		Atomic:           false,
		Endpoints:        []config.Endpoint{{URL: constantURL}},
		Name:             "handler1",
	}
)

func TestAPIQueryHandler(t *testing.T) {
	testCases := []struct {
		name           string
		requestHandler func() handlers.RequestHandler
		apiHandler     func() handlers.APIDataHandler[connecttypes.CurrencyPair, *big.Int]
		metrics        func() metrics.APIMetrics
		ids            []connecttypes.CurrencyPair
		atomic         bool
		responses      providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]
	}{
		{
			name: "no ids to query",
			requestHandler: func() handlers.RequestHandler {
				return mocks.NewRequestHandler(t)
			},
			apiHandler: func() handlers.APIDataHandler[connecttypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[connecttypes.CurrencyPair, *big.Int](t)
				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)
				return m
			},
			ids:    []connecttypes.CurrencyPair{},
			atomic: false,
			responses: providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
				Resolved:   map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{},
				UnResolved: map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{},
			},
		},
		{
			name: "single id to query with no errors and an atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Maybe().After(1 * time.Second)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[connecttypes.CurrencyPair, *big.Int] {
				expectedIDs := []connecttypes.CurrencyPair{btcusd}

				h := mocks.NewAPIDataHandler[connecttypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", expectedIDs).Return(constantURL, nil).Maybe()

				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				response := providertypes.NewGetResponse(
					resolved,
					nil,
				)

				h.On("ParseResponse", expectedIDs, newValidResponse()).Return(response).Maybe()

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", metrics.RedactedURL, mock.Anything).Maybe()
				m.On("AddHTTPStatusCode", "handler1", mock.Anything).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), providertypes.OK).Maybe()

				return m
			},
			ids:    []connecttypes.CurrencyPair{btcusd},
			atomic: true,
			responses: providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
				Resolved: map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				},
				UnResolved: map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{},
			},
		},
		{
			name: "single id to query with no errors and a non-atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Maybe().After(1 * time.Second)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[connecttypes.CurrencyPair, *big.Int] {
				expectedIDs := []connecttypes.CurrencyPair{btcusd}

				h := mocks.NewAPIDataHandler[connecttypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", expectedIDs).Return(constantURL, nil).Maybe()

				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				response := providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](
					resolved,
					nil,
				)

				h.On("ParseResponse", expectedIDs, newValidResponse()).Return(response).Maybe()

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", metrics.RedactedURL, mock.Anything).Maybe()
				m.On("AddHTTPStatusCode", "handler1", mock.Anything).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), providertypes.OK).Maybe()

				return m
			},
			ids:    []connecttypes.CurrencyPair{btcusd},
			atomic: false,
			responses: providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
				Resolved: map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				},
				UnResolved: map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{},
			},
		},
		{
			name: "single id to query with rate limit errors and atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newRateLimitResponse(), nil).Maybe().After(1 * time.Second)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[connecttypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[connecttypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []connecttypes.CurrencyPair{btcusd}).Return(constantURL, nil).Maybe()

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", metrics.RedactedURL, mock.Anything).Maybe()
				m.On("AddHTTPStatusCode", "handler1", mock.Anything).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), providertypes.ErrorRateLimitExceeded).Maybe()
				return m
			},
			ids:    []connecttypes.CurrencyPair{btcusd},
			atomic: true,
			responses: providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
				Resolved: map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{},
				UnResolved: map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{
					btcusd: {
						ErrorWithCode: providertypes.NewErrorWithCode(errors.ErrRateLimit, providertypes.ErrorRateLimitExceeded),
					},
				},
			},
		},
		{
			name: "single id to query with unexpected status code errors and atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newUnexpectedStatusCodeResponse(), nil).Maybe().After(1 * time.Second)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[connecttypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[connecttypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []connecttypes.CurrencyPair{btcusd}).Return(constantURL, nil).Maybe()

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", metrics.RedactedURL, mock.Anything).Maybe()
				m.On("AddHTTPStatusCode", "handler1", mock.Anything).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), mock.Anything).Maybe()

				return m
			},
			ids:    []connecttypes.CurrencyPair{btcusd},
			atomic: true,
			responses: providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
				Resolved: map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{},
				UnResolved: map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{
					btcusd: {
						ErrorWithCode: providertypes.NewErrorWithCode(errors.ErrUnexpectedStatusCodeWithCode(http.StatusInternalServerError), providertypes.ErrorUnknown),
					},
				},
			},
		},
		{
			name: "single id to query with failure to make request and atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(nil, fmt.Errorf("client has no rizz")).Maybe().After(1 * time.Second)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[connecttypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[connecttypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []connecttypes.CurrencyPair{btcusd}).Return(constantURL, nil).Maybe()

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", metrics.RedactedURL, mock.Anything).Maybe()
				m.On("AddHTTPStatusCode", "handler1", mock.Anything).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), providertypes.ErrorUnknown).Maybe()

				return m
			},
			ids:    []connecttypes.CurrencyPair{btcusd},
			atomic: true,
			responses: providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
				Resolved: map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{},
				UnResolved: map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{
					btcusd: {
						ErrorWithCode: providertypes.NewErrorWithCode(errors.ErrDoRequestWithErr(fmt.Errorf("client has no rizz")), providertypes.ErrorUnknown),
					},
				},
			},
		},
		{
			name: "multiple ids to query with no errors and atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Maybe().After(1 * time.Second)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[connecttypes.CurrencyPair, *big.Int] {
				expectedIDs := []connecttypes.CurrencyPair{btcusd, ethusd, atomusd}

				h := mocks.NewAPIDataHandler[connecttypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", expectedIDs).Return(constantURL, nil).Maybe()

				resolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
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
				response := providertypes.NewGetResponse(
					resolved,
					nil,
				)

				h.On("ParseResponse", expectedIDs, newValidResponse()).Return(response).Maybe()

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", metrics.RedactedURL, mock.Anything).Maybe()
				m.On("AddHTTPStatusCode", "handler1", mock.Anything).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), providertypes.OK).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(ethusd)), providertypes.OK).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(atomusd)), providertypes.OK).Maybe()

				return m
			},
			ids:    []connecttypes.CurrencyPair{btcusd, ethusd, atomusd},
			atomic: true,
			responses: providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
				Resolved: map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
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
				UnResolved: map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{},
			},
		},
		{
			name: "multiple ids to query with no errors and non-atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Maybe().After(1 * time.Second)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[connecttypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[connecttypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []connecttypes.CurrencyPair{btcusd}).Return(constantURL, nil).Maybe()
				h.On("CreateURL", []connecttypes.CurrencyPair{ethusd}).Return(constantURL, nil).Maybe()
				h.On("CreateURL", []connecttypes.CurrencyPair{atomusd}).Return(constantURL, nil).Maybe()

				btcResolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				btcResponse := providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](
					btcResolved,
					nil,
				)

				ethResolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					ethusd: {
						Value: big.NewInt(200),
					},
				}
				ethResponse := providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](
					ethResolved,
					nil,
				)

				atomResolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					atomusd: {
						Value: big.NewInt(300),
					},
				}
				atomResponse := providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](
					atomResolved,
					nil,
				)

				h.On("ParseResponse", []connecttypes.CurrencyPair{btcusd}, newValidResponse()).Return(btcResponse).Maybe()
				h.On("ParseResponse", []connecttypes.CurrencyPair{ethusd}, newValidResponse()).Return(ethResponse).Maybe()
				h.On("ParseResponse", []connecttypes.CurrencyPair{atomusd}, newValidResponse()).Return(atomResponse).Maybe()

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", metrics.RedactedURL, mock.Anything).Maybe()
				m.On("AddHTTPStatusCode", "handler1", mock.Anything).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), providertypes.OK).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(ethusd)), providertypes.OK).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(atomusd)), providertypes.OK).Maybe()

				return m
			},
			ids:    []connecttypes.CurrencyPair{btcusd, ethusd, atomusd},
			atomic: false,
			responses: providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
				Resolved: map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
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
				UnResolved: map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{},
			},
		},
		{
			name: "multiple ids to query with some errors and non-atomic handler",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Maybe().After(1 * time.Second)
				h.On("Do", mock.Anything, constantURL+"eth").Return(newRateLimitResponse(), nil).Maybe().After(1 * time.Second)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[connecttypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[connecttypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []connecttypes.CurrencyPair{btcusd}).Return(constantURL, nil).Maybe()
				h.On("CreateURL", []connecttypes.CurrencyPair{ethusd}).Return(constantURL+"eth", nil).Maybe()
				h.On("CreateURL", []connecttypes.CurrencyPair{atomusd}).Return(constantURL, nil).Maybe()

				btcResolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				btcResponse := providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](
					btcResolved,
					nil,
				)

				atomResolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					atomusd: {
						Value: big.NewInt(300),
					},
				}
				atomResponse := providertypes.NewGetResponse(
					atomResolved,
					nil,
				)

				h.On("ParseResponse", []connecttypes.CurrencyPair{btcusd}, newValidResponse()).Return(btcResponse).Maybe()
				h.On("ParseResponse", []connecttypes.CurrencyPair{atomusd}, newValidResponse()).Return(atomResponse).Maybe()

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", metrics.RedactedURL, mock.Anything).Maybe()
				m.On("AddHTTPStatusCode", "handler1", mock.Anything).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), providertypes.OK).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(ethusd)), providertypes.ErrorRateLimitExceeded).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(atomusd)), providertypes.OK).Maybe()

				return m
			},
			ids:    []connecttypes.CurrencyPair{btcusd, ethusd, atomusd},
			atomic: false,
			responses: providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
				Resolved: map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
					atomusd: {
						Value: big.NewInt(300),
					},
				},
				UnResolved: map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{
					ethusd: {
						ErrorWithCode: providertypes.NewErrorWithCode(errors.ErrRateLimit, providertypes.ErrorRateLimitExceeded),
					},
				},
			},
		},
		{
			name: "multiple ids to query with no errors and non-atomic handler and capacity on concurrent requests",
			requestHandler: func() handlers.RequestHandler {
				h := mocks.NewRequestHandler(t)

				// Delay the responses by 1 second to ensure that the requests are made sequentially.
				h.On("Do", mock.Anything, constantURL).Return(newValidResponse(), nil).Maybe().After(1 * time.Second)

				return h
			},
			apiHandler: func() handlers.APIDataHandler[connecttypes.CurrencyPair, *big.Int] {
				h := mocks.NewAPIDataHandler[connecttypes.CurrencyPair, *big.Int](t)

				h.On("CreateURL", []connecttypes.CurrencyPair{btcusd}).Return(constantURL, nil).Maybe()
				h.On("CreateURL", []connecttypes.CurrencyPair{ethusd}).Return(constantURL, nil).Maybe()
				h.On("CreateURL", []connecttypes.CurrencyPair{atomusd}).Return(constantURL, nil).Maybe()

				btcResolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				btcResponse := providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](
					btcResolved,
					nil,
				)

				ethResolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					ethusd: {
						Value: big.NewInt(200),
					},
				}
				ethResponse := providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](
					ethResolved,
					nil,
				)

				atomResolved := map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
					atomusd: {
						Value: big.NewInt(300),
					},
				}
				atomResponse := providertypes.NewGetResponse[connecttypes.CurrencyPair, *big.Int](
					atomResolved,
					nil,
				)

				h.On("ParseResponse", []connecttypes.CurrencyPair{btcusd}, newValidResponse()).Return(btcResponse).Maybe()
				h.On("ParseResponse", []connecttypes.CurrencyPair{ethusd}, newValidResponse()).Return(ethResponse).Maybe()
				h.On("ParseResponse", []connecttypes.CurrencyPair{atomusd}, newValidResponse()).Return(atomResponse).Maybe()

				return h
			},
			metrics: func() metrics.APIMetrics {
				m := mockmetrics.NewAPIMetrics(t)

				m.On("ObserveProviderResponseLatency", "handler1", metrics.RedactedURL, mock.Anything).Maybe()
				m.On("AddHTTPStatusCode", "handler1", mock.Anything).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(btcusd)), providertypes.OK).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(ethusd)), providertypes.OK).Maybe()
				m.On("AddProviderResponse", "handler1", strings.ToLower(fmt.Sprint(atomusd)), providertypes.OK).Maybe()

				return m
			},
			ids:    []connecttypes.CurrencyPair{btcusd, ethusd, atomusd},
			atomic: false,
			responses: providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int]{
				Resolved: map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int]{
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
				UnResolved: map[connecttypes.CurrencyPair]providertypes.UnresolvedResult{},
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

			handler, err := handlers.NewAPIQueryHandler[connecttypes.CurrencyPair, *big.Int](
				logger,
				apiCfg,
				tc.requestHandler(),
				tc.apiHandler(),
				tc.metrics(),
			)
			require.NoError(t, err)

			responseCh := make(chan providertypes.GetResponse[connecttypes.CurrencyPair, *big.Int], len(tc.ids))

			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			go func() {
				handler.Query(ctx, tc.ids, responseCh)
				close(responseCh)
			}()

			expectedResponses := tc.responses
			resolved := make(map[connecttypes.CurrencyPair]providertypes.ResolvedResult[*big.Int])
			unResolved := make(map[connecttypes.CurrencyPair]providertypes.UnresolvedResult)
			for resp := range responseCh {
				for id, result := range resp.Resolved {
					require.Equal(t, expectedResponses.Resolved[id], result)
					resolved[id] = result
				}

				for id, result := range resp.UnResolved {
					require.Equal(t, expectedResponses.UnResolved[id].Error(), result.Error())
					unResolved[id] = result
				}
			}

			// Ensure all responses are account for.
			require.Equal(t, len(tc.ids), len(resolved)+len(unResolved))
			require.Equal(t, len(expectedResponses.Resolved), len(resolved))
			require.Equal(t, len(expectedResponses.UnResolved), len(unResolved))
		})
	}
}

func TestAPIQueryHandlerWithBatchSize(t *testing.T) {
	cfg = config.APIConfig{
		Enabled:          true,
		Timeout:          500 * time.Millisecond,
		Interval:         250 * time.Millisecond,
		ReconnectTimeout: 250 * time.Millisecond,
		MaxQueries:       3,
		BatchSize:        2,
		Endpoints:        []config.Endpoint{{URL: constantURL}},
		Name:             "handler1",
	}

	pf := mocks.NewAPIFetcher[mmtypes.Ticker, *big.Int](t)

	handler, err := handlers.NewAPIQueryHandlerWithFetcher(
		zap.NewNop(),
		cfg,
		pf,
		metrics.NewNopAPIMetrics(),
	)
	require.NoError(t, err)

	t.Run("Query with batch-size correctly batches requests", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// create response channel
		responseCh := make(chan providertypes.GetResponse[mmtypes.Ticker, *big.Int], 3)

		// mock 3 executions to price-fetcher
		queriedTickers := map[string]bool{
			"BTC/USD":  false,
			"BTC1/USD": false,
			"BTC2/USD": false,
			"BTC3/USD": false,
			"BTC4/USD": false,
		}
		mtx := sync.Mutex{}
		pf.On("Fetch", mock.Anything, mock.Anything).Return(providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](nil, nil)).Run(func(args mock.Arguments) {
			// expect 2 executions w/ 2 arguments and 1 with 1 argument
			tickers := args.Get(1).([]mmtypes.Ticker)

			if !(len(tickers) == 1 || len(tickers) == 2) {
				t.Errorf("unexpected number of arguments: %d", len(args))
			}
			// mark tickers as queried
			for _, ticker := range tickers {
				mtx.Lock()
				if _, ok := queriedTickers[ticker.String()]; !ok {
					t.Errorf("unexpected ticker queried: %s", ticker.String())
				}

				queriedTickers[ticker.String()] = true
				mtx.Unlock()
			}
		})

		// query
		done := make(chan struct{})
		go func() {
			handler.Query(ctx, []mmtypes.Ticker{
				mmtypes.NewTicker("BTC", "USD", 8, 0, true),
				mmtypes.NewTicker("BTC1", "USD", 8, 0, true),
				mmtypes.NewTicker("BTC2", "USD", 8, 0, true),
				mmtypes.NewTicker("BTC3", "USD", 8, 0, true),
				mmtypes.NewTicker("BTC4", "USD", 8, 0, true),
			}, responseCh)
			close(done)
		}()

		// wait for response
		numResponses := 0
		for range responseCh {
			numResponses++
			if numResponses == 3 {
				break
			}
		}

		// close the handler
		cancel()
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("handler did not close")
		}

		// assert
		for ticker, queried := range queriedTickers {
			if !queried {
				t.Errorf("ticker not queried: %s", ticker)
			}
		}
	})
}

func newRateLimitResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Body:       io.NopCloser(strings.NewReader(`{"error": "rate limit exceeded"}`)),
	}
}

func newUnexpectedStatusCodeResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(strings.NewReader(`{"error": "unexpected error"}`)),
	}
}

func newValidResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(`{"result": "100"}`)),
	}
}
