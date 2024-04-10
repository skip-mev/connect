package json

import (
	"encoding/json"
	"fmt"
)

// IsValid checks if the given byte array is valid JSON.
// If the byte array is 0 length, this is a valid empty JSON object.
func IsValid(jsonBz []byte) error {
	if len(jsonBz) == 0 {
		return nil
	}

	var checkStruct map[string]interface{}
	if err := json.Unmarshal(jsonBz, &checkStruct); err != nil {
		return fmt.Errorf("unable to unmarshal string to json: %w", err)
	}

	return nil
}
