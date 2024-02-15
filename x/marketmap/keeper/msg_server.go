package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/marketmap/types"
)

// msgServer is the default implementation of the x/marketmap MsgService.
type msgServer struct {
	k Keeper
}

// NewMsgServer returns the default implementation of the x/marketmap message service.
func NewMsgServer(k Keeper) types.MsgServer {
	return &msgServer{k}
}

var _ types.MsgServer = (*msgServer)(nil)

// CreateMarket creates a market from the given message.
func (ms msgServer) CreateMarket(goCtx context.Context, msg *types.MsgCreateMarket) (*types.MsgCreateMarketResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("unable to process nil msg")
	}

	// Update the module's parameters.
	ctx := sdk.UnwrapSDKContext(goCtx)

	params, err := ms.k.GetParams(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get marketmap params: %w", err)
	}

	if msg.Signer != params.MarketAuthority {
		return nil, fmt.Errorf("request signer %s does not match module market authority %s", msg.Signer, params.MarketAuthority)
	}

	// check if market already exists
	aggCfgs, err := ms.k.GetAllAggregationConfigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get aggregation configs")
	}

	for _, cfg := range aggCfgs {
		if msg.Ticker.CurrencyPair == cfg.Ticker.CurrencyPair {
			return nil, fmt.Errorf("ticker %s already exists in marketmap", msg.Ticker.CurrencyPair.String())
		}
	}

	marketConfigs, err := ms.k.GetAllMarketConfigs(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get market configs")
	}

	// set market
	for providerName, offChainTicker := range msg.ProvidersToOffChainTickers {
		marketConfig, found := marketConfigs[providerName]
		if !found {
			// if not found, add new provider
			marketConfig = types.MarketConfig{
				Name: providerName,
			}
		}

		marketConfig.TickerConfigs[msg.Ticker.CurrencyPair.String()] = types.TickerConfig{
			Ticker:         msg.Ticker,
			OffChainTicker: offChainTicker,
		}

	}

	return &types.MsgCreateMarketResponse{}, nil
}

// Params updates the x/marketmap module's Params.
func (ms msgServer) Params(goCtx context.Context, msg *types.MsgParams) (*types.MsgParamsResponse, error) {
	if msg == nil {
		return nil, fmt.Errorf("unable to process nil msg")
	}

	// Update the module's parameters.
	ctx := sdk.UnwrapSDKContext(goCtx)

	if msg.Authority != ms.k.authority.String() {
		return nil, fmt.Errorf("request authority %s does not match module keeper authority %s", msg.Authority, ms.k.authority.String())
	}

	if err := ms.k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgParamsResponse{}, nil
}
