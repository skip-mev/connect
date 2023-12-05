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

type NoopOracleService struct{}

// NewNoopOracleService returns a new NoopOracleService.
func NewNoopOracleService() OracleService {
	return &NoopOracleService{}
}

// Start is a no-op.
func (s *NoopOracleService) Start(context.Context) error {
	return nil
}

// Stop is a no-op.
func (s *NoopOracleService) Stop(context.Context) error {
	return nil
}

// Prices is a no-op.
func (s *NoopOracleService) Prices(context.Context, *QueryPricesRequest) (*QueryPricesResponse, error) {
	return &QueryPricesResponse{}, nil
}
