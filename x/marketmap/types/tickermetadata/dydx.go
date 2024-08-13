package tickermetadata

import "encoding/json"

// DyDx is the Ticker.Metadata_JSON published to every Ticker in the x/marketmap module on dYdX.
type DyDx struct {
	// ReferencePrice gives a spot price for that Ticker at the point in time when the ReferencePrice was updated.
	// You should _not_ use this for up-to-date/instantaneous spot pricing data since it is updated infrequently.
	// The price is scaled by Ticker.Decimals.
	ReferencePrice uint64 `json:"reference_price"`
	// Liquidity gives a _rough_ estimate of the amount of liquidity in the Providers for a given Market.
	// It is _not_ updated in coordination with spot prices and only gives rough order of magnitude accuracy at the time
	// which the update for it is published.
	// The liquidity value stored here is USD denominated.
	Liquidity uint64 `json:"liquidity"`
	// AggregateIDs contains a list of AggregatorIDs associated with the ticker.
	// This field may not be populated if no aggregator currently indexes this Ticker.
	AggregateIDs []AggregatorID `json:"aggregate_ids"`
}

// NewDyDx returns a new DyDx instance.
func NewDyDx(referencePrice, liquidity uint64, aggregateIDs []AggregatorID) DyDx {
	return DyDx{
		ReferencePrice: referencePrice,
		Liquidity:      liquidity,
		AggregateIDs:   aggregateIDs,
	}
}

// MarshalDyDx returns the JSON byte encoding of the DyDx.
func MarshalDyDx(m DyDx) ([]byte, error) {
	return json.Marshal(m)
}

// DyDxFromJSONString returns a DyDx instance from a JSON string.
func DyDxFromJSONString(jsonString string) (DyDx, error) {
	var elem DyDx
	err := json.Unmarshal([]byte(jsonString), &elem)
	return elem, err
}

// DyDxFromJSONBytes returns a DyDx instance from JSON bytes.
func DyDxFromJSONBytes(jsonBytes []byte) (DyDx, error) {
	var elem DyDx
	err := json.Unmarshal(jsonBytes, &elem)
	return elem, err
}
