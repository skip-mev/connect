package proposals

import (
	cometabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/abci/strategies/codec"
	"github.com/skip-mev/connect/v2/abci/strategies/currencypair"
	"github.com/skip-mev/connect/v2/abci/ve"
)

// ValidateExtendedCommitInfo validates the extended commit info for a block. It first
// ensures that the vote extensions compose a super-majority of the signatures and
// voting power for the block. Then, it ensures that oracle vote extensions are correctly
// marshalled and contain valid prices.
func (h *ProposalHandler) ValidateExtendedCommitInfo(
	ctx sdk.Context,
	height int64,
	extendedCommitInfo cometabci.ExtendedCommitInfo,
) error {
	if err := h.validateVoteExtensionsFn(ctx, extendedCommitInfo); err != nil {
		h.logger.Error(
			"failed to validate vote extensions; vote extensions may not comprise a super-majority",
			"height", height,
			"err", err,
		)

		return err
	}

	// Validate all oracle vote extensions.
	for _, vote := range extendedCommitInfo.Votes {
		address := sdk.ConsAddress(vote.Validator.Address)
		// The vote extension are from the previous block.
		if err := validateVoteExtension(ctx, vote, h.voteExtensionCodec, h.currencyPairStrategy); err != nil {
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

// PruneAndValidateExtendedCommitInfo validates each vote-extension in the extended commit, and removes
// any vote-extensions that are invalid. Removal will effectively treat the validator's
// vote as absent.  This function performs all validation that ValidateExtendedCommitInfo performs.
func (h *ProposalHandler) PruneAndValidateExtendedCommitInfo(
	ctx sdk.Context, extendedCommitInfo cometabci.ExtendedCommitInfo,
) (cometabci.ExtendedCommitInfo, error) {
	// Validate all oracle vote extensions.
	for i, vote := range extendedCommitInfo.Votes {
		// validate the vote-extension
		if err := validateVoteExtension(ctx, vote, h.voteExtensionCodec, h.currencyPairStrategy); err != nil {
			h.logger.Info(
				"failed to validate vote extension - pruning vote",
				"err", err,
				"validator", vote.Validator.Address,
			)

			// failed to validate this vote-extension, mark it as absent in the original commit
			vote.BlockIdFlag = cometproto.BlockIDFlagAbsent
			vote.ExtensionSignature = nil
			vote.VoteExtension = nil
			extendedCommitInfo.Votes[i] = vote
		}
	}

	// validate after pruning
	if err := h.validateVoteExtensionsFn(ctx, extendedCommitInfo); err != nil {
		h.logger.Error(
			"failed to validate vote extensions; vote extensions may not comprise a super-majority",
			"err", err,
		)

		return cometabci.ExtendedCommitInfo{}, err
	}

	return extendedCommitInfo, nil
}

func validateVoteExtension(
	ctx sdk.Context,
	vote cometabci.ExtendedVoteInfo,
	voteExtensionCodec codec.VoteExtensionCodec,
	currencyPairStrategy currencypair.CurrencyPairStrategy,
) error {
	// vote is not voted for if VE is nil
	if vote.VoteExtension == nil && vote.ExtensionSignature == nil {
		return nil
	}

	voteExt, err := voteExtensionCodec.Decode(vote.VoteExtension)
	if err != nil {
		return err
	}

	// The vote extensions are from the previous block.
	if err := ve.ValidateOracleVoteExtension(ctx, voteExt, currencyPairStrategy); err != nil {
		return err
	}

	return nil
}
