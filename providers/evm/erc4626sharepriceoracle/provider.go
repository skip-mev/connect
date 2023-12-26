package erc4626sharepriceoracle

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

var _ providertypes.Provider[oracletypes.CurrencyPair, *big.Int] = (*ERC4626SharePriceProvider)(nil)

// ERC4626SharePriceProvider implements a provider that fetches data from the ERC4626SharePrice API.
type ERC4626SharePriceProvider struct {
	*base.BaseProvider[oracletypes.CurrencyPair, *big.Int]
}

// NewProvider returns a new ERC4626SharePriceOracle provider. It uses the provided API-key to
// make RPC calls to Alchemy. Note that only the Quote denom is used; the Quote/Base pair is
// naturally determined by the contract address, so be sure the configured addresses are
// correct.
func NewProvider(
	logger *zap.Logger,
	pairs []oracletypes.CurrencyPair,
	providerConfig config.ProviderConfig,
) (*ERC4626SharePriceProvider, error) {
	handler, err := NewERC4626SharePriceAPIHandler(logger, pairs, providerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create coingecko api handler: %w", err)
	}

	provider, err := base.NewProvider(logger, providerConfig, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to create base provider: %w", err)
	}

	return &ERC4626SharePriceProvider{
		BaseProvider: provider,
	}, nil
}
