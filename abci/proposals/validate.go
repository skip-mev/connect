package proposals

import (
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/abci/ve"
)

// ValidateExtendedCommitInfo validates the extended commit info for a block. It first
// ensures that the vote extensions compose a supermajority of the signatures and
// voting power for the block. Then, it ensures that oracle vote extensions are correctly
// marshalled and contain valid prices.
func (h *ProposalHandler) ValidateExtendedCommitInfo(
	ctx sdk.Context,
	height int64,
	extendedCommitInfo cometabci.ExtendedCommitInfo,
) error {
	if err := h.validateVoteExtensionsFn(ctx, extendedCommitInfo); err != nil {
		h.logger.Error(
			"failed to validate vote extensions; vote extensions may not comprise a supermajority",
			"height", height,
			"err", err,
		)

		return err
	}

	// Validate all oracle vote extensions.
	for _, vote := range extendedCommitInfo.Votes {
		address := sdk.ConsAddress{}
		if err := address.Unmarshal(vote.Validator.Address); err != nil {
			h.logger.Error(
				"failed to unmarshal validator address",
				"height", height,
			)

			return err
		}

		voteExt, err := h.voteExtensionCodec.Decode(vote.VoteExtension)
		if err != nil {
			return err
		}

		// The vote extension are from the previous block.
		if err := ve.ValidateOracleVoteExtension(ctx, voteExt, h.currencyPairStrategy); err != nil {
			h.logger.Error(
				"failed to validate oracle vote extension",
				"height", height,
				"validator", address.String(),
				"err", err,
			)

			return err
		}
	}

	return nil
}
