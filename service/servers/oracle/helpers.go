package oracle

import (
	"github.com/skip-mev/connect/v2/oracle/types"
)

func ToReqPrices(prices types.Prices) map[string]string {
	reqPrices := make(map[string]string, len(prices))

	for cp, price := range prices {
		intPrice, _ := price.Int(nil)
		reqPrices[cp] = intPrice.String()
	}

	return reqPrices
}
