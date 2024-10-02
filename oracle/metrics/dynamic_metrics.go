package metrics

import (
	"context"
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
}

func NewDynamicMetrics(ctx context.Context, cfg config.MetricsConfig, nc NodeClient) Metrics {
	impl := NewMetricsFromConfig(cfg, nc)
	dyn := &dynamicMetrics{
		cfg:         cfg,
		nc:          nc,
		impl:        impl,
		metricsType: determineMetricsType(impl),
	}
	if dyn.metricsType == NoOpMetricsType && !cfg.Telemetry.Disabled && nc != nil {
		go dyn.retrySwitchImpl(ctx)
	}
	return dyn
}

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
				// if we can successfully query the identifier, we can switch the metrics.
				// this means the node started, and we can actually start emitting real metrics.
				_, err := d.nc.DeriveNodeIdentifier()
				if err == nil {
					impl := NewMetricsFromConfig(d.cfg, d.nc)
					d.impl = impl
					d.metricsType = determineMetricsType(d.impl)
					d.impl.SetConnectBuildInfo()
					return
				}
			}
		}
	}()
}

func (d *dynamicMetrics) AddTick() {
	d.impl.AddTick()
}

func (d *dynamicMetrics) AddTickerTick(pairID string) {
	d.impl.AddTickerTick(pairID)
}

func (d *dynamicMetrics) UpdatePrice(name, pairID string, decimals uint64, price float64) {
	d.impl.UpdatePrice(name, pairID, decimals, price)
}

func (d *dynamicMetrics) UpdateAggregatePrice(pairID string, decimals uint64, price float64) {
	d.impl.UpdateAggregatePrice(pairID, decimals, price)
}

func (d *dynamicMetrics) AddProviderTick(providerName, pairID string, success bool) {
	d.impl.AddProviderTick(providerName, pairID, success)
}

func (d *dynamicMetrics) AddProviderCountForMarket(pairID string, count int) {
	d.impl.AddProviderCountForMarket(pairID, count)
}

func (d *dynamicMetrics) SetConnectBuildInfo() {
	d.impl.SetConnectBuildInfo()
}

func (d *dynamicMetrics) MissingPrices(pairIDs []string) {
	d.impl.MissingPrices(pairIDs)
}

func (d *dynamicMetrics) GetMissingPrices() []string {
	return d.impl.GetMissingPrices()
}
