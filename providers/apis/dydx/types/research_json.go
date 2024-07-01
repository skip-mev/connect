package types

// MarketParamIndex is the index into the research json market-structure under which
// the market's parameters are stored.
const MarketParamIndex = "params"

type ResearchJSONMarketParam struct {
	// Id is the unique identifier for the market
	ID uint32 `json:"id"`

	// Pair is the ticker symbol for the market
	Pair string `json:"ticker"`

	// Exponent is the number of decimal places to shift the price by
	Exponent float64 `json:"priceExponent"`

	// MinExchanges is the minimum number of exchanges that must provide data for the market
	MinExchanges uint32 `json:"minExchanges"`

	// MinPriceChangePpm is the minimum price change that must be observed for the market
	MinPriceChangePpm uint32 `json:"minPriceChange"`

	// ExchangeConfigJSON is the json object that contains the exchange configuration for the market
	ExchangeConfigJSON []ExchangeMarketConfigJson `json:"exchangeConfigJson"`
}

// ResearchJSON is the go-struct that encompasses the dydx research json, as hosted
// on [github](https://raw.githubusercontent.com/dydxprotocol/v4-web/main/public/configs/otherMarketData.json)
type ResearchJSON map[string]Params

type Params struct {
	ResearchJSONMarketParam `json:"params"`
	MetaData                `json:"meta"`
}

type MetaData struct {
	CMCID int `json:"cmcId"`
}
