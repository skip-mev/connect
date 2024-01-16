package cryptodotcom_test

import (
	"os"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/skip-mev/slinky/providers/websockets/cryptodotcom"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/require"
)

var (
	validJSON = `
{
    "markets": {
        "BITCOIN/USD": "BTCUSD-PERP",
        "ETHEREUM/USD": "ETHUSD-PERP",
        "SOLANA/USD": "SOLUSD-PERP"
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
		"BITCOIN/USD": "BTCUSD-PERP",
		"USD": "ETHUSD-PERP"
	},
	"production": true
}
`

	emptyMarketJSON = `
{
	"markets": {
		"BITCOIN/USD": "",
		"ETHEREUM/USD": "ETHUSD-PERP"
	},
	"production": true
}
`

	invalidJSON = `
{
	"markets": {
		"BITCOIN/USD": "BTCUSD-PERP",
	},
	"production": true
}
`

	duplicateMarketJSON = `
{
	"markets": {
		"BITCOIN/USD": "BTCUSD-PERP",
		"BITCOIN/USDT": "BTCUSD-PERP"
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
			f, err := os.CreateTemp("", "cryptodotcom_config")
			assert.NoError(t, err)
			defer os.Remove(f.Name())

			// Write the config as a toml file
			_, err = f.WriteString(tc.json)
			assert.NoError(t, err)

			// Read config from file
			_, err = cryptodotcom.ReadConfigFromFile(f.Name())
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
		config      cryptodotcom.Config
		expectedErr bool
	}{
		{
			name: "valid config",
			config: cryptodotcom.Config{
				Markets: map[string]string{
					"BITCOIN/USD":  "BTCUSD-PERP",
					"ETHEREUM/USD": "ETHUSD-PERP",
				},
				Production: true,
				Cache: map[oracletypes.CurrencyPair]string{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"):  "BTCUSD-PERP",
					oracletypes.NewCurrencyPair("ETHEREUM", "USD"): "ETHUSD-PERP",
				},
				ReverseCache: map[string]oracletypes.CurrencyPair{
					"BTCUSD-PERP": oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					"ETHUSD-PERP": oracletypes.NewCurrencyPair("ETHEREUM", "USD"),
				},
			},
			expectedErr: false,
		},
		{
			name: "missing currency pair in caches",
			config: cryptodotcom.Config{
				Markets: map[string]string{
					"BITCOIN/USD": "BTCUSD-PERP",
					"USD":         "ETHUSD-PERP",
				},
				Production: true,
			},
			expectedErr: true,
		},
		{
			name: "duplicate market",
			config: cryptodotcom.Config{
				Markets: map[string]string{
					"BITCOIN/USD":  "BTCUSD-PERP",
					"ETHEREUM/USD": "BTCUSD-PERP",
				},
				Production: true,
				Cache: map[oracletypes.CurrencyPair]string{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"):  "BTCUSD-PERP",
					oracletypes.NewCurrencyPair("ETHEREUM", "USD"): "BTCUSD-PERP",
				},
				ReverseCache: map[string]oracletypes.CurrencyPair{
					"BTCUSD-PERP": oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				},
			},
			expectedErr: true,
		},
		{
			name: "empty market",
			config: cryptodotcom.Config{
				Markets: map[string]string{
					"BITCOIN/USD":  "BTCUSD-PERP",
					"ETHEREUM/USD": "",
				},
				Production: true,
				Cache: map[oracletypes.CurrencyPair]string{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): "BTCUSD-PERP",
				},
				ReverseCache: map[string]oracletypes.CurrencyPair{
					"BTCUSD-PERP": oracletypes.NewCurrencyPair("BITCOIN", "USD"),
				},
			},
			expectedErr: true,
		},
		{
			name: "bad format for currency pair",
			config: cryptodotcom.Config{
				Markets: map[string]string{
					"BITCOIN/USD":  "BTCUSD-PERP",
					"ETHEREUM/USD": "ETHUSD-PERP",
				},
				Production: true,
				Cache: map[oracletypes.CurrencyPair]string{
					oracletypes.NewCurrencyPair("BITCOIN", "USD"): "BTCUSD-PERP",
					oracletypes.NewCurrencyPair("", "USD"):        "ETHUSD-PERP",
				},
				ReverseCache: map[string]oracletypes.CurrencyPair{
					"BTCUSD-PERP": oracletypes.NewCurrencyPair("BITCOIN", "USD"),
					"ETHUSD-PERP": oracletypes.NewCurrencyPair("", "USD"),
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
