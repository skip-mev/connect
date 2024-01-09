package testutils

import (
	"testing"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	handlers "github.com/skip-mev/slinky/providers/base/handlers"
	handlermocks "github.com/skip-mev/slinky/providers/base/handlers/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// CreateQueryHandlerWithGetResponses creates a mock query handler that returns the given responses every
// time it is invoked.
func CreateQueryHandlerWithGetResponses[K comparable, V any](
	t *testing.T,
	logger *zap.Logger,
	responses []providertypes.GetResponse[K, V],
) handlers.QueryHandler[K, V] {
	handler := handlermocks.NewQueryHandler[K, V](t)

	handler.On("Query", mock.Anything, mock.Anything, mock.Anything).Return().Run(func(args mock.Arguments) {
		responseCh := args.Get(2).(chan<- providertypes.GetResponse[K, V])

		for _, resp := range responses {
			logger.Debug("sending response", zap.String("response", resp.String()))
			responseCh <- resp
		}
	}).Maybe()

	return handler
}

// CreateQueryHandlerWithResponseFn creates a mock query handler that invokes the given function every time it is
// invoked. The function should utilize the response channel to send responses to the provider.
func CreateQueryHandlerWithResponseFn[K comparable, V any](
	t *testing.T,
	fn func(chan<- providertypes.GetResponse[K, V]),
) handlers.QueryHandler[K, V] {
	handler := handlermocks.NewQueryHandler[K, V](t)

	handler.On("Query", mock.Anything, mock.Anything, mock.Anything).Return().Run(func(args mock.Arguments) {
		responseCh := args.Get(2).(chan<- providertypes.GetResponse[K, V])
		fn(responseCh)
	}).Maybe()

	return handler
}

// CreateProviderWithGetResponses creates a new provider with the given responses.
func CreateProviderWithGetResponses[K comparable, V any](
	t *testing.T,
	logger *zap.Logger,
	cfg config.ProviderConfig,
	ids []K,
	responses []providertypes.GetResponse[K, V],
) providertypes.Provider[K, V] {
	handler := CreateQueryHandlerWithGetResponses[K, V](
		t,
		logger,
		responses,
	)

	provider, err := base.NewProvider[K, V](
		logger,
		cfg,
		handler,
		ids,
	)
	require.NoError(t, err)

	return provider
}

// CreateProviderWithResponseFn creates a new provider with the given response function.
func CreateProviderWithResponseFn[K comparable, V any](
	t *testing.T,
	logger *zap.Logger,
	cfg config.ProviderConfig,
	ids []K,
	fn func(chan<- providertypes.GetResponse[K, V]),
) providertypes.Provider[K, V] {
	handler := CreateQueryHandlerWithResponseFn[K, V](
		t,
		fn,
	)

	provider, err := base.NewProvider[K, V](
		logger,
		cfg,
		handler,
		ids,
	)
	require.NoError(t, err)

	return provider
}
