package types

//go:generate mockery --name QueryClient
type ClientWrapper interface {
	QueryClient
}
