package bybit_test

import (
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/skip-mev/slinky/providers/websockets/bybit"
)

var (
	goodConfig = `
{
	"supportedBases": {
		"BITCOIN": "BTC",
		"ETHEREUM": "ETH",
		"ATOM": "ATOM",
		"SOLANA": "SOL",
		"POLKADOT": "DOT",
		"DYDX": "DYDX"
	},
	"supportedQuotes": {
		"USD": "USD",
		"ETHEREUM": "ETH"
	}
}
	`
	noBasesConfig = `
{
	"supportedQuotes": {
		"USD": "USD",
		"ETHEREUM": "ETH"
	}
}
`
	noQuotesConfig = `
{
	"supportedBases": {
		"BITCOIN": "BTC",
		"ETHEREUM": "ETH",
		"ATOM": "ATOM",
		"SOLANA": "SOL",
		"POLKADOT": "DOT",
		"DYDX": "DYDX"
	}
}
	`
	malformedJSONConfig = `
{
	"supportedBases": {
		"BITCOIN": "BTC",
	},
	"supportedQuotes": {
		"USD": "USD",
	},
}
	`
	emptySupportedBaseKeyConfig = `
{
	"supportedBases": {
		"": "BTC"
	},
	"supportedQuotes": {
		"USD": "USD"
	}
}
`
	emptySupportedBaseValueConfig = `
{
	"supportedBases": {
		"BITCOIN": ""
	},
	"supportedQuotes": {
		"USD": "USD"
	}
}
`
	emptySupportedQuoteKeyConfig = `
{
	"supportedBases": {
		"BITCOIN": "BTC"
	},
	"supportedQuotes": {
		"": "USD"
	}
}
`
	emptySupportedQuoteValueConfig = `
{
	"supportedBases": {
		"BITCOIN": "BTC"
	},
	"supportedQuotes": {
		"USD": ""
	}
}
`
)

func TestReadConfigFromFile(t *testing.T) {
	testCases := []struct {
		name        string
		json        string
		expectedErr bool
	}{
		{
			name:        "good config",
			json:        goodConfig,
			expectedErr: false,
		},
		{
			name:        "no bases config",
			json:        noBasesConfig,
			expectedErr: true,
		},
		{
			name:        "no quotes config",
			json:        noQuotesConfig,
			expectedErr: true,
		},
		{
			name:        "malformed json config",
			json:        malformedJSONConfig,
			expectedErr: true,
		},
		{
			name:        "empty supported base key config",
			json:        emptySupportedBaseKeyConfig,
			expectedErr: true,
		},
		{
			name:        "empty supported base value config",
			json:        emptySupportedBaseValueConfig,
			expectedErr: true,
		},
		{
			name:        "empty supported quote key config",
			json:        emptySupportedQuoteKeyConfig,
			expectedErr: true,
		},
		{
			name:        "empty supported quote value config",
			json:        emptySupportedQuoteValueConfig,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temp file
			f, err := os.CreateTemp("", "bybit_config")
			assert.NoError(t, err)
			defer os.Remove(f.Name())

			// Write the config as a toml file
			_, err = f.WriteString(tc.json)
			assert.NoError(t, err)

			// Read config from file
			_, err = bybit.ReadConfigFromFile(f.Name())
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
