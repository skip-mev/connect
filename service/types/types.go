package types

import (
	"context"
	"time"

	"github.com/holiman/uint256"

	"github.com/skip-mev/slinky/x/oracle/types"
)

const (
	Transport = "tcp"
)

// Oracle is interface the OracleServer expects its underlying oracle to implement
//
//go:generate mockery --name Oracle --filename mock_oracle.go
type Oracle interface {
	IsRunning() bool
	GetLastSyncTime() time.Time
	GetPrices() map[types.CurrencyPair]*uint256.Int
	Start(ctx context.Context) error
	Stop()
}
