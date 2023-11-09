package sla

import (
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	slatypes "github.com/skip-mev/slinky/x/sla/types"
)

// getStatuses returns the price feed status updates for each currency pair.
func getStatuses(currencyPairs []oracletypes.CurrencyPair, prices map[string]string) map[oracletypes.CurrencyPair]slatypes.UpdateStatus {
	validatorUpdates := make(map[oracletypes.CurrencyPair]slatypes.UpdateStatus)

	for _, cp := range currencyPairs {
		if _, ok := prices[cp.String()]; !ok {
			validatorUpdates[cp] = slatypes.VoteWithoutPrice
		} else {
			validatorUpdates[cp] = slatypes.VoteWithPrice
		}
	}

	return validatorUpdates
}
