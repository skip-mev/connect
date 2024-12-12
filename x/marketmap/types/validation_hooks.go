package types

import (
	"context"
	"fmt"
)

// MarketValidationHook is a hook that is called for stateful validation of a market before
// some keeper operation is performed on it.
type MarketValidationHook func(ctx context.Context, market Market) error

// MarketValidationHooks is a type alias for an array of MarketValidationHook.
type MarketValidationHooks []MarketValidationHook

// ValidateMarket calls all validation hooks for the given market.
func (h MarketValidationHooks) ValidateMarket(ctx context.Context, market Market) error {
	for _, hook := range h {
		if err := hook(ctx, market); err != nil {
			return fmt.Errorf("failed validation hooks for market %s: %w", market.Ticker.String(), err)
		}
	}

	return nil
}

// DefaultDeleteMarketValidationHooks returns the default DeleteMarketValidationHook as an array.
func DefaultDeleteMarketValidationHooks() MarketValidationHooks {
	hooks := MarketValidationHooks{
		DefaultDeleteMarketValidationHook(),
	}

	return hooks
}

// DefaultDeleteMarketValidationHook returns the default DeleteMarketValidationHook for x/marketmap.
// This hook checks:
// - if the given market is enabled - error
// - if the given market is disabled - return nil.
func DefaultDeleteMarketValidationHook() MarketValidationHook {
	return func(_ context.Context, market Market) error {
		if market.Ticker.Enabled {
			return fmt.Errorf("market is enabled - cannot be deleted")
		}

		return nil
	}
}
