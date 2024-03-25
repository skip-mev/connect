package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof" //nolint: gosec

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/oracle/orchestrator"
	"github.com/skip-mev/slinky/oracle/types"
	oraclemath "github.com/skip-mev/slinky/pkg/math/oracle"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
	oracleserver "github.com/skip-mev/slinky/service/servers/oracle"
	promserver "github.com/skip-mev/slinky/service/servers/prometheus"
)

var (
	rootCmd = &cobra.Command{
		Use:   "oracle",
		Short: "Run the slinky oracle server.",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			return runOracle()
		},
	}

	oracleCfgPath     string
	marketCfgPath     string
	runPprof          bool
	profilePort       string
	chain             string
	updateLocalConfig bool
)

func init() {
	rootCmd.Flags().StringVarP(
		&oracleCfgPath,
		"oracle-config-path",
		"",
		"oracle.json",
		"Path to the oracle config file.",
	)
	rootCmd.Flags().StringVarP(
		&marketCfgPath,
		"market-config-path",
		"",
		"market.json",
		"Path to the market config file.",
	)
	rootCmd.Flags().BoolVarP(
		&runPprof,
		"run-pprof",
		"",
		false,
		"Run pprof server.",
	)
	rootCmd.Flags().StringVarP(
		&profilePort,
		"pprof-port",
		"",
		"6060",
		"Port for the pprof server to listen on.",
	)
	rootCmd.Flags().StringVarP(
		&chain,
		"chain",
		"",
		"",
		"The chain which the side car should run for {dydx, \"\"}. Empty (default) runs against a static market.json config.",
	)
	rootCmd.Flags().BoolVarP(
		&updateLocalConfig,
		"update-local-market-config",
		"",
		true,
		"Update the market map config when a new one is received; this will overwrite the existing config file.",
	)
}

// start the oracle-grpc server + oracle process, cancel on interrupt or terminate.
func main() {
	rootCmd.Execute()
}

func runOracle() error {
	// channel with width for either signal
	sigs := make(chan os.Signal, 1)

	// gracefully trigger close on interrupt or terminate signals
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.ReadOracleConfigFromFile(oracleCfgPath)
	if err != nil {
		return fmt.Errorf("failed to read oracle config file: %s", err.Error())
	}

	marketCfg, err := types.ReadMarketConfigFromFile(marketCfgPath)
	if err != nil {
		return fmt.Errorf("failed to read market config file: %s", err.Error())
	}

	var logger *zap.Logger
	if !cfg.Production {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return fmt.Errorf("failed to create logger: %s", err.Error())
		}
	} else {
		logger, err = zap.NewProduction()
		if err != nil {
			return fmt.Errorf("failed to create logger: %s", err.Error())
		}
	}

	metrics := oraclemetrics.NewMetricsFromConfig(cfg.Metrics)
	aggregator, err := oraclemath.NewMedianAggregator(
		logger,
		marketCfg,
		metrics,
	)
	if err != nil {
		return fmt.Errorf("failed to create data aggregator: %w", err)
	}

	// Define the orchestrator and oracle options. These determine how the orchestrator and oracle are created & executed.
	orchestratorOpts := []orchestrator.Option{
		orchestrator.WithLogger(logger),
		orchestrator.WithMarketMap(marketCfg),
		orchestrator.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),             // Replace with custom API query handler factory.
		orchestrator.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory), // Replace with custom websocket query handler factory.
		orchestrator.WithMarketMapperFactory(oraclefactory.MarketMapProviderFactory),
		orchestrator.WithAggregator(aggregator),
	}
	if updateLocalConfig {
		orchestratorOpts = append(orchestratorOpts, orchestrator.WithWriteTo(marketCfgPath))
	}
	oracleOpts := []oracle.Option{
		oracle.WithLogger(logger),
		oracle.WithUpdateInterval(cfg.UpdateInterval),
		oracle.WithMetrics(metrics),
		oracle.WithMaxCacheAge(cfg.MaxPriceAge),
		oracle.WithDataAggregator(aggregator),
	}

	// Create the orchestrator and start the orchestrator.
	orch, err := orchestrator.NewProviderOrchestrator(
		cfg,
		orchestratorOpts...,
	)
	if err != nil {
		return fmt.Errorf("failed to create provider orchestrator: %w", err)
	}

	if err := orch.Start(ctx); err != nil {
		return fmt.Errorf("failed to start provider orchestrator: %w", err)
	}
	defer orch.Stop()

	// Create the oracle and start the oracle server.
	oracleOpts = append(oracleOpts, oracle.WithProviders(orch.GetPriceProviders()))
	orc, err := oracle.New(oracleOpts...)
	if err != nil {
		return fmt.Errorf("failed to create oracle: %w", err)
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
			return fmt.Errorf("failed to start prometheus metrics: %w", err)
		}

		go ps.Start()

		// close server on shut-down
		go func() {
			<-ctx.Done()
			logger.Info("stopping prometheus metrics")
			ps.Close()
		}()
	}

	if runPprof {
		endpoint := fmt.Sprintf("%s:%s", cfg.Host, profilePort)
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
	return nil
}
