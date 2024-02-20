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

	// tickers is keyed by CurrencyPair string (BASE/QUOTE) and contains
	// the list of all Tickers.
	tickers collections.Map[types.TickerString, types.Ticker]

	// paths is keyed by CurrencyPair string (BASE/QUOTE) and contains
	// the list of all Paths.
	paths collections.Map[types.TickerString, types.Paths]

	// providers is keyed by CurrencyPair string (BASE/QUOTE) and contains
	// the list of all Providers.
	providers collections.Map[types.TickerString, types.Providers]

	// lastUpdated is the last block height the marketmap was updated.
	lastUpdated collections.Item[int64]

	// params is the module's parameters.
	params collections.Item[types.Params]
}

// NewKeeper initializes the keeper and its backing stores.
func NewKeeper(ss store.KVStoreService, cdc codec.BinaryCodec, authority sdk.AccAddress) Keeper {
	sb := collections.NewSchemaBuilder(ss)

	// Create the collections item that will track the module parameters.
	params := collections.NewItem(
		sb,
		types.ParamsPrefix,
		"params",
		codec.CollValue[types.Params](cdc),
	)

	return Keeper{
		cdc:         cdc,
		authority:   authority,
		tickers:     collections.NewMap(sb, types.TickersPrefix, "tickers", types.TickersCodec, codec.CollValue[types.Ticker](cdc)),
		paths:       collections.NewMap(sb, types.PathsPrefix, "paths", types.TickersCodec, codec.CollValue[types.Paths](cdc)),
		providers:   collections.NewMap(sb, types.ProvidersPrefix, "providers", types.TickersCodec, codec.CollValue[types.Providers](cdc)),
		lastUpdated: collections.NewItem[int64](sb, types.LastUpdatedPrefix, "last_updated", types.LastUpdatedCodec),
		params:      params,
	}
}

// SetLastUpdated sets the lastUpdated field to the current block height.
func (k *Keeper) SetLastUpdated(ctx sdk.Context) error {
	return k.lastUpdated.Set(ctx, ctx.BlockHeight())
}

// GetLastUpdated gets the last block-height the market map was updated.
func (k *Keeper) GetLastUpdated(ctx sdk.Context) (int64, error) {
	return k.lastUpdated.Get(ctx)
}

// GetAllTickers returns the set of Ticker objects currently stored in state.
func (k *Keeper) GetAllTickers(ctx sdk.Context) ([]types.Ticker, error) {
	iter, err := k.tickers.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	tickers, err := iter.Values()
	if err != nil {
		return nil, err
	}
	return tickers, err
}

// GetAllTickersMap returns the set of Ticker objects currently stored in state
// as a map[TickerString] -> Tickers.
func (k *Keeper) GetAllTickersMap(ctx sdk.Context) (map[string]types.Ticker, error) {
	iter, err := k.tickers.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	keyValues, err := iter.KeyValues()
	if err != nil {
		return nil, err
	}

	m := make(map[string]types.Ticker, len(keyValues))
	for _, keyValue := range keyValues {
		m[string(keyValue.Key)] = keyValue.Value
	}

	return m, nil
}

// CreateTicker initializes a new Ticker.
// The Ticker.String corresponds to a market, and must be unique.
func (k *Keeper) CreateTicker(ctx sdk.Context, ticker types.Ticker) error {
	// Check if Ticker already exists for the provider
	alreadyExists, err := k.tickers.Has(ctx, types.TickerString(ticker.String()))
	if err != nil {
		return err
	}
	if alreadyExists {
		return types.NewTickerAlreadyExistsError(types.TickerString(ticker.String()))
	}
	// Create the config
	return k.tickers.Set(ctx, types.TickerString(ticker.String()), ticker)
}

// GetAllProvidersMap returns the set of Providers objects currently stored in state
// as a map[TickerString] -> Providers.
func (k *Keeper) GetAllProvidersMap(ctx sdk.Context) (map[string]types.Providers, error) {
	iter, err := k.providers.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	keyValues, err := iter.KeyValues()
	if err != nil {
		return nil, err
	}

	m := make(map[string]types.Providers, len(keyValues))
	for _, keyValue := range keyValues {
		m[string(keyValue.Key)] = keyValue.Value
	}

	return m, nil
}

// CreateProviders initializes a new providers.
// The Ticker.String corresponds to a market, and must be unique.
func (k *Keeper) CreateProviders(ctx sdk.Context, providers types.Providers, ticker types.Ticker) error {
	// Check if MarketConfig already exists for the provider
	alreadyExists, err := k.providers.Has(ctx, types.TickerString(ticker.String()))
	if err != nil {
		return err
	}
	if alreadyExists {
		return types.NewTickerAlreadyExistsError(types.TickerString(ticker.String()))
	}
	// Create the config
	return k.providers.Set(ctx, types.TickerString(ticker.String()), providers)
}

// GetAllPathsMap returns the set of Paths objects currently stored in state
// as a map[TickerString] -> Paths.
func (k *Keeper) GetAllPathsMap(ctx sdk.Context) (map[string]types.Paths, error) {
	iter, err := k.paths.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	keyValues, err := iter.KeyValues()
	if err != nil {
		return nil, err
	}

	m := make(map[string]types.Paths, len(keyValues))
	for _, keyValue := range keyValues {
		m[string(keyValue.Key)] = keyValue.Value
	}

	return m, nil
}

// CreatePaths initializes a new Paths.
// The Ticker.String corresponds to a market, and must be unique.
func (k *Keeper) CreatePaths(ctx sdk.Context, paths types.Paths, ticker types.Ticker) error {
	// Check if MarketConfig already exists for the provider
	alreadyExists, err := k.paths.Has(ctx, types.TickerString(ticker.String()))
	if err != nil {
		return err
	}
	if alreadyExists {
		return types.NewTickerAlreadyExistsError(types.TickerString(ticker.String()))
	}
	// Create the config
	return k.paths.Set(ctx, types.TickerString(ticker.String()), paths)
}

// CreateMarket sets the ticker, paths, and providers for a given market.  It also
// sets the LastUpdated field to the current block height.
func (k *Keeper) CreateMarket(ctx sdk.Context, ticker types.Ticker, paths types.Paths, providers types.Providers) error {
	if err := k.CreateTicker(ctx, ticker); err != nil {
		return err
	}

	if err := k.CreatePaths(ctx, paths, ticker); err != nil {
		return err
	}

	if err := k.CreateProviders(ctx, providers, ticker); err != nil {
		return err
	}

	return k.SetLastUpdated(ctx)
}

// SetParams sets the x/marketmap module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	return k.params.Set(ctx, params)
}

// GetParams returns the x/marketmap module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	return k.params.Get(ctx)
}
