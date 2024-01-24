package handlers_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	wserrors "github.com/skip-mev/slinky/providers/base/websocket/errors"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	handlermocks "github.com/skip-mev/slinky/providers/base/websocket/handlers/mocks"
	"github.com/skip-mev/slinky/providers/base/websocket/metrics"
	mockmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	logger  = zap.NewExample()
	btcusd  = oracletypes.NewCurrencyPair("BTC", "USD")
	ethusd  = oracletypes.NewCurrencyPair("ETH", "USD")
	atomusd = oracletypes.NewCurrencyPair("ATOM", "USD")

	name        = "sirmoggintonwebsocket"
	testMessage = []byte("gib me money")

	cfg = config.WebSocketConfig{
		Name:                "sirmoggintonwebsocket",
		WSS:                 "ws://localhost:8080",
		Enabled:             true,
		MaxBufferSize:       1024,
		ReconnectionTimeout: 5 * time.Second,
		ReadBufferSize:      config.DefaultReadBufferSize,
		WriteBufferSize:     config.DefaultWriteBufferSize,
		HandshakeTimeout:    config.DefaultHandshakeTimeout,
		EnableCompression:   config.DefaultEnableCompression,
		ReadTimeout:         config.DefaultReadTimeout,
		WriteTimeout:        config.DefaultWriteTimeout,
	}
)

func TestWebSocketQueryHandler(t *testing.T) {
	testCases := []struct {
		name        string
		connHandler func() handlers.WebSocketConnHandler
		dataHandler func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int]
		metrics     func() metrics.WebSocketMetrics
		ids         []oracletypes.CurrencyPair
		responses   providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]
	}{
		{
			name: "fails to dial the websocket",
			connHandler: func() handlers.WebSocketConnHandler {
				connHandler := handlermocks.NewWebSocketConnHandler(t)

				connHandler.On("Dial", mock.Anything).Return(fmt.Errorf("no rizz alert")).Once()

				return connHandler
			},
			dataHandler: func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] {
				dataHandler := handlermocks.NewWebSocketDataHandler[oracletypes.CurrencyPair, *big.Int](t)
				return dataHandler
			},
			metrics: func() metrics.WebSocketMetrics {
				m := mockmetrics.NewWebSocketMetrics(t)

				m.On("AddWebSocketConnectionStatus", name, metrics.DialErr).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Unhealthy).Return().Once()

				return m
			},
			ids: []oracletypes.CurrencyPair{btcusd},
			responses: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
				UnResolved: map[oracletypes.CurrencyPair]error{
					btcusd: wserrors.ErrDial,
				},
			},
		},
		{
			name: "fails to create subscriptions",
			connHandler: func() handlers.WebSocketConnHandler {
				connHandler := handlermocks.NewWebSocketConnHandler(t)

				connHandler.On("Dial", mock.Anything).Return(nil).Once()
				connHandler.On("Heartbeat", mock.Anything).Return(nil).Maybe()

				return connHandler
			},
			dataHandler: func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] {
				dataHandler := handlermocks.NewWebSocketDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				dataHandler.On("CreateMessages", mock.Anything).Return(nil, fmt.Errorf("no rizz alert")).Once()

				return dataHandler
			},
			metrics: func() metrics.WebSocketMetrics {
				m := mockmetrics.NewWebSocketMetrics(t)

				m.On("AddWebSocketConnectionStatus", name, metrics.DialSuccess).Return().Once()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageErr).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Unhealthy).Return().Once()

				return m
			},
			ids: []oracletypes.CurrencyPair{btcusd},
			responses: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
				UnResolved: map[oracletypes.CurrencyPair]error{
					btcusd: wserrors.ErrCreateMessages,
				},
			},
		},
		{
			name: "fails to write to the websocket on initial subscription",
			connHandler: func() handlers.WebSocketConnHandler {
				connHandler := handlermocks.NewWebSocketConnHandler(t)

				connHandler.On("Dial", mock.Anything).Return(nil).Once()
				connHandler.On("Write", mock.Anything).Return(fmt.Errorf("no rizz alert")).Once()
				connHandler.On("Heartbeat", mock.Anything).Return(nil).Maybe()

				return connHandler
			},
			dataHandler: func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] {
				dataHandler := handlermocks.NewWebSocketDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				dataHandler.On("CreateMessages", mock.Anything).Return([]handlers.WebsocketEncodedMessage{testMessage}, nil).Once()

				return dataHandler
			},
			metrics: func() metrics.WebSocketMetrics {
				m := mockmetrics.NewWebSocketMetrics(t)

				m.On("AddWebSocketConnectionStatus", name, metrics.DialSuccess).Return().Once()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteErr).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Unhealthy).Return().Once()

				return m
			},
			ids: []oracletypes.CurrencyPair{btcusd},
			responses: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
				UnResolved: map[oracletypes.CurrencyPair]error{
					btcusd: wserrors.ErrWrite,
				},
			},
		},
		{
			name: "fails to read from the websocket",
			connHandler: func() handlers.WebSocketConnHandler {
				connHandler := handlermocks.NewWebSocketConnHandler(t)

				connHandler.On("Dial", mock.Anything).Return(nil).Once()
				connHandler.On("Write", mock.Anything).Return(nil).Once()
				connHandler.On("Read").Return(nil, fmt.Errorf("no rizz alert")).Maybe().After(time.Second)
				connHandler.On("Close").Return(nil).Once()
				connHandler.On("Heartbeat", mock.Anything).Return(nil).Maybe()

				return connHandler
			},
			dataHandler: func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] {
				dataHandler := handlermocks.NewWebSocketDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				dataHandler.On("CreateMessages", mock.Anything).Return([]handlers.WebsocketEncodedMessage{testMessage}, nil).Once()

				return dataHandler
			},
			metrics: func() metrics.WebSocketMetrics {
				m := mockmetrics.NewWebSocketMetrics(t)

				m.On("AddWebSocketConnectionStatus", name, metrics.DialSuccess).Return().Once()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Healthy).Return().Once()

				m.On("AddWebSocketConnectionStatus", name, metrics.ReadErr).Return().Maybe()
				m.On("ObserveWebSocketLatency", name, mock.Anything).Return().Maybe()

				m.On("AddWebSocketConnectionStatus", name, metrics.CloseSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Unhealthy).Return().Once()

				return m
			},
			ids:       []oracletypes.CurrencyPair{btcusd},
			responses: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
		},
		{
			name: "fails to parse the response from the websocket",
			connHandler: func() handlers.WebSocketConnHandler {
				connHandler := handlermocks.NewWebSocketConnHandler(t)

				connHandler.On("Dial", mock.Anything).Return(nil).Once()
				connHandler.On("Write", mock.Anything).Return(nil).Once()
				connHandler.On("Read").Return(testMessage, nil).Maybe().After(time.Second)
				connHandler.On("Close").Return(nil).Once()
				connHandler.On("Heartbeat", mock.Anything).Return(nil).Maybe()

				return connHandler
			},
			dataHandler: func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] {
				dataHandler := handlermocks.NewWebSocketDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				dataHandler.On("CreateMessages", mock.Anything).Return([]handlers.WebsocketEncodedMessage{testMessage}, nil).Once()
				dataHandler.On("HandleMessage", mock.Anything).Return(
					providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, nil),
					nil,
					fmt.Errorf("no rizz alert"),
				).Maybe()

				return dataHandler
			},
			metrics: func() metrics.WebSocketMetrics {
				m := mockmetrics.NewWebSocketMetrics(t)

				m.On("AddWebSocketConnectionStatus", name, metrics.DialSuccess).Return().Once()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Healthy).Return().Once()

				m.On("AddWebSocketConnectionStatus", name, metrics.ReadSuccess).Return().Maybe()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.HandleMessageErr).Return().Maybe()
				m.On("ObserveWebSocketLatency", name, mock.Anything).Return().Maybe()

				m.On("AddWebSocketConnectionStatus", name, metrics.CloseSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Unhealthy).Return().Once()

				return m
			},
			ids:       []oracletypes.CurrencyPair{btcusd},
			responses: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
		},
		{
			name: "pseudo heart beat update message with no response",
			connHandler: func() handlers.WebSocketConnHandler {
				connHandler := handlermocks.NewWebSocketConnHandler(t)

				connHandler.On("Dial", mock.Anything).Return(nil).Once()
				connHandler.On("Write", mock.Anything).Return(nil).Maybe()
				connHandler.On("Read").Return(testMessage, nil).Maybe().After(time.Second)
				connHandler.On("Close").Return(nil).Once()
				connHandler.On("Heartbeat", mock.Anything).Return(nil).Maybe()

				return connHandler
			},
			dataHandler: func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] {
				dataHandler := handlermocks.NewWebSocketDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				dataHandler.On("CreateMessages", mock.Anything).Return([]handlers.WebsocketEncodedMessage{testMessage}, nil).Once()
				dataHandler.On("HandleMessage", mock.Anything).Return(
					providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, nil),
					[]handlers.WebsocketEncodedMessage{[]byte("hearb eat")},
					nil,
				).Maybe()

				return dataHandler
			},
			metrics: func() metrics.WebSocketMetrics {
				m := mockmetrics.NewWebSocketMetrics(t)

				m.On("AddWebSocketConnectionStatus", name, metrics.DialSuccess).Return().Once()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Healthy).Return().Once()

				m.On("AddWebSocketConnectionStatus", name, metrics.ReadSuccess).Return().Maybe()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.HandleMessageSuccess).Return().Maybe()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Maybe()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteSuccess).Return().Maybe()
				m.On("ObserveWebSocketLatency", name, mock.Anything).Return().Maybe()

				m.On("AddWebSocketConnectionStatus", name, metrics.CloseSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Unhealthy).Return().Once()

				return m
			},
			ids:       []oracletypes.CurrencyPair{btcusd},
			responses: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
		},
		{
			name: "fails to send the update message to the websocket",
			connHandler: func() handlers.WebSocketConnHandler {
				connHandler := handlermocks.NewWebSocketConnHandler(t)

				connHandler.On("Dial", mock.Anything).Return(nil).Once()
				connHandler.On("Read").Return(testMessage, nil).Maybe().After(time.Second)
				connHandler.On("Close").Return(nil).Once()
				connHandler.On("Write", mock.Anything).Return(nil).Once()
				connHandler.On("Write", mock.Anything).Return(fmt.Errorf("no rizz alert")).Maybe()
				connHandler.On("Heartbeat", mock.Anything).Return(nil).Maybe()

				return connHandler
			},
			dataHandler: func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] {
				dataHandler := handlermocks.NewWebSocketDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				dataHandler.On("CreateMessages", mock.Anything).Return([]handlers.WebsocketEncodedMessage{testMessage}, nil).Once()
				dataHandler.On("HandleMessage", mock.Anything).Return(
					providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, nil),
					[]handlers.WebsocketEncodedMessage{[]byte("hearb eat")},
					nil,
				).Maybe()

				return dataHandler
			},
			metrics: func() metrics.WebSocketMetrics {
				m := mockmetrics.NewWebSocketMetrics(t)

				m.On("AddWebSocketConnectionStatus", name, metrics.DialSuccess).Return().Once()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.CloseSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Healthy).Return().Once()

				m.On("AddWebSocketConnectionStatus", name, metrics.ReadSuccess).Return().Maybe()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.HandleMessageSuccess).Return().Maybe()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Maybe()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteErr).Return().Maybe()
				m.On("ObserveWebSocketLatency", name, mock.Anything).Return().Maybe()

				m.On("AddWebSocketConnectionStatus", name, metrics.CloseSuccess).Return().Maybe()
				m.On("AddWebSocketConnectionStatus", name, metrics.Unhealthy).Return().Once()

				return m
			},
			ids:       []oracletypes.CurrencyPair{btcusd},
			responses: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
		},
		{
			name: "fails to close the websocket",
			connHandler: func() handlers.WebSocketConnHandler {
				connHandler := handlermocks.NewWebSocketConnHandler(t)

				connHandler.On("Dial", mock.Anything).Return(nil).Once()
				connHandler.On("Read").Return(testMessage, nil).Maybe().After(time.Second)
				connHandler.On("Close").Return(fmt.Errorf("no rizz alert")).Once()
				connHandler.On("Write", mock.Anything).Return(nil).Maybe()
				connHandler.On("Heartbeat", mock.Anything).Return(nil).Maybe()

				return connHandler
			},
			dataHandler: func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] {
				dataHandler := handlermocks.NewWebSocketDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				dataHandler.On("CreateMessages", mock.Anything).Return([]handlers.WebsocketEncodedMessage{testMessage}, nil).Once()
				dataHandler.On("HandleMessage", mock.Anything).Return(
					providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, nil),
					[]handlers.WebsocketEncodedMessage{[]byte("hearb eat")},
					nil,
				).Maybe()

				return dataHandler
			},
			metrics: func() metrics.WebSocketMetrics {
				m := mockmetrics.NewWebSocketMetrics(t)

				m.On("AddWebSocketConnectionStatus", name, metrics.DialSuccess).Return().Once()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Healthy).Return().Once()

				m.On("AddWebSocketConnectionStatus", name, metrics.ReadSuccess).Return().Maybe()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.HandleMessageSuccess).Return().Maybe()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Maybe()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteSuccess).Return().Maybe()
				m.On("ObserveWebSocketLatency", name, mock.Anything).Return().Maybe()

				m.On("AddWebSocketConnectionStatus", name, metrics.CloseErr).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Unhealthy).Return().Once()

				return m
			},
			ids:       []oracletypes.CurrencyPair{btcusd},
			responses: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{},
		},
		{
			name: "returns a single response with no update message",
			connHandler: func() handlers.WebSocketConnHandler {
				connHandler := handlermocks.NewWebSocketConnHandler(t)

				connHandler.On("Dial", mock.Anything).Return(nil).Once()
				connHandler.On("Read").Return(testMessage, nil).Maybe().After(time.Second)
				connHandler.On("Close").Return(nil).Once()
				connHandler.On("Write", mock.Anything).Return(nil).Once()
				connHandler.On("Heartbeat", mock.Anything).Return(nil).Maybe()

				return connHandler
			},
			dataHandler: func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] {
				dataHandler := handlermocks.NewWebSocketDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				dataHandler.On("CreateMessages", mock.Anything).Return([]handlers.WebsocketEncodedMessage{testMessage}, nil).Once()

				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				response := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)
				dataHandler.On("HandleMessage", mock.Anything).Return(
					response,
					nil,
					nil,
				).Maybe()

				return dataHandler
			},
			metrics: func() metrics.WebSocketMetrics {
				m := mockmetrics.NewWebSocketMetrics(t)

				m.On("AddWebSocketConnectionStatus", name, metrics.DialSuccess).Return().Once()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Healthy).Return().Once()

				m.On("AddWebSocketConnectionStatus", name, metrics.ReadSuccess).Return().Maybe()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.HandleMessageSuccess).Return().Maybe()
				m.On("ObserveWebSocketLatency", name, mock.Anything).Return().Maybe()

				m.On("AddWebSocketConnectionStatus", name, metrics.CloseSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Unhealthy).Return().Once()

				return m
			},
			ids: []oracletypes.CurrencyPair{btcusd},
			responses: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
				Resolved: map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				},
			},
		},
		{
			name: "returns a single response with an update message",
			connHandler: func() handlers.WebSocketConnHandler {
				connHandler := handlermocks.NewWebSocketConnHandler(t)

				connHandler.On("Dial", mock.Anything).Return(nil).Once()
				connHandler.On("Read").Return(testMessage, nil).Maybe().After(time.Second)
				connHandler.On("Close").Return(nil).Once()
				connHandler.On("Write", mock.Anything).Return(nil).Maybe()
				connHandler.On("Heartbeat", mock.Anything).Return(nil).Maybe()

				return connHandler
			},
			dataHandler: func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] {
				dataHandler := handlermocks.NewWebSocketDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				dataHandler.On("CreateMessages", mock.Anything).Return([]handlers.WebsocketEncodedMessage{testMessage}, nil).Once()

				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				response := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)
				dataHandler.On("HandleMessage", mock.Anything).Return(
					response,
					[]handlers.WebsocketEncodedMessage{[]byte("hearb eat")},
					nil,
				).Maybe()

				return dataHandler
			},
			metrics: func() metrics.WebSocketMetrics {
				m := mockmetrics.NewWebSocketMetrics(t)

				m.On("AddWebSocketConnectionStatus", name, metrics.DialSuccess).Return().Once()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Healthy).Return().Once()

				m.On("AddWebSocketConnectionStatus", name, metrics.ReadSuccess).Return().Maybe()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.HandleMessageSuccess).Return().Maybe()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Maybe()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteSuccess).Return().Maybe()
				m.On("ObserveWebSocketLatency", name, mock.Anything).Return().Maybe()

				m.On("AddWebSocketConnectionStatus", name, metrics.CloseSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Unhealthy).Return().Once()

				return m
			},
			ids: []oracletypes.CurrencyPair{btcusd},
			responses: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
				Resolved: map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				},
			},
		},
		{
			name: "returns multiple responses with no update message",
			connHandler: func() handlers.WebSocketConnHandler {
				connHandler := handlermocks.NewWebSocketConnHandler(t)

				connHandler.On("Dial", mock.Anything).Return(nil).Once()
				connHandler.On("Read").Return(testMessage, nil).Maybe().After(time.Second)
				connHandler.On("Close").Return(nil).Once()
				connHandler.On("Write", mock.Anything).Return(nil).Once()
				connHandler.On("Heartbeat", mock.Anything).Return(nil).Maybe()

				return connHandler
			},
			dataHandler: func() handlers.WebSocketDataHandler[oracletypes.CurrencyPair, *big.Int] {
				dataHandler := handlermocks.NewWebSocketDataHandler[oracletypes.CurrencyPair, *big.Int](t)

				dataHandler.On("CreateMessages", mock.Anything).Return([]handlers.WebsocketEncodedMessage{testMessage}, nil).Once()

				resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
				}
				resolved2 := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					ethusd: {
						Value: big.NewInt(200),
					},
				}
				unresolved := map[oracletypes.CurrencyPair]error{
					atomusd: wserrors.ErrHandleMessage,
				}

				response1 := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)
				dataHandler.On("HandleMessage", mock.Anything).Return(
					response1,
					nil,
					nil,
				).Once()

				response2 := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved2, nil)
				dataHandler.On("HandleMessage", mock.Anything).Return(
					response2,
					nil,
					nil,
				).Once()

				response3 := providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](nil, unresolved)
				dataHandler.On("HandleMessage", mock.Anything).Return(
					response3,
					nil,
					nil,
				).Maybe()

				return dataHandler
			},
			metrics: func() metrics.WebSocketMetrics {
				m := mockmetrics.NewWebSocketMetrics(t)

				m.On("AddWebSocketConnectionStatus", name, metrics.DialSuccess).Return().Once()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.CreateMessageSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.WriteSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Healthy).Return().Once()

				m.On("AddWebSocketConnectionStatus", name, metrics.ReadSuccess).Return().Maybe()
				m.On("AddWebSocketDataHandlerStatus", name, metrics.HandleMessageSuccess).Return().Maybe()
				m.On("ObserveWebSocketLatency", name, mock.Anything).Return().Maybe()

				m.On("AddWebSocketConnectionStatus", name, metrics.CloseSuccess).Return().Once()
				m.On("AddWebSocketConnectionStatus", name, metrics.Unhealthy).Return().Once()

				return m
			},
			ids: []oracletypes.CurrencyPair{btcusd, ethusd},
			responses: providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int]{
				Resolved: map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
					btcusd: {
						Value: big.NewInt(100),
					},
					ethusd: {
						Value: big.NewInt(200),
					},
				},
				UnResolved: map[oracletypes.CurrencyPair]error{
					atomusd: wserrors.ErrHandleMessage,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, err := handlers.NewWebSocketQueryHandler[oracletypes.CurrencyPair, *big.Int](
				logger,
				cfg,
				tc.dataHandler(),
				tc.connHandler(),
				tc.metrics(),
			)
			require.NoError(t, err)

			responseCh := make(chan providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], 20)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			handler.Start(ctx, tc.ids, responseCh)
			cancel()
			close(responseCh)

			expectedResponses := tc.responses
			seenResponses := make(map[oracletypes.CurrencyPair]bool)
			for resp := range responseCh {
				for id, result := range resp.Resolved {
					if _, ok := seenResponses[id]; ok {
						continue
					}

					require.Equal(t, expectedResponses.Resolved[id], result)
					delete(expectedResponses.Resolved, id)
					seenResponses[id] = true
				}

				for id, err := range resp.UnResolved {
					if _, ok := seenResponses[id]; ok {
						continue
					}

					require.True(t, errors.Is(err, expectedResponses.UnResolved[id]))
					delete(expectedResponses.UnResolved, id)
					seenResponses[id] = true
				}
			}

			// Ensure all responses are account for.
			require.Empty(t, expectedResponses.Resolved)
			require.Empty(t, expectedResponses.UnResolved)
		})
	}
}
