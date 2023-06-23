package types

import (
	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/oracle/types"
)

func ToReqPrices(prices map[types.CurrencyPair]*uint256.Int) map[string]string {
	reqPrices := make(map[string]string, len(prices))

	for k, v := range prices {
		reqPrices[k.String()] = v.String()
	}

	return reqPrices
}
