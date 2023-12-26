package erc4626

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

var _ providertypes.Provider[oracletypes.CurrencyPair, *big.Int] = (*ERC4626Provider)(nil)

// ERC4626Provider implements a provider that fetches data from the ERC4626 API.
type ERC4626Provider struct { //nolint
	*base.BaseProvider[oracletypes.CurrencyPair, *big.Int]
}

// NewProvider returns a new ERC4626 provider. It uses the provided API-key to make RPC calls to Alchemy.
// Note that only the Quote denom is used; the Base denom is naturally determined by the contract address.
func NewProvider(
	logger *zap.Logger,
	pairs []oracletypes.CurrencyPair,
	providerCfg config.ProviderConfig,
) (*ERC4626Provider, error) {
	handler, err := NewERC4626APIHandler(logger, pairs, providerCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create erc4626 api handler: %w", err)
	}

	provider, err := base.NewProvider(logger, providerCfg, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to create base provider: %w", err)
	}

	return &ERC4626Provider{
		BaseProvider: provider,
	}, nil
}
