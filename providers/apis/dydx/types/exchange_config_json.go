package types

// ExchangeConfigJson demarshals the exchange configuration json for a particular market.
// The result is a list of parameters that define how the market is resolved on
// each supported exchange.
//
// This struct stores data in an intermediate form as it's being assigned to various
// `ExchangeMarketConfig` objects, which are keyed by exchange id. These objects are not kept
// past the time the `GetAllMarketParams` API response is parsed, and do not contain an id
// because the id is expected to be known at the time the object is in use.
type ExchangeConfigJson struct { //nolint
	Exchanges []ExchangeMarketConfigJson `json:"exchanges"`
}

// ExchangeMarketConfigJson captures per-exchange information for resolving a market, including
// the ticker and conversion details. It demarshals JSON parameters from the chain for a
// particular market on a specific exchange.
type ExchangeMarketConfigJson struct { //nolint
	ExchangeName   string `json:"exchangeName"`
	Ticker         string `json:"ticker"`
	AdjustByMarket string `json:"adjustByMarket,omitempty"`
	Invert         bool   `json:"invert,omitempty"`
}
