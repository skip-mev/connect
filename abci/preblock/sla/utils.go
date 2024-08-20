package sla

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/abci/strategies/currencypair"
	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

// getStatuses returns the price feed status updates for each currency pair.
func getStatuses(ctx sdk.Context, currencyPairIDStrategy currencypair.CurrencyPairStrategy, currencyPairs []slinkytypes.CurrencyPair, prices map[uint64][]byte) map[slinkytypes.CurrencyPair]slatypes.UpdateStatus {
	validatorUpdates := make(map[slinkytypes.CurrencyPair]slatypes.UpdateStatus)

	for _, cp := range currencyPairs {
		id, err := currencyPairIDStrategy.ID(ctx, cp)
		if err != nil {
			continue
		}

		if _, ok := prices[id]; !ok {
			validatorUpdates[cp] = slatypes.VoteWithoutPrice
		} else {
			validatorUpdates[cp] = slatypes.VoteWithPrice
		}
	}

	return validatorUpdates
}
