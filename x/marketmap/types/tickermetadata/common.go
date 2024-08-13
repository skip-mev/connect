package tickermetadata

import "encoding/json"

type AggregatorID struct {
	// Venue is the name of the aggregator for which the ID is valid.
	// E.g. `coingecko`, `cmc`
	Venue string `json:"venue"`
	// ID is the string ID of the Ticker's Base denom in the aggregator.
	ID string `json:"ID"`
}

// NewAggregatorID returns a new AggregatorID instance.
func NewAggregatorID(venue, id string) AggregatorID {
	return AggregatorID{
		Venue: venue,
		ID:    id,
	}
}

// MarshalAggregatorID returns the JSON byte encoding of the AggregatorID.
func MarshalAggregatorID(m AggregatorID) ([]byte, error) {
	return json.Marshal(m)
}

// AggregatorIDFromJSONString returns an AggregatorID instance from a JSON string.
func AggregatorIDFromJSONString(jsonString string) (AggregatorID, error) {
	var elem AggregatorID
	err := json.Unmarshal([]byte(jsonString), &elem)
	return elem, err
}

// AggregatorIDFromJSONBytes returns an AggregatorID instance from JSON bytes.
func AggregatorIDFromJSONBytes(jsonBytes []byte) (AggregatorID, error) {
	var elem AggregatorID
	err := json.Unmarshal(jsonBytes, &elem)
	return elem, err
}
