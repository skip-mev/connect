package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	oraclemath "github.com/skip-mev/slinky/pkg/math/oracle"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
	oracleserver "github.com/skip-mev/slinky/service/servers/oracle"
	promserver "github.com/skip-mev/slinky/service/servers/prometheus"
)

var (
	host          = flag.String("host", "localhost", "host for the grpc-service to listen on")
	port          = flag.String("port", "8080", "port for the grpc-service to listen on")
	oracleCfgPath = flag.String("oracle-config-path", "oracle_config.toml", "path to the oracle config file")
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

	cfg, err := config.ReadOracleConfigFromFile(*oracleCfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read oracle config file: %s\n", err.Error())
		return
	}

	var logger *zap.Logger
	if !cfg.Production {
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

	// Create the API and websocket query handler factories. These are used to create the
	// data collection handlers for the providers. To read more about what these query handlers
	// do, see the documentation for the `providers/base` package.
	apiFactory := oraclefactory.APIQueryHandlerFactory()      // Replace with custom API factory.
	wsFactory := oraclefactory.WebSocketQueryHandlerFactory() // Replace with custom websocket factory.

	// Create the providers using the default provider factory.
	generator, err := oraclefactory.NewDefaultProviderFactory(
		logger,
		apiFactory,
		wsFactory,
	)
	if err != nil {
		logger.Error("failed to create provider factory", zap.Error(err))
		return
	}

	providers, err := generator.Factory()(cfg)
	if err != nil {
		logger.Error("failed to create providers", zap.Error(err))
		return
	}

	// Create the conversion market aggregator.
	aggregator, err := oraclemath.NewMedianAggregator(logger, cfg.Market)
	if err != nil {
		logger.Error("failed to create median aggregator", zap.Error(err))
		return
	}

	// Create the oracle.
	oracle, err := oracle.New(
		oracle.WithUpdateInterval(cfg.UpdateInterval),
		oracle.WithProviders(providers),                        // Replace with custom providers.
		oracle.WithAggregateFunction(aggregator.AggregateFn()), // Replace with custom aggregation function.
		oracle.WithMetricsConfig(cfg.Metrics),
		oracle.WithLogger(logger),
	)
	if err != nil {
		logger.Error("failed to create oracle", zap.Error(err))
		return
	}

	// create server
	srv := oracleserver.NewOracleServer(oracle, logger)

	// cancel oracle on interrupt or terminate
	go func() {
		<-sigs
		logger.Info(
			"received interrupt or terminate signal, closing oracle",
		)

		cancel()
	}()

	// start prometheus metrics
	if cfg.Metrics.Enabled {
		logger.Info("starting prometheus metrics", zap.String("address", cfg.Metrics.PrometheusServerAddress))
		ps, err := promserver.NewPrometheusServer(cfg.Metrics.PrometheusServerAddress, logger)
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
