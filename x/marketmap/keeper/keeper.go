package keeper

import (
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

	// collections
	// marketConfigs is keyed by provider name and provides the MarketConfig for each given provider
	marketConfigs collections.Map[types.MarketProvider, types.MarketConfig]
	// aggregationConfigs is keyed by CurrencyPair string (BASE/QUOTE) and contains the PathsConfig used
	// to do price aggregation for a given canonical Ticker
	aggregationConfigs collections.Map[types.TickerString, types.PathsConfig]
}

// NewKeeper initializes the keeper and its backing stores.
func NewKeeper(ss store.KVStoreService, cdc codec.BinaryCodec, authority sdk.AccAddress) Keeper {
	sb := collections.NewSchemaBuilder(ss)

	return Keeper{
		cdc:                cdc,
		authority:          authority,
		marketConfigs:      collections.NewMap(sb, types.MarketConfigsPrefix, "market_configs", types.MarketProviderCodec, codec.CollValue[types.MarketConfig](cdc)),
		aggregationConfigs: collections.NewMap(sb, types.AggregationConfigsPrefix, "aggregation_configs", types.TickerStringCodec, codec.CollValue[types.PathsConfig](cdc)),
	}
}

// GetAllMarketConfigs returns the set of MarketConfig objects currently stored in state.
func (k Keeper) GetAllMarketConfigs(ctx sdk.Context) ([]types.MarketConfig, error) {
	iter, err := k.marketConfigs.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	configs, err := iter.Values()
	if err != nil {
		return nil, err
	}
	return configs, err
}

// GetAllAggregationConfigs returns all PathsConfig objects currently in x/marketmap state.
// The keys are omitted since the PathsConfig object contains a Ticker which effectively identifies
// which pair each config refers to.
func (k Keeper) GetAllAggregationConfigs(ctx sdk.Context) ([]types.PathsConfig, error) {
	iter, err := k.aggregationConfigs.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}
	configs, err := iter.Values()
	if err != nil {
		return nil, err
	}
	return configs, err
}

// GetMarketMap returns an AggregateMarketConfig object which effectively contains the entire state of the module.
func (k Keeper) GetMarketMap(ctx sdk.Context) (*types.AggregateMarketConfig, error) {
	marketMap := &types.AggregateMarketConfig{
		MarketConfigs: make(map[string]types.MarketConfig),
		TickerConfigs: make(map[string]types.PathsConfig),
	}
	aggregationCfgs, err := k.GetAllAggregationConfigs(ctx)
	if err != nil {
		return nil, err
	}
	for _, pathConfig := range aggregationCfgs {
		marketMap.TickerConfigs[slinkytypes.CurrencyPair{Base: pathConfig.Ticker.Base, Quote: pathConfig.Ticker.Quote}.String()] = pathConfig
	}
	marketConfigs, err := k.GetAllMarketConfigs(ctx)
	if err != nil {
		return nil, err
	}
	for _, marketCfg := range marketConfigs {
		marketMap.MarketConfigs[marketCfg.Name] = marketCfg
	}
	return marketMap, nil
}

// CreateAggregationConfig initializes a new PathsConfig.
// The combination of pathsConfig.Ticker.Base and pathsConfig.Ticker.Quote provide a unique key which is used to
// validate against duplication.
func (k Keeper) CreateAggregationConfig(ctx sdk.Context, pathsConfig types.PathsConfig) error {
	// Construct the key for the PathsConfig
	configKey := types.TickerString(slinkytypes.CurrencyPair{Base: pathsConfig.Ticker.Base, Quote: pathsConfig.Ticker.Quote}.String())
	// Check if AggregationConfig already exists for the Ticker
	alreadyExists, err := k.aggregationConfigs.Has(ctx, configKey)
	if err != nil {
		return err
	}
	if alreadyExists {
		return types.NewAggregationConfigAlreadyExistsError(configKey)
	}
	// Create the config
	return k.aggregationConfigs.Set(ctx, configKey, pathsConfig)
}

// CreateMarketConfig initializes a new MarketConfig.
// The marketConfig.Name corresponds to a price provider, and must be unique.
func (k Keeper) CreateMarketConfig(ctx sdk.Context, marketConfig types.MarketConfig) error {
	// Check if MarketConfig already exists for the provider
	alreadyExists, err := k.marketConfigs.Has(ctx, types.MarketProvider(marketConfig.Name))
	if err != nil {
		return err
	}
	if alreadyExists {
		return types.NewMarketConfigAlreadyExistsError(types.MarketProvider(marketConfig.Name))
	}
	// Create the config
	return k.marketConfigs.Set(ctx, types.MarketProvider(marketConfig.Name), marketConfig)
}
