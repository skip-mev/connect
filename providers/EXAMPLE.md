# Provider Walkthrough

Let's walk through an example provider to see how it is implemented.

## Binance API Provider

The [Binance API provider](./apis/binance/README.md) can be used as a reference for how to implement an API provider.

### Implementation

The package is laid out as follows:

* **binance/**: Contains the source code of the package.
  * `api_handler.go`: Main implementation of the `PriceAPIDataHandler` interface.
  * `utils.go`: Helper functions and configuration of the expected types to receive from the Binance API.


The logic below defines our implementation of the Binance API provider and its constructor.

```go
// APIHandler implements the PriceAPIDataHandler interface for Binance.
// for more information about the Binance API, refer to the following link:
// https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#public-api-endpoints
type APIHandler struct {
	// api is the config for the Binance API.
	api config.APIConfig
	// cache maintains the latest set of tickers seen by the handler.
	cache types.ProviderTickers
}

// NewAPIHandler returns a new Binance PriceAPIDataHandler.
func NewAPIHandler(
	api config.APIConfig,
) (types.PriceAPIDataHandler, error) {
	if api.Name != Name {
		return nil, fmt.Errorf("expected api config name %s, got %s", Name, api.Name)
	}

	if !api.Enabled {
		return nil, fmt.Errorf("api config for %s is not enabled", Name)
	}

	if err := api.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid api config for %s: %w", Name, err)
	}

	return &APIHandler{
		api:   api,
		cache: types.NewProviderTickers(),
	}, nil
}
```

The `CreateURL()` function implementation for the Binance APIHandler creates the URL needed for querying
the Binance API.  The URL is created by appending a list of desired tickers to the base URL of `https://api.binance.com/api/v3/ticker/price?symbols`.
The URL and logic for appending ticker IDs were chosen based on the [Binance API documentation](https://binance-docs.github.io/apidocs/spot/en/#symbol-price-ticker).

```go
// CreateURL returns the URL that is used to fetch data from the Binance API for the
// given tickers.
func (h *APIHandler) CreateURL(
	tickers []types.ProviderTicker,
) (string, error) {
	var tickerStrings string
	for _, ticker := range tickers {
		tickerStrings += fmt.Sprintf("%s%s%s%s", Quotation, ticker.GetOffChainTicker(), Quotation, Separator)
		h.cache.Add(ticker)
	}

	if len(tickerStrings) == 0 {
		return "", fmt.Errorf("empty url created. invalid or no ticker were provided")
	}

	return fmt.Sprintf(
		h.api.Endpoints[0].URL,
		LeftBracket,
		strings.TrimSuffix(tickerStrings, Separator),
		RightBracket,
	), nil
}
```

The `ParseResponse()` function implementation for the Binance APIHandler handles the response returned from the Binance API.
The function:
* Decodes the response to a known type (using JSON parsing)
* Resolves the response tickers to the requested tickers
* Converts the returned price for each ticker to a `*big.Float` for internal oracle use

```go
// ParseResponse parses the response from the Binance API and returns a GetResponse. Each
// of the tickers supplied will get a response or an error.
func (h *APIHandler) ParseResponse(
    tickers []types.ProviderTicker,
    resp *http.Response,
) types.PriceResponse {
    // Parse the response into a BinanceResponse.
    result, err := Decode(resp)
    if err != nil {
        return types.NewPriceResponseWithErr(
            tickers,
            providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToDecode),
        )
    }

    var (
        resolved   = make(types.ResolvedPrices)
        unresolved = make(types.UnResolvedPrices)
    )

    for _, data := range result {
        // Filter out the responses that are not expected.
        ticker, ok := h.cache.FromOffChainTicker(data.Symbol)
        if !ok {
            continue
        }

        price, err := math.Float64StringToBigFloat(data.Price)
        if err != nil {
			wErr := fmt.Errorf("failed to convert price %s to big.Float: %w", data.Price, err)
			unresolved[ticker] = providertypes.UnresolvedResult{
                ErrorWithCode: providertypes.NewErrorWithCode(wErr, providertypes.ErrorFailedToParsePrice),
            }
            continue
        }

        resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
    }

    // Add currency pairs that received no response to the unresolved map.
    for _, ticker := range tickers {
        _, resolvedOk := resolved[ticker]
        _, unresolvedOk := unresolved[ticker]

        if !resolvedOk && !unresolvedOk {
            unresolved[ticker] = providertypes.UnresolvedResult{
                ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no response"), providertypes.ErrorNoResponse),
            }
        }
    }

    return types.NewPriceResponse(resolved, unresolved)
}
```

### Wiring to the Oracle

To wire the provider to the oracle, it must be added to the corresponding [Oracle Factory](./factories/README.md).
Here, the provider will be added to `APIQueryHandlerFactory` since we are making an API provider.  The factories
are pre-wired to the oracle, so when it starts up, it will now have registered the given new provider.