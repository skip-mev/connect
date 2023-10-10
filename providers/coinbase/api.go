package coinbase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/skip-mev/slinky/aggregator"
	"github.com/skip-mev/slinky/providers"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// NameToSymbol is a map of currency names to their symbols.
var NameToSymbol = map[string]string{
	"BITCOIN":  "BTC",
	"COSMOS":   "ATOM",
	"ETHEREUM": "ETH",
	"USD":      "USD",
	"POLKADOT": "DOT",
	"POLYGON":  "MATIC",
}

// getPriceForPair returns the spot price of a currency pair. In practice,
// this should not be used because price data should come from an aggregated
// price feed - API that uses a TWAP, TVWAP, or median price.
//
// Response format:
//
//	{
//	  "data": {
//	    "amount": "1020.25",
//	    "currency": "USD"
//	  }
//	}
func getPriceForPair(ctx context.Context, pair oracletypes.CurrencyPair) (*aggregator.QuotePrice, error) {
	baseSymbol, ok := NameToSymbol[pair.Base]
	if !ok {
		return nil, fmt.Errorf("invalid base currency %s", pair.Base)
	}

	quoteSymbol, ok := NameToSymbol[pair.Quote]
	if !ok {
		return nil, fmt.Errorf("invalid quote currency %s", pair.Quote)
	}

	url := getSpotPriceEndpoint(baseSymbol, quoteSymbol)

	// make the request to url and unmarshal the response into respMap
	respMap := make(map[string]map[string]string)

	if err := providers.GetWithContext(ctx, url, func(body []byte) error {
		return json.Unmarshal(body, &respMap)
	}); err != nil {
		return nil, err
	}

	data, ok := respMap["data"]
	if !ok {
		return nil, fmt.Errorf("failed to parse response")
	}

	amount, ok := data["amount"]
	if !ok {
		return nil, fmt.Errorf("failed to parse response")
	}

	price, err := providers.Float64StringToUint256(amount, pair.Decimals())
	if err != nil {
		return nil, err
	}

	return &aggregator.QuotePrice{
		Price:     price,
		Timestamp: time.Now(),
	}, nil
}

// getSpotPriceEndpoint is the Coinbase endpoint for getting the spot price of a
// currency pair.
func getSpotPriceEndpoint(base, quote string) string {
	return fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-%s/spot", base, quote)
}
