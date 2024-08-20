package keeper

import (
	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

// Keeper defines a new keeper for the price feed SLA module. This module
// tracks the current SLAs and the corresponding price feed updates. Each
// price feed is associated with an SLA, validator, and currency pair. The
// currency pairs utilized by the x/sla module are defined in the x/oracle
// module.
type Keeper struct {
	cdc codec.BinaryCodec

	// State management variables
	storeService store.KVStoreService
	schema       collections.Schema

	// slas is a map of (sla ID -> SLA) that is used to track the SLAs that are
	// currently in the x/sla module's state.
	slas collections.Map[string, slatypes.PriceFeedSLA]

	// priceFeeds is a map of (sla ID, currency pair, consensus address -> price feed)
	priceFeeds collections.Map[collections.Triple[string, string, []byte], slatypes.PriceFeed]

	// currencyPairs is a map of (currency pair string -> currency pair) that is used to
	// track the currency pairs that are currently in the x/sla module's state. This set
	// of currency pairs is used to remove stale price feeds.
	currencyPairs collections.Map[string, slinkytypes.CurrencyPair]

	// params is the module's parameters.
	params collections.Item[slatypes.Params]

	// authority is the address that is authorized to add new SLAs.
	authority sdk.AccAddress

	// stakingKeeper is utilized to retrieve validator information.
	stakingKeeper slatypes.StakingKeeper

	// slashingKeeper is utilized to slash validators that do not meet the SLA.
	slashingKeeper slatypes.SlashingKeeper
}

// NewKeeper returns a new keeper for the price feed SLAs. The keeper is
// responsible for maintaining the current set of SLAs and the corresponding
// price feed updates.
func NewKeeper(
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	authority sdk.AccAddress,
	stakingKeeper slatypes.StakingKeeper,
	slashingKeeper slatypes.SlashingKeeper,
) *Keeper {
	schemaBuilder := collections.NewSchemaBuilder(storeService)

	// Create the collections map that will track the SLAs.
	slas := collections.NewMap(
		schemaBuilder,
		slatypes.KeyPrefixSLA,
		"slas",
		collections.StringKey,
		codec.CollValue[slatypes.PriceFeedSLA](cdc),
	)

	// Create the price feed map that will track the price feed updates.
	priceFeeds := collections.NewMap(
		schemaBuilder,
		slatypes.KeyPrefixPriceFeeds,
		"price_feeds",
		collections.TripleKeyCodec(collections.StringKey, collections.StringKey, collections.BytesKey),
		codec.CollValue[slatypes.PriceFeed](cdc),
	)

	// Create the collections map that will track the currency pairs.
	currencyPairs := collections.NewMap(
		schemaBuilder,
		slatypes.KeyPrefixCurrencyPairs,
		"currency_pairs",
		collections.StringKey,
		codec.CollValue[slinkytypes.CurrencyPair](cdc),
	)

	// Create the collections item that will track the module parameters.
	params := collections.NewItem(
		schemaBuilder,
		slatypes.KeyPrefixParams,
		"params",
		codec.CollValue[slatypes.Params](cdc),
	)

	// Build the schema and return the keeper.
	schema, err := schemaBuilder.Build()
	if err != nil {
		panic(err)
	}

	return &Keeper{
		cdc:            cdc,
		storeService:   storeService,
		schema:         schema,
		slas:           slas,
		priceFeeds:     priceFeeds,
		currencyPairs:  currencyPairs,
		params:         params,
		authority:      authority,
		stakingKeeper:  stakingKeeper,
		slashingKeeper: slashingKeeper,
	}
}

// Logger returns the keeper's logger.
func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/sla")
}

// SetParams sets the x/sla module's parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params slatypes.Params) error {
	return k.params.Set(ctx, params)
}

// GetParams returns the x/sla module's parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (slatypes.Params, error) {
	return k.params.Get(ctx)
}

// SetCurrencyPairs sets the x/sla module's currency pairs. Note, this function
// is primarily used to remove stale price feeds.
func (k *Keeper) SetCurrencyPairs(ctx sdk.Context, currencyPairs map[slinkytypes.CurrencyPair]struct{}) error {
	// Remove all currency pairs that are currently in the x/sla module's state.
	if err := k.currencyPairs.Clear(ctx, nil); err != nil {
		return err
	}

	for cp := range currencyPairs {
		if err := k.currencyPairs.Set(ctx, cp.String(), cp); err != nil {
			return err
		}
	}

	return nil
}

// GetCurrencyPairs returns the x/sla module's currency pairs.
func (k *Keeper) GetCurrencyPairs(ctx sdk.Context) (map[slinkytypes.CurrencyPair]struct{}, error) {
	iterator, err := k.currencyPairs.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	currencyPairs := make(map[slinkytypes.CurrencyPair]struct{})
	for ; iterator.Valid(); iterator.Next() {
		cp, err := iterator.Value()
		if err != nil {
			return nil, err
		}

		currencyPairs[cp] = struct{}{}
	}

	return currencyPairs, nil
}
