package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/skip-mev/slinky/providers/apis/coinmarketcap"
	"github.com/skip-mev/slinky/providers/apis/marketmap"

	_ "net/http/pprof" //nolint: gosec

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	cmdconfig "github.com/skip-mev/slinky/cmd/slinky/config"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"

	"github.com/skip-mev/slinky/cmd/build"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/pkg/log"
	oraclemath "github.com/skip-mev/slinky/pkg/math/oracle"
	oraclefactory "github.com/skip-mev/slinky/providers/factories/oracle"
	mmservicetypes "github.com/skip-mev/slinky/service/clients/marketmap/types"
	oracleserver "github.com/skip-mev/slinky/service/servers/oracle"
	promserver "github.com/skip-mev/slinky/service/servers/prometheus"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
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

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version of the oracle.",
		Args:  cobra.NoArgs,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(build.Build)
		},
	}

	oracleCfgPath       string
	marketCfgPath       string
	marketMapProvider   string
	updateMarketCfgPath string
	runPprof            bool
	profilePort         string
	logLevel            string
	fileLogLevel        string
	writeLogsTo         string
	marketMapEndPoint   string
	maxLogSize          int
	maxBackups          int
	maxAge              int
	disableCompressLogs bool
	disableRotatingLogs bool
	useCMCOnly          bool
)

const (
	DefaultLegacyConfigPath = "./oracle.json"
)

func init() {
	rootCmd.Flags().StringVarP(
		&marketMapProvider,
		"marketmap-provider",
		"",
		marketmap.Name,
		"MarketMap provider to use (marketmap_api, dydx_api, dydx_migration_api).",
	)
	rootCmd.Flags().StringVarP(
		&oracleCfgPath,
		"oracle-config",
		"",
		"",
		"Path to the oracle config file.",
	)
	rootCmd.Flags().StringVarP(
		&marketCfgPath,
		"market-config-path",
		"",
		"",
		"Path to the market config file. If you supplied a node URL in your config, this will not be required.",
	)
	rootCmd.Flags().StringVarP(
		&updateMarketCfgPath,
		"update-market-config-path",
		"",
		"",
		"Path where the current market config will be written. Overwrites any pre-existing file. Requires an http-node-url/marketmap provider in your oracle.json config.",
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
		&logLevel,
		"log-std-out-level",
		"",
		"info",
		"Log level (debug, info, warn, error, dpanic, panic, fatal).",
	)
	rootCmd.Flags().StringVarP(
		&fileLogLevel,
		"log-file-level",
		"",
		"info",
		"Log level for the file logger (debug, info, warn, error, dpanic, panic, fatal).",
	)
	rootCmd.Flags().StringVarP(
		&writeLogsTo,
		"log-file",
		"",
		"sidecar.log",
		"Write logs to a file.",
	)
	rootCmd.Flags().IntVarP(
		&maxLogSize,
		"log-max-size",
		"",
		100,
		"Maximum size in megabytes before log is rotated.",
	)
	rootCmd.Flags().IntVarP(
		&maxBackups,
		"log-max-backups",
		"",
		1,
		"Maximum number of old log files to retain.",
	)
	rootCmd.Flags().IntVarP(
		&maxAge,
		"log-max-age",
		"",
		3,
		"Maximum number of days to retain an old log file.",
	)
	rootCmd.Flags().BoolVarP(
		&disableCompressLogs,
		"log-file-disable-compression",
		"",
		false,
		"Compress rotated log files.",
	)
	rootCmd.Flags().BoolVarP(
		&disableRotatingLogs,
		"log-disable-file-rotation",
		"",
		false,
		"Disable writing logs to a file.",
	)
	rootCmd.Flags().StringVarP(
		&marketMapEndPoint,
		"market-map-endpoint",
		"",
		"",
		"Use a custom listen-to endpoint for market-map (overwrites what is provided in oracle-config).",
	)
	rootCmd.Flags().BoolVarP(
		&useCMCOnly,
		"use-cmc-only",
		"",
		false,
		"Use CoinMarketCap only for price data. **This should only be used for testing.**. Only works if you pair this with the --market-config-path flag.",
	)
	rootCmd.MarkFlagsMutuallyExclusive("update-market-config-path", "market-config-path")
	rootCmd.MarkFlagsMutuallyExclusive("market-map-endpoint", "market-config-path")

	rootCmd.AddCommand(versionCmd)
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

	// Set up logging.
	logCfg := log.NewDefaultConfig()
	logCfg.StdOutLogLevel = logLevel
	logCfg.FileOutLogLevel = fileLogLevel
	logCfg.DisableRotating = disableRotatingLogs
	logCfg.WriteTo = writeLogsTo
	logCfg.MaxSize = maxLogSize
	logCfg.MaxBackups = maxBackups
	logCfg.MaxAge = maxAge
	logCfg.Compress = !disableCompressLogs

	// Build logger.
	logger := log.NewLogger(logCfg)
	defer logger.Sync()

	var cfg config.OracleConfig
	var err error

	cfg, err = cmdconfig.ReadOracleConfigWithOverrides(oracleCfgPath, marketMapProvider)
	if err != nil {
		return fmt.Errorf("failed to get oracle config: %w", err)
	}

	// overwrite endpoint
	if marketMapEndPoint != "" {
		cfg, err = overwriteMarketMapEndpoint(cfg, marketMapEndPoint)
		if err != nil {
			return fmt.Errorf("failed to overwrite market endpoint %s: %w", marketMapEndPoint, err)
		}
	}

	var marketCfg mmtypes.MarketMap
	if marketCfgPath != "" {
		marketCfg, err = mmtypes.ReadMarketMapFromFile(marketCfgPath)
		if err != nil {
			return fmt.Errorf("failed to read market config file: %w", err)
		}

		if useCMCOnly {
			marketCfg = filterToOnlyCMCMarkets(marketCfg)
		}
	}

	logger.Info(
		"successfully read in configs",
		zap.String("oracle_config_path", oracleCfgPath),
		zap.String("market_config_path", marketCfgPath),
	)

	metrics := oraclemetrics.NewMetricsFromConfig(cfg.Metrics)
	aggregator, err := oraclemath.NewIndexPriceAggregator(
		logger,
		marketCfg,
		metrics,
	)
	if err != nil {
		return fmt.Errorf("failed to create data aggregator: %w", err)
	}

	// Define the oracle options. These determine how the oracle is created & executed.
	oracleOpts := []oracle.Option{
		oracle.WithLogger(logger),
		oracle.WithMarketMap(marketCfg),
		oracle.WithPriceAPIQueryHandlerFactory(oraclefactory.APIQueryHandlerFactory),             // Replace with custom API query handler factory.
		oracle.WithPriceWebSocketQueryHandlerFactory(oraclefactory.WebSocketQueryHandlerFactory), // Replace with custom websocket query handler factory.
		oracle.WithMarketMapperFactory(oraclefactory.MarketMapProviderFactory),
		oracle.WithMetrics(metrics),
	}
	if updateMarketCfgPath != "" {
		oracleOpts = append(oracleOpts, oracle.WithWriteTo(updateMarketCfgPath))
	}

	// Create the oracle and start the oracle.
	orc, err := oracle.New(
		cfg,
		aggregator,
		oracleOpts...,
	)
	if err != nil {
		return fmt.Errorf("failed to create oracle: %w", err)
	}
	go func() {
		if err := orc.Start(ctx); err != nil {
			logger.Fatal("failed to start oracle", zap.Error(err))
		}
	}()
	defer orc.Stop()

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

	// start server (blocks).
	if err := srv.StartServer(ctx, cfg.Host, cfg.Port); err != nil {
		logger.Error("stopping server", zap.Error(err))
	}
	return nil
}

func overwriteMarketMapEndpoint(cfg config.OracleConfig, overwrite string) (config.OracleConfig, error) {
	for providerName, provider := range cfg.Providers {
		if provider.Type == mmservicetypes.ConfigType {
			provider.API.Endpoints = []config.Endpoint{
				{
					URL: overwrite,
				},
			}
			cfg.Providers[providerName] = provider
			return cfg, cfg.ValidateBasic()
		}
	}

	return cfg, fmt.Errorf("no market-map provider found in config")
}

// filterToOnlyCMCMarkets is a helper function that filters out all markets that are not from CoinMarketCap. It
// mutates the marketmap to only include CoinMarketCap markets.
func filterToOnlyCMCMarkets(marketmap mmtypes.MarketMap) mmtypes.MarketMap {
	res := mmtypes.MarketMap{
		Markets: make(map[string]mmtypes.Market),
	}

	// Filter out all markets that are not from CoinMarketCap.
	for _, market := range marketmap.Markets {
		var meta metaDataJson
		if err := json.Unmarshal([]byte(market.Ticker.Metadata_JSON), &meta); err != nil {
			continue
		}

		var id string
		for _, aggregateID := range meta.AggregateIDs {
			if aggregateID.Venue == "coinmarketcap" {
				id = aggregateID.ID
				break
			}
		}

		if len(id) == 0 {
			continue
		}

		resTicker := market.Ticker
		resTicker.MinProviderCount = 1

		providers := []mmtypes.ProviderConfig{
			{
				Name:           coinmarketcap.Name,
				OffChainTicker: id,
			},
		}

		res.Markets[resTicker.CurrencyPair.String()] = mmtypes.Market{
			Ticker:          resTicker,
			ProviderConfigs: providers,
		}
	}

	return res
}

// Ref:
// {\"reference_price\":0,\"liquidity\":0,\"aggregate_ids\":[{\"venue\":\"coinmarketcap\",\"ID\":\"4030\"}]}
type metaDataJson struct {
	AggregateIDs []aggregateIDsJson `json:"aggregate_ids"`
}
type aggregateIDsJson struct {
	Venue string `json:"venue"`
	ID    string `json:"ID"`
}
