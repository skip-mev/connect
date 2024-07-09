package osmosis_test

import (
	"testing"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/defi/osmosis"
)

func TestMultiGRPCClient(t *testing.T) {
	cfg := osmosis.DefaultAPIConfig
	cfg.Endpoints = []config.Endpoint{
		{
			URL: "http://localhost:8899",
		},
		{
			URL: "http://localhost:8899/",
			Authentication: config.Authentication{
				APIKey:       "test",
				APIKeyHeader: "X-API-Key",
			},
		},
		{
			URL: "http://localhost:8899/",
		},
	}
}
