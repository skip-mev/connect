package static

import (
	"encoding/json"
	"fmt"
)

// MetaData is the per-ticker specific metadata that is used to configure the static provider.
type MetaData struct {
	Price float64 `json:"price"`
}

// FromJSON unmarshals the JSON data into a MetaData struct.
func (m *MetaData) FromJSON(jsonStr string) error {
	err := json.Unmarshal([]byte(jsonStr), m)
	return err
}

// MustToJSON marshals the MetaData struct into a JSON string.
func (m *MetaData) MustToJSON() string {
	bz, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Errorf("failed to marshal metadata: %w", err))
	}
	return string(bz)
}
