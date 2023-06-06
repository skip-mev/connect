package service

// OracleService defines the service all clients must implement.
type OracleService interface {
	OracleServer

	Start() error
	Stop() error
}
