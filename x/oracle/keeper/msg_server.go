package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/oracle/types"
)

// msgServer is the default implementation of the x/oracle MsgService.
type msgServer struct {
	k Keeper
}

// NewMsgServer returns the default implementation of the x/oracle message service.
func NewMsgServer(k Keeper) types.MsgServer {
	return &msgServer{k}
}

var _ types.MsgServer = (*msgServer)(nil)

// AddCurrencyPairs takes a set of currency pairs to be added, and adds them to the module's state. This method fails if any of the
// currency-pairs fail to be set, if the message is invalid, or if the signer is not the authority account of the module. If any of the CurrencyPairs
// to be added already exist in the module, they will be skipped.
func (m *msgServer) AddCurrencyPairs(goCtx context.Context, req *types.MsgAddCurrencyPairs) (*types.MsgAddCurrencyPairsResponse, error) {
	// check the validity of the message
	if req == nil {
		return nil, fmt.Errorf("message cannot be empty")
	}

	if m.k.mmKeeper != nil {
		return nil, fmt.Errorf("x/oracle message server is disabled when using x/marketmap")
	}

	// check that the authority of the message is the authority of the module
	if req.Authority != m.k.authority.String() {
		return nil, fmt.Errorf("message validation failed: authority %s is not module authority %s", req.Authority, m.k.authority)
	}

	// finally, add all currency pairs in message to state
	ctx := sdk.UnwrapSDKContext(goCtx)
	for _, cp := range req.CurrencyPairs {
		// only set the currency-pair if it does not already exist in state
		if !m.k.HasCurrencyPair(ctx, cp) {
			// set to state, initial nonce will be zero (no price updates have been made for this CurrencyPair)
			err := m.k.CreateCurrencyPair(ctx, cp)
			if err != nil {
				return nil, fmt.Errorf("error creating currency pair state: %w", err)
			}
		}
	}

	return &types.MsgAddCurrencyPairsResponse{}, nil
}

// RemoveCurrencyPairs takes a set of CurrencyPairs to remove. CurrencyPairs given are represented by string identifiers of CurrencyPairs
// i.e `cp.String()`. For each CurrencyPair in the message, remove the Nonce / QuotePrice data for that CurrencyPair, if a CurrencyPair is
// given that is not currently tracked, skip, and continue removing CurrencyPairs.
func (m *msgServer) RemoveCurrencyPairs(goCtx context.Context, req *types.MsgRemoveCurrencyPairs) (*types.MsgRemoveCurrencyPairsResponse, error) {
	// check validity of message
	if req == nil {
		return nil, fmt.Errorf("message cannot be empty")
	}

	if m.k.mmKeeper != nil {
		return nil, fmt.Errorf("x/oracle message server is disabled when using x/marketmap")
	}

	// check that the authority of the message is the authority of the module
	if req.Authority != m.k.authority.String() {
		return nil, fmt.Errorf("message validation failed: authority %s is not module authority %s", req.Authority, m.k.authority)
	}

	// remove all currency-pairs in msg from state
	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, id := range req.CurrencyPairIds {
		// get cp from identifier string
		cp, err := connecttypes.CurrencyPairFromString(id)
		if err != nil {
			return nil, fmt.Errorf("error retrieving CurrencyPair from request: %w", err)
		}

		// delete the currency pair from state
		if err := m.k.RemoveCurrencyPair(ctx, cp); err != nil {
			return nil, fmt.Errorf("error removing currency pair from state: %w", err)
		}
	}

	return nil, nil
}
