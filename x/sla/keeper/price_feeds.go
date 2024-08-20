package keeper

import (
	"cosmossdk.io/collections"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

// PriceFeedCB is a callback function that can be used to process price feeds
// as they are received.
type PriceFeedCB func(priceFeed slatypes.PriceFeed) error

// SetPriceFeed adds a price feed to the x/sla module's state.
func (k *Keeper) SetPriceFeed(
	ctx sdk.Context,
	priceFeed slatypes.PriceFeed,
) error {
	key := collections.Join3(priceFeed.ID, priceFeed.CurrencyPair.String(), priceFeed.Validator)
	return k.priceFeeds.Set(ctx, key, priceFeed)
}

// GetPriceFeed returns the price feed with the given ID from the x/sla module's
// state.
func (k *Keeper) GetPriceFeed(
	ctx sdk.Context,
	slaID string,
	cp slinkytypes.CurrencyPair,
	consAddress sdk.ConsAddress,
) (slatypes.PriceFeed, error) {
	key := collections.Join3(slaID, cp.String(), consAddress.Bytes())
	return k.priceFeeds.Get(ctx, key)
}

// GetAllPriceFeeds returns the set of price feeds that are currently in the
// x/sla module's state for a given SLA.
func (k *Keeper) GetAllPriceFeeds(ctx sdk.Context, slaID string) ([]slatypes.PriceFeed, error) {
	feeds := make([]slatypes.PriceFeed, 0)
	cb := func(feed slatypes.PriceFeed) error {
		feeds = append(feeds, feed)
		return nil
	}

	if err := k.iteratePriceFeeds(ctx, slaID, cb); err != nil {
		return nil, err
	}

	return feeds, nil
}

// RemovePriceFeed removes a price feed from the x/sla module's state. Note,
// if the price feed does not exist, this function will not return an error.
func (k *Keeper) RemovePriceFeed(
	ctx sdk.Context,
	slaID string,
	cp slinkytypes.CurrencyPair,
	consAddress sdk.ConsAddress,
) error {
	key := collections.Join3(slaID, cp.String(), consAddress.Bytes())
	return k.priceFeeds.Remove(ctx, key)
}

// RemovePriceFeedByCurrencyPair removes all price feeds that track
// a given currency pair from the x/sla module's state for a given sla.
func (k *Keeper) RemovePriceFeedByCurrencyPair(
	ctx sdk.Context,
	slaID string,
	cp slinkytypes.CurrencyPair,
) error {
	prefix := collections.NewSuperPrefixedTripleRange[string, string, []byte](slaID, cp.String())
	return k.priceFeeds.Clear(ctx, prefix)
}

// RemovePriceFeedsBySLA removes all price feeds that track a given SLA
// from the x/sla module's state.
func (k *Keeper) RemovePriceFeedsBySLA(ctx sdk.Context, slaID string) error {
	prefix := collections.NewPrefixedTripleRange[string, string, []byte](slaID)
	return k.priceFeeds.Clear(ctx, prefix)
}

// ContainsPriceFeed returns true if the x/sla module's state contains
// a price feed with the given sla ID, currency pair, and validator.
func (k *Keeper) ContainsPriceFeed(
	ctx sdk.Context,
	slaID string,
	cp slinkytypes.CurrencyPair,
	validator sdk.ConsAddress,
) (bool, error) {
	key := collections.Join3(slaID, cp.String(), validator.Bytes())
	return k.priceFeeds.Has(ctx, key)
}

// iteratePriceFeeds iterates over the set of price feeds that
// are currently in the x/sla module's state and belong to a given SLA.
func (k *Keeper) iteratePriceFeeds(ctx sdk.Context, slaID string, cb PriceFeedCB) error {
	prefix := collections.NewPrefixedTripleRange[string, string, []byte](slaID)
	iterator, err := k.priceFeeds.Iterate(ctx, prefix)
	if err != nil {
		return err
	}
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		priceFeed, err := iterator.Value()
		if err != nil {
			return err
		}

		if err := cb(priceFeed); err != nil {
			return err
		}
	}

	return nil
}
