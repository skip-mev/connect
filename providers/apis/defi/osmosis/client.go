package osmosis

import (
	"context"

	"github.com/osmosis-labs/osmosis/v25/x/poolmanager/client/queryproto"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base/api/metrics"
)

var _ GRPCCLient = &GRPCCLientImpl{}

// GRPCClient is the expected interface for an osmosis grpc client.
//
//go:generate mockery --name GRPCCLient --output ./mocks/ --case underscore
type GRPCCLient interface {
	SpotPrice(grpcCtx context.Context,
		req *queryproto.SpotPriceRequest,
	) (*queryproto.SpotPriceResponse, error)
}

type GRPCCLientImpl struct {
	api        config.APIConfig
	apiMetrics metrics.APIMetrics

	pmClient queryproto.QueryClient
}

func (c *GRPCCLientImpl) SpotPrice(grpcCtx context.Context, req *queryproto.SpotPriceRequest) (*queryproto.SpotPriceResponse, error) {
	// TODO implement me
	panic("implement me")
}
