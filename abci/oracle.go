package abci

import (
	"fmt"

	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/abci/types"
	oracletypes "github.com/skip-mev/slinky/oracle/types"
)

// AggregateOracleData ingresses extended commit info which contains all of the
// vote extensions each validator extended in the previous block. Each vote extension
// corresponds to the oracle data that the validator is providing for the current
// block. However, it is important to note that
//  1. The vote extension may be nil, in which case the validator is not providing
//     any oracle data for the current block. This could have occurred because the
//     validator was offline, or its local oracle service was down.
//  2. The vote extension may contain prices updates for only a subset of currency pairs.
//     This could have occurred because the price providers for the validator were
//     offline, or the price providers did not provide a price update for a given
//     currency pair.
//
// In order for a currency pair to be included in the final oracle price, the currency
// pair must be provided by a supermajority (2/3+) of validators. This is enforced by the
// price aggregator.
func (h *ProposalHandler) AggregateOracleData(
	ctx sdk.Context,
	extendedCommitInfo cometabci.ExtendedCommitInfo,
) ([]byte, error) {
	// Reset the price aggregator.
	h.priceAggregator.ResetProviderPrices()

	// Iterate through all vote extensions and consolidate all price info before
	// aggregating.
	for _, commitInfo := range extendedCommitInfo.Votes {
		address := &sdk.ValAddress{}
		if err := address.Unmarshal(commitInfo.Validator.Address); err != nil {
			h.logger.Debug(
				"failed to unmarshal validator address",
				"err", err,
			)

			continue
		}

		// Retrieve the oracle data from the vote extension.
		oracleData, err := h.GetOracleDataFromVE(commitInfo.VoteExtension)
		if err != nil {
			h.logger.Debug(
				"failed to get oracle data from vote extension",
				"validator", address.String(),
				"err", err,
			)

			continue
		}

		// Add the validator data to the price aggregator.
		h.AddOracleDataToAggregator(address.String(), oracleData)
	}

	// Compute the final prices for each currency pair.
	h.priceAggregator.UpdatePrices()

	// Build the oracle transaction and return it.
	return h.BuildOracleTx(ctx, extendedCommitInfo, h.priceAggregator.GetPrices())
}

// BuildOracleTx marshals the final oracle prices, vote extensions and more into bytes so that
// it can be included in a proposal. Note that this is not an actual transaction.
func (h *ProposalHandler) BuildOracleTx(
	ctx sdk.Context,
	extendedCommitInfo cometabci.ExtendedCommitInfo,
	prices map[oracletypes.CurrencyPair]*uint256.Int,
) ([]byte, error) {
	// Convert the prices to a string map of currency pair -> price.
	priceMap := make(map[string]string)
	for currencyPair, price := range prices {
		priceMap[currencyPair.String()] = price.String()
	}

	// Inject the extended commit info into the proposal which contains all vote extensions.
	commitInfoBz, err := extendedCommitInfo.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal commit info: %w", err)
	}

	// Create the injected oracle data.
	aggregatedOracleData := &types.OracleData{
		Prices:             priceMap,
		ExtendedCommitInfo: commitInfoBz,
		Height:             ctx.BlockHeight(),
		Timestamp:          ctx.BlockHeader().Time,
	}

	// Marshal the oracle data and attempt to include it in the proposal.
	bz, err := aggregatedOracleData.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal oracle data: %w", err)
	}

	return bz, nil
}

// AddOracleDataToAggregator consolidates the oracle data from a single validator
// into the price aggregator. The oracle data is provided in the form of a vote
// extension. The vote extension contains the prices for each currency pair that
// the validator is providing for the current block. In the case where the vote
// extension is nil, or price info is not contained within the vote extension,
// the oracle data is not added to the price aggregator.
func (h *ProposalHandler) AddOracleDataToAggregator(address string, oracleData *types.OracleVoteExtension) {
	if len(oracleData.Prices) == 0 {
		return
	}

	// Format all of the prices into a map of currency pair -> price.
	prices := make(map[oracletypes.CurrencyPair]oracletypes.QuotePrice)
	for asset, priceString := range oracleData.Prices {
		// Convert the price to a uint256.Int. All price feeds are expected to be
		// in the form of a string hex before conversion.
		price, err := uint256.FromHex(priceString)
		if err != nil {
			continue
		}

		// Convert the asset into a currency pair.
		currencyPair, err := oracletypes.NewCurrencyPairFromString(asset)
		if err != nil {
			continue
		}

		prices[currencyPair] = oracletypes.QuotePrice{
			Price:     price,
			Timestamp: oracleData.Timestamp,
		}
	}

	// insert the prices into the price aggregator.
	h.priceAggregator.SetProviderPrices(address, prices)
}

// GetOracleDataFromVE inputs the raw vote extension bytes and returns the
// oracle data contained within. In the case where the vote extension is nil,
// or price info is not contained within the vote extension, an error is returned.
func (h *ProposalHandler) GetOracleDataFromVE(voteExtension []byte) (*types.OracleVoteExtension, error) {
	if len(voteExtension) == 0 {
		return nil, fmt.Errorf("vote extension is nil")
	}

	oracleData := &types.OracleVoteExtension{}
	if err := oracleData.Unmarshal(voteExtension); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vote extension: %w", err)
	}

	return oracleData, nil
}
