package oracle

import (
	"math/big"

	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

func ToReqPrices(prices map[mmtypes.Ticker]*big.Int) map[string]string {
	reqPrices := make(map[string]string, len(prices))

	for cp, price := range prices {
		reqPrices[cp.String()] = price.String()
	}

	return reqPrices
}
