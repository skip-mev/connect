package keeper

import (
	"context"
	"errors"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/marketmap/types"
)

// Keeper is the module's keeper implementation.
type Keeper struct {
	cdc codec.BinaryCodec

	// module authority
	authority sdk.AccAddress

	// registered hooks
	hooks types.MarketMapHooks

	// markets is keyed by CurrencyPair string (BASE/QUOTE) and contains
	// the list of all Markets.
	markets collections.Map[types.TickerString, types.Market]

	// lastUpdated is the last block height the marketmap was updated.
	lastUpdated collections.Item[uint64]

	// params is the module's parameters.
	params collections.Item[types.Params]

	// deleteValidationHooks are called by the keeper before any deletion call is performed.
	deleteMarketValidationHooks types.MarketValidationHooks
}

// NewKeeper initializes the keeper and its backing stores.
func NewKeeper(ss store.KVStoreService, cdc codec.BinaryCodec, authority sdk.AccAddress, opts ...Option) *Keeper {
	sb := collections.NewSchemaBuilder(ss)

	// Create the collections item that will track the module parameters.
	params := collections.NewItem[types.Params](
		sb,
		types.ParamsPrefix,
		"params",
		codec.CollValue[types.Params](cdc),
	)

	k := &Keeper{
		cdc:                         cdc,
		authority:                   authority,
		markets:                     collections.NewMap(sb, types.MarketsPrefix, "markets", types.TickersCodec, codec.CollValue[types.Market](cdc)),
		lastUpdated:                 collections.NewItem[uint64](sb, types.LastUpdatedPrefix, "last_updated", types.LastUpdatedCodec),
		params:                      params,
		hooks:                       &types.NoopMarketMapHooks{},
		deleteMarketValidationHooks: types.DefaultDeleteMarketValidationHooks(),
	}

	// apply options to default initialized keeper
	for _, opt := range opts {
		opt(k)
	}

	return k
}

// SetDeleteMarketValidationHooks sets the MarketValidationHooks for deletion in the keeper.
func (k *Keeper) SetDeleteMarketValidationHooks(hooks types.MarketValidationHooks) {
	k.deleteMarketValidationHooks = hooks
}

// SetLastUpdated sets the lastUpdated field to the current block height.
func (k *Keeper) SetLastUpdated(ctx context.Context, height uint64) error {
	return k.lastUpdated.Set(ctx, height)
}

// GetLastUpdated gets the last block-height the market map was updated.
func (k *Keeper) GetLastUpdated(ctx context.Context) (uint64, error) {
	return k.lastUpdated.Get(ctx)
}

// GetMarket returns a market from the store by its currency pair string ID.
func (k *Keeper) GetMarket(ctx context.Context, tickerStr string) (types.Market, error) {
	return k.markets.Get(ctx, types.TickerString(tickerStr))
}

// setMarket sets a market.
func (k *Keeper) setMarket(ctx context.Context, market types.Market) error {
	return k.markets.Set(ctx, types.TickerString(market.Ticker.String()), market)
}

// EnableMarket sets the Enabled field of a Market Ticker to true.
func (k *Keeper) EnableMarket(ctx context.Context, tickerStr string) error {
	market, err := k.GetMarket(ctx, tickerStr)
	if err != nil {
		return err
	}

	market.Ticker.Enabled = true

	return k.setMarket(ctx, market)
}

// DisableMarket sets the Enabled field of a Market Ticker to false.
func (k *Keeper) DisableMarket(ctx context.Context, tickerStr string) error {
	market, err := k.GetMarket(ctx, tickerStr)
	if err != nil {
		return err
	}

	market.Ticker.Enabled = false

	return k.setMarket(ctx, market)
}

// GetAllMarkets returns the set of Market objects currently stored in state
// as a map[TickerString] -> Markets.
func (k *Keeper) GetAllMarkets(ctx context.Context) (map[string]types.Market, error) {
	iter, err := k.markets.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	keyValues, err := iter.KeyValues()
	if err != nil {
		return nil, err
	}
	m := make(map[string]types.Market, len(keyValues))
	for _, keyValue := range keyValues {
		m[string(keyValue.Key)] = keyValue.Value
	}

	return m, nil
}

// CreateMarket initializes a new Market.
// The Ticker.String corresponds to a market, and must be unique.
func (k *Keeper) CreateMarket(ctx context.Context, market types.Market) error {
	// Check if Ticker already exists for the provider
	alreadyExists, err := k.markets.Has(ctx, types.TickerString(market.Ticker.String()))
	if err != nil {
		return err
	}
	if alreadyExists {
		return types.NewMarketAlreadyExistsError(types.TickerString(market.Ticker.String()))
	}
	// Create the config
	return k.setMarket(ctx, market)
}

// UpdateMarket updates a Market.
// The Ticker.String corresponds to a market, and exists uniquely.
func (k *Keeper) UpdateMarket(ctx context.Context, market types.Market) error {
	// Check if Ticker already exists for the provider
	alreadyExists, err := k.markets.Has(ctx, types.TickerString(market.Ticker.String()))
	if err != nil {
		return err
	}
	if !alreadyExists {
		return types.NewMarketDoesNotExistsError(types.TickerString(market.Ticker.String()))
	}
	// Create the config
	return k.setMarket(ctx, market)
}

// DeleteMarket removes a Market.  If the market does not exist, this is a no-op and nil is returned.
// If the market exists, all DeleteMarketValidationHooks are called on the market before deletion.
// Additionally, returns true if the market was deleted.
func (k *Keeper) DeleteMarket(ctx context.Context, tickerStr string) (bool, error) {
	market, err := k.GetMarket(ctx, tickerStr)
	switch {
	case errors.Is(err, collections.ErrNotFound):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("failed to get market for ticker %s: %w", tickerStr, err)
	}

	if err := k.deleteMarketValidationHooks.ValidateMarket(ctx, market); err != nil {
		return false, err
	}

	err = k.markets.Remove(ctx, types.TickerString(market.Ticker.String()))
	if err != nil {
		return false, err
	}

	return true, nil
}

// HasMarket checks if a market exists in the store.
func (k *Keeper) HasMarket(ctx context.Context, tickerStr string) (bool, error) {
	return k.markets.Has(ctx, types.TickerString(tickerStr))
}

// SetParams sets the x/marketmap module's parameters.
func (k *Keeper) SetParams(ctx context.Context, params types.Params) error {
	return k.params.Set(ctx, params)
}

// GetParams returns the x/marketmap module's parameters.
func (k *Keeper) GetParams(ctx context.Context) (types.Params, error) {
	return k.params.Get(ctx)
}

// ValidateState is called after keeper modifications have been made to the market map to verify that
// the aggregate of all updates has led to a valid state.
func (k *Keeper) ValidateState(ctx sdk.Context, updates []types.Market) error {
	for _, market := range updates {
		if err := k.IsMarketValid(ctx, market); err != nil {
			return err
		}
	}

	return nil
}

// IsMarketValid checks if a market is valid by statefully checking if each of the currency pairs
// specified by its provider configs are valid and in state.
func (k *Keeper) IsMarketValid(ctx sdk.Context, market types.Market) error {
	// check that all markets already exist in the keeper store:
	for _, providerConfig := range market.ProviderConfigs {
		if providerConfig.NormalizeByPair != nil {
			norm, err := k.markets.Get(ctx, types.TickerString(providerConfig.NormalizeByPair.String()))
			if err != nil {
				return fmt.Errorf("unable to get normalize market %s for market %s: %w",
					providerConfig.NormalizeByPair.String(), market.Ticker.String(), err)
			}

			// if the new market is enabled, its normalize by market must also be enabled
			if market.Ticker.Enabled && !norm.Ticker.Enabled {
				return fmt.Errorf("needed normalize market %s for market %s is not enabled",
					providerConfig.NormalizeByPair.String(), market.Ticker.String())
			}
		}
	}

	return nil
}
