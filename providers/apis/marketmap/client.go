package marketmap

import (
	"github.com/skip-mev/slinky/oracle/config"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// NewGRPCClient returns a new GRPC client for MarketMap module.
func NewGRPCClient(
	config config.APIConfig,
) (mmtypes.QueryClient, error) {
	kacp := keepalive.ClientParameters{
		Time:                config.Interval, // send pings every 10 seconds if there is no activity
		Timeout:             config.Timeout,  // wait a second for ping ack before considering the connection dead
		PermitWithoutStream: true,            // send pings even without active streams
	}

	// TODO: Do we want to ignore proxy settings?
	conn, err := grpc.Dial(
		config.URL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*4)), // Set max receive message size to 4MB
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(1024*1024*4)), // Set max send message size to 4MB
		grpc.WithKeepaliveParams(kacp),
	)
	if err != nil {
		return nil, err
	}

	return mmtypes.NewQueryClient(conn), nil
}
