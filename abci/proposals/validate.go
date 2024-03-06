package proposals

import (
	"fmt"

	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/skip-mev/slinky/abci/strategies/codec"
	"github.com/skip-mev/slinky/abci/strategies/currencypair"
	"github.com/skip-mev/slinky/abci/ve"
)

// PruneExtendedCommitInfo validates each vote-extension in the extended commit, and removes
// any vote-extensions that are invalid. Removal will effectively treat the validator's
// vote as absent.
func (h *ProposalHandler) PruneExtendedCommitInfo(
	ctx sdk.Context, extendedCommitInfo cometabci.ExtendedCommitInfo,
) (cometabci.ExtendedCommitInfo, error) {
	// branch the context for verifying vote-extensions on verify-vote extension state
	verifyVECtx, err := h.appState.VerifyVoteExtensionState(ctx)
	if err != nil {
		h.logger.Error(
			"failed to retrieve verify vote extension state",
			"err", err,
		)
	}

	// tally total voting-power of valid votes
	var (
		votingPowerInCommit, totalVotingPower uint64
	)

	// Validate all oracle vote extensions.
	for i, vote := range extendedCommitInfo.Votes {
		// sum the total voting power
		totalVotingPower += uint64(vote.Validator.Power)
		// ignore non-commit votes
		if vote.BlockIdFlag != cometproto.BlockIDFlagCommit {
			continue
		}
		
		// validate the vote-extension
		if err := validateVoteExtension(verifyVECtx, vote, h.voteExtensionCodec, h.currencyPairStrategy); err != nil {
			h.logger.Error(
				"failed to validate vote extension",
				"err", err,
				"validator", vote.Validator.Address,
			)

			// failed to validate this vote-extension, mark it as absent in the original commit
			vote.BlockIdFlag = cometproto.BlockIDFlagAbsent
			vote.ExtensionSignature = []byte{}
			vote.VoteExtension = []byte{}
			extendedCommitInfo.Votes[i] = vote
		} else {
			// tally the voting power of this vote
			votingPowerInCommit += uint64(vote.Validator.Power)
		}
	}

	// ensure that the valid vote-extensions compose a super-majority of voting-power
	if requiredVP := (totalVotingPower * 2/3) + 1; votingPowerInCommit < requiredVP {
		h.logger.Error(
			"vote extensions do not compose a supermajority",
			"voting-power-in-commit", votingPowerInCommit,
			"required-voting-power", requiredVP,
		)

		return cometabci.ExtendedCommitInfo{}, fmt.Errorf(
			"vote extensions do not compose a supermajority, expected: %d, got: %d", requiredVP, votingPowerInCommit,
		)
	}

	return extendedCommitInfo, nil
}

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

	// branch the context for verifying vote-extensions on verify-vote extension state
	verifyVECtx, err := h.appState.VerifyVoteExtensionState(ctx)
	if err != nil {
		h.logger.Error(
			"failed to retrieve verify vote extension state",
			"height", height,
			"err", err,
		)

		return err
	}

	// Validate all oracle vote extensions.
	for _, vote := range extendedCommitInfo.Votes {

		// validate the vote-extension with the verify-vote extension state
		if err := validateVoteExtension(verifyVECtx, vote, h.voteExtensionCodec, h.currencyPairStrategy); err != nil {
			h.logger.Error(
				"failed to validate vote extension",
				"height", ctx.BlockHeight(),
				"err", err,
				"validator", vote.Validator.Address,
			)

			return err
		}
	}

	return nil
}

func validateVoteExtension(
	ctx sdk.Context,
	vote cometabci.ExtendedVoteInfo,
	voteExtensioncodec codec.VoteExtensionCodec,
	currencyPairStrategy currencypair.CurrencyPairStrategy,
) error {
	address := sdk.ConsAddress{}
	if err := address.Unmarshal(vote.Validator.Address); err != nil {
		return err
	}

	voteExt, err := voteExtensioncodec.Decode(vote.VoteExtension)
	if err != nil {
		return err
	}

	// The vote extensions are from the previous block.
	if err := ve.ValidateOracleVoteExtension(ctx, voteExt, currencyPairStrategy); err != nil {
		return err
	}

	return nil
}
