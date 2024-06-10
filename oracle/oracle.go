package oracle

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math/oracle"
	ssync "github.com/skip-mev/slinky/pkg/sync"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
	mmclienttypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var _ Oracle = (*OracleImpl)(nil)

// Oracle defines the expected interface for an oracle. It is consumed by the oracle server.
//
//go:generate mockery --name Oracle --filename mock_oracle.go
type Oracle interface {
	IsRunning() bool
	GetLastSyncTime() time.Time
	GetPrices() types.Prices
	GetMarketMap() mmtypes.MarketMap
	Start(ctx context.Context) error
	Stop()
}

// OracleImpl is a stateful orchestrator that is responsible for maintaining
// all providers that the oracle is using. This includes initializing the providers,
// creating the provider specific market map, and enabling/disabling the providers based
// on the oracle configuration and market map.
//
// TODO(tyler): rewrite this comment.
type OracleImpl struct {
	mut     sync.Mutex
	logger  *zap.Logger
	closer  *ssync.Closer
	running atomic.Bool

	// -------------------Lifecycle Fields-------------------//
	//
	// mainCtx is the main context for the provider orchestrator.
	mainCtx context.Context
	// mainCancel is the main context cancel function.
	mainCancel context.CancelFunc
	// wg is the wait group for the provider orchestrator.
	wg sync.WaitGroup

	// -------------------Stateful Fields-------------------//
	//
	// providers is a map of all providers that the oracle is using.
	providers map[string]ProviderState
	// mmProvider is the market map provider. Specifically this provider is responsible
	// for making requests for the latest market map data.
	mmProvider *mmclienttypes.MarketMapProvider
	// aggregator is the price aggregator.
	aggregator *oracle.IndexPriceAggregator
	// lastPriceSync is the last time the oracle successfully updated its prices.
	lastPriceSync time.Time

	// -------------------Oracle Configuration Fields-------------------//
	//
	// cfg is the oracle configuration.
	cfg config.OracleConfig
	// marketMap is the market map that the oracle is using.
	marketMap mmtypes.MarketMap
	// writeTo is a path to write the market map to.
	writeTo string

	// -------------------Provider Constructor Fields-------------------//
	//
	// priceAPIFactory factory is a factory function that creates price API query handlers.
	priceAPIFactory types.PriceAPIQueryHandlerFactory
	// priceWSFactory is a factory function that creates price websocket query handlers.
	priceWSFactory types.PriceWebSocketQueryHandlerFactory
	// marketMapperFactory is a factory function that creates market map providers.
	marketMapperFactory mmclienttypes.MarketMapFactory

	// -------------------Metrics Fields-------------------//
	//
	// wsMetrics is the web socket metrics.
	wsMetrics wsmetrics.WebSocketMetrics
	// apiMetrics is the API metrics.
	apiMetrics apimetrics.APIMetrics
	// providerMetrics is the provider metrics.
	providerMetrics providermetrics.ProviderMetrics
	// metrics are the base metrics of the oracle.
	metrics oraclemetrics.Metrics
}

// ProviderState is the state of a provider. This includes the provider implementation,
// the provider specific market map, and whether the provider is enabled.
type ProviderState struct {
	// Provider is the price provider implementation.
	Provider *types.PriceProvider
	// Cfg is the provider configuration.
	//
	// TODO: Deprecate this once we have synchronous configuration updates.
	Cfg config.ProviderConfig
}

// New returns a new Oracle.
func New(
	cfg config.OracleConfig,
	opts ...Option,
) (Oracle, error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, err
	}

	orc := &OracleImpl{
		cfg:             cfg,
		closer:          ssync.NewCloser(),
		providers:       make(map[string]ProviderState),
		logger:          zap.NewNop(),
		wsMetrics:       wsmetrics.NewWebSocketMetricsFromConfig(cfg.Metrics),
		apiMetrics:      apimetrics.NewAPIMetricsFromConfig(cfg.Metrics),
		providerMetrics: providermetrics.NewProviderMetricsFromConfig(cfg.Metrics),
	}

	for _, opt := range opts {
		opt(orc)
	}

	return orc, nil
}

// GetProviderState returns all providers and their state.
func (o *OracleImpl) GetProviderState() map[string]ProviderState {
	o.mut.Lock()
	defer o.mut.Unlock()

	return o.providers
}

// GetPriceProviders returns all price providers.
func (o *OracleImpl) GetPriceProviders() []*types.PriceProvider {
	o.mut.Lock()
	defer o.mut.Unlock()

	providers := make([]*types.PriceProvider, 0, len(o.providers))
	for _, state := range o.providers {
		providers = append(providers, state.Provider)
	}

	return providers
}

// GetMarketMapProvider returns the market map provider.
func (o *OracleImpl) GetMarketMapProvider() *mmclienttypes.MarketMapProvider {
	o.mut.Lock()
	defer o.mut.Unlock()

	return o.mmProvider
}

// GetMarketMap returns the market map.
func (o *OracleImpl) GetMarketMap() mmtypes.MarketMap {
	o.mut.Lock()
	defer o.mut.Unlock()

	return o.marketMap
}

func (o *OracleImpl) GetLastSyncTime() time.Time {
	return o.lastPriceSync
}

func (o *OracleImpl) GetPrices() types.Prices {
	return o.aggregator.GetPrices()
}
