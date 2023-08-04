package service

import (
	"context"
)

// OracleService defines the service all clients must implement.
//
//go:generate mockery --name OracleService --filename mock_oracle_service.go
type OracleService interface {
	OracleServer
	Start(context.Context) error
	Stop(context.Context) error
}
