package osmosis

import (
	gammtypes "github.com/osmosis-labs/osmosis/v25/x/gamm/types"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
)

type GRPCCLient struct {
	api        config.APIConfig
	apiMetrics metrics.APIMetrics

	gammClient gammtypes.QueryClient
}
