package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/oracle/types"
)

// queryServer is the default implementation of the x/oracle QueryService.
type queryServer struct {
	k Keeper
}

// NewQueryServer returns an implementation of the x/oracle QueryServer.
func NewQueryServer(k Keeper) types.QueryServer {
	return queryServer{
		k,
	}
}

var _ types.QueryServer = queryServer{}

// GetAllCurrencyPairs returns the set of all currency pairs that the module is tracking QuotePrices for.
// It returns an error to the caller if there are no CurrencyPairs being tracked by the module.
func (q queryServer) GetAllCurrencyPairs(ctx context.Context, _ *types.GetAllCurrencyPairsRequest) (*types.GetAllCurrencyPairsResponse, error) {
	// get all currency pairs from state
	cps := q.k.GetAllCurrencyPairs(sdk.UnwrapSDKContext(ctx))
	return &types.GetAllCurrencyPairsResponse{
		CurrencyPairs: cps,
	}, nil
}

// GetPrice gets the QuotePrice and the nonce for the QuotePrice for a given CurrencyPair. The request contains a
// CurrencyPairSelector (either the stringified CurrencyPair, or the CurrencyPair itself). If the request is nil this method fails.
// If the selector is an incorrectly formatted string this method fails. If the QuotePrice / Nonce do not exist for this CurrencyPair, this method fails.
func (q queryServer) GetPrice(goCtx context.Context, req *types.GetPriceRequest) (_ *types.GetPriceResponse, err error) {
	// fail on nil requests
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	cp, err := connecttypes.CurrencyPairFromString(req.CurrencyPair)
	if err != nil {
		return nil, fmt.Errorf("invalid currency pair: %w", err)
	}

	if err = cp.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid currency pair: %w", err)
	}

	// unwrap ctx
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get the QuotePrice + nonce for the given CurrencyPair
	qpn, err := q.k.GetPriceWithNonceForCurrencyPair(ctx, cp)
	if err != nil {
		return nil, fmt.Errorf("no price / nonce reported for CurrencyPair: %s, the module is not tracking this CurrencyPair", cp.String())
	}

	id, ok := q.k.GetIDForCurrencyPair(ctx, cp)
	if !ok {
		return nil, fmt.Errorf("no ID found for CurrencyPair: %s", cp.String())
	}

	decimals, err := q.k.GetDecimalsForCurrencyPair(ctx, cp)
	if err != nil {
		return nil, err
	}

	// return the QuotePrice + Nonce
	return &types.GetPriceResponse{
		Price:    &qpn.QuotePrice,
		Nonce:    qpn.Nonce(),
		Decimals: decimals,
		Id:       id,
	}, nil
}

// GetPrices gets the array of the QuotePrice and the nonce for the QuotePrice for a given CurrencyPairs.
func (q queryServer) GetPrices(goCtx context.Context, req *types.GetPricesRequest) (_ *types.GetPricesResponse, err error) {
	var cp connecttypes.CurrencyPair

	// fail on nil requests
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	prices := make([]types.GetPriceResponse, 0, len(req.CurrencyPairIds))
	for _, cid := range req.CurrencyPairIds {
		cp, err = connecttypes.CurrencyPairFromString(cid)
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling CurrencyPairID: %w", err)
		}

		// unwrap ctx
		ctx := sdk.UnwrapSDKContext(goCtx)

		// get the QuotePrice + nonce for the given CurrencyPair
		qpn, err := q.k.GetPriceWithNonceForCurrencyPair(ctx, cp)
		if err != nil {
			return nil, fmt.Errorf("no price / nonce reported for CurrencyPair: %v, the module is not tracking this CurrencyPair", cp)
		}

		id, ok := q.k.GetIDForCurrencyPair(ctx, cp)
		if !ok {
			return nil, fmt.Errorf("no ID found for CurrencyPair: %v", cp)
		}

		decimals, err := q.k.GetDecimalsForCurrencyPair(ctx, cp)
		if err != nil {
			return nil, err
		}

		prices = append(prices, types.GetPriceResponse{
			Price:    &qpn.QuotePrice,
			Nonce:    qpn.Nonce(),
			Decimals: decimals,
			Id:       id,
		})
	}

	return &types.GetPricesResponse{
		Prices: prices,
	}, nil
}

func (q queryServer) GetCurrencyPairMapping(ctx context.Context, _ *types.GetCurrencyPairMappingRequest) (*types.GetCurrencyPairMappingResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	pairs, err := q.k.GetCurrencyPairMapping(sdkCtx)
	if err != nil {
		return nil, err
	}
	return &types.GetCurrencyPairMappingResponse{CurrencyPairMapping: pairs}, nil
}
