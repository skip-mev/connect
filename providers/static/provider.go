package static

import (
	"fmt"
	"math/big"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/providers/base"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

var _ providertypes.Provider[oracletypes.CurrencyPair, *big.Int] = (*StaticMockProvider)(nil)

// StaticMockProvider implements a provider that returns static values for all pairs.
type StaticMockProvider struct { //nolint
	*base.BaseProvider[oracletypes.CurrencyPair, *big.Int]
}

// NewProvider returns a new static mock provider.
func NewProvider(
	logger *zap.Logger,
	pairs []oracletypes.CurrencyPair,
	providerConfig config.ProviderConfig,
) (*StaticMockProvider, error) {
	handler, err := NewStaticMockAPIHandler(logger, pairs, providerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create static mock api handler: %w", err)
	}

	provider, err := base.NewProvider(logger, providerConfig, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to create base provider: %w", err)
	}

	return &StaticMockProvider{
		BaseProvider: provider,
	}, nil
}
