package abci

import (
	"fmt"

	"cosmossdk.io/math"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/abci/types"
	oracleservice "github.com/skip-mev/slinky/oracle/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// CheckOracleData checks the validity of the oracle data in the proposal by re-running
// the same aggregation logic on the vote extensions that was run in the prepare proposal
// handler. The oracle data is valid if:
//  1. The oracle/vote extension data is present in the proposal and is not nil.
//  2. The vote extensions included compose of a supermajority of signatures (2/3+). This
//     is enforced by the validateVoteExtensionsFn which can be replaced by the application.
//  3. The number of prices in the oracle data included in a proposal matches the number of prices
//     the proposer calculates off of the vote extensions.
//  4. The prices for each asset in the oracle data included in a proposal matches the prices for
//     each asset the proposer calculates off of the vote extensions.
func (h *ProposalHandler) CheckOracleData(ctx sdk.Context, req *cometabci.RequestProcessProposal) (types.OracleData, error) {
	h.logger.Info("verifying oracle data included in proposal")

	// There must be at least one slot in the proposal for the oracle data.
	if len(req.Txs) < NumInjectedTxs {
		return types.OracleData{}, fmt.Errorf("invalid number of transactions in proposal; expected at least %d txs", NumInjectedTxs)
	}

	// Retrieve the oracle info from the proposal. This cannot be empty as we have to at least
	// verify that vote extensions were included and that they are valid.
	oracleInfoBytes := req.Txs[OracleInfoIndex]
	if len(oracleInfoBytes) == 0 {
		return types.OracleData{}, fmt.Errorf("oracle data is nil")
	}

	proposalOracleData := types.OracleData{}
	if err := proposalOracleData.Unmarshal(oracleInfoBytes); err != nil {
		return types.OracleData{}, fmt.Errorf("failed to unmarshal oracle data: %w", err)
	}

	extendedCommitInfo := cometabci.ExtendedCommitInfo{}
	if err := extendedCommitInfo.Unmarshal(proposalOracleData.ExtendedCommitInfo); err != nil {
		return types.OracleData{}, fmt.Errorf("failed to unmarshal extended commit info: %w", err)
	}

	// Verify that the vote extensions included in the proposal are valid.
	if err := h.validateVoteExtensionsFn(ctx, req.Height, extendedCommitInfo); err != nil {
		return types.OracleData{}, fmt.Errorf("failed to validate vote extensions: %w", err)
	}

	// Verify that the oracle price info provided by the proposer matches the vote extensions
	// included in the proposal.
	oracleData, err := h.VerifyOraclePrices(ctx, proposalOracleData, extendedCommitInfo)
	if err != nil {
		return types.OracleData{}, err
	}

	return oracleData, err
}

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
// price aggregator but can be replaced by the application.
func (h *ProposalHandler) AggregateOracleData(
	ctx sdk.Context,
	extendedCommitInfo cometabci.ExtendedCommitInfo,
) (types.OracleData, error) {
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
	return h.BuildOracleData(ctx, extendedCommitInfo, h.priceAggregator.GetPrices())
}

// BuildOracleData combinens all of the price information and vote extensions
// into a single oracle data object.
func (h *ProposalHandler) BuildOracleData(
	ctx sdk.Context,
	extendedCommitInfo cometabci.ExtendedCommitInfo,
	prices map[oracletypes.CurrencyPair]*uint256.Int,
) (types.OracleData, error) {
	h.logger.Info(
		"aggregating oracle data",
		"num_votes", len(extendedCommitInfo.Votes),
		"num_prices", len(prices),
	)

	// Convert the prices to a string map of currency pair -> price.
	priceMap := make(map[string]string)
	for currencyPair, price := range prices {
		priceMap[currencyPair.ToString()] = price.String()
	}

	// Inject the extended commit info into the proposal which contains all vote extensions.
	commitInfoBz, err := extendedCommitInfo.Marshal()
	if err != nil {
		return types.OracleData{}, fmt.Errorf("failed to marshal commit info: %w", err)
	}

	// Create the injected oracle data.
	aggregatedOracleData := types.OracleData{
		Prices:             priceMap,
		ExtendedCommitInfo: commitInfoBz,
		Height:             ctx.BlockHeight(),
		Timestamp:          ctx.BlockHeader().Time,
	}

	return aggregatedOracleData, nil
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
	prices := make(map[oracletypes.CurrencyPair]oracleservice.QuotePrice)
	for asset, priceString := range oracleData.Prices {
		// Convert the price to a uint256.Int. All price feeds are expected to be
		// in the form of a string hex before conversion.
		price, err := uint256.FromHex(priceString)
		if err != nil {
			continue
		}

		// Convert the asset into a currency pair.
		currencyPair, err := oracletypes.CurrencyPairFromString(asset)
		if err != nil {
			continue
		}

		prices[currencyPair] = oracleservice.QuotePrice{
			Price:     price,
			Timestamp: oracleData.Timestamp,
		}
	}

	// insert the prices into the price aggregator.
	h.logger.Debug(
		"adding oracle prices to aggregator",
		"num_prices", len(prices),
		"validator_address", address,
	)

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

// VerifyOraclePrices verifies that the oracle prices provided by the proposer are valid. The
// oracle prices are valid if:
//  1. The number of prices in the oracle data included in a proposal matches the number of prices
//     the proposer calculates off of the vote extensions.
//  2. The prices for each asset in the oracle data included in a proposal matches the prices for
//     each asset the proposer calculates off of the vote extensions.
//
// The same exact aggregation logic that was run in the prepare proposal handler is
// run here to verify the oracle data.
func (h *ProposalHandler) VerifyOraclePrices(
	ctx sdk.Context,
	proposalOracleData types.OracleData,
	extendedCommitInfo cometabci.ExtendedCommitInfo,
) (types.OracleData, error) {
	// Process the oracle info by re-running the same aggregation logic
	// that was run in the prepare proposal handler.
	oracleData, err := h.AggregateOracleData(ctx, extendedCommitInfo)
	if err != nil {
		return types.OracleData{}, err
	}

	// Invariant 1: The number of prices in calculated oracle data must match the number of prices
	// in the proposal oracle data.
	if len(oracleData.Prices) != len(proposalOracleData.Prices) {
		return types.OracleData{}, fmt.Errorf("invalid number of prices in oracle data")
	}

	// Invariant 2: The prices for each asset in the calculated oracle data must match the prices
	// for each asset in the proposal oracle data.
	for asset, priceStr := range oracleData.Prices {
		proposalPriceStr, ok := proposalOracleData.Prices[asset]
		if !ok {
			return types.OracleData{}, fmt.Errorf("missing asset %s in oracle data", asset)
		}

		price, err := uint256.FromHex(priceStr)
		if err != nil {
			return types.OracleData{}, fmt.Errorf("invalid price %s for asset %s", priceStr, asset)
		}

		proposalPrice, err := uint256.FromHex(proposalPriceStr)
		if err != nil {
			return types.OracleData{}, fmt.Errorf("invalid proposal price %s for asset %s", proposalPrice, asset)
		}

		if !price.Eq(proposalPrice) {
			return types.OracleData{}, fmt.Errorf("price mismatch for asset %s", asset)
		}
	}

	return oracleData, nil
}

// WriteOracleData writes the oracle data to state for the supported assets.
func (h *ProposalHandler) WriteOracleData(ctx sdk.Context, oracleData types.OracleData) error {
	if len(oracleData.Prices) == 0 {
		return nil
	}

	// convert the OracleData prices map to a format that is compatible with the oracle module.
	modulePrices := h.toModulePrices(ctx, oracleData)

	// Get the currency pairs currently supported by the oracle module.
	currencyPairs := h.oracleKeeper.GetAllCurrencyPairs(ctx)
	for _, cp := range currencyPairs {
		// Check if there is a price update for the given currency pair.
		quotePrice, ok := modulePrices[cp.ToString()]
		if !ok {
			h.logger.Debug(
				"no price update",
				"currency_pair", cp.ToString(),
			)

			continue
		}

		// Write the currency pair price info to state.
		if err := h.writeCurrencyPairToState(ctx, cp, quotePrice); err != nil {
			return err
		}

		h.logger.Info(
			"price data written to state",
			"currency_pair", cp.ToString(),
			"price", quotePrice.Price,
		)
	}

	return nil
}

// writeCurrencyPairToState writes the currency pair price info to state. There are two different types of writes
// that can occur:
//  1. Prepare Proposal Mode: The currency pair price info is written to only the prepare proposal state.
//  2. Process Proposal Mode: The currency pair price info is written to both the prepare and finalize proposal state.
//
// We do not write oracle data to finalize block state in prepare proposal mode because that will directly
// apply state changes for a proposal that may be rejected. When current mode is process proposal, we write to the current
// context (ProcessProposalState) and finalize context (FinalizeBlockState). This is the desired behavior as otherwise
// there may be inconsistency in whether transactions are successfully executed in the different modes (Prepare,
// Process, Finalize).
func (h *ProposalHandler) writeCurrencyPairToState(
	ctx sdk.Context,
	cp oracletypes.CurrencyPair,
	quotePrice oracletypes.QuotePrice,
) error {
	if err := h.oracleKeeper.SetPriceForCurrencyPair(ctx, cp, quotePrice); err != nil {
		return fmt.Errorf("failed to write oracle data to state: %s", err)
	}

	if ctx.ExecMode() == sdk.ExecModePrepareProposal {
		return nil
	}

	finalizeCtx := h.baseApp.GetFinalizeBlockStateCtx()
	if err := h.oracleKeeper.SetPriceForCurrencyPair(finalizeCtx, cp, quotePrice); err != nil {
		return fmt.Errorf("failed to write oracle data to finalize state: %w", err)
	}

	return nil
}

func (h *ProposalHandler) toModulePrices(ctx sdk.Context, oracleData types.OracleData) map[string]oracletypes.QuotePrice {
	modulePrices := make(map[string]oracletypes.QuotePrice)

	// for each CurrencyPair in the oracleData, convert to a oracletypes.CurrencyPair + format as string
	for currencyPairStr, priceStr := range oracleData.Prices {
		// convert big int string to big int
		price, err := uint256.FromHex(priceStr)
		if err != nil {
			h.logger.Debug(
				"failed to convert price string to big int",
				"currency_pair", currencyPairStr,
				"price", priceStr,
				"err", err,
			)

			continue
		}

		// convert big int to sdk int
		quotePrice := oracletypes.QuotePrice{
			Price:          math.NewIntFromBigInt(price.ToBig()),
			BlockTimestamp: ctx.BlockHeader().Time,
			BlockHeight:    uint64(ctx.BlockHeight()),
		}

		// set the quote price for the currency pair
		modulePrices[currencyPairStr] = quotePrice
	}

	return modulePrices
}
