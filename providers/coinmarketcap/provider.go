package coinmarketcap

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

var _ providertypes.Provider[oracletypes.CurrencyPair, *big.Int] = (*CoinMarketCapProvider)(nil)

// CoinMarketCapProvider implements a provider that fetches data from the CoinMarketCap API.
type CoinMarketCapProvider struct { //nolint
	*base.BaseProvider[oracletypes.CurrencyPair, *big.Int]
}

// NewProvider returns a new CoinMarketCap provider. It uses the provided API-key in the
// header of outgoing requests to CoinMarketCap's API.
func NewProvider(
	logger *zap.Logger,
	pairs []oracletypes.CurrencyPair,
	providerConfig config.ProviderConfig,
) (*CoinMarketCapProvider, error) {
	handler, err := NewCoinMarketCapAPIHandler(logger, pairs, providerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create coinmarketcap api handler: %w", err)
	}

	provider, err := base.NewProvider(logger, providerConfig, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to create base provider: %w", err)
	}

	return &CoinMarketCapProvider{
		BaseProvider: provider,
	}, nil
}
