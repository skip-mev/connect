package abci

import sdk "github.com/cosmos/cosmos-sdk/types"

// VoteExtensionsEnabled determines if vote extensions are enabled for the current block.
func VoteExtensionsEnabled(ctx sdk.Context) bool {
	cp := ctx.ConsensusParams()
	if cp.Abci == nil || cp.Abci.VoteExtensionsEnableHeight == 0 {
		return false
	}

	// Per the cosmos sdk, the first block should not utilize the latest finalize block state. This means
	// vote extensions should NOT be making state changes.
	//
	// Ref: https://github.com/cosmos/cosmos-sdk/blob/2100a73dcea634ce914977dbddb4991a020ee345/baseapp/baseapp.go#L488-L495
	if ctx.BlockHeight() <= 1 {
		return false
	}

	// We do a +1 here because the vote extensions are enabled at height h
	// but a proposer will only receive vote extensions in height h+1.
	return cp.Abci.VoteExtensionsEnableHeight+1 < ctx.BlockHeight()
}
