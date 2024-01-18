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
	providerCfg config.ProviderConfig,
) (*binance.APIHandler, error) {
	if providerCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, providerCfg.Name)
	}

	cfg, err := binance.ReadBinanceConfigFromFile(providerCfg.Path)
	if err != nil {
		return nil, err
	}

	return &binance.APIHandler{
		Config:  cfg,
		BaseURL: BaseURL,
	}, nil
}
