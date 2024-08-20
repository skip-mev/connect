package testutils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/base"
	"github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
	handlermocks "github.com/skip-mev/connect/v2/providers/base/websocket/handlers/mocks"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

// CreateWebSocketQueryHandlerWithGetResponses creates a mock query handler that returns the given responses every
// time it is invoked.
func CreateWebSocketQueryHandlerWithGetResponses[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	timeout time.Duration,
	logger *zap.Logger,
	responses []providertypes.GetResponse[K, V],
) handlers.WebSocketQueryHandler[K, V] {
	t.Helper()

	handler := handlermocks.NewWebSocketQueryHandler[K, V](t)

	handler.On("Copy").Return(handler).Maybe()
	handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		responseCh := args.Get(2).(chan<- providertypes.GetResponse[K, V])

		for _, resp := range responses {
			logger.Debug("sending response", zap.String("response", resp.String()))
			select {
			case <-ctx.Done():
				return
			case responseCh <- resp:
				time.Sleep(timeout)
			}
		}
	}).Maybe()

	return handler
}

// CreateWebSocketQueryHandlerWithResponseFn creates a mock query handler that invokes the given function every time it is
// invoked. The function should utilize the response channel to send responses to the provider.
func CreateWebSocketQueryHandlerWithResponseFn[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	fn func(context.Context, chan<- providertypes.GetResponse[K, V]),
) handlers.WebSocketQueryHandler[K, V] {
	t.Helper()

	handler := handlermocks.NewWebSocketQueryHandler[K, V](t)

	handler.On("Copy").Return(handler).Maybe()
	handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		responseCh := args.Get(2).(chan<- providertypes.GetResponse[K, V])
		fn(ctx, responseCh)
	}).Maybe()

	return handler
}

// CreateWebSocketProviderWithGetResponses creates a new websocket provider with the given responses.
func CreateWebSocketProviderWithGetResponses[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	timeout time.Duration,
	ids []K,
	cfg config.ProviderConfig,
	logger *zap.Logger,
	responses []providertypes.GetResponse[K, V],
) *base.Provider[K, V] {
	t.Helper()

	handler := CreateWebSocketQueryHandlerWithGetResponses[K, V](
		t,
		timeout,
		logger,
		responses,
	)

	p, err := base.NewProvider[K, V](
		base.WithName[K, V](cfg.Name),
		base.WithWebSocketQueryHandler[K, V](handler),
		base.WithWebSocketConfig[K, V](cfg.WebSocket),
		base.WithLogger[K, V](logger),
		base.WithIDs[K, V](ids),
	)
	require.NoError(t, err)

	return p
}

// CreateWebSocketProviderWithResponseFn creates a new websocket provider with the given response function.
func CreateWebSocketProviderWithResponseFn[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	cfg config.ProviderConfig,
	logger *zap.Logger,
	fn func(context.Context, chan<- providertypes.GetResponse[K, V]),
) *base.Provider[K, V] {
	t.Helper()

	handler := CreateWebSocketQueryHandlerWithResponseFn[K, V](
		t,
		fn,
	)

	p, err := base.NewProvider[K, V](
		base.WithName[K, V](cfg.Name),
		base.WithWebSocketQueryHandler[K, V](handler),
		base.WithWebSocketConfig[K, V](cfg.WebSocket),
		base.WithLogger[K, V](logger),
	)
	require.NoError(t, err)

	return p
}
