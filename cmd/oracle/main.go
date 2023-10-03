package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"cosmossdk.io/log"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/service/server"
)

var (
	host    = flag.String("host", "localhost", "host for the grpc-service to listen on")
	port    = flag.String("port", "8080", "port for the grpc-service to listen on")
	cfgPath = flag.String("config", "config.toml", "path to config file")
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

	// create oracle
	cfg, err := config.ReadConfigFromFile(*cfgPath)
	if err != nil {
		logger.Error("failed to read config file", "err", err)
		return
	}

	o, err := oracle.NewOracleFromConfig(logger, cfg)
	if err != nil {
		logger.Error("failed to create oracle from config", "err", err)
		return
	}

	// create server
	srv := server.NewOracleServer(o, logger)

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
	if err := config.StartPrometheusServer(cfg.Metrics, logger); err != nil {
		logger.Error("failed to start prometheus metrics", "err", err)
		return
	}

	// start oracle + server, and wait for either to finish
	if err := srv.StartServer(ctx, *host, *port); err != nil {
		logger.Error("stopping server", "err", err)
	}
}
