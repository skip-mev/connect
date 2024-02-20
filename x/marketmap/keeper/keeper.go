package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

// Keeper is the module's keeper implementation.
type Keeper struct {
	cdc codec.BinaryCodec

	// module authority
	authority sdk.AccAddress

	// markets is keyed by CurrencyPair string (BASE/QUOTE) and contains
	// the list of all Tickers.
	markets collections.Map[types.TickerString, types.Ticker]

	// lastUpdated is the last block height the marketmap was updated.
	lastUpdated collections.Item[int64]
}

// NewKeeper initializes the keeper and its backing stores.
func NewKeeper(ss store.KVStoreService, cdc codec.BinaryCodec, authority sdk.AccAddress) Keeper {
	sb := collections.NewSchemaBuilder(ss)

	return Keeper{
		cdc:         cdc,
		authority:   authority,
		markets:     collections.NewMap(sb, types.TickersPrefix, "markets", types.TickersCodec, codec.CollValue[types.Ticker](cdc)),
		lastUpdated: collections.NewItem[int64](sb, types.LastUpdatedPrefix, "last_updated", types.LastUpdatedCodec),
	}
}

// SetLastUpdated sets the lastUpdated field to the current block height.
func (k Keeper) SetLastUpdated(ctx sdk.Context) error {
	return k.lastUpdated.Set(ctx, ctx.BlockHeight())
}

// GetLastUpdated gets the last block-height the market map was updated.
func (k Keeper) GetLastUpdated(ctx sdk.Context) (int64, error) {
	return k.lastUpdated.Get(ctx)
}

// GetAllTickers returns the set of Ticker objects currently stored in state.
func (k Keeper) GetAllTickers(ctx sdk.Context) ([]types.Ticker, error) {
	iter, err := k.markets.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	tickers, err := iter.Values()
	if err != nil {
		return nil, err
	}
	return tickers, err
}

// CreateTicker initializes a new Ticker.
// The Ticker.String corresponds to a market, and must be unique.
func (k Keeper) CreateTicker(ctx sdk.Context, ticker types.Ticker) error {
	// Check if MarketConfig already exists for the provider
	alreadyExists, err := k.markets.Has(ctx, types.TickerString(ticker.String()))
	if err != nil {
		return err
	}
	if alreadyExists {
		return types.NewTickerAlreadyExistsError(types.TickerString(ticker.String()))
	}
	// Create the config
	err = k.markets.Set(ctx, types.TickerString(ticker.String()), ticker)
	if err != nil {
		return err
	}

	return k.SetLastUpdated(ctx)
}
