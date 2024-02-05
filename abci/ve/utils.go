package ve

import (
	"fmt"

	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/abci/strategies/currencypair"
	slinkyabci "github.com/skip-mev/slinky/abci/types"
	vetypes "github.com/skip-mev/slinky/abci/ve/types"
)

// ValidateOracleVoteExtension validates the vote extension provided by a validator.
func ValidateOracleVoteExtension(
	ctx sdk.Context,
	ve vetypes.OracleVoteExtension,
	strategy currencypair.CurrencyPairStrategy,
) error {
	// Verify prices are valid.
	for id, bz := range ve.Prices {
		// Ensure that the price bytes are not too long.
		if len(bz) > slinkyabci.MaximumPriceSize {
			return fmt.Errorf("price bytes are too long: %d", len(bz))
		}

		// Ensure that the currency pair ID is valid.
		cp, err := strategy.FromID(ctx, id)
		if err != nil {
			return fmt.Errorf("invalid currency pair ID: %d", id)
		}

		// Ensure that the price bytes are valid.
		if _, err := strategy.GetDecodedPrice(ctx, cp, bz); err != nil {
			return fmt.Errorf("invalid price bytes: %s", err)
		}
	}

	return nil
}

// VoteExtensionsEnabled determines if vote extensions are enabled for the current block. If
// vote extensions are enabled at height h, then a proposer will receive vote extensions
// in height h+1. This is primarily utilized by any module that needs to make state changes
// based on whether or not they were included in a proposal.
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

	return cp.Abci.VoteExtensionsEnableHeight < ctx.BlockHeight()
}

// ValidateVoteExtensionsFn defines the function for validating vote extensions. This
// function is not explicitly used to validate the oracle data but rather that
// the signed vote extensions included in the proposal are valid and provide
// a supermajority of vote extensions for the current block. This method is
// expected to be used in PrepareProposal and ProcessProposal.
type ValidateVoteExtensionsFn func(
	ctx sdk.Context,
	height int64,
	extInfo cometabci.ExtendedCommitInfo,
) error

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
