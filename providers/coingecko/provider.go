package coingecko

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

var _ providertypes.Provider[oracletypes.CurrencyPair, *big.Int] = (*CoinGeckoProvider)(nil)

// CoinGeckoProvider implements a provider that fetches data from the CoinGecko API.
type CoinGeckoProvider struct { //nolint
	*base.BaseProvider[oracletypes.CurrencyPair, *big.Int]
}

// NewProvider returns a new CoinGecko provider.
func NewProvider(
	logger *zap.Logger,
	pairs []oracletypes.CurrencyPair,
	providerConfig config.ProviderConfig,
) (*CoinGeckoProvider, error) {
	handler, err := NewCoinGeckoAPIHandler(logger, pairs, providerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create coingecko api handler: %w", err)
	}

	provider, err := base.NewProvider(logger, providerConfig, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to create base provider: %w", err)
	}

	return &CoinGeckoProvider{
		BaseProvider: provider,
	}, nil
}
