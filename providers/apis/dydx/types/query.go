package types

// QueryAllMarketParamsResponse is response type for the Query/Params
// `AllMarketParams` RPC method.
type QueryAllMarketParamsResponse struct {
	MarketParams []MarketParam `protobuf:"bytes,1,rep,name=market_params,json=marketParams,proto3" json:"market_params"`
}

// MarketParam represents the x/prices configuration for markets, including
// representing price values, resolving markets on individual exchanges, and
// generating price updates. This configuration is specific to the quote
// currency.
type MarketParam struct {
	// Unique, sequentially-generated value.
	Id uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"` //nolint
	// The human-readable name of the market pair (e.g. `BTC-USD`).
	Pair string `protobuf:"bytes,2,opt,name=pair,proto3" json:"pair,omitempty"`
	// Static value. The exponent of the price.
	// For example if `Exponent == -5` then a `Value` of `1,000,000,000`
	// represents “$10,000`. Therefore `10 ^ Exponent` represents the smallest
	// price step (in dollars) that can be recorded.
	Exponent int32 `protobuf:"zigzag32,3,opt,name=exponent,proto3" json:"exponent,omitempty"`
	// The minimum number of exchanges that should be reporting a live price for
	// a price update to be considered valid.
	MinExchanges uint32 `protobuf:"varint,4,opt,name=min_exchanges,json=minExchanges,proto3" json:"min_exchanges,omitempty"`
	// The minimum allowable change in `price` value that would cause a price
	// update on the network. Measured as `1e-6` (parts per million).
	MinPriceChangePpm uint32 `protobuf:"varint,5,opt,name=min_price_change_ppm,json=minPriceChangePpm,proto3" json:"min_price_change_ppm,omitempty"`
	// A string of json that encodes the configuration for resolving the price
	// of this market on various exchanges.
	ExchangeConfigJson string `protobuf:"bytes,6,opt,name=exchange_config_json,json=exchangeConfigJson,proto3" json:"exchange_config_json,omitempty"` //nolint
}
