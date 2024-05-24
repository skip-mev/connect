package simapp

import (
	"github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/abci/preblock/oracle"
)

// PreBlocker wraps the provided preblocker with the oracle preblocker.
func PreBlocker(appPreBlock sdk.PreBlocker, oracleHandler *oracle.PreBlockHandler) sdk.PreBlocker {
	if oracleHandler == nil {
		panic("nil oracleHandler")
	}

	return func(ctx sdk.Context, block *types.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		resp, err := appPreBlock(ctx, block)
		if err != nil {
			return resp, err
		}

		oraclePreBlocker := oracleHandler.PreBlocker()
		_, err = oraclePreBlocker(ctx, block)
		if err != nil {
			return &sdk.ResponsePreBlock{}, err
		}

		return resp, nil
	}
}
