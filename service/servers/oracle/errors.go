package oracle

import "errors"

var (
	ErrNilRequest       = errors.New("request cannot be nil")
	ErrOracleNotRunning = errors.New("oracle is not running")
	ErrContextCancelled = errors.New("context cancelled")
)
