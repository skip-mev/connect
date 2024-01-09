package client

import (
	"math/big"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/oracle/config"
	providertypes "github.com/skip-mev/slinky/providers/types"
	"github.com/skip-mev/slinky/service"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

// NewOracleService reads a config and instantiates either a grpc-client / local-client from a config
// and returns a new OracleService.
func NewOracleService(
	logger *zap.Logger,
	oracleCfg config.OracleConfig,
	metricsCfg config.MetricsConfig,
	factory providertypes.ProviderFactory[oracletypes.CurrencyPair, *big.Int],
	aggregateFn aggregator.AggregateFn[string, map[oracletypes.CurrencyPair]*big.Int],
) (service.OracleService, error) {
	if !oracleCfg.Enabled {
		return service.NewNoopOracleService(), nil
	}

	var oracleService service.OracleService
	if oracleCfg.InProcess {
		providers, err := factory(logger, oracleCfg, metricsCfg.OracleMetrics)
		if err != nil {
			return nil, err
		}

		oracle, err := oracle.New(
			oracleCfg,
			oracle.WithLogger(logger),
			oracle.WithMetricsConfig(metricsCfg.OracleMetrics),
			oracle.WithProviders(providers),
			oracle.WithAggregateFunction(aggregateFn),
		)
		if err != nil {
			return nil, err
		}

		oracleService = NewLocalClient(oracle, oracleCfg.ClientTimeout)
	} else {
		oracleService = NewGRPCClient(oracleCfg.RemoteAddress, oracleCfg.ClientTimeout)
	}

	return oracleService, nil
}
