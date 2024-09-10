package http

import "net"

func IsValidAddress(address string) bool {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return false
	}

	if host == "" || port == "" {
		return false
	}

	return true
}
