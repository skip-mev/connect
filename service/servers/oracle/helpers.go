package oracle

import (
	"github.com/skip-mev/slinky/oracle/types"
)

func ToReqPrices(prices types.TickerPrices) map[string]string {
	reqPrices := make(map[string]string, len(prices))

	for cp, price := range prices {
		reqPrices[cp.String()] = price.String()
	}

	return reqPrices
}
