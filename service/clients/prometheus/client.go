package prometheus

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"go.uber.org/zap"

	libhttp "github.com/skip-mev/connect/v2/pkg/http"
)

type Client struct {
	v1.API
	logger *zap.Logger
}

func NewClient(address string, logger *zap.Logger) (Client, error) {
	if logger == nil {
		logger = zap.NewNop()
	}

	logger = logger.With(zap.String("service", "prometheus_client"))
	logger.Info("creating prometheus client", zap.String("address", address))

	// get the prometheus server address
	if address == "" || !libhttp.IsValidAddress(address) {
		return Client{}, fmt.Errorf("invalid prometheus server address: %s", address)
	}

	const httpPrefix = "http://"
	if !strings.HasPrefix(address, httpPrefix) {
		address = httpPrefix + address
	}

	// Create a Prometheus API client
	client, err := api.NewClient(api.Config{
		Address: address, // Address of your Prometheus server
	})
	if err != nil {
		return Client{}, fmt.Errorf("failed to create prometheus client: %w", err)
	}

	// Create a new Prometheus API v1 client
	return Client{
		API:    v1.NewAPI(client),
		logger: logger,
	}, nil
}
