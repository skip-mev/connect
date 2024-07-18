package grpc

import (
	"fmt"
	"net/url"

	grpc "google.golang.org/grpc"
)

// NewClient is a wrapper around the `grpc.NewClient` function. Which strips the
// (`http` / `https`) schemes from the URL, and returns a new client using a
// plain url (<address>:<host>) as the target.
func NewClient(
	target string,
	opts ...grpc.DialOption,
) (conn *grpc.ClientConn, err error) {
	// We need to strip the protocol / scheme from the URL
	ip, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	// create a new client
	return grpc.NewClient(
		fmt.Sprintf("%s:%s", ip.Hostname(), ip.Port()),
		opts...,
	)
}
