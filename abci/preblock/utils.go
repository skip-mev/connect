package preblock

import (
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NoOpPreBlocker is a no-op preblocker. This should only be used for testing.
func NoOpPreBlocker() sdk.PreBlocker {
	return func(_ sdk.Context, _ *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		return &sdk.ResponsePreBlock{}, nil
	}
}
