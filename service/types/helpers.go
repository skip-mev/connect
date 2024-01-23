package types

import (
	"math/big"

	"github.com/skip-mev/slinky/x/oracle/types"
)

func ToReqPrices(prices map[types.CurrencyPair]*big.Int) map[string]string {
	reqPrices := make(map[string]string, len(prices))

	for cp, price := range prices {
		reqPrices[cp.String()] = price.String()
	}

	return reqPrices
}
