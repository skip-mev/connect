package coinbase

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

var _ providertypes.Provider[oracletypes.CurrencyPair, *big.Int] = (*CoinBaseProvider)(nil)

// CoinBaseProvider implements a provider that fetches data from the Coinbase API.
type CoinBaseProvider struct { //nolint
	*base.BaseProvider[oracletypes.CurrencyPair, *big.Int]
}

// NewProvider returns a new Coinbase provider.
//
// THIS PROVIDER SHOULD NOT BE USED IN PRODUCTION. IT IS ONLY MEANT FOR TESTING.
func NewProvider(
	logger *zap.Logger,
	pairs []oracletypes.CurrencyPair,
	providerConfig config.ProviderConfig,
) (*CoinBaseProvider, error) {
	handler, err := NewCoinBaseAPIHandler(logger, pairs, providerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create coinbase api handler: %w", err)
	}

	provider, err := base.NewProvider(logger, providerConfig, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to create base provider: %w", err)
	}

	return &CoinBaseProvider{
		BaseProvider: provider,
	}, nil
}
