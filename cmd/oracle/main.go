package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof" //nolint: gosec

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/oracle/constants"
	"github.com/skip-mev/slinky/oracle/orchestrator"
	"github.com/skip-mev/slinky/oracle/types"
	oraclemath "github.com/skip-mev/slinky/pkg/math/oracle"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
	oracleserver "github.com/skip-mev/slinky/service/servers/oracle"
	promserver "github.com/skip-mev/slinky/service/servers/prometheus"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

var (
	oracleCfgPath     = flag.String("oracle-config-path", "oracle_config.json", "path to the oracle config file")
	marketCfgPath     = flag.String("market-config-path", "market_config.json", "path to the market config file")
	runPprof          = flag.Bool("run-pprof", false, "run pprof server")
	profilePort       = flag.String("pprof-port", "6060", "port for the pprof server to listen on")
	chain             = flag.String("chain-id", "", "the chain id for which the side car should run for (ex. dydx-mainnet-1)")
	updateLocalConfig = flag.Bool("update-local-market-config", true, "update the market map config when a new one is received; this will overwrite the existing config file.")
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

	marketCfg, err := types.ReadMarketConfigFromFile(*marketCfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read market config file: %s\n", err.Error())
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

	// Define the orchestrator and oracle options. These determine how the orchestrator and oracle are created & executed.
	orchestratorOpts := []orchestrator.Option{
		orchestrator.WithLogger(logger),
		orchestrator.WithMarketMap(marketCfg),
		orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),             // Replace with custom API query handler factory.
		orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory), // Replace with custom websocket query handler factory.
	}
	oracleOpts := []oracle.Option{
		oracle.WithLogger(logger),
		oracle.WithUpdateInterval(cfg.UpdateInterval),
		oracle.WithMetricsConfig(cfg.Metrics),
		oracle.WithMaxCacheAge(cfg.MaxPriceAge),
	}

	if *chain == constants.DYDXMainnet.ID || *chain == constants.DYDXTestnet.ID {
		customOrchestratorOps, customOracleOpts, err := dydxOptions(logger, marketCfg)
		if err != nil {
			logger.Error("failed to create dydx orchestrator and oracle options", zap.Error(err))
			return
		}

		orchestratorOpts = append(orchestratorOpts, customOrchestratorOps...)
		oracleOpts = append(oracleOpts, customOracleOpts...)
	}

	// Create the orchestrator and start the orchestrator.
	orch, err := orchestrator.NewProviderOrchestrator(
		cfg,
		orchestratorOpts...,
	)
	if err != nil {
		logger.Error("failed to create provider orchestrator", zap.Error(err))
		return
	}

	if err := orch.Start(ctx); err != nil {
		logger.Error("failed to start provider orchestrator", zap.Error(err))
		return
	}
	defer orch.Stop()

	// Create the oracle and start the oracle server.
	oracleOpts = append(oracleOpts, oracle.WithProviders(orch.GetPriceProviders()))
	orc, err := oracle.New(oracleOpts...)
	if err != nil {
		logger.Error("failed to create oracle", zap.Error(err))
		return
	}
	srv := oracleserver.NewOracleServer(orc, logger)

	// cancel oracle on interrupt or terminate
	go func() {
		<-sigs
		logger.Info("received interrupt or terminate signal; closing oracle")

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

	if *runPprof {
		endpoint := fmt.Sprintf("%s:%s", cfg.Host, *profilePort)
		// Start pprof server
		go func() {
			logger.Info("Starting pprof server", zap.String("endpoint", endpoint))
			if err := http.ListenAndServe(endpoint, nil); err != nil { //nolint: gosec
				logger.Error("pprof server failed", zap.Error(err))
			}
		}()
	}

	// start oracle + server, and wait for either to finish
	if err := srv.StartServer(ctx, cfg.Host, cfg.Port); err != nil {
		logger.Error("stopping server", zap.Error(err))
	}
}

// dydxOptions specifies the custom orchestrator and oracle options for dYdX.
func dydxOptions(
	logger *zap.Logger,
	marketCfg mmtypes.MarketMap,
) ([]orchestrator.Option, []oracle.Option, error) {
	// dYdX uses the median index price aggregation strategy.
	aggregator, err := oraclemath.NewMedianAggregator(
		logger,
		marketCfg,
	)
	if err != nil {
		return nil, nil, err
	}

	// The oracle must be configured with the median index price aggregator.
	customOracleOpts := []oracle.Option{
		oracle.WithDataAggregator(aggregator),
	}

	// Additionally, dYdX requires a custom market map provider that fetches market params from the chain.
	customOrchestratorOps := []orchestrator.Option{
		orchestrator.WithMarketMapperFactory(oraclefactory.DefaultDYDXMarketMapProvider),
		orchestrator.WithAggregator(aggregator),
	}
	if *updateLocalConfig {
		customOrchestratorOps = append(customOrchestratorOps, orchestrator.WithWriteTo(*marketCfgPath))
	}

	return customOrchestratorOps, customOracleOpts, nil
}
