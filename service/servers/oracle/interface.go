package oracle

import (
	"context"

	"github.com/skip-mev/slinky/service/servers/oracle/types"
)

// OracleService defines the service all clients must implement.
//
//go:generate mockery --name OracleService --filename mock_oracle_service.go
type OracleService interface { //nolint
	types.OracleServer

	Start(context.Context) error
	Stop(context.Context) error
}
