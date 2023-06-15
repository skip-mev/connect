package client

import (
	"context"
	"errors"

	"github.com/skip-mev/slinky/oracle"
	"github.com/skip-mev/slinky/service"
)

var _ service.OracleService = (*LocalClient)(nil)

// LocalClient defines an implementation of a local, i.e. in-process, oracle client.
// This client can be used in ABCI++ calls where the application wants the oracle
// process to be run in-process. The client must be started upon app construction
// and stopped upon app shutdown/cleanup.
type LocalClient struct {
	oracle *oracle.Oracle
}

func NewLocalClient(o *oracle.Oracle) *LocalClient {
	return &LocalClient{
		oracle: o,
	}
}

func (c *LocalClient) Prices(_ context.Context, req *service.QueryPricesRequest) (*service.QueryPricesResponse, error) {
	if req == nil {
		return nil, ErrorNilRequest
	}

	prices := c.oracle.GetPrices()

	resp := &service.QueryPricesResponse{
		Prices:    make(map[string]string, len(prices)),
		Timestamp: c.oracle.GetLastSyncTime(),
	}

	for k, v := range prices {
		resp.Prices[k.String()] = v.String()
	}

	return resp, nil
}

// Note: Start(ctx) is a blocking call, so the caller will need to run it in a
// goroutine.
func (c *LocalClient) Start(ctx context.Context) error {
	if c.oracle.IsRunning() {
		return errors.New("oracle already running")
	}

	return c.oracle.Start(ctx)
}

func (c *LocalClient) Stop(_ context.Context) error {
	c.oracle.Stop()
	return nil
}
