package okx_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/alecthomas/assert/v2"
	"github.com/skip-mev/slinky/providers/websockets/okx"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

var (
	validJSON = `
{
	"markets": {
		"BITCOIN/USD": "BTC-USD",
		"ETHEREUM/USD": "ETH-USD",
		"SOLANA/USD": "SOL-USD"
	},
	"production": true
}
	`

	emptyJSON = `
{
	"markets": {},
	"production": true
}
`

	invalidCPJSON = `
{
	"markets": {
		"BITCOIN/USD": "BTC-USD",
		"USD": "ETH-USD"
	},
	"production": true
}
`

	emptyMarketJSON = `
{
	"markets": {
		"BITCOIN/USD": "",
		"ETHEREUM/USD": "ETH-USD"
	},
	"production": true
}
`

	invalidJSON = `
{
	"markets": {
		"BITCOIN/USD": "BTC-USD",
	},
	"production": true
}
`

	duplicateMarketJSON = `
{
	"markets": {
		"BITCOIN/USD": "BTC-USD",
		"BITCOIN/USDT": "BTC-USD"
	},
	"production": true
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
			name:        "invalid currency pair",
			json:        invalidCPJSON,
			expectedErr: true,
		},
		{
			name:        "empty market",
			json:        emptyMarketJSON,
			expectedErr: true,
		},
		{
			name:        "invalid json",
			json:        invalidJSON,
			expectedErr: true,
		},
		{
			name:        "duplicate market",
			json:        duplicateMarketJSON,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temp file
			f, err := os.CreateTemp("", "okx_config")
			assert.NoError(t, err)
			defer os.Remove(f.Name())

			// Write the config as a toml file
			_, err = f.WriteString(tc.json)
			assert.NoError(t, err)

			// Read config from file
			_, err = okx.ReadConfigFromFile(f.Name())
			if tc.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateBsic(t *testing.T) {
	testCases := []struct {
		name        string
		config      okx.Config
		expectedErr bool
	}{
		{
			name: "valid config",
			config: okx.Config{
				Markets: map[string]string{
					"BITCOIN/USD":  "BTC-USD",
					"ETHEREUM/USD": "ETH-USD",
				},
				Production: true,
				Cache: map[oracletypes.CurrencyPair]string{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"):  "BTC-USD",
					oracletypes.NewCurrencyPair("ETHEREUM", "USD"): "ETH-USD",
				},
				ReverseCache: map[string]oracletypes.CurrencyPair{
					"BTC-USD": oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					"ETH-USD": oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
				},
			},
			expectedErr: false,
		},
		{
			name: "missing currency pair in caches",
			config: okx.Config{
				Markets: map[string]string{
					"BITCOIN/USD": "BTC-USD",
					"USD":         "ETH-USD",
				},
				Production: true,
			},
			expectedErr: true,
		},
		{
			name: "duplicate market",
			config: okx.Config{
				Markets: map[string]string{
					"BITCOIN/USD":  "BTC-USD",
					"ETHEREUM/USD": "BTC-USD",
				},
				Production: true,
				Cache: map[oracletypes.CurrencyPair]string{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"):  "BTC-USD",
					oracletypes.NewCurrencyPair("ETHEREUM", "USD"): "BTC-USD",
				},
				ReverseCache: map[string]oracletypes.CurrencyPair{
					"BTC-USD": oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				},
			},
			expectedErr: true,
		},
		{
			name: "empty market",
			config: okx.Config{
				Markets: map[string]string{
					"BITCOIN/USD":  "BTC-USD",
					"ETHEREUM/USD": "",
				},
				Production: true,
				Cache: map[oracletypes.CurrencyPair]string{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): "BTC-USD",
				},
				ReverseCache: map[string]oracletypes.CurrencyPair{
					"BTC-USD": oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				},
			},
			expectedErr: true,
		},
		{
			name: "bad format for currency pair",
			config: okx.Config{
				Markets: map[string]string{
					"BITCOIN/USD":  "BTC-USD",
					"ETHEREUM/USD": "ETH-USD",
				},
				Production: true,
				Cache: map[oracletypes.CurrencyPair]string{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): "BTC-USD",
					oracletypes.NewCurrencyPair("", "USD"):        "ETH-USD",
				},
				ReverseCache: map[string]oracletypes.CurrencyPair{
					"BTC-USD": oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					"ETH-USD": oracletypes.NewCurrencyPair("", "USD"),
				},
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.ValidateBasic()
			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
