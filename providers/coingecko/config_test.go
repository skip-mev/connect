package coingecko_test

import (
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"

	"github.com/skip-mev/slinky/providers/coingecko"
)

var (
	goodConfig = `
{
	"apiKey": "",
	"supportedBases": {
		"BITCOIN": "bitcoin",
		"ETHEREUM": "ethereum",
		"ATOM": "cosmos",
		"SOLANA": "solana",
		"POLKADOT": "polkadot",
		"DYDX": "dydx-chain"
	},
	"supportedQuotes": {
		"USD": "usd",
		"ETHEREUM": "eth"
	}
}
	`
	noBasesConfig = `
{
	"apiKey": "",
	"supportedQuotes": {
		"USD": "usd",
		"ETHEREUM": "eth"
	}
}
`
	noQuotesConfig = `
{
	"apiKey": "",
	"supportedBases": {
		"BITCOIN": "bitcoin",
		"ETHEREUM": "ethereum",
		"ATOM": "cosmos",
		"SOLANA": "solana",
		"POLKADOT": "polkadot",
		"DYDX": "dydx-chain"
	}
}
	`
	malformedJSONConfig = `
{
	"apiKey": "",
	"supportedBases": {
		"BITCOIN": "bitcoin",
	},
	"supportedQuotes": {
		"USD": "usd",
	},
}
	`
	emptySupportedBaseKeyConfig = `
{
	"apiKey": "",
	"supportedBases": {
		"": "bitcoin"
	},
	"supportedQuotes": {
		"USD": "usd"
	}
}
`
	emptySupportedBaseValueConfig = `
{
	"apiKey": "",
	"supportedBases": {
		"BITCOIN": ""
	},
	"supportedQuotes": {
		"USD": "usd"
	}
}
`
	emptySupportedQuoteKeyConfig = `
{
	"apiKey": "",
	"supportedBases": {
		"BITCOIN": "bitcoin"
	},
	"supportedQuotes": {
		"": "usd"
	}
}
`
	emptySupportedQuoteValueConfig = `
{
	"apiKey": "",
	"supportedBases": {
		"BITCOIN": "bitcoin"
	},
	"supportedQuotes": {
		"USD": ""
	}
}
`
)

func TestReadCoinGeckoConfigFromFile(t *testing.T) {
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
			f, err := os.CreateTemp("", "coingecko_config")
			assert.NoError(t, err)
			defer os.Remove(f.Name())

			// Write the config as a toml file
			_, err = f.WriteString(tc.json)
			assert.NoError(t, err)

			// Read config from file
			_, err = coingecko.ReadCoinGeckoConfigFromFile(f.Name())
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
