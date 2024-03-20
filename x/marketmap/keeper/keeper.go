package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/types"
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
}

// NewKeeper initializes the keeper and its backing stores.
func NewKeeper(ss store.KVStoreService, cdc codec.BinaryCodec, authority sdk.AccAddress) *Keeper {
	sb := collections.NewSchemaBuilder(ss)

	// Create the collections item that will track the module parameters.
	params := collections.NewItem[types.Params](
		sb,
		types.ParamsPrefix,
		"params",
		codec.CollValue[types.Params](cdc),
	)

	return &Keeper{
		cdc:         cdc,
		authority:   authority,
		markets:     collections.NewMap(sb, types.MarketsPrefix, "markets", types.TickersCodec, codec.CollValue[types.Market](cdc)),
		lastUpdated: collections.NewItem[uint64](sb, types.LastUpdatedPrefix, "last_updated", types.LastUpdatedCodec),
		params:      params,
	}
}

// SetLastUpdated sets the lastUpdated field to the current block height.
func (k *Keeper) SetLastUpdated(ctx sdk.Context, height uint64) error {
	return k.lastUpdated.Set(ctx, height)
}

// GetLastUpdated gets the last block-height the market map was updated.
func (k *Keeper) GetLastUpdated(ctx sdk.Context) (uint64, error) {
	return k.lastUpdated.Get(ctx)
}

// GetMarket returns a market from the store by its currency pair string ID.
func (k *Keeper) GetMarket(ctx sdk.Context, tickerStr string) (types.Market, error) {
	return k.markets.Get(ctx, types.TickerString(tickerStr))
}

// GetAllMarkets returns the set of Market objects currently stored in state.
func (k *Keeper) GetAllMarkets(ctx sdk.Context) ([]types.Market, error) {
	iter, err := k.markets.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	markets, err := iter.Values()
	if err != nil {
		return nil, err
	}
	return markets, err
}

// GetAllMarketsMap returns the set of Market objects currently stored in state
// as a map[TickerString] -> Markets.
func (k *Keeper) GetAllMarketsMap(ctx sdk.Context) (map[string]types.Market, error) {
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

// createMarket initializes a new Market.
// The Ticker.String corresponds to a market, and must be unique.
func (k *Keeper) createMarket(ctx sdk.Context, market types.Market) error {
	// Check if Ticker already exists for the provider
	alreadyExists, err := k.markets.Has(ctx, types.TickerString(market.Ticker.String()))
	if err != nil {
		return err
	}
	if alreadyExists {
		return types.NewMarketAlreadyExistsError(types.TickerString(market.Ticker.String()))
	}
	// Create the config
	return k.markets.Set(ctx, types.TickerString(market.Ticker.String()), market)
}

// updateMarket updates a Market.
// The Ticker.String corresponds to a market, and exist unique.
func (k *Keeper) updateMarket(ctx sdk.Context, market types.Market) error {
	// Check if Ticker already exists for the provider
	alreadyExists, err := k.markets.Has(ctx, types.TickerString(market.Ticker.String()))
	if err != nil {
		return err
	}
	if !alreadyExists {
		return types.NewMarketDoesNotExistsError(types.TickerString(market.Ticker.String()))
	}
	// Create the config
	return k.markets.Set(ctx, types.TickerString(market.Ticker.String()), market)
}

// CreateMarket sets the ticker, paths, and providers for a given market.  It also
// sets the LastUpdated field to the current block height.
func (k *Keeper) CreateMarket(ctx sdk.Context, market types.Market) error {
	if err := k.createMarket(ctx, market); err != nil {
		return err
	}

	return k.SetLastUpdated(ctx, uint64(ctx.BlockHeight()))
}

// UpdateMarket updaters the ticker, paths, and providers for a given market.  It also
// sets the LastUpdated field to the current block height.
func (k *Keeper) UpdateMarket(ctx sdk.Context, market types.Market) error {
	if err := k.updateMarket(ctx, market); err != nil {
		return err
	}

	return k.SetLastUpdated(ctx, uint64(ctx.BlockHeight()))
}

// SetParams sets the x/marketmap module's parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	return k.params.Set(ctx, params)
}

// GetParams returns the x/marketmap module's parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	return k.params.Get(ctx)
}

// ValidateState is called after keeper modifications have been made to the market map to verify that
// the aggregate of all updates has led to a valid state.
func (k *Keeper) ValidateState(ctx sdk.Context, creates []types.Market, updates []types.Market) error {
	for _, market := range append(creates, updates...) {
		// check that all paths already exist in the keeper store:
		for _, path := range market.Paths.Paths {
			for _, op := range path.Operations {
				cp := op.CurrencyPair
				if op.Invert {
					cp = slinkytypes.CurrencyPair{
						Base:  cp.Quote,
						Quote: cp.Base,
					}
				}

				has, err := k.markets.Has(ctx, types.TickerString(cp.String()))
				if err != nil {
					return err
				}

				if !has {
					return fmt.Errorf("currency pair %s in path %s does not exist", cp.String(), path.ShowRoute())
				}
			}
		}
	}

	return nil
}
