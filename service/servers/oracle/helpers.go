package oracle

import (
	"math/big"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
)

func ToReqPrices(prices map[slinkytypes.CurrencyPair]*big.Int) map[string]string {
	reqPrices := make(map[string]string, len(prices))

	for cp, price := range prices {
		reqPrices[cp.String()] = price.String()
	}

	return reqPrices
}
