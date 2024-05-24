package simapp

import (
	"github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/abci/preblock/oracle"
)

// PreBlocker wraps the provider preblocker with the oracle pre blocker.
func PreBlocker(appPreBlock sdk.PreBlocker, oracleHandler *oracle.PreBlockHandler) sdk.PreBlocker {
	return func(context sdk.Context, block *types.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		oraclePreBlocker := oracleHandler.PreBlocker()
		_, err := oraclePreBlocker(context, block)
		if err != nil {
			return nil, err
		}

		return appPreBlock(context, block)
	}
}
