package binanceus

// NOTE: All the documentation for this file can be located on the Binance github
// API documentation: https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#symbol-price-ticker. This
// API does not require a subscription to use (i.e. No API key is required).

const (
	// BaseURL is the base URL of the Binance US API. This includes the base and quote
	// currency pairs that need to be inserted into the URL.
	BaseURL = "https://api.binance.us/api/v3/ticker/price?symbols=%s%s%s"
)
