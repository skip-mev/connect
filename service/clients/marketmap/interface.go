package marketmap

import (
	"context"
)

// MarketMapClient defines the interface that will be utilized by the oracle side-car
// to query the marketmap module. The marketmap module is responsible for maintaining
// the cannonical mapping of providers -> assets as well as all of the conversion logic
// to resolve price feeds to a common set of currency pairs.
//
//go:generate mockery --name MarketMapClient --filename mock_marketmap_client.go
type MarketMapClient interface { //nolint
	// Start starts the marketmap client. This should connect to the remote marketmap
	// service and return an error if the connection fails.
	Start(ctx context.Context) error

	// Stop stops the marketmap client.
	Stop() error
}

// NoOpClient is a no-op implementation of the MarketMapClient interface. This
// implementation is used when the marketmap service is disabled or not utilized.
type NoOpClient struct{}

// Start is a no-op.
func (NoOpClient) Start(context.Context) error {
	return nil
}

// Stop is a no-op.
func (NoOpClient) Stop() error {
	return nil
}
