package testutils

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	"github.com/skip-mev/slinky/providers/base/api/handlers"
	handlermocks "github.com/skip-mev/slinky/providers/base/api/handlers/mocks"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// CreateResponseFromJSON creates a http response from a json string.
func CreateResponseFromJSON(m string) *http.Response {
	jsonBlob := bytes.NewReader([]byte(m))
	return &http.Response{Body: io.NopCloser(jsonBlob)}
}

// CreateAPIQueryHandlerWithGetResponses creates a mock query handler that returns the given responses every
// time it is invoked.
func CreateAPIQueryHandlerWithGetResponses[K comparable, V any](
	t *testing.T,
	logger *zap.Logger,
	responses []providertypes.GetResponse[K, V],
) handlers.APIQueryHandler[K, V] {
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

// CreateAPIQueryHandlerWithResponseFn creates a mock query handler that invokes the given function every time it is
// invoked. The function should utilize the response channel to send responses to the provider.
func CreateAPIQueryHandlerWithResponseFn[K comparable, V any](
	t *testing.T,
	fn func(chan<- providertypes.GetResponse[K, V]),
) handlers.APIQueryHandler[K, V] {
	handler := handlermocks.NewQueryHandler[K, V](t)

	handler.On("Query", mock.Anything, mock.Anything, mock.Anything).Return().Run(func(args mock.Arguments) {
		responseCh := args.Get(2).(chan<- providertypes.GetResponse[K, V])
		fn(responseCh)
	}).Maybe()

	return handler
}

// CreateProviderWithGetResponses creates a new provider with the given responses.
func CreateAPIProviderWithGetResponses[K comparable, V any](
	t *testing.T,
	logger *zap.Logger,
	cfg config.ProviderConfig,
	ids []K,
	responses []providertypes.GetResponse[K, V],
) providertypes.Provider[K, V] {
	handler := CreateAPIQueryHandlerWithGetResponses[K, V](
		t,
		logger,
		responses,
	)

	provider, err := base.NewProvider[K, V](
		cfg,
		base.WithAPIQueryHandler[K, V](handler),
		base.WithLogger[K, V](logger),
		base.WithIDs[K, V](ids),
	)
	require.NoError(t, err)

	return provider
}

// CreateAPIProviderWithResponseFn creates a new provider with the given response function.
func CreateAPIProviderWithResponseFn[K comparable, V any](
	t *testing.T,
	logger *zap.Logger,
	cfg config.ProviderConfig,
	ids []K,
	fn func(chan<- providertypes.GetResponse[K, V]),
) providertypes.Provider[K, V] {
	handler := CreateAPIQueryHandlerWithResponseFn[K, V](
		t,
		fn,
	)

	provider, err := base.NewProvider[K, V](
		cfg,
		base.WithAPIQueryHandler[K, V](handler),
		base.WithLogger[K, V](logger),
		base.WithIDs[K, V](ids),
	)
	require.NoError(t, err)

	return provider
}
