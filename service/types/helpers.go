package types

import (
	"github.com/holiman/uint256"

	"github.com/skip-mev/slinky/x/oracle/types"
)

func ToReqPrices(prices map[types.CurrencyPair]*uint256.Int) map[string]string {
	reqPrices := make(map[string]string, len(prices))

	for cp, price := range prices {
		reqPrices[cp.ToString()] = price.String()
	}

	return reqPrices
}
