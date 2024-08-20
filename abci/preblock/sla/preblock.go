package sla

import (
	"fmt"

	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	voteaggregator "github.com/skip-mev/connect/v2/abci/strategies/aggregator"
	compression "github.com/skip-mev/connect/v2/abci/strategies/codec"
	"github.com/skip-mev/connect/v2/abci/strategies/currencypair"
	"github.com/skip-mev/connect/v2/abci/ve"
	slakeeper "github.com/skip-mev/connect/v2/x/sla/keeper"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

// PreBlockHandler is responsible for aggregating information about
// oracle price feeds that each validator is including via their vote extensions.
// This handler is run before any transactions are executed/finalized for
// a given block. The handler check's the vote extensions included in each
// transaction and updates the price feed incentives for each validator.
type PreBlockHandler struct {
	// Expected keepers required by the handler.
	oracleKeeper  OracleKeeper
	stakingKeeper StakingKeeper
	slaKeeper     Keeper

	// currencyPairIDStrategy is the strategy used for generating / retrieving
	// IDs for currency-pairs
	currencyPairIDStrategy currencypair.CurrencyPairStrategy

	// voteExtensionCodec is the codec used for encoding / decoding vote extensions.
	// This is used to decode vote extensions included in transactions.
	voteExtensionCodec compression.VoteExtensionCodec

	// extendedCommitCodec is the codec used for encoding / decoding extended
	// commit messages. This is used to decode extended commit messages included
	// in transactions.
	extendedCommitCodec compression.ExtendedCommitCodec
}

// NewSLAPreBlockHandler returns a new PreBlockHandler.
func NewSLAPreBlockHandler(
	oracleKeeper OracleKeeper,
	stakingKeeper StakingKeeper,
	slaKeeper Keeper,
	strategy currencypair.CurrencyPairStrategy,
	voteExtCodec compression.VoteExtensionCodec,
	extendedCommitCodec compression.ExtendedCommitCodec,
) *PreBlockHandler {
	return &PreBlockHandler{
		oracleKeeper:           oracleKeeper,
		stakingKeeper:          stakingKeeper,
		slaKeeper:              slaKeeper,
		currencyPairIDStrategy: strategy,
		voteExtensionCodec:     voteExtCodec,
		extendedCommitCodec:    extendedCommitCodec,
	}
}

// PreBlocker is called by the base app before the block is finalized. Specifically, this
// function is called before any transactions are executed, but after oracle prices have
// been written to the oracle store. This function will retrieve all of the vote extensions
// that were included in the block, determine which currency pairs each validator included
// prices for, and update the price feed incentives for each validator. Enforcements of SLAs
// is done in the SLA BeginBlocker.
func (h *PreBlockHandler) PreBlocker() sdk.PreBlocker {
	return func(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		if req == nil {
			ctx.Logger().Error(
				"received nil RequestFinalizeBlock in SLA preblocker",
				"height", ctx.BlockHeight(),
			)

			return &sdk.ResponsePreBlock{}, fmt.Errorf("received nil RequestFinalizeBlock in SLA preblocker: height %d", ctx.BlockHeight())
		}

		if !ve.VoteExtensionsEnabled(ctx) {
			ctx.Logger().Info(
				"skipping sla price feed preblocker because vote extensions are not enabled",
				"height", ctx.BlockHeight(),
			)

			return &sdk.ResponsePreBlock{}, nil
		}

		ctx.Logger().Info(
			"executing price feed sla pre-block hook",
			"height", ctx.BlockHeight(),
		)

		// Retrieve all vote extensions that were included in the block. This
		// returns a list of validators and the price updates that they made.
		votes, err := voteaggregator.GetOracleVotes(req.Txs, h.voteExtensionCodec, h.extendedCommitCodec)
		if err != nil {
			ctx.Logger().Error(
				"failed to get extended commit info from proposal",
				"height", ctx.BlockHeight(),
				"num_txs", len(req.Txs),
				"err", err,
			)

			return &sdk.ResponsePreBlock{}, err
		}

		// Create a mapping of price updates by status for each validator.
		updates, err := h.GetUpdates(ctx, votes)
		if err != nil {
			ctx.Logger().Error(
				"failed to get price feed updates",
				"height", ctx.BlockHeight(),
				"votes", votes,
				"err", err,
			)

			return &sdk.ResponsePreBlock{}, err
		}

		// Update all of the price feeds for each validator.
		if err := h.slaKeeper.UpdatePriceFeeds(ctx, updates); err != nil {
			ctx.Logger().Error(
				"failed to update price feeds",
				"height", ctx.BlockHeight(),
				"votes", votes,
				"err", err,
			)
		}

		return &sdk.ResponsePreBlock{}, nil
	}
}

// GetUpdates returns a mapping of every validator's price feed status updates. This function
// will iterate through the active set of validators, determine which currency pairs they
// included prices for via their vote extensions, and return a mapping of each validator's
// status updates.
func (h *PreBlockHandler) GetUpdates(ctx sdk.Context, votes []voteaggregator.Vote) (slakeeper.PriceFeedUpdates, error) {
	updates := slakeeper.NewPriceFeedUpdates()

	// Retrieve all bonded validators which will be considered for updates.
	bondedValidators, err := h.stakingKeeper.GetBondedValidatorsByPower(ctx)
	if err != nil {
		ctx.Logger().Error("failed to get bonded validators", "err", err)
		return updates, err
	}

	currencyPairs := h.oracleKeeper.GetAllCurrencyPairs(ctx)
	for _, cp := range currencyPairs {
		updates.CurrencyPairs[cp] = struct{}{}
	}

	// Initialize an empty status for each of the bonded validators.
	for _, validator := range bondedValidators {
		consAddrBz, err := validator.GetConsAddr()
		if err != nil {
			ctx.Logger().Error("failed to get consensus address", "err", err)
			return updates, err
		}
		consAddress := sdk.ConsAddress(consAddrBz)

		// Initialize the status for each currency pair to NoVote.
		validator := slakeeper.NewValidatorUpdate(consAddress)
		for _, cp := range currencyPairs {
			validator.Updates[cp] = slatypes.NoVote
		}

		updates.ValidatorUpdates[consAddress.String()] = validator
	}

	// Determine the price feed status updates for each validator that included
	// their vote extension in the block.
	for _, vote := range votes {
		valUpdates := getStatuses(ctx, h.currencyPairIDStrategy, currencyPairs, vote.OracleVoteExtension.Prices)

		ctx.Logger().Debug(
			"retrieved status updates by validator",
			"validator", vote.ConsAddress.String(),
			"updates", valUpdates,
		)

		validator := updates.ValidatorUpdates[vote.ConsAddress.String()]
		validator.Updates = valUpdates
		updates.ValidatorUpdates[vote.ConsAddress.String()] = validator
	}

	return updates, nil
}
