package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/skip-mev/connect/v2/oracle/config"
)

type ImplType int

const (
	OracleMetricsType ImplType = iota
	NoOpMetricsType
	UnknownMetricsType
)

// determines the underlying impl of the metrics.
func determineMetricsType(m Metrics) ImplType {
	switch m.(type) {
	case *OracleMetricsImpl:
		return OracleMetricsType
	case *noOpOracleMetrics:
		return NoOpMetricsType
	default:
		return UnknownMetricsType
	}
}

var _ Metrics = &dynamicMetrics{}

// dynamicMetrics is a type that can change its internal metrics impl on the fly.
// this is useful for when a connect instance is started before a node. we can't be sure which one starts first,
// so we need to be able to switch when the node comes online.
type dynamicMetrics struct {
	cfg         config.MetricsConfig
	nc          NodeClient
	impl        Metrics
	metricsType ImplType
	mu          sync.RWMutex // Add a mutex for concurrent access
}

func NewDynamicMetrics(ctx context.Context, cfg config.MetricsConfig, nc NodeClient) Metrics {
	impl := NewMetricsFromConfig(cfg, nc)
	dyn := &dynamicMetrics{
		cfg:         cfg,
		nc:          nc,
		impl:        impl,
		metricsType: determineMetricsType(impl),
	}
	// we only want to kick off the routine of attempting to switch if we're a noop metrics, telemetry is enabled,
	// and we have a node client.
	if dyn.metricsType == NoOpMetricsType && !cfg.Telemetry.Disabled && nc != nil {
		dyn.retrySwitchImpl(ctx)
	}
	return dyn
}

// retrySwitchImpl kicks off a go routine that attempts to contact a node every 3 seconds for 10 mins.
// if it gets a response, it will switch its internal metrics impl.
func (d *dynamicMetrics) retrySwitchImpl(ctx context.Context) {
	go func() {
		retryDuration := time.NewTimer(10 * time.Minute)
		ticker := time.NewTicker(3 * time.Second)

		for {
			select {
			case <-ctx.Done():
				return
			case <-retryDuration.C:
				return
			case <-ticker.C:
				_, err := d.nc.DeriveNodeIdentifier()
				if err == nil {
					impl := NewMetricsFromConfig(d.cfg, d.nc)
					d.mu.Lock()
					d.impl = impl
					d.metricsType = determineMetricsType(d.impl)
					d.mu.Unlock()
					d.SetConnectBuildInfo()
					return
				}
			}
		}
	}()
}

func (d *dynamicMetrics) AddTick() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	d.impl.AddTick()
}

func (d *dynamicMetrics) AddTickerTick(pairID string) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	d.impl.AddTickerTick(pairID)
}

func (d *dynamicMetrics) UpdatePrice(name, pairID string, decimals uint64, price float64) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	d.impl.UpdatePrice(name, pairID, decimals, price)
}

func (d *dynamicMetrics) UpdateAggregatePrice(pairID string, decimals uint64, price float64) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	d.impl.UpdateAggregatePrice(pairID, decimals, price)
}

func (d *dynamicMetrics) AddProviderTick(providerName, pairID string, success bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	d.impl.AddProviderTick(providerName, pairID, success)
}

func (d *dynamicMetrics) AddProviderCountForMarket(pairID string, count int) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	d.impl.AddProviderCountForMarket(pairID, count)
}

func (d *dynamicMetrics) SetConnectBuildInfo() {
	d.mu.RLock()
	defer d.mu.RUnlock()
	d.impl.SetConnectBuildInfo()
}

func (d *dynamicMetrics) MissingPrices(pairIDs []string) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	d.impl.MissingPrices(pairIDs)
}

func (d *dynamicMetrics) GetMissingPrices() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.impl.GetMissingPrices()
}
