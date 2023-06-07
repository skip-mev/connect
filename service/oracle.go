package service

import (
	"context"
)

// OracleService defines the service all clients must implement.
type OracleService interface {
	OracleServer

	Start(context.Context) error
	Stop(context.Context) error
}
