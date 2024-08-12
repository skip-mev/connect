package tickermetadata

import "encoding/json"

// CoreMetadata is the Ticker.Metadata_JSON published to every Ticker in the x/marketmap module on core markets.
type CoreMetadata struct {
	// AggregateIDs contains a list of AggregatorIDs associated with the ticker.
	// This field may not be populated if no aggregator currently indexes this Ticker.
	AggregateIDs []AggregatorID `json:"aggregate_ids"`
}

// NewCoreMetadata returns a new CoreMetadata instance.
func NewCoreMetadata(aggregateIDs []AggregatorID) CoreMetadata {
	return CoreMetadata{
		AggregateIDs: aggregateIDs,
	}
}

// MarshalCoreMetadata returns the JSON byte encoding of the CoreMetadata.
func MarshalCoreMetadata(m CoreMetadata) ([]byte, error) {
	return json.Marshal(m)
}

// CoreMetadataFromJSONString returns a CoreMetadata instance from a JSON string.
func CoreMetadataFromJSONString(jsonString string) (CoreMetadata, error) {
	var elem CoreMetadata
	err := json.Unmarshal([]byte(jsonString), &elem)
	return elem, err
}

// CoreMetadataFromJSONBytes returns a CoreMetadata instance from JSON bytes.
func CoreMetadataFromJSONBytes(jsonBytes []byte) (CoreMetadata, error) {
	var elem CoreMetadata
	err := json.Unmarshal(jsonBytes, &elem)
	return elem, err
}
