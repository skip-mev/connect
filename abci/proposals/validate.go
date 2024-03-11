package proposals

import (
	"fmt"

	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkyabci "github.com/skip-mev/slinky/abci/types"
	servicemetrics "github.com/skip-mev/slinky/service/metrics"

	"github.com/skip-mev/slinky/abci/ve"
)

// ValidateExtendedCommitInfoPrepare validates the extended commit info for a block. It first
// ensures that the vote extensions compose a supermajority of the signatures and
// voting power for the block. Then, it ensures that oracle vote extensions are correctly
// marshalled and contain valid prices.  This function is to be run in PrepareProposal.
func (h *ProposalHandler) ValidateExtendedCommitInfoPrepare(
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

// ValidateExtendedCommitInfoProcess validates the extended commit info for a block. It first
// ensures that the vote extensions compose a supermajority of the signatures and
// voting power for the block. Then, it ensures that oracle vote extensions are correctly
// marshalled and contain valid prices. This function contains extra validation to be run in
// ProcessProposal.
func (h *ProposalHandler) ValidateExtendedCommitInfoProcess(
	ctx sdk.Context,
	req *cometabci.RequestProcessProposal,
	extendedCommitInfo cometabci.ExtendedCommitInfo,
) error {
	if req == nil {
		return slinkyabci.NilRequestError{
			Handler: servicemetrics.ProcessProposal,
		}
	}

	if err := h.validateVoteExtensionsFn(ctx, extendedCommitInfo); err != nil {
		h.logger.Error(
			"failed to validate vote extensions; vote extensions may not comprise a supermajority",
			"height", req.Height,
			"err", err,
		)

		return err
	}

	if len(extendedCommitInfo.Votes) != len(req.ProposedLastCommit.Votes) {
		h.logger.Error(
			"mismatched length in encoded extended commit info and proposed last commit",
			"height", req.Height,
			"extended commit length", len(extendedCommitInfo.Votes),
			"proposed last commit length", len(req.ProposedLastCommit.Votes),
		)

		return fmt.Errorf("mismatched length in encoded extended commit info and proposed last commit")
	}

	// Validate all oracle vote extensions.  And cross-reference them with the ProposedLastCommit
	for _, vote := range extendedCommitInfo.Votes {
		address := sdk.ConsAddress{}
		if err := address.Unmarshal(vote.Validator.Address); err != nil {
			h.logger.Error(
				"failed to unmarshal validator address",
				"height", req.Height,
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
				"height", req.Height,
				"height", req.Height,
				"validator", address.String(),
				"err", err,
			)

			return err
		}
	}

	return nil
}
