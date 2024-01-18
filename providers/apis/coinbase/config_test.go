package coinbase_test

import (
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/skip-mev/slinky/providers/apis/coinbase"
)

var (
	validJSON = `
{
	"symbolMap": {
		"BITCOIN": "BTC",
		"USD": "USD",
		"ETHEREUM": "ETH",
		"ATOM": "ATOM",
		"SOLANA": "SOL",
		"POLKADOT": "DOT",
		"DYDX": "DYDX"
	}
}
`
	emptyJSON = `
{
	"symbolMap": {}
}
`
	invalidJSON = `
{
	"symbol_map": {
		"BITCOIN": "BTC",
		"USD": "USD"
	}
}
`
	emptyKeyJSON = `
{
	"symbolMap": {
		"": "BTC",
		"USD": "USD"
	}
}
`
	emptyValueJSON = `
{
	"symbolMap": {
		"BITCOIN": "",
		"USD": "USD"
	}
}
`
)

func TestReadCoinbaseConfigFromFile(t *testing.T) {
	testCases := []struct {
		name        string
		json        string
		expectedErr bool
	}{
		{
			name:        "valid json",
			json:        validJSON,
			expectedErr: false,
		},
		{
			name:        "empty json",
			json:        emptyJSON,
			expectedErr: true,
		},
		{
			name:        "invalid json",
			json:        invalidJSON,
			expectedErr: true,
		},
		{
			name:        "empty key json",
			json:        emptyKeyJSON,
			expectedErr: true,
		},
		{
			name:        "empty value json",
			json:        emptyValueJSON,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temp file
			f, err := os.CreateTemp("", "coinbase_config")
			assert.NoError(t, err)
			defer os.Remove(f.Name())

			// Write the config as a toml file
			_, err = f.WriteString(tc.json)
			assert.NoError(t, err)

			// Read config from file
			_, err = coinbase.ReadCoinbaseConfigFromFile(f.Name())
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
