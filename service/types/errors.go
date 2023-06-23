package types

import (
	"errors"
)

// Client sentinel errors.
var (
	ErrorNilRequest       = errors.New("request cannot be nil")
	ErrorProviderNotFound = errors.New("provider not found")
	ErrorOracleNotRunning = errors.New("oracle is not running")
	ErrorContextCancelled = errors.New("context cancelled")
)
