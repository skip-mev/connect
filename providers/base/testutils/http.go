package testutils

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/base"
	"github.com/skip-mev/connect/v2/providers/base/api/handlers"
	handlermocks "github.com/skip-mev/connect/v2/providers/base/api/handlers/mocks"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

// CreateResponseFromJSON creates a http response from a json string.
func CreateResponseFromJSON(m string) *http.Response {
	jsonBlob := bytes.NewReader([]byte(m))
	return &http.Response{Body: io.NopCloser(jsonBlob)}
}

// CreateAPIQueryHandlerWithGetResponses creates a mock query handler that returns the given responses every
// time it is invoked.
func CreateAPIQueryHandlerWithGetResponses[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	logger *zap.Logger,
	responses []providertypes.GetResponse[K, V],
	timeout time.Duration,
) handlers.APIQueryHandler[K, V] {
	t.Helper()

	handler := handlermocks.NewQueryHandler[K, V](t)

	handler.On("Query", mock.Anything, mock.Anything, mock.Anything).Return().Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		responseCh := args.Get(2).(chan<- providertypes.GetResponse[K, V])

		time.Sleep(timeout)

		for _, resp := range responses {
			select {
			case <-ctx.Done():
				return
			case responseCh <- resp:
				logger.Debug("sending response", zap.String("response", resp.String()))
			}
		}
	}).Maybe()

	return handler
}

// CreateAPIQueryHandlerWithResponseFn creates a mock query handler that invokes the given function every time it is
// invoked. The function should utilize the response channel to send responses to the provider.
func CreateAPIQueryHandlerWithResponseFn[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	fn func(context.Context, chan<- providertypes.GetResponse[K, V]),
) handlers.APIQueryHandler[K, V] {
	t.Helper()

	handler := handlermocks.NewQueryHandler[K, V](t)

	handler.On("Query", mock.Anything, mock.Anything, mock.Anything).Return().Run(func(args mock.Arguments) {
		ctx := args.Get(0).(context.Context)
		responseCh := args.Get(2).(chan<- providertypes.GetResponse[K, V])
		fn(ctx, responseCh)
	}).Maybe()

	return handler
}

// CreateAPIProviderWithGetResponses creates a new provider with the given responses.
func CreateAPIProviderWithGetResponses[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	logger *zap.Logger,
	cfg config.ProviderConfig,
	ids []K,
	responses []providertypes.GetResponse[K, V],
	timeout time.Duration,
) *base.Provider[K, V] {
	t.Helper()

	handler := CreateAPIQueryHandlerWithGetResponses[K, V](
		t,
		logger,
		responses,
		timeout,
	)

	provider, err := base.NewProvider[K, V](
		base.WithName[K, V](cfg.Name),
		base.WithAPIQueryHandler[K, V](handler),
		base.WithAPIConfig[K, V](cfg.API),
		base.WithLogger[K, V](logger),
		base.WithIDs[K, V](ids),
	)
	require.NoError(t, err)

	return provider
}

// CreateAPIProviderWithResponseFn creates a new provider with the given response function.
func CreateAPIProviderWithResponseFn[K providertypes.ResponseKey, V providertypes.ResponseValue](
	t *testing.T,
	logger *zap.Logger,
	cfg config.ProviderConfig,
	ids []K,
	fn func(context.Context, chan<- providertypes.GetResponse[K, V]),
) *base.Provider[K, V] {
	t.Helper()

	handler := CreateAPIQueryHandlerWithResponseFn[K, V](
		t,
		fn,
	)

	provider, err := base.NewProvider[K, V](
		base.WithName[K, V](cfg.Name),
		base.WithAPIQueryHandler[K, V](handler),
		base.WithAPIConfig[K, V](cfg.API),
		base.WithLogger[K, V](logger),
		base.WithIDs[K, V](ids),
	)
	require.NoError(t, err)

	return provider
}
