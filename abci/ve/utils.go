package ve

import (
	"fmt"

	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	vetypes "github.com/skip-mev/slinky/abci/ve/types"
)

const (
	// MaximumPriceSize defines the maximum size of a price in bytes. This allows
	// up to 32 bytes for the price and 1 byte for the sign (positive/negative).
	MaximumPriceSize = 33
)

// ValidateOracleVoteExtension validates the vote extension provided by a validator.
func ValidateOracleVoteExtension(ve vetypes.OracleVoteExtension) error {
	// Verify prices are valid.
	for _, bz := range ve.Prices {
		// validate the price bytes
		if len(bz) > MaximumPriceSize {
			return fmt.Errorf("price bytes are too long: %d", len(bz))
		}
	}

	return nil
}

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

type (
	// ValidateVoteExtensionsFn defines the function for validating vote extensions. This
	// function is not explicitly used to validate the oracle data but rather that
	// the signed vote extensions included in the proposal are valid and provide
	// a supermajority of vote extensions for the current block. This method is
	// expected to be used in PrepareProposal and ProcessProposal.
	ValidateVoteExtensionsFn func(
		ctx sdk.Context,
		height int64,
		extInfo cometabci.ExtendedCommitInfo,
	) error
)

// NewDefaultValidateVoteExtensionsFn returns a new DefaultValidateVoteExtensionsFn.
func NewDefaultValidateVoteExtensionsFn(chainID string, validatorStore baseapp.ValidatorStore) ValidateVoteExtensionsFn {
	return func(ctx sdk.Context, height int64, info cometabci.ExtendedCommitInfo) error {
		return baseapp.ValidateVoteExtensions(ctx, validatorStore, height, chainID, info)
	}
}

// NoOpValidateVoteExtensions is a no-op validation method (purely used for testing).
func NoOpValidateVoteExtensions(
	_ sdk.Context,
	_ int64,
	_ cometabci.ExtendedCommitInfo,
) error {
	return nil
}
