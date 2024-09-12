package oracle

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/config"
	oraclemetrics "github.com/skip-mev/connect/v2/oracle/metrics"
	"github.com/skip-mev/connect/v2/oracle/types"
	apimetrics "github.com/skip-mev/connect/v2/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/connect/v2/providers/base/metrics"
	wsmetrics "github.com/skip-mev/connect/v2/providers/base/websocket/metrics"
	mmclienttypes "github.com/skip-mev/connect/v2/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var _ Oracle = (*OracleImpl)(nil)

// OracleImpl maintains providers and the state provided by them. This includes pricing data and market map updates.
type OracleImpl struct { //nolint:revive
	mut     sync.RWMutex
	logger  *zap.Logger
	running atomic.Bool

	// -------------------Lifecycle Fields-------------------//
	//
	// mainCtx is the main context for the oracle.
	mainCtx context.Context
	// mainCancel is the main context cancel function.
	mainCancel context.CancelFunc
	// wg is the wait group for the oracle.
	wg sync.WaitGroup

	// -------------------Stateful Fields-------------------//
	//
	// priceProviders is a map of all price providers that the oracle is using.
	priceProviders map[string]ProviderState
	// mmProvider is the market map provider. Specifically this provider is responsible
	// for making requests for the latest market map data.
	mmProvider *mmclienttypes.MarketMapProvider
	// aggregator is the price aggregator.
	aggregator PriceAggregator
	// lastPriceSync is the last time the oracle successfully updated its prices.
	lastPriceSync time.Time

	// -------------------Oracle Configuration Fields-------------------//
	//
	// cfg is the oracle configuration.
	cfg config.OracleConfig
	// marketMap is the market map that the oracle is using.
	marketMap mmtypes.MarketMap
	// lastUpdated is the field in the marketmap module tracking the last block at which an update was posted
	lastUpdated uint64
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
	aggregator PriceAggregator,
	opts ...Option,
) (Oracle, error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, err
	}
	if aggregator == nil {
		return nil, errors.New("aggregator is required")
	}

	orc := &OracleImpl{
		cfg:             cfg,
		aggregator:      aggregator,
		priceProviders:  make(map[string]ProviderState), // this will be initialized via the Init method.
		logger:          zap.NewNop(),
		wsMetrics:       wsmetrics.NewWebSocketMetricsFromConfig(cfg.Metrics),
		apiMetrics:      apimetrics.NewAPIMetricsFromConfig(cfg.Metrics),
		providerMetrics: providermetrics.NewProviderMetricsFromConfig(cfg.Metrics),
		metrics:         oraclemetrics.NewNopMetrics(),
	}

	for _, opt := range opts {
		opt(orc)
	}

	return orc, nil
}

// GetProviderState returns all providers and their state. This method is used for testing purposes only.
func (o *OracleImpl) GetProviderState() map[string]ProviderState {
	o.mut.Lock()
	defer o.mut.Unlock()

	return o.priceProviders
}

// GetMarketMap returns the market map.
func (o *OracleImpl) GetMarketMap() mmtypes.MarketMap {
	o.mut.Lock()
	defer o.mut.Unlock()

	return o.marketMap
}

func (o *OracleImpl) GetMarketMapProvider() *mmclienttypes.MarketMapProvider {
	return o.mmProvider
}

func (o *OracleImpl) GetLastSyncTime() time.Time {
	o.mut.RLock()
	defer o.mut.RUnlock()
	return o.lastPriceSync
}

func (o *OracleImpl) GetPrices() types.Prices {
	return o.aggregator.GetPrices()
}
