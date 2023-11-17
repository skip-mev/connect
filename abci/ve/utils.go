package ve

import (
	"fmt"

	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/holiman/uint256"

	"github.com/skip-mev/slinky/abci/ve/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// ValidateOracleVoteExtension validates the vote extension provided by a validator.
func ValidateOracleVoteExtension(voteExtension []byte, height int64) error {
	if len(voteExtension) == 0 {
		return nil
	}

	voteExt := types.OracleVoteExtension{}
	if err := voteExt.Unmarshal(voteExtension); err != nil {
		return fmt.Errorf("failed to unmarshal vote extension: %w", err)
	}

	// The height of the vote extension must match the height of the request.
	if voteExt.Height != height {
		return fmt.Errorf(
			"vote extension height does not match request height; expected: %d, got: %d",
			height,
			voteExt.Height,
		)
	}

	// Verify tickers and prices are valid.
	for currencyPair, price := range voteExt.Prices {
		if _, err := oracletypes.CurrencyPairFromString(currencyPair); err != nil {
			return fmt.Errorf("invalid ticker in oracle vote extension %s: %w", currencyPair, err)
		}

		if _, err := uint256.FromHex(price); err != nil {
			return fmt.Errorf("invalid price in oracle vote extension %s: %w", currencyPair, err)
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
