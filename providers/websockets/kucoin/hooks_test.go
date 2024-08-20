package kucoin_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/oracle/config"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	apimocks "github.com/skip-mev/connect/v2/providers/base/api/handlers/mocks"
	"github.com/skip-mev/connect/v2/providers/base/testutils"
	wshandlers "github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
	"github.com/skip-mev/connect/v2/providers/websockets/kucoin"
)

var postURL = fmt.Sprintf("%s%s", kucoin.URL, kucoin.BulletPublicEndpoint)

func TestPreDialHook(t *testing.T) {
	testCases := []struct {
		name           string
		requestHandler func() apihandlers.RequestHandler
		expectedConfig *config.WebSocketConfig
		expectErr      bool
	}{
		{
			name: "request handler fails to make the request",
			requestHandler: func() apihandlers.RequestHandler {
				h := apimocks.NewRequestHandler(t)

				h.On("Do", mock.Anything, postURL).Return(nil, fmt.Errorf("error")).Once()

				return h
			},
			expectErr: true,
		},
		{
			name: "request handler returns a non-200 status code",
			requestHandler: func() apihandlers.RequestHandler {
				h := apimocks.NewRequestHandler(t)

				resp := http.Response{
					StatusCode: http.StatusForbidden,
				}
				h.On("Do", mock.Anything, postURL).Return(&resp, nil).Once()

				return h
			},
			expectErr: true,
		},
		{
			name: "request handler cannot decode the response",
			requestHandler: func() apihandlers.RequestHandler {
				h := apimocks.NewRequestHandler(t)

				resp := testutils.CreateResponseFromJSON("invalid json,")
				resp.StatusCode = http.StatusOK

				h.On("Do", mock.Anything, postURL).Return(resp, nil).Once()

				return h
			},
			expectErr: true,
		},
		{
			name: "bullet public response is not successful",
			requestHandler: func() apihandlers.RequestHandler {
				h := apimocks.NewRequestHandler(t)

				resp := testutils.CreateResponseFromJSON(`{"code": "error"}`)
				resp.StatusCode = http.StatusOK

				h.On("Do", mock.Anything, postURL).Return(resp, nil).Once()

				return h
			},
			expectErr: true,
		},
		{
			name: "bullet public response contains no servers",
			requestHandler: func() apihandlers.RequestHandler {
				h := apimocks.NewRequestHandler(t)

				resp := testutils.CreateResponseFromJSON(`{"code": "200000", "data": {}}`)
				resp.StatusCode = http.StatusOK

				h.On("Do", mock.Anything, postURL).Return(resp, nil).Once()

				return h
			},
			expectErr: true,
		},
		{
			name: "bullet public response does not contain a websocket protocol",
			requestHandler: func() apihandlers.RequestHandler {
				h := apimocks.NewRequestHandler(t)

				jsonResp := `
				{
					"code": "200000", 
					"data": {
						"instanceServers": 
							[{
								"endpoint": "wss://ws-api-spot.kucoin.com/",
								"encrypt": true,
								"protocol": "api",
								"pingInterval": 18000,
								"pingTimeout": 10000
							}]
						}
				}`
				resp := testutils.CreateResponseFromJSON(jsonResp)
				resp.StatusCode = http.StatusOK

				h.On("Do", mock.Anything, postURL).Return(resp, nil).Once()

				return h
			},
			expectErr: true,
		},
		{
			name: "bullet public response contains a websocket protocol",
			requestHandler: func() apihandlers.RequestHandler {
				h := apimocks.NewRequestHandler(t)

				jsonResp := `
				{
					"code": "200000",
					"data": {
						"token": "abcd",
						"instanceServers": 
							[{
								"endpoint": "wss://ws-api-spot.kucoin.com/",
								"encrypt": true,
								"protocol": "websocket",
								"pingInterval": 20000,
								"pingTimeout": 15000
							}]
						}
				}`
				resp := testutils.CreateResponseFromJSON(jsonResp)
				resp.StatusCode = http.StatusOK

				h.On("Do", mock.Anything, postURL).Return(resp, nil).Once()

				return h
			},
			expectedConfig: &config.WebSocketConfig{
				Enabled: true,
				Endpoints: []config.Endpoint{
					{
						URL: "wss://ws-api-spot.kucoin.com/?token=abcd", // this is dynamically generated.
					},
				},
				ReadTimeout:  15 * time.Second, // this is dynamically generated
				PingInterval: 20 * time.Second, // this is dynamically generated
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			connHandler, err := wshandlers.NewWebSocketHandlerImpl(kucoin.DefaultWebSocketConfig)
			require.NoError(t, err)

			hook := kucoin.PreDialHook(kucoin.DefaultAPIConfig, tc.requestHandler())
			err = hook(connHandler)
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			cfg := connHandler.GetConfig()
			require.Equal(t, tc.expectedConfig.Endpoints, cfg.Endpoints)
			require.Equal(t, tc.expectedConfig.ReadTimeout, cfg.ReadTimeout)
			require.Equal(t, tc.expectedConfig.PingInterval, cfg.PingInterval)
		})
	}
}
