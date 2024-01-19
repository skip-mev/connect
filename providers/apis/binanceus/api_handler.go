package binanceus

import (
	"fmt"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/apis/binance"
)

const (
	// Name is the name of the provider.
	Name = "binanceus"
)

// NewBinanceUSAPIHandler returns a new Binance US API handler.
func NewBinanceUSAPIHandler(
	cfg config.ProviderConfig,
) (*binance.APIHandler, error) {
	if cfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, cfg.Name)
	}

	return &binance.APIHandler{
		ProviderConfig: cfg,
		BaseURL:        BaseURL,
	}, nil
}
