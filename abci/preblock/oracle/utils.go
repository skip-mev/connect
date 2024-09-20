package oracle

import (
	"math/big"

	cometabci "github.com/cometbft/cometbft/abci/types"
	cometproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	servicemetrics "github.com/skip-mev/connect/v2/service/metrics"
)

// recordPrice records all the given prices per ticker, and reports them as a float64.
func (h *PreBlockHandler) recordPrices(prices map[connecttypes.CurrencyPair]*big.Int) {
	for ticker, price := range prices {
		floatPrice, _ := price.Float64()
		h.metrics.ObservePriceForTicker(ticker, floatPrice)
	}
}

// recordValidatorReports takes the commit decided for this block, and for each validator in the commit, records
// whether their vote was included in the commit, whether they reported a price for each currency-pair, and if so
// the price they reported.
func (h *PreBlockHandler) recordValidatorReports(ctx sdk.Context, decidedCommit cometabci.CommitInfo) {
	pricesToReport := h.keeper.GetAllCurrencyPairs(ctx)

	// iterate over each validator in the commit
	for _, vote := range decidedCommit.Votes {
		var nilVote bool
		validator := sdk.ConsAddress(vote.Validator.Address)
		// if the validator voted nil, record that status
		if vote.BlockIdFlag != cometproto.BlockIDFlagCommit {
			nilVote = true
		}
		// iterate over each currency-pair, and record whether the validator reported a price for it
		validatorPrices := h.pa.GetPricesForValidator(validator)
		for _, cp := range pricesToReport {
			// if the validator reported a nil-vote, record that and skip
			if nilVote {
				h.metrics.AddValidatorReportForTicker(validator.String(), cp, servicemetrics.Absent)
				continue
			}

			// otherwise, check if the validator reported a price for the currency-pair
			price, ok := validatorPrices[cp]
			if !ok || price == nil {
				h.metrics.AddValidatorReportForTicker(validator.String(), cp, servicemetrics.MissingPrice)
				continue
			}

			// if the validator reported a price, record that price
			floatPrice, _ := price.Float64()
			h.metrics.AddValidatorReportForTicker(validator.String(), cp, servicemetrics.WithPrice)
			h.metrics.AddValidatorPriceForTicker(validator.String(), cp, floatPrice)
		}
	}
}
