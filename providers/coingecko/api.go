package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"github.com/skip-mev/slinky/providers"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	apiKeyHeader = "x-cg-pro-api-key" // #nosec G101
	baseURL      = "https://api.coingecko.com/api/v3"
	apiURL       = "https://pro-api.coingecko.com/api/v3"

	pairPriceRequest = "/simple/price?ids=%s&vs_currencies=%s"
	precision        = "&precision=18"
)

// getPriceForPair returns the price of a currency pair. The price is fetched
// from the CoinGecko API in a single request for all pairs. Since the CoinGecko
// response will match some base denoms to quote denoms that should not be supported,
// we filter out pairs that are not supported by the provider.
//
// Response format:
//
//	{
//	  "cosmos": {
//	    "usd": 11.35
//	  },
//	  "bitcoin": {
//	    "usd": 10000
//	  }
//	}
func (p *Provider) getPrices(ctx context.Context) (map[oracletypes.CurrencyPair]*big.Int, error) {
	url := p.getPriceEndpoint(p.bases, p.quotes)

	// make the request to url and unmarshal the response into respMap
	respMap := make(map[string]map[string]float64)

	// if an API key is set, add it to the request
	var reqFn providers.ReqFn
	if p.config.APIKey != "" {
		reqFn = func(req *http.Request) {
			req.Header.Set(apiKeyHeader, p.config.APIKey)
		}
	}

	if err := providers.GetWithContextAndHeader(ctx, url, func(body []byte) error {
		return json.Unmarshal(body, &respMap)
	}, reqFn); err != nil {
		return nil, err
	}

	prices := make(map[oracletypes.CurrencyPair]*big.Int)

	for _, pair := range p.pairs {
		base := strings.ToLower(pair.Base)
		quote := strings.ToLower(pair.Quote)

		if _, ok := respMap[base]; !ok {
			continue
		}

		if _, ok := respMap[base][quote]; !ok {
			continue
		}

		prices[pair] = providers.Float64ToBigInt(respMap[base][quote], pair.Decimals())
	}

	return prices, nil
}

// getPriceEndpoint is the CoinGecko endpoint for getting the price of a
// currency pair.
func (p *Provider) getPriceEndpoint(base, quote string) string {
	if p.config.APIKey != "" {
		return fmt.Sprintf(apiURL+pairPriceRequest+precision, base, quote)
	}
	return fmt.Sprintf(baseURL+pairPriceRequest+precision, base, quote)
}
