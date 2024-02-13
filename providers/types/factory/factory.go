package factory

import (
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	apihandlers "github.com/skip-mev/slinky/providers/base/api/handlers"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	wshandlers "github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
)

type (
	// ProviderFactory inputs the oracle configuration and returns a set of providers. Developers
	// can implement their own provider factory to create their own providers.
	ProviderFactory[K providertypes.ResponseKey, V providertypes.ResponseValue] func(
		config.OracleConfig,
	) ([]providertypes.Provider[K, V], error)

	// APIQueryHandlerFactory inputs the oracle configuration and returns a API Query Handler.
	APIQueryHandlerFactory[K providertypes.ResponseKey, V providertypes.ResponseValue] func(
		*zap.Logger,
		config.ProviderConfig,
		apimetrics.APIMetrics,
	) (apihandlers.APIQueryHandler[K, V], error)

	// WebSocketQueryHandlerFactory inputs the oracle configuration and returns a WebSocket Query Handler.
	WebSocketQueryHandlerFactory[K providertypes.ResponseKey, V providertypes.ResponseValue] func(
		*zap.Logger,
		config.ProviderConfig,
		wsmetrics.WebSocketMetrics,
	) (wshandlers.WebSocketQueryHandler[K, V], error)
)
