package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

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

	oracleCfg, err := config.ReadOracleConfigFromFile(*oracleCfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read oracle config file: %s\n", err.Error())
		return
	}

	if !oracleCfg.Enabled {
		fmt.Fprintf(os.Stderr, "oracle is not enabled\n")
		return
	}

	var logger *zap.Logger
	if !oracleCfg.Production {
		logger, err = zap.NewDevelopment()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create logger: %s\n", err.Error())
			return
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create logger: %s\n", err.Error())
			return
		}
	}

	metricsCfg, err := config.ReadMetricsConfigFromFile(*metricsCfgPath)
	if err != nil {
		logger.Error("failed to read metrics config file", zap.Error(err))
		return
	}

	// This can be replaced with a custom provider factory. See the simapp package for an example.
	// Construct all of the web socket providers.
	wsProviders, err := simapp.DefaultWebSocketProviderFactory()(logger, oracleCfg, metricsCfg.OracleMetrics)
	if err != nil {
		logger.Error("failed to create providers using the factory", zap.Error(err))
		return
	}

	// This can be replaced with a custom provider factory. See the simapp package for an example.
	// Construct all of the API providers.
	apiProviders, err := simapp.DefaultAPIProviderFactory()(logger, oracleCfg, metricsCfg.OracleMetrics)
	if err != nil {
		logger.Error("failed to create providers using the factory", zap.Error(err))
		return
	}

	// Create the oracle.
	oracle, err := oracle.New(
		oracleCfg,
		oracle.WithProviders(append(wsProviders, apiProviders...)), // Replace with custom providers.
		oracle.WithAggregateFunction(aggregator.ComputeMedian()),   // Replace with custom aggregation function.
		oracle.WithMetricsConfig(metricsCfg.OracleMetrics),
		oracle.WithLogger(logger),
	)
	if err != nil {
		logger.Error("failed to create oracle", zap.Error(err))
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
	if metricsCfg.OracleMetrics.Enabled {
		logger.Info("starting prometheus metrics", zap.String("address", metricsCfg.PrometheusServerAddress))
		ps, err := oraclemetrics.NewPrometheusServer(metricsCfg.PrometheusServerAddress, logger)
		if err != nil {
			logger.Error("failed to start prometheus metrics", zap.Error(err))
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
		logger.Error("stopping server", zap.Error(err))
	}
}
