package oracle

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/oracle/types"
	apimetrics "github.com/skip-mev/slinky/providers/base/api/metrics"
	providermetrics "github.com/skip-mev/slinky/providers/base/metrics"
	wsmetrics "github.com/skip-mev/slinky/providers/base/websocket/metrics"
	mmclienttypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
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
	priceProviders map[string]*PriceProviderState
	// mmProviders is a map of all market map providers that the oracle is using.
	mmProviders map[string]*MarketMapProviderState
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

// PriceProviderState is the state of a provider. This includes the provider implementation,
// and the provider specific market map.
type PriceProviderState struct {
	// Provider is the price provider implementation.
	Provider *types.PriceProvider
	// Tickers are the provider's tickers.
	Tickers types.ProviderTickers
}

// MarketMapProviderState is the state of the market map provider.
type MarketMapProviderState struct {
	// Provider is the market map provider implementation.
	Provider *mmclienttypes.MarketMapProvider
	// MarketMap is the provider's market map.
	MarketMap mmtypes.MarketMap
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
		priceProviders:  make(map[string]*PriceProviderState), // this will be initialized via the Init method.
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

// GetPriceProvidersState returns all providers and their state. This method is used for testing purposes only.
func (o *OracleImpl) GetPriceProvidersState() map[string]*PriceProviderState {
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

// GetMarketMapProvider returns the market map provider.
func (o *OracleImpl) GetMarketMapProvidersState() map[string]*MarketMapProviderState {
	o.mut.Lock()
	defer o.mut.Unlock()

	return o.mmProviders
}

// GetLastSyncTime returns the last time the oracle successfully updated its prices.
func (o *OracleImpl) GetLastSyncTime() time.Time {
	o.mut.RLock()
	defer o.mut.RUnlock()

	return o.lastPriceSync
}

// GetPrices returns the current prices.
func (o *OracleImpl) GetPrices() types.Prices {
	return o.aggregator.GetPrices()
}
