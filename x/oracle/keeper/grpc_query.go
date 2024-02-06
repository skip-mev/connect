package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/oracle/types"
)

// queryServer is the default implementation of the x/oracle QueryService
type queryServer struct {
	k Keeper
}

// NewQueryServer returns an implementation of the x/oracle QueryServer
func NewQueryServer(k Keeper) types.QueryServer {
	return queryServer{
		k,
	}
}

var _ types.QueryServer = queryServer{}

// GetAllCurrencyPairs returns the set of all currency pairs that the module is tracking QuotePrices for.
// It returns an error to the caller if there are no CurrencyPairs being tracked by the module
func (q queryServer) GetAllCurrencyPairs(ctx context.Context, _ *types.GetAllCurrencyPairsRequest) (*types.GetAllCurrencyPairsResponse, error) {
	// get all currency pairs from state
	cps := q.k.GetAllCurrencyPairs(sdk.UnwrapSDKContext(ctx))

	// if no currency pairs exist in the module, return an error to indicate to caller
	if len(cps) == 0 {
		return &types.GetAllCurrencyPairsResponse{}, nil
	}

	return &types.GetAllCurrencyPairsResponse{
		CurrencyPairs: cps,
	}, nil
}

// GetPrice gets the QuotePrice and the nonce for the QuotePrice for a given CurrencyPair. The request contains a
// CurrencyPairSelector (either the stringified CurrencyPair, or the CurrencyPair itself). If the request is nil this method fails.
// If the selector is an incorrectly formatted string this method fails. If the QuotePrice / Nonce do not exist for this CurrencyPair, this method fails.
func (q queryServer) GetPrice(goCtx context.Context, req *types.GetPriceRequest) (*types.GetPriceResponse, error) {
	// fail on nil requests
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// unwrap ctx
	ctx := sdk.UnwrapSDKContext(goCtx)

	cps, err := q.k.currencyPairs.Get(ctx, req.CurrencyPairId)
	if err != nil {
		return nil, fmt.Errorf("no price / nonce reported for CurrencyPair: %s, the module is not tracking this CurrencyPair", req.CurrencyPairId)
	}

	// return the QuotePrice + Nonce
	return &types.GetPriceResponse{
		Price:    cps.Price,
		Nonce:    cps.Nonce,
		Decimals: cps.Decimals,
		Id:       cps.Id,
	}, nil
}

// GetPrices gets the array of the QuotePrice and the nonce for the QuotePrice for a given CurrencyPairs.
func (q queryServer) GetPrices(goCtx context.Context, req *types.GetPricesRequest) (*types.GetPricesResponse, error) {
	// fail on nil requests
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// unwrap ctx
	ctx := sdk.UnwrapSDKContext(goCtx)

	prices := make([]types.GetPriceResponse, 0, len(req.CurrencyPairIds))
	for _, cpID := range req.CurrencyPairIds {
		cps, err := q.k.currencyPairs.Get(ctx, cpID)
		if err != nil {
			return nil, fmt.Errorf("no price / nonce reported for CurrencyPair: %s, the module is not tracking this CurrencyPair", cpID)
		}

		prices = append(prices, types.GetPriceResponse{
			Price:    cps.Price,
			Nonce:    cps.Nonce,
			Decimals: cps.Decimals,
			Id:       cps.Id,
		})
	}

	return &types.GetPricesResponse{
		Prices: prices,
	}, nil
}
