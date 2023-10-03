package config

import (
	"fmt"
	"net"
	"net/http"

	"cosmossdk.io/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/skip-mev/slinky/service/metrics"
)

// StartPrometheusServer starts a prometheus server if the metrics are enabled and
// address is set, and valid. This method will spawn an http server in a go-routine,
// that will handle requests to /metrics and serve the metrics registered in the DefaultRegisterer.
func StartPrometheusServer(cfg Metrics, log log.Logger) error {
	log.Info("starting prometheus server", "cfg", cfg)
	if !cfg.AppMetrics.Enabled && !cfg.OracleMetrics.Enabled {
		return nil
	}

	// get the prometheus server address
	if cfg.PrometheusServerAddress == "" || !isValidAddress(cfg.PrometheusServerAddress) {
		return fmt.Errorf("invalid prometheus server address: %s", cfg.PrometheusServerAddress)
	}

	// create server for DefaultRegisterer
	http.Handle("/metrics", promhttp.Handler())

	// serve
	go func() {
		// TODO: Do we need more security here / configuration options here?
		if err := http.ListenAndServe(cfg.PrometheusServerAddress, nil); err != nil { //nolint: gosec
			log.Info("failed to start prometheus server", "err", err)
			panic(err)
		}
	}()
	return nil
}

func isValidAddress(address string) bool {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return false
	}

	if host == "" || port == "" {
		return false
	}

	return true
}

// NewServiceMetricsFromConfig returns a new Metrics implementation based on the config. The Metrics
// returned is safe to be used in the client, and in the Oracle used by the PreFinalizeBlockHandler.
// If the metrics are not enabled, a nop implementation is returned.
func NewServiceMetricsFromConfig(cfg Metrics) (metrics.Metrics, sdk.ConsAddress, error) {
	if !cfg.AppMetrics.Enabled {
		return metrics.NewNopMetrics(), nil, nil
	}

	// ensure that the metrics are enabled
	if err := cfg.AppMetrics.ValidateBasic(); err != nil {
		return nil, nil, err
	}

	// get the cons address
	consAddress, err := cfg.AppMetrics.ConsAddress()
	if err != nil {
		return nil, nil, err
	}

	// create the metrics
	metrics := metrics.NewMetrics()
	return metrics, consAddress, nil
}
