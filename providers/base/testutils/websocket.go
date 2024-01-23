package testutils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	handlermocks "github.com/skip-mev/slinky/providers/base/websocket/handlers/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// CreateWebSocketQueryHandlerWithGetResponses creates a mock query handler that returns the given responses every
// time it is invoked.
func CreateWebSocketQueryHandlerWithGetResponses[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	timeout time.Duration,
	logger *zap.Logger,
	responses []providertypes.GetResponse[K, V],
) handlers.WebSocketQueryHandler[K, V] {
	handler := handlermocks.NewWebSocketQueryHandler[K, V](t)

	handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		responseCh := args.Get(2).(chan<- providertypes.GetResponse[K, V])

		for _, resp := range responses {
			logger.Debug("sending response", zap.String("response", resp.String()))
			responseCh <- resp
			time.Sleep(timeout)
		}
	}).Maybe()

	return handler
}

// CreateWebSocketQueryHandlerWithResponseFn creates a mock query handler that invokes the given function every time it is
// invoked. The function should utilize the response channel to send responses to the provider.
func CreateWebSocketQueryHandlerWithResponseFn[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	fn func(chan<- providertypes.GetResponse[K, V]),
) handlers.WebSocketQueryHandler[K, V] {
	handler := handlermocks.NewWebSocketQueryHandler[K, V](t)

	handler.On("Start", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		responseCh := args.Get(2).(chan<- providertypes.GetResponse[K, V])
		fn(responseCh)
	}).Maybe()

	return handler
}

// CreateWebSocketProviderWithGetResponses creates a new web socket provider with the given responses.
func CreateWebSocketProviderWithGetResponses[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	timeout time.Duration,
	cfg config.ProviderConfig,
	logger *zap.Logger,
	responses []providertypes.GetResponse[K, V],
) providertypes.Provider[K, V] {
	handler := CreateWebSocketQueryHandlerWithGetResponses[K, V](
		t,
		timeout,
		logger,
		responses,
	)

	p, err := base.NewProvider[K, V](
		cfg,
		base.WithWebSocketQueryHandler[K, V](handler),
		base.WithLogger[K, V](logger),
	)
	require.NoError(t, err)

	return p
}

// CreateWebSocketProviderWithResponseFn creates a new web socket provider with the given response function.
func CreateWebSocketProviderWithResponseFn[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	cfg config.ProviderConfig,
	logger *zap.Logger,
	fn func(chan<- providertypes.GetResponse[K, V]),
) providertypes.Provider[K, V] {
	handler := CreateWebSocketQueryHandlerWithResponseFn[K, V](
		t,
		fn,
	)

	p, err := base.NewProvider[K, V](
		cfg,
		base.WithWebSocketQueryHandler[K, V](handler),
		base.WithLogger[K, V](logger),
	)
	require.NoError(t, err)

	return p
}
