package marketmap

import (
	"github.com/skip-mev/slinky/oracle/config"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NewGRPCClient returns a new GRPC client for MarketMap module.
func NewGRPCClient(
	config config.APIConfig,
) (mmtypes.QueryClient, error) {
	// TODO: Do we want to ignore proxy settings?
	conn, err := grpc.Dial(
		config.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return mmtypes.NewQueryClient(conn), nil
}
