package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

type (
	// PriceFeedUpdates is utilized to track price feed updates as they are received in the
	// SLA module's preblock handler.
	PriceFeedUpdates struct {
		// ValidatorUpdates is a map of validator updates. The key is the validator's consensus
		// address and the value is the validator's price feed updates.
		ValidatorUpdates map[string]ValidatorUpdate

		// CurrencyPairs is a set of currency pairs that are supported by the network.
		CurrencyPairs map[slinkytypes.CurrencyPair]struct{}
	}

	// ValidatorUpdate is utilized to map a validator's price feed updates.
	ValidatorUpdate struct {
		// ConsAddress is the validator's consensus address.
		ConsAddress sdk.ConsAddress

		// Updates is a map of price feed updates. The key is the currency pair and the value is
		Updates map[slinkytypes.CurrencyPair]slatypes.UpdateStatus
	}
)

// NewPriceFeedUpdates returns a new PriceFeedUpdates.
func NewPriceFeedUpdates() PriceFeedUpdates {
	return PriceFeedUpdates{
		ValidatorUpdates: make(map[string]ValidatorUpdate),
		CurrencyPairs:    make(map[slinkytypes.CurrencyPair]struct{}),
	}
}

// NewValidatorUpdate returns a new ValidatorUpdate.
func NewValidatorUpdate(consAddress sdk.ConsAddress) ValidatorUpdate {
	return ValidatorUpdate{
		ConsAddress: consAddress,
		Updates:     make(map[slinkytypes.CurrencyPair]slatypes.UpdateStatus),
	}
}

// UpdatePriceFeeds will update the price feed incentives for all given updates. The
// updates parameter is constructed in the preblock handler and contains all price feed
// updates for the current block for every validator and currency pair. The validators included
// are the ones in the active set from the previous block. There are a few cases that need to be
// handled:
// 1. A new validator is added to the active set.
// 2. A validator is removed from the active set.
// 3. A new currency pair is added to the network.
// 4. A currency pair is removed from the network.
// 5. A currency pair is updated.
func (k *Keeper) UpdatePriceFeeds(ctx sdk.Context, updates PriceFeedUpdates) error {
	slas, err := k.GetSLAs(ctx)
	if err != nil {
		return err
	}

	// Determine the set of currency pairs that are currently stored in the x/sla module's state.
	// but are not supported by the network anymore.
	cpsToRemove, err := k.GetCurrencyPairs(ctx)
	if err != nil {
		return err
	}
	for cp := range cpsToRemove {
		if _, ok := updates.CurrencyPairs[cp]; ok {
			delete(cpsToRemove, cp)
		}
	}

	// Update the currency pairs that are currently stored in the x/sla module's state.
	if err := k.SetCurrencyPairs(ctx, updates.CurrencyPairs); err != nil {
		return err
	}

	// Update the price feeds for each SLA.
	for _, sla := range slas {
		for cp := range cpsToRemove {
			if err := k.RemovePriceFeedByCurrencyPair(ctx, sla.ID, cp); err != nil {
				return err
			}
		}

		if err := k.UpdatePriceFeedsForSLA(ctx, sla, updates); err != nil {
			return err
		}
	}

	return nil
}

// UpdatePriceFeedsForSLA will update the price feeds for given SLA.
func (k *Keeper) UpdatePriceFeedsForSLA(ctx sdk.Context, sla slatypes.PriceFeedSLA, updates PriceFeedUpdates) error {
	for _, validator := range updates.ValidatorUpdates {
		for cp, status := range validator.Updates {
			contains, err := k.ContainsPriceFeed(ctx, sla.ID, cp, validator.ConsAddress)
			if err != nil {
				return err
			}

			if contains {
				if err := k.updatePriceFeedWithStatus(ctx, sla, cp, validator.ConsAddress, status); err != nil {
					return err
				}
			} else {
				if err := k.initPriceFeedWithStatus(ctx, sla, cp, validator.ConsAddress, status); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

// updatePriceFeedWithStatus will update the price feed with the given status and add it to the
// x/sla module's state.
func (k *Keeper) updatePriceFeedWithStatus(
	ctx sdk.Context,
	sla slatypes.PriceFeedSLA,
	cp slinkytypes.CurrencyPair,
	validator sdk.ConsAddress,
	status slatypes.UpdateStatus,
) error {
	feed, err := k.GetPriceFeed(ctx, sla.ID, cp, validator)
	if err != nil {
		return err
	}

	if err := feed.SetUpdate(status); err != nil {
		return err
	}

	return k.SetPriceFeed(ctx, feed)
}

// initPriceFeedWithStatus will initialize a price feed with the given status and add it to the
// x/sla module's state.
func (k *Keeper) initPriceFeedWithStatus(
	ctx sdk.Context,
	sla slatypes.PriceFeedSLA,
	cp slinkytypes.CurrencyPair,
	validator sdk.ConsAddress,
	status slatypes.UpdateStatus,
) error {
	feed, err := slatypes.NewPriceFeed(uint(sla.MaximumViableWindow), validator, cp, sla.ID)
	if err != nil {
		return err
	}

	if err := feed.SetUpdate(status); err != nil {
		return err
	}

	return k.SetPriceFeed(ctx, feed)
}
