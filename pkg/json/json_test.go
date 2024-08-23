package json_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/pkg/json"
)

func TestIsValid(t *testing.T) {
	type testCase struct {
		name      string
		bz        []byte
		expectErr bool
	}

	testCases := []testCase{
		{
			name:      "valid empty",
			bz:        []byte{},
			expectErr: false,
		},
		{
			name:      "valid basic JSON",
			bz:        []byte(`{"key": "value"}`),
			expectErr: false,
		},
		{
			name:      "invalid basic JSON missing quotation",
			bz:        []byte(`{"key": "value}`),
			expectErr: true,
		},
		{
			name: "valid JSON array",
			bz: []byte(`{
					"arr": [
						{
							"key1": "value1a",
							"key2": "value2a",
							"key3": "value3a"
						},
						{
							"key1": "value1b",
							"key2": "value2b",
							"key3": "value3b"
						}
					]
				}`),
			expectErr: false,
		},
		{
			name: "invalid JSON array - extra comma",
			bz: []byte(`{
					"arr": [
						{
							"key1": "value1a",
							"key2": "value2a",
							"key3": "value3a",
						},
						{
							"key1": "value1b",
							"key2": "value2b",
							"key3": "value3b"
						}
					]
				}`),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := json.IsValid(tc.bz)
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
