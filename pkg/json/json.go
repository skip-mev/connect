package json

import "encoding/json"

// IsValid checks if the given byte array is valid JSON.
// If the string is 0 length, this is a valid empty JSON object.
func IsValid(jsonBz []byte) error {
	if len(jsonBz) == 0 {
		return nil
	}

	var checkStruct map[string]interface{}
	return json.Unmarshal(jsonBz, &checkStruct)
}
