package client

import (
	"context"

	"cosmossdk.io/log"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/service"
)

// NewOracleServiceFromConfig reads a config and instantiates either a grpc-client / local-client from a config
// and returns a new OracleService.
func NewOracleServiceFromConfig(cfg config.Config, l log.Logger) (service.OracleService, error) {
	var oracleService service.OracleService

	if cfg.InProcess {
		// retrieve oracle from a config
		oracle, err := oracle.NewOracleFromConfig(l, &cfg)
		if err != nil {
			return nil, err
		}

		oracleService = NewLocalClient(oracle, cfg.Timeout)
	} else {
		oracleService = NewGRPCClient(cfg.RemoteAddress, cfg.Timeout)
	}

	// start the service in a go-routine
	go oracleService.Start(context.Background())
	return oracleService, nil
}
