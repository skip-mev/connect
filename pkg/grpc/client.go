package grpc

import (
	"fmt"
	"net"
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
	// check if this is a host:port URI, if so continue,
	// otherwise, parse the URL and extract the host and port
	host, port, err := net.SplitHostPort(target)
	if err != nil {
		// parse the URL
		ip, err := url.Parse(target)
		if err != nil {
			return nil, err
		}

		// extract the host and port
		host, port = ip.Hostname(), ip.Port()
	}

	// create a new client
	return grpc.NewClient(
		fmt.Sprintf("%s:%s", host, port),
		opts...,
	)
}
