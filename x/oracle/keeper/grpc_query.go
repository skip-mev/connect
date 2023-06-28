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
func (q queryServer) GetPrice(goCtx context.Context, req *types.GetPriceRequest) (_ *types.GetPriceResponse, err error) {
	var cp types.CurrencyPair

	// fail on nil requests
	if req == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}
	// determine what type the selector from the response is giving
	switch cpI := req.CurrencyPairSelector.(type) {

	case *types.GetPriceRequest_CurrencyPairId:
		// retrieve the currency pair from the stringified ID, and fail if incorrectly formatted
		cp, err = types.CurrencyPairFromString(cpI.CurrencyPairId)

		if err != nil {
			return nil, fmt.Errorf("error unmarshalling CurrencyPairID: %v", err)
		}

	case *types.GetPriceRequest_CurrencyPair:
		// retrieve CurrencyPair directly from selector
		if cpI.CurrencyPair == nil {
			return nil, fmt.Errorf("currency Pair cannot be nil")
		}
		cp = *cpI.CurrencyPair

	default:
		// fail if any other type of CurrencyPairSelector is given
		return nil, fmt.Errorf("invalid CurrencyPairSelector given in request (consult documentation)")
	}

	// unwrap ctx
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get the QuotePrice + nonce for the given CurrencyPair
	qpn, err := q.k.GetPriceWithNonceForCurrencyPair(ctx, cp)
	if err != nil {
		return nil, fmt.Errorf("no price / nonce reported for CurrencyPair: %v, the module is not tracking this CurrencyPair", cp)
	}

	// return the QuotePrice + Nonce
	return &types.GetPriceResponse{
		Price:    &qpn.QuotePrice,
		Nonce:    qpn.Nonce(),
		Decimals: uint64(cp.Decimals()),
	}, nil
}
