package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckMarketMapEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "Valid gRPC endpoint",
			endpoint: "example.com:8080",
			wantErr:  false,
		},
		{
			name:     "Valid IP address endpoint",
			endpoint: "192.168.1.1:9090",
			wantErr:  false,
		},
		{
			name:     "HTTP endpoint",
			endpoint: "http://example.com:8080",
			wantErr:  true,
			errMsg:   `expected gRPC endpoint but got HTTP endpoint "http://example.com:8080". Please provide a gRPC endpoint (e.g. some.host:9090)`,
		},
		{
			name:     "HTTPS endpoint",
			endpoint: "https://example.com:8080",
			wantErr:  true,
			errMsg:   `expected gRPC endpoint but got HTTP endpoint "https://example.com:8080". Please provide a gRPC endpoint (e.g. some.host:9090)`,
		},
		{
			name:     "Missing port",
			endpoint: "example.com",
			wantErr:  true,
			errMsg:   `invalid gRPC endpoint "example.com". Must specify port (e.g. example.com:9090)`,
		},
		{
			name:     "Invalid port format",
			endpoint: "example.com:port",
			wantErr:  true,
			errMsg:   `invalid gRPC endpoint "example.com:port". Must specify port (e.g. example.com:9090)`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isValidGRPCEndpoint(tt.endpoint)
			if tt.wantErr {
				require.EqualError(t, err, tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
