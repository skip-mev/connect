package client

import (
	"math/big"

	"cosmossdk.io/log"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	oraclemetrics "github.com/skip-mev/slinky/oracle/metrics"
	"github.com/skip-mev/slinky/service"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// NewOracleService reads a config and instantiates either a grpc-client / local-client from a config
// and returns a new OracleService.
func NewOracleService(
	logger log.Logger,
	oracleCfg config.OracleConfig,
	metricsCfg config.MetricsConfig,
	factory oracle.ProviderFactory,
	aggregateFn aggregator.AggregateFn[string, map[oracletypes.CurrencyPair]*big.Int],
) (service.OracleService, error) {
	var (
		oracleService service.OracleService
		metrics       = oraclemetrics.NewNopMetrics()
	)

	if !oracleCfg.Enabled {
		return service.NewNoopOracleService(), nil
	}

	if oracleCfg.InProcess {
		if metricsCfg.OracleMetrics.Enabled {
			metrics = oraclemetrics.NewMetrics()
		}

		oracle, err := oracle.New(
			logger,
			oracleCfg,
			factory,
			aggregateFn,
			metrics,
		)
		if err != nil {
			return nil, err
		}

		oracleService = NewLocalClient(oracle, oracleCfg.Timeout)
	} else {
		oracleService = NewGRPCClient(oracleCfg.RemoteAddress, oracleCfg.Timeout)
	}

	return oracleService, nil
}
