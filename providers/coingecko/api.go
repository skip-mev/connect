package coingecko

import (
	"fmt"
)

const (
	apiKeyHeader = "x-cg-pro-api-key" // #nosec G101
	baseURL      = "https://api.coingecko.com/api/v3"
	apiURL       = "https://pro-api.coingecko.com/api/v3"

	pairPriceRequest = "/simple/price?ids=%s&vs_currencies=%s"
	precision        = "&precision=18"
)

// getPriceEndpoint is the CoinGecko endpoint for getting the price of a
// currency pair.
func (h *CoinGeckoAPIHandler) getPriceEndpoint(base, quote string) string {
	if h.config.APIKey != "" {
		return fmt.Sprintf(apiURL+pairPriceRequest+precision, base, quote)
	}
	return fmt.Sprintf(baseURL+pairPriceRequest+precision, base, quote)
}
