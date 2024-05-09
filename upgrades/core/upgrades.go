package core

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/skip-mev/slinky/upgrades"
	marketmapkeeper "github.com/skip-mev/slinky/x/marketmap/keeper"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
	oraclekeeper "github.com/skip-mev/slinky/x/oracle/keeper"

	upgradetypes "cosmossdk.io/x/upgrade/types"
)

const Name = "initialize slinky state"

var _ upgrades.Upgrade = &InitializeUpgrade{}

type InitializeUpgrade struct {
	params  marketmaptypes.Params
	markets marketmaptypes.Markets
}

func NewInitializeUpgrade(params marketmaptypes.Params, markets marketmaptypes.Markets) *InitializeUpgrade {
	return &InitializeUpgrade{
		params:  params,
		markets: markets,
	}
}

func (i *InitializeUpgrade) CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ *oraclekeeper.Keeper,
	marketMapKeeper *marketmapkeeper.Keeper,
	_ codec.Codec,
	handler upgradetypes.UpgradeHandler,
) upgradetypes.UpgradeHandler {
	if oraclekeeper == nil

	return func(c context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx := sdk.UnwrapSDKContext(c)

		ctx.Logger().Info("Starting module migrations...")
		vm, err := mm.RunMigrations(ctx, configurator, vm)
		if err != nil {
			return vm, err
		}

		ctx.Logger().Info("Setting marketmap params...")
		err = setMarketMapParams(ctx, marketMapKeeper, marketmaptypes.DefaultParams())
		if err != nil {
			return nil, err
		}

		ctx.Logger().Info("Setting marketmap and oracle state...")
		err = setMarketState(ctx, marketMapKeeper)
		if err != nil {
			return nil, err
		}

		ctx.Logger().Info(fmt.Sprintf("Migration {%s} applied", Name))
		return handler(c, plan, vm)
	}
}

func setMarketMapParams(ctx sdk.Context, marketmapKeeper *marketmapkeeper.Keeper, params marketmaptypes.Params) error {
	if err := params.ValidateBasic(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	return marketmapKeeper.SetParams(ctx, params)
}

func setMarketState(ctx sdk.Context, mmKeeper *marketmapkeeper.Keeper, markets marketmaptypes.Markets) error {
	// markets, err := marketmaptypes.ReadMarketsFromFile("markets.json")
	//if err != nil {
	//	return err
	//}

	for _, market := range markets {
		err = mmKeeper.CreateMarket(ctx, market)
		if err != nil {
			return err
		}

		err = mmKeeper.Hooks().AfterMarketCreated(ctx, market)
		if err != nil {
			return err
		}

	}
	return nil
}
