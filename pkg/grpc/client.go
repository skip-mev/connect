package grpc

import (
	"strings"

	grpc "google.golang.org/grpc"
)

// NewClient is a wrapper around the `grpc.NewClient` function. Which strips the
// (`http` / `https`) schemes from the URL, and returns a new client using a
// plain url (<address>:<host>) as the target.
func NewClient(
	target string,
	opts ...grpc.DialOption,
) (conn *grpc.ClientConn, err error) {
	// strip the scheme from the target
	target = strings.TrimPrefix(strings.TrimPrefix(target, "http://"), "https://")

	// create a new client
	return grpc.NewClient(
		target,
		opts...,
	)
}
