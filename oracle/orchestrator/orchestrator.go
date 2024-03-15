package orchestrator

import (
	"context"
	"fmt"
	"maps"
	"math/big"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers/base"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// ProviderOrchestrator is a stateful orchestrator that is responsible for maintaining
// all of the providers that the oracle is using. This includes initializing the providers,
// creating the provider specific market map, and enabling/disabling the providers based
// on the oracle configuration and market map.
type ProviderOrchestrator struct {
	mut    sync.Mutex
	logger *zap.Logger

	// -------------------Lifecycle Fields-------------------//
	//
	// mainCtx is the main context for the provider orchestrator.
	mainCtx context.Context
	// mainCancel is the main context cancel function.
	mainCancel context.CancelFunc
	// errGroup is the error group for the provider orchestrator.
	errGroup *errgroup.Group

	// -------------------State Fields-------------------//
	//
	// providers is a map of all of the providers that the oracle is using.
	providers map[string]ProviderState
	// marketmapper is the market map provider. Specifically this provider is responsible
	// for making requests for the latest market map data.
	mapper MapperState

	// -------------------Oracle Configuration Fields-------------------//
	//
	// cfg is the oracle configuration.
	cfg config.OracleConfig
	// marketMap is the market map that the oracle is using.
	marketMap mmtypes.MarketMap

	// -------------------Price Provider Constructor Fields-------------------//
	//
	// apiQueryHandler factory is a factory function that creates API query handlers.
	apiQueryHandlerFactory types.PriceAPIQueryHandlerFactory
	// webSocketQueryHandlerFactory is a factory function that creates websocket query
	// handlers.
	webSocketQueryHandlerFactory types.PriceWebSocketQueryHandlerFactory

	// -------------------Market Mapper Provider Constructor Fields-------------------//
	//

	// -------------------Metrics Fields-------------------//
	// wsMetrics is the web socket metrics.
	wsMetrics wsmetrics.WebSocketMetrics
	// apiMetrics is the API metrics.
	apiMetrics apimetrics.APIMetrics
	// providerMetrics is the provider metrics.
	providerMetrics providermetrics.ProviderMetrics
}

// NewProviderOrchestrator returns a new provider orchestrator.
func NewProviderOrchestrator(
	cfg config.OracleConfig,
	opts ...Option,
) (*ProviderOrchestrator, error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, err
	}

	orchestrator := &ProviderOrchestrator{
		cfg:             cfg,
		providers:       make(map[string]ProviderState),
		logger:          zap.NewNop(),
		wsMetrics:       wsmetrics.NewWebSocketMetricsFromConfig(cfg.Metrics),
		apiMetrics:      apimetrics.NewAPIMetricsFromConfig(cfg.Metrics),
		providerMetrics: providermetrics.NewProviderMetricsFromConfig(cfg.Metrics),
	}

	for _, opt := range opts {
		opt(orchestrator)
	}

	return orchestrator, nil
}

// Init initializes the all of the providers that are configured via the oracle config.
// Specifically, this will:
//
// 1. This will initialize the provider.
// 2. Create the provider specific market map, if configured with a marketmap.
// 3. Enable the provider if the provider is included in the oracle config and marketmap.
func (o *ProviderOrchestrator) Init() error {
	o.mut.Lock()
	defer o.mut.Unlock()

	for _, providerCfg := range o.cfg.Providers {
		// Initialize the provider.
		state, err := o.CreateProviderState(providerCfg)
		if err != nil {
			o.logger.Error("failed to create provider state", zap.Error(err))
			return err
		}

		// Add the provider to the orchestrator.
		o.providers[providerCfg.Name] = state
		o.logger.Info(
			"created provider state",
			zap.String("provider", providerCfg.Name),
			zap.Bool("enabled", state.Enabled),
			zap.Int("num_tickers", len(state.Market.GetTickers())),
		)
	}

	return nil
}

// CreateProviderState creates a provider state for the given provider. This constructs the
// query handler, based on the provider's type and configuration. The provider state is then
// enabled/disabled based on whether the provider is configured to support any of the tickers.
func (o *ProviderOrchestrator) CreateProviderState(
	cfg config.ProviderConfig,
) (ProviderState, error) {
	// Create the provider market map. This creates the tickers the provider is configured to
	// support.
	market, err := types.ProviderMarketMapFromMarketMap(cfg.Name, o.marketMap)
	if err != nil {
		return ProviderState{}, fmt.Errorf("failed to create %s's provider market map: %w", cfg.Name, err)
	}

	// Select the query handler based on the provider's configuration.
	var provider *types.PriceProvider
	switch {
	case cfg.API.Enabled:
		queryHandler, err := o.createAPIQueryHandler(cfg, market)
		if err != nil {
			return ProviderState{}, fmt.Errorf("failed to create %s's api query handler: %w", cfg.Name, err)
		}

		provider, err = types.NewPriceProvider(
			base.WithName[mmtypes.Ticker, *big.Int](cfg.Name),
			base.WithLogger[mmtypes.Ticker, *big.Int](o.logger),
			base.WithAPIQueryHandler(queryHandler),
			base.WithAPIConfig[mmtypes.Ticker, *big.Int](cfg.API),
			base.WithIDs[mmtypes.Ticker, *big.Int](market.GetTickers()),
			base.WithMetrics[mmtypes.Ticker, *big.Int](o.providerMetrics),
		)
		if err != nil {
			return ProviderState{}, fmt.Errorf("failed to create %s's provider: %w", cfg.Name, err)
		}
	case cfg.WebSocket.Enabled:
		queryHandler, err := o.createWebSocketQueryHandler(cfg, market)
		if err != nil {
			return ProviderState{}, fmt.Errorf("failed to create %s's web socket query handler: %w", cfg.Name, err)
		}

		provider, err = types.NewPriceProvider(
			base.WithName[mmtypes.Ticker, *big.Int](cfg.Name),
			base.WithLogger[mmtypes.Ticker, *big.Int](o.logger),
			base.WithWebSocketQueryHandler(queryHandler),
			base.WithWebSocketConfig[mmtypes.Ticker, *big.Int](cfg.WebSocket),
			base.WithIDs[mmtypes.Ticker, *big.Int](market.GetTickers()),
			base.WithMetrics[mmtypes.Ticker, *big.Int](o.providerMetrics),
		)
		if err != nil {
			return ProviderState{}, fmt.Errorf("failed to create %s's provider: %w", cfg.Name, err)
		}
	default:
		return ProviderState{}, fmt.Errorf("provider %s has no enabled query handlers", cfg.Name)
	}

	return ProviderState{
		Provider: provider,
		Market:   market,
		Enabled:  len(market.GetTickers()) > 0,
		Cfg:      cfg,
	}, nil
}

// GetProviderState returns all of the providers and their state.
func (o *ProviderOrchestrator) GetProviderState() map[string]ProviderState {
	o.mut.Lock()
	defer o.mut.Unlock()

	// Copy the providers
	providers := make(map[string]ProviderState, len(o.providers))
	maps.Copy(providers, o.providers)

	return providers
}

// GetProviders returns all of the providers.
func (o *ProviderOrchestrator) GetProviders() []*types.PriceProvider {
	o.mut.Lock()
	defer o.mut.Unlock()

	providers := make([]*types.PriceProvider, 0)
	for _, state := range o.providers {
		providers = append(providers, state.Provider)
	}

	return providers
}
