package factory

import (
	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	"github.com/skip-mev/connect/v2/providers/base"
	apihandlers "github.com/skip-mev/connect/v2/providers/base/api/handlers"
	apimetrics "github.com/skip-mev/connect/v2/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/connect/v2/providers/base/metrics"
	wshandlers "github.com/skip-mev/connect/v2/providers/base/websocket/handlers"
	wsmetrics "github.com/skip-mev/connect/v2/providers/base/websocket/metrics"
	providertypes "github.com/skip-mev/connect/v2/providers/types"
)

type (
	// ProviderFactory inputs the oracle configuration and returns a set of providers. Developers
	// can implement their own provider factory to create their own providers.
	ProviderFactory[K providertypes.ResponseKey, V providertypes.ResponseValue] func(
		config.OracleConfig,
	) ([]providertypes.Provider[K, V], error)

	// BaseProviderFactory inputs the provider configuration and returns a set of base providers. The factory
	// should case on all the different provider configurations and return the appropriate base providers.
	BaseProviderFactory[K providertypes.ResponseKey, V providertypes.ResponseValue] func(
		logger *zap.Logger,
		cfg config.OracleConfig,
		wsMetrics wsmetrics.WebSocketMetrics,
		apiMetrics apimetrics.APIMetrics,
		providerMetrics providermetrics.ProviderMetrics,
	) ([]*base.Provider[K, V], error)

	// APIQueryHandlerFactory inputs the provider configuration and returns a API Query Handler. The
	// factory should case on all the different provider configurations and return the appropriate
	// API Query Handler.
	APIQueryHandlerFactory[K providertypes.ResponseKey, V providertypes.ResponseValue] func(
		*zap.Logger,
		config.ProviderConfig,
		apimetrics.APIMetrics,
	) (apihandlers.APIQueryHandler[K, V], error)

	// WebSocketQueryHandlerFactory inputs the provider configuration and returns a WebSocket Query Handler.
	// The factory should case on all the different provider configurations and return the appropriate
	// WebSocket Query Handler.
	WebSocketQueryHandlerFactory[K providertypes.ResponseKey, V providertypes.ResponseValue] func(
		*zap.Logger,
		config.ProviderConfig,
		wsmetrics.WebSocketMetrics,
	) (wshandlers.WebSocketQueryHandler[K, V], error)
)
