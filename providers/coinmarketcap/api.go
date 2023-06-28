package coinmarketcap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/providers"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	headerFieldKey = "X-CMC_PRO_API_KEY"
)

// getPriceForPair gets the the price of base in terms of quote, and returns the response scaled by cp.Decimals
// Response format:
//
//	{
//		"status": {
//		  "timestamp": "2023-06-28T23:17:55.178Z",
//	    ...
//		},
//		"data": {
//		  "BTC": [
//			{
//			  ...
//			  "quote": {
//				"USD": {
//				  "price": 30172.950080795657,
//				  ...
//				}
//			  }
//			}
//		  ]
//		}
//	  }
func (p *Provider) getPriceForPair(ctx context.Context, pair oracletypes.CurrencyPair) (types.QuotePrice, error) {
	p.logger.Info("Fetching price for pair", "pair", pair)

	// make request to coinmarketcap api w/ X-CMC_PRO_API_KEY header set to api-key
	base, quote := p.getSymbolForTokenName(pair.Base), p.getSymbolForTokenName(pair.Quote)
	p.logger.Info("Fetching price for pair", "pair", pair, "base", base, "quote", quote, "apikey", p.apiKey)

	// make request to coinmarketcap api w/ X-CMC_PRO_API_KEY header set to api-key
	var resp map[string]interface{}
	if err := providers.GetWithContextAndHeader(
		ctx,
		getPriceEndpoint(base, quote),

		func(body []byte) error {
			return json.Unmarshal(body, &resp)
		},

		func(req *http.Request) {
			req.Header.Add(headerFieldKey, p.apiKey)
		},
	); err != nil {
		return types.QuotePrice{}, err
	}

	// unmarshal request body to get price
	price, err := unmarshalRequest(resp, pair.Base, pair.Quote)
	if err != nil {
		return types.QuotePrice{}, err
	}

	return types.QuotePrice{
		Price:     providers.Float64ToUint256(price, pair.Decimals()),
		Timestamp: time.Now(),
	}, nil
}

// getPriceEndpoint returns the endpoint to fetch prices from.
func getPriceEndpoint(base, quote string) string {
	return fmt.Sprintf("https://pro-api.coinmarketcap.com/v2/cryptocurrency/quotes/latest?symbol=%s&convert=%s", base, quote)
}

// unmarshalRequest unmarshals the response from the coinmarketcap api, and retrieves the units of quote for a unit of base.
func unmarshalRequest(resp map[string]interface{}, base, quote string) (float64, error) {
	// get the price + data
	data, ok := resp["data"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("failed to get data")
	}

	// get the base currency
	baseDataArr, ok := data[strings.ToUpper(base)].([]interface{})
	if !ok || len(baseDataArr) == 0 {
		return 0, fmt.Errorf("failed to get base data")
	}

	// get the first element of the base-data array
	baseData, ok := baseDataArr[0].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("failed to get base data")
	}

	// get the quote data
	quoteData, ok := baseData["quote"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("failed to get quote data")
	}

	// get the quote data for the quote currency
	quoteData, ok = quoteData[strings.ToUpper(quote)].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("failed to get quote data")
	}

	// get the price
	price, ok := quoteData["price"].(float64)
	if !ok {
		return 0, fmt.Errorf("failed to get price")
	}

	return price, nil
}
