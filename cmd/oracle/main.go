package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"cosmossdk.io/log"
	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/service/server"
	"github.com/skip-mev/slinky/tests/simapp"
)

var (
	host           = flag.String("host", "localhost", "host for the grpc-service to listen on")
	port           = flag.String("port", "8080", "port for the grpc-service to listen on")
	oracleCfgPath  = flag.String("oracle-config-path", "oracle_config.toml", "path to the oracle config file")
	metricsCfgPath = flag.String("metrics-config-path", "metrics_config.toml", "path to the metrics config file")
)

// start the oracle-grpc server + oracle process, cancel on interrupt or terminate.
func main() {
	// channel with width for either signal
	sigs := make(chan os.Signal, 1)

	// gracefully trigger close on interrupt or terminate signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// parse flags
	flag.Parse()

	logger := log.NewLogger(os.Stderr)

	oracleCfg, err := config.ReadOracleConfigFromFile(*oracleCfgPath)
	if err != nil {
		logger.Error("failed to read oracle config file", "err", err)
		return
	}

	metricsCfg, err := config.ReadMetricsConfigFromFile(*metricsCfgPath)
	if err != nil {
		logger.Error("failed to read metrics config file", "err", err)
		return
	}

	metrics := oraclemetrics.NewNopMetrics()
	if metricsCfg.OracleMetrics.Enabled {
		metrics = oraclemetrics.NewMetrics()
	}

	// Create the oracle.
	oracle, err := oracle.New(
		logger,
		oracleCfg,
		simapp.DefaultProviderFactory(), // Replace with custom provider factory
		aggregator.ComputeMedian(),      // Replace with custom aggregator
		metrics,
	)
	if err != nil {
		logger.Error("failed to create oracle", "err", err)
		return
	}

	// create server
	srv := server.NewOracleServer(oracle, logger)

	// cancel oracle on interrupt or terminate
	go func() {
		<-sigs
		logger.Info(
			"received interrupt or terminate signal, closing oracle",
		)

		cancel()
	}()

	// start prometheus metrics
	logger.Info("starting prometheus metrics")
	if metricsCfg.OracleMetrics.Enabled {
		ps, err := oraclemetrics.NewPrometheusServer(metricsCfg.PrometheusServerAddress, logger)
		if err != nil {
			logger.Error("failed to start prometheus metrics", "err", err)
			return
		}

		go ps.Start()

		// close server on shut-down
		go func() {
			<-ctx.Done()
			logger.Info("stopping prometheus metrics")
			ps.Close()
		}()
	}

	// start oracle + server, and wait for either to finish
	if err := srv.StartServer(ctx, *host, *port); err != nil {
		logger.Error("stopping server", "err", err)
	}
}
