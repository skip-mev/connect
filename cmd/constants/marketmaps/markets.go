package marketmaps

import (
	"encoding/json"
	"errors"
	"fmt"

	mmtypes "github.com/skip-mev/connect/v2/x/marketmap/types"
)

var (
	// CoinMarketCapMarketMap is used to initialize the CoinMarketCap market map. This only includes
	// the markets that are supported by CoinMarketCap.
	CoinMarketCapMarketMap mmtypes.MarketMap
	// CoinMarketCapMarketMapJSON is the JSON representation of the CoinMarketCap MarketMap that can be used
	// to initialize for a genesis state or used by the sidecar as as static market map.
	CoinMarketCapMarketMapJSON = `
{
    "markets": {
	  "W/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "W",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "29587"
          }
        ]
      },
	  "TON/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "TON",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "11419"
          }
        ]
      },
	  "ZRO/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ZRO",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "26997"
          }
        ]
      },
	  "CHZ/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "CHZ",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "4066"
          }
        ]
      },
	  "ZK/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ZK",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "24091"
          }
        ]
      },
	  "BODEN/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "BODEN",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "29687"
          }
        ]
      },
	  "ETHFI/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ETHFI",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "29814"
          }
        ]
      },
      "KHAI/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "KHAI",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "30948"
          }
        ]
      },
      "WAFFLES/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "WAFFLES",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "31442"
          }
        ]
      },
      "HEGE/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "HEGE",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "31044"
          }
        ]
      },
      "WUF/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "WUF",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "30683"
          }
        ]
      },
      "CHAT/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "CHAT",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "29478"
          }
        ]
      },
      "BEER/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "BEER",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "31337"
          }
        ]
      },
      "MANEKI/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "MANEKI",
            "Quote": "USD"
          },
          "decimals": 11,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "30912"
          }
        ]
      },
      "SLERF/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "SLERF",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "29920"
          }
        ]
      },
      "MYRO/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "MYRO",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "28382"
          }
        ]
      },
      "RAY/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "RAY",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "8526"
          }
        ]
      },
      "WIF/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "WIF",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "28752"
          }
        ]
      },
      "MICHI/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "MICHI",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "30943"
          }
        ]
      },
      "MEW/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "MEW",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "30126"
          }
        ]
      },
      "PONKE/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "PONKE",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "29150"
          }
        ]
      },
      "BOME/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "BOME",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "29870"
          }
        ]
      },
      "DJT/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "DJT",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "31891"
          }
        ]
      },
      "POPCAT/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "POPCAT",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "28782"
          }
        ]
      },
      "AAVE/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "AAVE",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "7278"
          }
        ]
      },
      "ADA/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ADA",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "2010"
          }
        ]
      },
      "AEVO/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "AEVO",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "29676"
          }
        ]
      },
      "ALGO/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ALGO",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "4030"
          }
        ]
      },
      "APE/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "APE",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "18876"
          }
        ]
      },
      "APT/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "APT",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "21794"
          }
        ]
      },
      "ARB/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ARB",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "11841"
          }
        ]
      },
      "ARKM/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ARKM",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "27565"
          }
        ]
      },
      "ASTR/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ASTR",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "12885"
          }
        ]
      },
      "ATOM/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ATOM",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "3794"
          }
        ]
      },
      "AVAX/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "AVAX",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "5805"
          }
        ]
      },
      "AXL/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "AXL",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "17799"
          }
        ]
      },
      "BCH/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "BCH",
            "Quote": "USD"
          },
          "decimals": 7,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "1831"
          }
        ]
      },
      "BLUR/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "BLUR",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "23121"
          }
        ]
      },
      "BNB/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "BNB",
            "Quote": "USD"
          },
          "decimals": 7,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "1839"
          }
        ]
      },
      "BONK/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "BONK",
            "Quote": "USD"
          },
          "decimals": 14,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "23095"
          }
        ]
      },
      "BTC/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "BTC",
            "Quote": "USD"
          },
          "decimals": 5,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "1"
          }
        ]
      },
      "BUBBA/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "BUBBA",
            "Quote": "USD"
          },
          "decimals": 12,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "31411"
          }
        ]
      },
      "COMP/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "COMP",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "5692"
          }
        ]
      },
      "CRV/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "CRV",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "6538"
          }
        ]
      },
      "DOGE/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "DOGE",
            "Quote": "USD"
          },
          "decimals": 11,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "74"
          }
        ]
      },
      "DOT/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "DOT",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "6636"
          }
        ]
      },
      "DYDX/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "DYDX",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "28324"
          }
        ]
      },
      "DYM/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "DYM",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "28932"
          }
        ]
      },
      "TREMP/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "TREMP",
            "Quote": "USD"
          },
          "decimals": 11,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "29717"
          }
        ]
      },
      "MOG/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "MOG",
            "Quote": "USD"
          },
          "decimals": 11,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "27659"
          }
        ]
      },
      "MOTHER/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "MOTHER",
            "Quote": "USD"
          },
          "decimals": 11,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "31510"
          }
        ]
      },
      "EOS/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "EOS",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "1765"
          }
        ]
      },
      "ETC/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ETC",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "1321"
          }
        ]
      },
      "ETH/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ETH",
            "Quote": "USD"
          },
          "decimals": 6,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "1027"
          }
        ]
      },
      "FET/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "FET",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "3773"
          }
        ]
      },
      "FIL/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "FIL",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "2280"
          }
        ]
      },
      "GRT/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "GRT",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "6719"
          }
        ]
      },
      "HBAR/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "HBAR",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "4642"
          }
        ]
      },
      "ICP/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ICP",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "8916"
          }
        ]
      },
      "IMX/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "IMX",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "10603"
          }
        ]
      },
      "INJ/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "INJ",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "7226"
          }
        ]
      },
      "JTO/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "JTO",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "28541"
          }
        ]
      },
      "JUP/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "JUP",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "29210"
          }
        ]
      },
      "LDO/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "LDO",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "8000"
          }
        ]
      },
      "LINK/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "LINK",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "1975"
          }
        ]
      },
      "LTC/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "LTC",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "2"
          }
        ]
      },
      "MANA/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "MANA",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "1966"
          }
        ]
      },
      "POL/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "POL",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "28321"
          }
        ]
      },
      "MKR/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "MKR",
            "Quote": "USD"
          },
          "decimals": 6,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "1518"
          }
        ]
      },
      "NEAR/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "NEAR",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "6535"
          }
        ]
      },
      "NTRN/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "NTRN",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "26680"
          }
        ]
      },
      "OP/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "OP",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "11840"
          }
        ]
      },
      "ORDI/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "ORDI",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "25028"
          }
        ]
      },
      "PEPE/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "PEPE",
            "Quote": "USD"
          },
          "decimals": 16,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "24478"
          }
        ]
      },
      "PYTH/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "PYTH",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "28177"
          }
        ]
      },
      "RUNE/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "RUNE",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "4157"
          }
        ]
      },
      "SEI/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "SEI",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "23149"
          }
        ]
      },
      "SHIB/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "SHIB",
            "Quote": "USD"
          },
          "decimals": 15,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "5994"
          }
        ]
      },
      "SNX/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "SNX",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "2586"
          }
        ]
      },
      "SOL/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "SOL",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "5426"
          }
        ]
      },
      "STRK/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "STRK",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "22691"
          }
        ]
      },
      "STX/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "STX",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "4847"
          }
        ]
      },
      "SUI/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "SUI",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "20947"
          }
        ]
      },
      "TIA/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "TIA",
            "Quote": "USD"
          },
          "decimals": 8,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "22861"
          }
        ]
      },
      "TRX/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "TRX",
            "Quote": "USD"
          },
          "decimals": 11,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "1958"
          }
        ]
      },
      "UNI/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "UNI",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "7083"
          }
        ]
      },
      "USDT/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "USDT",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "825"
          }
        ]
      },
      "WLD/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "WLD",
            "Quote": "USD"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "13502"
          }
        ]
      },
      "WOO/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "WOO",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "7501"
          }
        ]
      },
      "XLM/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "XLM",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "512"
          }
        ]
      },
      "XRP/USD": {
        "ticker": {
          "currency_pair": {
            "Base": "XRP",
            "Quote": "USD"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": true
        },
        "provider_configs": [
          {
            "name": "coinmarketcap_api",
            "off_chain_ticker": "52"
          }
        ]
      }
    }
}
	`

	// RaydiumMarketMap is used to initialize the Raydium market map. This only includes
	// the markets that are supported by Raydium.
	RaydiumMarketMap mmtypes.MarketMap
	// RaydiumMarketMapJSON is the JSON representation of the Raydium MarketMap that can be used
	// to initialize for a genesis state or used by the sidecar as as static market map.
	RaydiumMarketMapJSON = `
	{
		"markets": {
		  "SOL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SOL",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 2,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "SOL-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "SOLUSD"
			  }
			]
		  },
		  "TRUMP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TRUMP",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "TRUMP/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"GwUn3JYQ5PC6P9cB9HsatJN7y7aw3BXrNqHPpoHtWyKF\",\"token_decimals\":8},\"quote_token_vault\":{\"token_vault_address\":\"5DNPt6WYgj3x7EmU4Fyqe3jDYPk2HMAB21H5N4Ggbev9\",\"token_decimals\":9},\"amm_info_address\":\"7Lco4QdQLaW6M4sxVhWe8BHjrykyzjcjGTo4a6qYGABK\",\"open_orders_address\":\"FAWLdBB8kmWZQ74KpYAYN3YaEW31Si8qrwuQPauFSoma\"}"
			  }
			]
		  },
		  "BAZINGA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BAZINGA",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/BAZINGA",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"invert": true,
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"GHVSjfv2kEpiMSTXsxP1S9KZKNzqa4rG8u3qVGVvNiEU\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"GKpbD62S56ZFtTBR5C1tJE4ZgaPesu5jhuLhkY3BXXKb\",\"token_decimals\":6},\"amm_info_address\":\"BhQgvhYpYVccRt5wJnxi13waXNaC3dJVcX6TjTNY9kee\",\"open_orders_address\":\"DbK9zkkFDh9aHfV3TVbbDjrXtFtdecSbsGwfofW4KzvC\"}"
			  }
			]
		  },
		  "BENDOG/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BENDOG",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "BENDOG/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"2Pza1YUczgc4RWLhAgdXSJh4oYUspvhhAiSecFDd7ZJ3\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"2BFpMzi33JtpY4CGUjY7x5JPApy6f2AdkuLZsd1QGqRv\",\"token_decimals\":9},\"amm_info_address\":\"47857wX96Tb4Ud3M3ka949iVRFmUqS33KLBxoVsqgfLK\",\"open_orders_address\":\"H1FPc9WQpA3GPnXMmzSjtt6gMuYuyDqYndBscaHNyCbv\"}"
			  }
			]
		  },
		  "BUBBA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BUBBA",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "BUBBA/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9sDmGRF3AmwrwAh75vne5dx7aVt7aKyHFeqLNP4ecXwh\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"J4Tisz2Ca8baDwvnrG5zNJBmE7U9DP47hgw1BdVUj7M1\",\"token_decimals\":9},\"amm_info_address\":\"8xtsif8mhNfpiHg3QNk24NW6X6wosoWWbCAhRTiUGW2n\",\"open_orders_address\":\"Dxjq2KVoJsrsHPViZyZWT89BDQexqQuC5uaEXd4ugqRG\"}"
			  }
			]
		  },
		  "CATGPT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CATGPT",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "CATGPT/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"AYAmCRPotwZprbNpPQ1hVGSEpbgWUgWHUbjnjt4bfLo1\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"ok2NPhxx2q3tn8XL289m5irCyGntURLQNhtisLowQ7y\",\"token_decimals\":9},\"amm_info_address\":\"92NvJRnTxkaiHcfRd72B8h1SHyj5ZGtMoeFAQvCdB3vB\",\"open_orders_address\":\"DXA5jH1r5c9QeAZxAYQb6emFGY2eBb3ZgMjSDuTNFZ6n\"}"
			  }
			]
		  },
		  "CHEESE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CHEESE",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/CHEESE",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"invert": true,
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"ELVeWgH6XnWtkLZHWWTPV2RRLCTHbVR6hCu4XMWbtS9M\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"5HmRek4a2d43r3Hot18AfTTjkRwgHwWQf6qvSDgWwSju\",\"token_decimals\":6},\"amm_info_address\":\"C4ZHt1fPtb6CLcUkivhnnNtxBfxYoJq6x8HEZpUexQvR\",\"open_orders_address\":\"7uBMeaAWzhfdHFnd9QVNS274wMmfmZYxmjtTnGWUf32A\"}"
			  }
			]
		  },
		  "CHUD/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CHUD",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/CHUD",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"invert": true,
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"4ESAztYfBqzZzqqLzKxzB7RDuQMi5Ho2Vrfjb1AdHwiG\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"BvcAm78ro36qESWkQorH1qnnJP21DJR2yaKwbsTNnMqx\",\"token_decimals\":6},\"amm_info_address\":\"3Cv7Z8KN21M6Ur6nTeKiKEPcjTuEAg87Ciu6cC3gRnw6\",\"open_orders_address\":\"3xr3BZ5EDnFiA2B2PGEo87nC9gDDUBERZPzVHHqQDPRw\"}"
			  }
			]
		  },
		  "COK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "COK",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "COK/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"C4hk6k49gotrWP1b9j2ejPcPo4Lq59jVmfGwB2YYYGds\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"6mk1jhhWr6yeYxQkcrkia2wLHFyuy1LW6Xmj2MmwJ2x5\",\"token_decimals\":9},\"amm_info_address\":\"1D5GHSzrcaSXLtUYxSCg4vWHdKGd7hFnasYPiPFYFGX\",\"open_orders_address\":\"F8F7FGDKfqVEC4qpnVjigZHB8kijTx8qqpmc1fX8s1dY\"}"
			  }
			]
		  },
		  "COST/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "COST",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "COST/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"FgQifwxmQfjhHvh2ggVxQwb9qwRwHrxwwxxQXASLAnVH\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"9dDGmEfXJXgjcMAp516c5eUd1eBRW3ZgKg6diyBmd1xh\",\"token_decimals\":9},\"amm_info_address\":\"GQdUPA8cUV8WsqEdCfDQtphvztocNCoSBGo1wARtaAXK\",\"open_orders_address\":\"4MQHW9GXiypDCGgjgEGKYB6pLiPPF7v38ki9VpaiUvni\"}"
			  }
			]
		  },
		  "DEVIN/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DEVIN",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/DEVIN",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"invert": true,
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"G8Uvd8VWWiVvPCdpAikiNryRfoNfSpL6EUcZAdcrb68D\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"4atVfrWD9J6jwksaJLaTFRiwkxS4DyZ9MnFTK1vh9Hih\",\"token_decimals\":6},\"amm_info_address\":\"2cZQ71uDTBwFZT456koEwfZDLSV736hT688A18sD3n4M\",\"open_orders_address\":\"AqYvf3WyRQAjHyiM6HtsTUqwfzNQjsvghYynjbsjetR6\"}"
			  }
			]
		  },
		  "EGG/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "EGG",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "EGG/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"AgzrwnLSRntNBMaQNXF8kmRy3Zg3VC6thq2KA4zNrFop\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"BzoTmVRmsi1Ukw4wT9SuhtXjBmwFvgQK6t7LNX5b1uE8\",\"token_decimals\":9},\"amm_info_address\":\"BUPTMMrUnfeE5mk6L6ZUgrHtZo2qyeR44s9L6UGg1oQB\",\"open_orders_address\":\"A8Vq3pY7s43KsJFxe6UQpS6RVZCGVLyHUdqc4Q2HuLeo\"}"
			  }
			]
		  },
		  "FALX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "FALX",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "FALX/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"6TtniYJPdHJ764d3rUvk8SokmzyZMYCgqUjSgUAnjma2\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"2eH2DVXQvs5qWwDQjgiSPsdZ19KQsj84RKRfCeorsGem\",\"token_decimals\":9},\"amm_info_address\":\"2hPp2aKd6T6HZmMQW2LkqH7R1wLZDjzZ1bZjhj5nrhrV\",\"open_orders_address\":\"3jf9f9VJdUXQQha6nHJkZVxqBTW5oJUNHYuKLyfDDuMM\"}"
			  }
			]
		  },
		  "GOL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "GOL",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "GOL/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"ACuv7Frh33MqLZtv5aKjz4uh2ZZFauFQfm2t23Wk2Gkr\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"GDY1yj4VyCi4Wa57PFVVvEY7pDVwtjwMTmpgda3NQC8h\",\"token_decimals\":9},\"amm_info_address\":\"E3E5grXmLfETytkBKVBHCLt7FcRAfQLLXftJqSYF1noJ\",\"open_orders_address\":\"8GwoCLwbYKWXyqcJeo3TzcVGoWdRFuz7Qqn9ByZQ4d9s\"}"
			  }
			]
		  },
		  "GUMMY/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "GUMMY",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "GUMMY/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"33AQRrPaZTckDJQd5DZstiwi11tcMVryu63V8rAHFF7N\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"4DHoHzkMHYqJeNDpkdeL6AGDymLFjJnS4SRsJHoT52Bm\",\"token_decimals\":9},\"amm_info_address\":\"FMiecMsYhPdBf94zZKa7i6inK1GX7aypLf7QewNz1i6w\",\"open_orders_address\":\"FSv96pMp3x5XwFdYgqXUY47o7nSKhA6tvCHX1UZZPWnv\"}"
			  }
			]
		  },
		  "HABIBI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "HABIBI",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "HABIBI/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"AseQmyRFtmH2KGcBtsnDmVGiH68WP32KEak7VshLddr5\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"6Vj4gTxdkkhS2DFgyLuoAv1c1iTWZuYnZpw7bhK8oeQj\",\"token_decimals\":9},\"amm_info_address\":\"2ukgjDC99Nk34RfRjWjCoHAuQLtLnz8TLcBrDQk3f2ay\",\"open_orders_address\":\"Ap3oiBWsLbDFwcigjktNvt2WjPQnLReRd28wtqJE1yDF\"}"
			  }
			]
		  },
		  "HAMMY/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "HAMMY",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "HAMMY/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"324NgHgEDyU9d7TE9dkAkB2GNtqxdEU4PsYRTDL68qoR\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"4kPJL1LmempALPjjwMWSo6JRBjmKQY7HX3edozqmJBPe\",\"token_decimals\":9},\"amm_info_address\":\"X131b3frGn4b8ue51EyvrnzWuTuBGoM93uRYrNteEFy\",\"open_orders_address\":\"9WZDqKjvpyoAShnp3Dg1725uyo2aQtgp8z7GG9XdB5NM\"}"
			  }
			]
		  },
		  "KITTY/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "KITTY",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "KITTY/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"E1hKKYHebq4naKVxG191vL19Lm6afCP6sXneBdskSqcc\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"Ac2VgEf8eSXxECL2a8wH4z13TbDw8sk5zhoTnXSe3Zbh\",\"token_decimals\":9},\"amm_info_address\":\"StJ9GP9KKVsbvtEtBDSjWNL9jpgybCjyHAwYyTe4SpW\",\"open_orders_address\":\"9sDic5pic3Q4HYRtu2cv8W2WB7Kt5EUGsdmvghvPobPv\"}"
			  }
			]
		  },
		  "MANEKI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MANEKI",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MANEKI/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"5FeTzLNqwrvSzexFujeV62a2v4kmQUrBnCQjJANStMXj\",\"token_decimals\":5},\"quote_token_vault\":{\"token_vault_address\":\"2kjCeDKKK9pCiDqfsbS72q81RZiUnSwoaruuwz1avUWn\",\"token_decimals\":9},\"amm_info_address\":\"2aPsSVxFw6dGRqWWUKfwujN6WVoyxuhjJaPzYaJvGDDR\",\"open_orders_address\":\"9pd9FFJfVjY1aG9dh96ArJB5F2HAyfj2XryjVTHbJhc9\"}"
			  }
			]
		  },
		  "PAJAMAS/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PAJAMAS",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "PAJAMAS/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"2amsF7CaXcxBDU39e8H8Cm4EFJWJqhWhJ4TBgFFvkbMQ\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"F7vbaUrc9z1CMWDqVtpQCpSQ8m5k5s3WkAf7NVAHdemD\",\"token_decimals\":9},\"amm_info_address\":\"BqricZnjjtFg8wuTbckV6NZcTstuR7BZtKJtzH8oV3eK\",\"open_orders_address\":\"8eSiN9JD5WYJVznfu4EWwPnEDMSvwfSx12NyXdhhkUJ9\"}"
			  }
			]
		  },
		  "PUNDU/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PUNDU",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "PUNDU/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"FQoagUqLxpNq69dpqFLrKm1gySC92NLKMkVgtdHWMKtt\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"GpnU5CAFUyyodrtXTBfvWK7ewsmcJGNksr4Fe49AvpM5\",\"token_decimals\":9},\"amm_info_address\":\"7yEXWTjLyXwBEjMhNwP9dWVJp8G68JvY9KXGT83sDCaM\",\"open_orders_address\":\"CmQ1XbjSeg5opqnJ1nnVf3TBzXUA2KbeWgysCVwKaCSN\"}"
			  }
			]
		  },
		  "SELFIE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SELFIE",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/SELFIE",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"invert": true,
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"HGDoCdba9yPpKvyYptWv747mG2ti8oVr8Cz88gV9TMdW\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"HmEomvDg2BjV8bvdb1DWH52WEju6KDTpG5CZBBqW2Zgb\",\"token_decimals\":6},\"amm_info_address\":\"Dfk133hHxjAA1yPryNkoPERGJ5DMpUtm79YeY1p1Wiyh\",\"open_orders_address\":\"7xGGsWHaXoPw4mJaJKoUatrQbUSVyy3TjvniWVxTBfbc\"}"
			  }
			]
		  },
		  "SLOTH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SLOTH",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SLOTH/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"Bs7VsZxQYHndLFnfDRRmJ4D44gCoTv7vNoDF2s5s11cV\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"5xQzJAvJ7Ut4qoTwiKECnaMDUhZFivx96EFomcBbUShq\",\"token_decimals\":9},\"amm_info_address\":\"7mtJbVNEtejYmCLRriwQhymZdzn4wGRFTvTZ5721b4BD\",\"open_orders_address\":\"A7k1mZQNNNKCakhHZN9bQqLzowDmApHTb4564uw5tAVU\"}"
			  }
			]
		  },
		  "SPEED/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SPEED",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/SPEED",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"invert": true,
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"H1Lj6xGsMnbPKZcGAtJbnNWf8W72XEY3FCb9JMyWH9jq\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"Eiwa11wqf1T8nEgiH1DadeNJ56hXw7iobFExVYerBi9V\",\"token_decimals\":6},\"amm_info_address\":\"81E14MT778WvYKYpavCWoaqzSiT87VhboZqXmYgbABan\",\"open_orders_address\":\"6FAtx1SHxfdpWEvTH6zsDT8vb93cRno95ELBE1sUYmDP\"}"
			  }
			]
		  },
		  "SPIKE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SPIKE",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/SPIKE",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"invert": true,
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"12ZLC7u75vJr9nAuzuJ6yFECq3JAgFKFWRcEhjEwvrF8\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"EqGTncTywZiYeA74Zog9hgEtY7DapEuKd7QVDX38yP3k\",\"token_decimals\":6},\"amm_info_address\":\"Gk4uCFPHUMriPVGNaAFr6v2YB491ViZtdMpGNsJAWfTe\",\"open_orders_address\":\"4BYG6KxdExgvEd8Mp8rB2iH3VDQis3GkaHqHhPojx7QM\"}"
			  }
			]
		  },
		  "STACKS/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "STACKS",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "STACKS/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"AoKWNc3zKu3hfprVR4ZH9jPod4NyUZ2nzkU3gYGpwpfZ\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"4TBThQTWQ15sXNQJfRjp2gxQ8BKarn7u4FmEy1D5guCL\",\"token_decimals\":9},\"amm_info_address\":\"HdVFDQRgQiHBfBjz4oGuFLZy7m7qacLMQCiMC2y6QQAt\",\"open_orders_address\":\"9s6zTUAXti5r4k8AcpFszujdfqyPRJNpcPUQiT2rTsFp\"}"
			  }
			]
		  },
		  "MOTHER/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MOTHER",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/MOTHER",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"invert": true,
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"8uQwXPi1sWwUTVbDBnjznmf4mV44CETiNAh3UENvHejV\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"8ZcrNqaDbqy1H4R2DtmGnuZnJ6TKGSsaGmyRGQQeELLv\",\"token_decimals\":6},\"amm_info_address\":\"HcPgh6B2yHNvT6JsEmkrHYT8pVHu9Xiaoxm4Mmn2ibWw\",\"open_orders_address\":\"1z3rLR6AwR8gjVZ8ArqHp9kdaPrNwPvCdrZ6jGy6wwF\"}"
			  }
			]
		  },
		  "$RETIRE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "$RETIRE",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "$RETIRE/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"2NJXbbLxfygPGusgTyxbFngLaodRRRCpXxeo1pv7M5XQ\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"8PcfqMPandh18bYZJKEvjuQRH5bcH4Y6TZzGUSWEPeYG\",\"token_decimals\":9},\"amm_info_address\":\"CQQDXt2M6Cx1J8N3cYsSmPiD7fcLdU5RpVtRbs9WaCXZ\",\"open_orders_address\":\"7uUa9ELGG1NQymmbNiUwYmkbrzhDzAdb1u9yHPNfnuZk\"}"
			  }
			]
		  },
		  "$WIF/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "$WIF",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "$WIF/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"7UYZ4vX13mmGiopayLZAduo8aie77yZ3o8FMzTeAX8uJ\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"7e9ExBAvDvuJP3GE6eKL5aSMi4RfXv3LkQaiNZBPmffR\",\"token_decimals\":9},\"amm_info_address\":\"EP2ib6dYdEeqD8MfE2ezHCxX3kP3K2eLKkirfPm5eyMx\",\"open_orders_address\":\"6jeayPbLeJq9o6zXbCtLsEJuPyPFyojWoH55xrksfsoL\"}"
			  }
			]
		  },
		  "ANDY/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ANDY",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "ANDY/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"pVCniSexvGFdDTjYuzoSKXDkoTqFjRhqJpmzzYMs7tY\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"2tt4w3C9hKjtzVgPqqa9Apbxz5qEKEAWodJEbbmfpwEm\",\"token_decimals\":9},\"amm_info_address\":\"7KqUHWB67F1qpqRn5sHLtpdZaXkb8pWdKcqWb6zwKoQY\",\"open_orders_address\":\"A4wdbaSW9LgBkyj1kvzmbcXuJQ8tbcet3u4ny6UUGepN\"}"
			  }
			]
		  },
		  "ANSEM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ANSEM",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "ANSEM/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"5cM8YfCvkALkout2h9WxmYrc5e61YmvLgLLWsrZgumgK\",\"token_decimals\":1},\"quote_token_vault\":{\"token_vault_address\":\"B86KHBLhVVQQnsgbn6SDJR43NSKbqmpxZsrTai45yrMy\",\"token_decimals\":9},\"amm_info_address\":\"7xGQkpvqrqCNKwangJaj6h8KFqMu3RC9PRGYkAXhw2kw\",\"open_orders_address\":\"Dt51heAAxkx8hVok4zhotDC1J7pnpFrRr5ALgnwRhUY8\"}"
			  }
			]
		  },
		  "ATR/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "ATR",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "ATR/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"4rm2qWHwGZGj9rWoRvU7m3FDdsZJV11wuHmczw27C3Wc\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"Ahy4NhTyBaXZfsGjq4DDxGaMosBkAjaanGYdfeZjuDzP\",\"token_decimals\":6},\"amm_info_address\":\"2Ky6BskrcKNCJSrP4X6bgrPPe1erBArBAhyZi2C8nPwy\",\"open_orders_address\":\"DmY92VBBtKaNX2ZitPNZDX1G7VSrkpyeGwizXbCJc5ed\"}"
			  }
			]
		  },
		  "BAG/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BAG",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "BAG/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"7eLwyCqfhxKLsKeFwcN4JdfspKK22rSC4uQHNy3zWNPB\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"Cr7Yo8Uf5f8pzMsY3ZwgDFNx85nb3UDvPfQxuWG4acxc\",\"token_decimals\":9},\"amm_info_address\":\"Bv7mM5TwLxsukrRrwzEc6TFAj22GAdVCcH5ViAZFNZC\",\"open_orders_address\":\"Du6ZaABu8cxmCAvwoGMixZgZuw57cCQc8xE8yRenaxL4\"}"
			  }
			]
		  },
		  "BODEN/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BODEN",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "BODEN/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"54zedUwxuSnmHHYg9oY1AfykeBDaCF6ZFZDW3ym2Nea4\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"DzpiXKsTUCacKyahLBUC5sfjj2fiWbwCpiCPEgyS3zDC\",\"token_decimals\":9},\"amm_info_address\":\"6UYbX1x8YUcFj8YstPYiZByG7uQzAq2s46ZWphUMkjg5\",\"open_orders_address\":\"9ndGwmmTcFLut1TNjWFA8pDvcrxgmqPEJTZ2Y3jTipva\"}"
			  }
			]
		  },
		  "BOME/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BOME",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "BOME/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"FBba2XsQVhkoQDMfbNLVmo7dsvssdT39BMzVc2eFfE21\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"GuXKCb9ibwSeRSdSYqaCL3dcxBZ7jJcj6Y7rDwzmUBu9\",\"token_decimals\":9},\"amm_info_address\":\"DSUvc5qf5LJHHV5e2tD184ixotSnCnwj7i4jJa4Xsrmt\",\"open_orders_address\":\"38p42yoKFWgxw2LCbB96wAKa2LwAxiBArY3fc3eA9yWv\"}"
			  }
			]
		  },
		  "CHAT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CHAT",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "CHAT/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"FKCCPsYDgEoVpEhyE2XMFAXq5zWFrWHgpQjVEKQk1C54\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"9APPnbdEXbJpktfKSGbbdgYvj6K3ZFRDFwQUabFw6CHP\",\"token_decimals\":9},\"amm_info_address\":\"9kLGUEFwEuFzn9txDfGJ3FimGp9LjMtNPp4GvMLfkZSY\",\"open_orders_address\":\"G9fse9D2feKdSjy4eLDQfuuBfxQDqektwNMG9smVBJr9\"}"
			  }
			]
		  },
		  "CHEEMS/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "CHEEMS",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "CHEEMS/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"HRi4eJ8xWFG4hsv7FA9L7YnPLSxhQR4U7CCXKyZvcLXe\",\"token_decimals\":4},\"quote_token_vault\":{\"token_vault_address\":\"4gnEBvHQEx4nLcf9qWk1Wsxh9V1GFFDf4MfwEYmFo8hm\",\"token_decimals\":6},\"amm_info_address\":\"3FRFbRMvUjjufZBAjbdGcb1PYYyH9MJAyRQ7WuHuXTXe\",\"open_orders_address\":\"FZ6NVUWGcok9vGTBiHs4BMcvqTqEH1YCzpd3WPRqSpp7\"}"
			  }
			]
		  },
		  "CHONKY/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CHONKY",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "CHONKY/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9rJqiGuRG971HCpapVNJtN4ho2fKMhkPiZRhQCAohonU\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"6Fbs4sCBH5jZp1v2Rr6zVdU44Fn4Vv9iPhY6eihjfdbz\",\"token_decimals\":9},\"amm_info_address\":\"E61pEDMEwf8iUHFhmGn3Wcj5P32DPjKDgo1UNjjaNrg1\",\"open_orders_address\":\"2PP6pkjc3QcQcB5qr6xg6gD4AWtjzDZzBDNebMQYFJqP\"}"
			  }
			]
		  },
		  "DUKO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DUKO",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "DUKO/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"HD7aZ6YrqAxVbGNAMEKxozcW1ZDU7pbKfd7XMmZtxyzk\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"A9J2mXPXfRZ7Sh2ymUgCJM4p9iUjZBcyAfrz49PoBBN4\",\"token_decimals\":9},\"amm_info_address\":\"BGS69Ju7DRRVxw9b2B5TnrMLzVdJcscV8UtKywqNsgwx\",\"open_orders_address\":\"FoBQDGey332Ppv1KiTow8z9oZP8n6mEPLyhedPdG1nUG\"}"
			  }
			]
		  },
		  "FKETH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "FKETH",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "FKETH/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"DNh9pRMiRX6zwWuzsXtmxmXLdbAGNuYg4dmmnzrpL871\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"47fCbQKnJYaMbPaPSyUrPXPUahizhYwAbwXnEcKN1vwD\",\"token_decimals\":9},\"amm_info_address\":\"HvAUaYpykFbUmzyGUCPbLR2nKA43cXspfxYNyYT2mw7j\",\"open_orders_address\":\"5VEfGLutckRLb3sFj9US8Dz4sFQ29xNMigpdxpTFj1bj\"}"
			  }
			]
		  },
		  "GMEOW/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "GMEOW",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "GMEOW/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9crtLiYfxYVYQ9sCfWix9vAMPJyBXCcMzCPXZ5isPFxB\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"ANLamVN6Df4HqC6YevQskovddsjhkqBqHsyLzhDibFEj\",\"token_decimals\":9},\"amm_info_address\":\"BbbuEcbX1JhxDTjnDXvi5CuQFR77FW7AbqnuyvZEmAFK\",\"open_orders_address\":\"8LdXZ1TUMBmZ59WJXt3YPCcifjRmvTdMcRvzpBBoUgyQ\"}"
			  }
			]
		  },
		  "HARAMBE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "HARAMBE",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "HARAMBE/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"5f9Fgcp2C9vdrp75GspNKBjzdaxq5uiqpLVkgtWKpDZZ\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"Apgp3SzNB5VpVWbK5q2ucBvCJEsf1gqXL4iUAqvD9pgB\",\"token_decimals\":9},\"amm_info_address\":\"2BJKy9pnzTDvMPdHJhv8qbWejKiLzebD7i2taTyJxAze\",\"open_orders_address\":\"BPv68DZUMxpqvfRye2JoeK1GRkkGs5PEUycmx5b448x2\"}"
			  }
			]
		  },
		  "KHAI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "KHAI",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "KHAI/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"6g4rhxVTrN6SrtNvimq4QiU8yA5XScvwL6wxaMkegrtJ\",\"token_decimals\":8},\"quote_token_vault\":{\"token_vault_address\":\"7p2PrGMyeetNRqTKFraL7eYo2TbU3apWz6vfqrZFiPcG\",\"token_decimals\":9},\"amm_info_address\":\"ECbK6PSMZ5yQaUYBocsXaVrax2fWADw2ijTqLGPtt9sC\",\"open_orders_address\":\"2DaRg4UycKL9GSVfARBDrcensb89WD5PyyFX9NrMunLc\"}"
			  }
			]
		  },
		  "LIGMA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "LIGMA",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "LIGMA/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"6aefaP7C6eKbW2taLqmyHinYH4ZMyY2G6MdqNu6PvfbL\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"G3kKDmmxwpjt4NVdQgdvgiuFxFsAsC1hSv4PVg63cKwM\",\"token_decimals\":9},\"amm_info_address\":\"2AvbJNsiDD99yYYgk7eGKXuKU4PRVibQYFv6xfNA5Fce\",\"open_orders_address\":\"2RvH4R5h7r3f6BJfZpPqzcFqUC46RziUko5HFh2KYqff\"}"
			  }
			]
		  },
		  "LIKE/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "LIKE",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "LIKE/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"8LoHX6f6bMdQVs4mThoH2KwX2dQDSkqVFADi4ZjDQv9T\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"2Fwm8M8vuPXEXxvKz98VdawDxsK9W8uRuJyJhvtRdhid\",\"token_decimals\":6},\"amm_info_address\":\"GmaDNMWsTYWjaXVBjJTHNmCWAKU6cn5hhtWWYEZt4odo\",\"open_orders_address\":\"Crn5beRFeyj4Xw13E2wdJ9YkkLLEZzKYmtTV4LFDx3MN\"}"
			  }
			]
		  },
		  "MARVIN/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MARVIN",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MARVIN/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"3vLTpZneGAzceAXKu2HuesT4rt6ksRJ3Q9WvjUmwksqA\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"3rWa7PstXZk4ZaEhLamfoqMVozwq7hfXEDqyNbHcL4uK\",\"token_decimals\":9},\"amm_info_address\":\"CtR5r6zp6NhXfR2iVSRmSir5TwYnj6fcYJNuehDBbDyF\",\"open_orders_address\":\"AqVCnWScTmiKbvpmaDrAGLZbTCypHg4s37UAokKRCgVH\"}"
			  }
			]
		  },
		  "MEW/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MEW",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MEW/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"4HqAjFKuQX7tnXEkxsjGQha3G4bcgER8qPXRahn9gj8F\",\"token_decimals\":5},\"quote_token_vault\":{\"token_vault_address\":\"BhNdEGJef9jSqT1iCEkFZ2bYZCdpC1vuiWtqDt87vBVp\",\"token_decimals\":9},\"amm_info_address\":\"879F697iuDJGMevRkRcnW21fcXiAeLJK1ffsw2ATebce\",\"open_orders_address\":\"CV3Gq5M2R7KRU5ey4LpnZYRR7r7vzKoV9Bt4mZ8P6bSB\"}"
			  }
			]
		  },
		  "MONK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MONK",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MONK/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"CeLZLhC2nScSpsKqRL1eRr3L3LLfjDzakZLCUKcUHW1m\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"Be6DM12uTWtxHMSRJeah3J5PRP4CumR28Yy2qpLQFTE7\",\"token_decimals\":9},\"amm_info_address\":\"5mYCuaXmqW1McUe18Ry6gbWUQhtk1f4GxJ9j7vRj34HJ\",\"open_orders_address\":\"F1KFumMDuNonPprwUxarH6bbip9TR5wfsCZzRuKM8XBM\"}"
			  }
			]
		  },
		  "MOUTAI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MOUTAI",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MOUTAI/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"4VkmzH14ETcNhSQLTK6AtL1ZP8UmvWpbNCgokDVfiCcD\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"HSQPdDCxtGo4fTHeZuBGWtQUqHgRsgdz4BVhTCCAtsTv\",\"token_decimals\":9},\"amm_info_address\":\"578CbhKnpAW5NjbmYku6qSaesZZLy3xwFQ8UkDANzd91\",\"open_orders_address\":\"FCQvrj9mrWN5XsPHDSfKf17i8xbzLxW3Esor7nw42nsp\"}"
			  }
			]
		  },
		  "MPLX/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "MPLX",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MPLX/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"5zdZza5N2TzV7cPtLeqCZQQRYCCHFVdXWLMeJo75DK24\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"7mwk7ZBiTtrKkKC5o34gpFBSCabEvLkp2fLjGNz43PyM\",\"token_decimals\":6},\"amm_info_address\":\"F3enDARqxRtTDuD8RxYLFXUEfPgph9VQLE6HXmWDqTSS\",\"open_orders_address\":\"ANiHc2R9qkDJxj9kLr6Bfwkw51niaK8q22CkuFTTUbTd\"}"
			  }
			]
		  },
		  "MUMU/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MUMU",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "MUMU/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"2Re1H89emr8hNacyDTrm1NU8VEhuwaJX7JwcdDqy5Q6g\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"4VPXFMpndqZhME27vMqtkuGtBo7hVTA9kEvo87zbjXsA\",\"token_decimals\":9},\"amm_info_address\":\"FvMZrD1qC66Zw8VPrW15xN1N5owUPqpQgNQ5oH18mR4E\",\"open_orders_address\":\"BjWyTUxXSNXN1GNzwR7iRhqmdc3XukYpWFfqy1o94DF2\"}"
			  }
			]
		  },
		  "NICK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "NICK",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "NICK/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"FDvQiwbJVHdSZE3ngZ6WCkJfciFTdg958W7bxyKU2PJ9\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"z6ZCZESyof3ZgCJ23hY31f1SSD33gQgyVRfMB8wP9iM\",\"token_decimals\":9},\"amm_info_address\":\"2FAZebHiPTxoTskLcV4EFJiqpXgourmP4qa9rt9qiU1o\",\"open_orders_address\":\"34SvJTo8VDcQQGcoAVpc3rJrCdPooLFWj83nUQYadzeF\"}"
			  }
			]
		  },
		  "NINJA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "NINJA",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "NINJA/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"5Mmie9Drh6RGMy8X8UQ3egyBi4Hvva1TR778bf77ViCV\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"HxVbv76N8EoRGEPJsKdtWCu3mz7ZXJi8dbZy8kM3QL3i\",\"token_decimals\":9},\"amm_info_address\":\"B8sv1wiDf9VydktLmBDHUn4W5UFkFAcbgu7igK6Fn2sW\",\"open_orders_address\":\"6DMN41zTWap2S2e2fzz5bKKD8DToxRiandvZ74h7FC6s\"}"
			  }
			]
		  },
		  "NOS/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "NOS",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "NOS/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9Gs4LvFZw18EBLrSmZbQBw4G2SpTu4bJRCWH1Dz33cUZ\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"FqKU4BxbabPd1tcZAVVv8JkdUWmdz32CocRM856gA3Lw\",\"token_decimals\":6},\"amm_info_address\":\"Cj7kD2VmzwSrwKBieuYYbjPEvr8gwhNi76KUESbGDNfF\",\"open_orders_address\":\"DEM1Zse8UWfKEk9dH1Jkjzepdb9DSaxMZ8uDe34rmbE5\"}"
			  }
			]
		  },
		  "NUB/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "NUB",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "NUB/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9uNqUwneLXbQ6YKndciL5aBXTLJhwpyDXkZmaBbWfwWz\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"75DrZt3zmGSFfKaYDm7yHLKMrr35Wy8ffBNN1143PWbj\",\"token_decimals\":9},\"amm_info_address\":\"83G6VzJzLRCnHBsLATj94VCpRimyyqwuN6ZfL11McADL\",\"open_orders_address\":\"CLXBUkh3hMKNDRUZFFKS721Q1NJb11oHrYvV66QMBcVv\"}"
			  }
			]
		  },
		  "PENG/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PENG",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "PENG/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"2g5q7fBGKZm2CXix8JjK4ZFdBTHQ1LerxkseBTqWuDdD\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"GmLJXUzjQAAU86a91hKesg5P9pKb6p9AZaGBEZLaDySD\",\"token_decimals\":9},\"amm_info_address\":\"AxBDdiMK9hRPLMPM7k6nCPC1gRARgXQHNejfP2LvNGr6\",\"open_orders_address\":\"9E5VWkY1UsbhkXW4Lk1YovkVouWMG57CuCNXUmecrGpC\"}"
			  }
			]
		  },
		  "PONKE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PONKE",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "PONKE/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"D7rw7fyEzo9EQcozjqAHJwbdbywGcSLw1at5MioZtMZ4\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"8DcvfWidQ53a3SCBrWxBWL2UU5zEBAKEypApiiCCFu2Y\",\"token_decimals\":9},\"amm_info_address\":\"5uTwG3y3F5cx4YkodgTjWEHDrX5HDKZ5bZZ72x8eQ6zE\",\"open_orders_address\":\"ECoptgCPMxXXWtxv3Zg2V3E7SpWp6SGqKqj32AcdWRQK\"}"
			  }
			]
		  },
		  "POPCAT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "POPCAT",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "POPCAT/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"4Vc6N76UBu26c3jJDKBAbvSD7zPLuQWStBk7QgVEoeoS\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"n6CwMY77wdEftf2VF6uPvbusYoraYUci3nYBPqH1DJ5\",\"token_decimals\":9},\"amm_info_address\":\"FRhB8L7Y9Qq41qZXYLtC2nw8An1RJfLLxRF2x9RwLLMo\",\"open_orders_address\":\"4ShRqC2n3PURN7EiqmB8X4WLR81pQPvGLTPjL9X8SNQp\"}"
			  }
			]
		  },
		  "POPCAT/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "POPCAT",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "POPCAT/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"Ej1CZHUcHYcqAx3pJXUvqCTs3diVmEWSfozvQQLsQkyU\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"pTJo5c4ynoxxRAgDyWgQKasR8dgqQHP3CSeoXzDgZvZ\",\"token_decimals\":6},\"amm_info_address\":\"HBS7a3br8GMMWuqVa7VB3SMFa7xVi1tSFdoF5w4ZZ3kS\",\"open_orders_address\":\"9Q1E7B4w5Vhb5RjbmojpEuZbMZ944m9vDZZJoxBcGBRS\"}"
			  }
			]
		  },
		  "RAY/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "RAY",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "RAY/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"Em6rHi68trYgBFyJ5261A2nhwuQWfLcirgzZZYoRcrkX\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"3mEFzHsJyu2Cpjrz6zPmTzP7uoLFj9SbbecGVzzkL1mJ\",\"token_decimals\":9},\"amm_info_address\":\"AVs9TA4nWDzfPJE9gGVNJMVhcQy3V9PGazuz33BfG2RA\",\"open_orders_address\":\"6Su6Ea97dBxecd5W92KcVvv6SzCurE2BXGgFe9LNGMpE\"}"
			  }
			]
		  },
		  "RAY/USDC": {
			"ticker": {
			  "currency_pair": {
				"Base": "RAY",
				"Quote": "USDC"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "RAY/USDC",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"FdmKUE4UMiJYFK5ogCngHzShuVKrFXBamPWcewDr31th\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"Eqrhxd7bDUCH3MepKmdVkgwazXRzY6iHhEoBpY7yAohk\",\"token_decimals\":6},\"amm_info_address\":\"6UmmUiYoBjSrhakAobJw8BvkmJtDVxaeBtbt7rxWo1mg\",\"open_orders_address\":\"CSCS9J8eVQ4vnWfWCx59Dz8oLGtcdQ5R53ea4V9o2eUp\"}"
			  }
			]
		  },
		  "RAY/USDT": {
			"ticker": {
			  "currency_pair": {
				"Base": "RAY",
				"Quote": "USDT"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "RAY/USDT",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"3wqhzSB9avepM9xMteiZnbJw75zmTBDVmPFLTQAGcSMN\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"5GtSbKJEPaoumrDzNj4kGkgZtfDyUceKaHrPziazALC1\",\"token_decimals\":6},\"amm_info_address\":\"DVa7Qmb5ct9RCpaU7UTpSaf3GVMYz17vNVU67XpdCRut\",\"open_orders_address\":\"8MSPLj4c2hi1fZGDARvxLXQp1ooDQ8iGnWXbGdwvZxUQ\"}"
			  }
			]
		  },
		  "RETARDIO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "RETARDIO",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "RETARDIO/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"HXzTvbuKKPyNMmLKJb8vaSUaRZsVS2J2AAsDuDm36rNC\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"HNcAAdLKHSRnwdmmWCYnP5Zcd11sfGpAoCuWFtugt2ma\",\"token_decimals\":9},\"amm_info_address\":\"5eLRsN6qDQTQSBF8KdW4B8mVpeeAzHCCwaDptzMyszxH\",\"open_orders_address\":\"5TcDuBbtU8Q6LagcM8wfw1Ux2MWgCC6Q1FY22FVDZnXX\"}"
			  }
			]
		  },
		  "SLERF/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SLERF",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SLERF/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9A2ZsPz5Zg6jKN4o4KRMjTVPmkH51wYWFLmt4KBRy1Rq\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"5Zumc1SYPmQ89nqwXqzogeuhdJ85iEMpSk35A4P87pmD\",\"token_decimals\":9},\"amm_info_address\":\"AgFnRLUScRD2E4nWQxW73hdbSN7eKEUb2jHX7tx9YTYc\",\"open_orders_address\":\"FT5Ptk37g5r6D9BKt3hne8ovHZ1g56oJBvuZRwn3zS3j\"}"
			  }
			]
		  },
		  "SMOLE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SMOLE",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SMOLE/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"VDZ9kwvKRbqhNdsoRZyLVzAAQMbGY9akHbtM6YugViS\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"HiLcngHP5y1Jno53tuuNeFHKWhyyZp3XuxtKPszD6rG2\",\"token_decimals\":9},\"amm_info_address\":\"5EgCcjkuE42YyTZY4QG8qTioUwNh6agTvJuNRyEqcqV1\",\"open_orders_address\":\"FeKBjZ5rBvHPyppHf11qjYxwaQuiympppCTQ5pC6om3F\"}"
			  }
			]
		  },
		  "SOL/EPIK": {
			"ticker": {
			  "currency_pair": {
				"Base": "SOL",
				"Quote": "EPIK"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/EPIK",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"9aGBQqKRyC5bbrZsnZJJtp59EqJj7vBkgV3HehgKEu5y\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"ANpMJb9ToMGNivLEdmBNBC2Qcf5ASaZkEdmUddV1FUZB\",\"token_decimals\":6},\"amm_info_address\":\"AZaaQaRhp1ys9VaJBRZYbmPz3JSBSp7m8cSSrLBn4BP9\",\"open_orders_address\":\"FjCKdnpN1t262QGGn6chWYRtoSaY6fuYxyKoqhinyGEK\"}"
			  }
			]
		  },
		  "SOL/HOBBES": {
			"ticker": {
			  "currency_pair": {
				"Base": "SOL",
				"Quote": "HOBBES"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/HOBBES",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"4J3cZy8G2qm6MBVGfeXhYETZvbRThv9TPPeY3p83QYLb\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"7DejqJN5iRcuUhR7C1Vif3SbjTXKCzkpyS3AxC28tiaF\",\"token_decimals\":6},\"amm_info_address\":\"9VffBiow5r5YQzgK56rirEWpu45gZGrDWzm9JUt6zL9G\",\"open_orders_address\":\"9q3x5stYdC6xuxdNjYQCRVktdHZiVrqWw2qcShrAQB2b\"}"
			  }
			]
		  },
		  "SOL/MAIL": {
			"ticker": {
			  "currency_pair": {
				"Base": "SOL",
				"Quote": "MAIL"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "SOL/MAIL",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"BtJXP2iWPGW2x3EsscHtCuFLBECRCfrxsJ2SDi9jh96C\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"GmHpq7Wgri9TFitGj89quvwRc1ALhe7dePM6VBAiqxrC\",\"token_decimals\":6},\"amm_info_address\":\"9DYGj7g2b5fipk9wGsUhxdv5zpfTsoGzCiS29vH8Cfrs\",\"open_orders_address\":\"CLxCfhwK9SQAPLHu3KpMTLxvunG4JwRYxd5wY3tuNqQL\"}"
			  }
			]
		  },
		  "TOOKER/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TOOKER",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "TOOKER/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"Cki9WdL3sCoNY3cLmfG4iqSbvB8g1Fr9tw8qa5tP1m3Y\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"2vTTh5pGbzc6guAJmt78XnTcXVBEZEWmGBkXkSNZwN59\",\"token_decimals\":9},\"amm_info_address\":\"3vGHsKVKNapB4hSapzKNwtiJ6DA8Ytd9SsMFSoAk154B\",\"open_orders_address\":\"5dzcxMHjuNU5LZyEXBhoWWKuxw51Z3626TTf2FTfLJjb\"}"
			  }
			]
		  },
		  "TREMP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TREMP",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "TREMP/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"G2XNC6Rt2G7JZQWhqpJriYwZyxd2L52KSDbDNBCYCpvx\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"9DfnSR9h3hrmgy5pjqBP3SrVQRWPfjSqZZBrNNYGoyaN\",\"token_decimals\":9},\"amm_info_address\":\"5o9kGvozArYNWfbYTZD1WDRkPkkDr6LdpQbUUqM57nFJ\",\"open_orders_address\":\"kTgLvRcrvhxJy9KZFureP8fU5L11BzFrRvUEUa1joai\"}"
			  }
			]
		  },
		  "TRUMP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TRUMP",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "TRUMP/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"GwUn3JYQ5PC6P9cB9HsatJN7y7aw3BXrNqHPpoHtWyKF\",\"token_decimals\":8},\"quote_token_vault\":{\"token_vault_address\":\"5DNPt6WYgj3x7EmU4Fyqe3jDYPk2HMAB21H5N4Ggbev9\",\"token_decimals\":9},\"amm_info_address\":\"7Lco4QdQLaW6M4sxVhWe8BHjrykyzjcjGTo4a6qYGABK\",\"open_orders_address\":\"FAWLdBB8kmWZQ74KpYAYN3YaEW31Si8qrwuQPauFSoma\"}"
			  }
			]
		  },
		  "USDC/$MYRO": {
			"ticker": {
			  "currency_pair": {
				"Base": "USDC",
				"Quote": "$MYRO"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "USDC/$MYRO",
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"AMtPGYQS873njs35mD9MAAMKoospEuzNHPy7LQuuKo4A\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"BUvMbqP311JDU4ZGnf1rSZJLjutTU9VpNLEos393TYyW\",\"token_decimals\":9},\"amm_info_address\":\"5WGYajM1xtLy3QrLHGSX4YPwsso3jrjEsbU1VivUErzk\",\"open_orders_address\":\"2w1mZXi8PNqUz4gbezu4xvPzcGogDmVdLXGLhTpPbczd\"}"
			  }
			]
		  },
		  "VCAT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "VCAT",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "VCAT/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"8gNjWm2wGubUiAXT9cXBeoY2NqAFkrnvKkh9J3gHZ7Wn\",\"token_decimals\":9},\"quote_token_vault\":{\"token_vault_address\":\"7G9HpLoYVhcBsg7ZEy928iUuzetJFK4AWBcfaCQTMp72\",\"token_decimals\":9},\"amm_info_address\":\"m9uiXqNMAxP7BNmtRf4NwkeKExjT5Hc6ftyjtNJB6E5\",\"open_orders_address\":\"9qnrgiBExCuugfUrK4Hgb3treG1YbBsNRqyk4cGDPzMF\"}"
			  }
			]
		  },
		  "VONSPEED/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "VONSPEED",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "VONSPEED/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"89CwpgTSsCc9u187kKvQQo6VAL5gKZViVub4eaNXfrtu\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"2fEchHP5r5wks9PPN1C2S39wPCe6Ni4247oMMTausc2a\",\"token_decimals\":9},\"amm_info_address\":\"34yiYY6kZmnVfmdQGtv2HugiNB5g1DcMDuc2VdckidB7\",\"open_orders_address\":\"CzfQYWMuLpnxCNMmQbwezkkivLD8Mn7W26MDNLgPtdS7\"}"
			  }
			]
		  },
		  "WHALES/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "WHALES",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "raydium_api",
				"off_chain_ticker": "WHALES/SOL",
				"normalize_by_pair": {
				  "Base": "SOL",
				  "Quote": "USD"
				},
				"metadata_JSON": "{\"base_token_vault\":{\"token_vault_address\":\"P9uSMnNEGHDP7Dhu7fKWfRViAGGHjEMv6urC8c2qG4k\",\"token_decimals\":6},\"quote_token_vault\":{\"token_vault_address\":\"CktEbT37HFRtwXVjwPEVfXHdcTAnqnmCvkgAw9SEN7zf\",\"token_decimals\":9},\"amm_info_address\":\"DczmyvnV8hR7d8zvy6bAoc2itZbFvLAx9iG2D7gyyt9e\",\"open_orders_address\":\"5JAwqabcp6KnfUe88RiaMgdpE3nw6CQu4NyAfbGmNEz2\"}"
			  }
			]
		  }
		}
	  }
	`

	// CoreMarketMap is used to initialize the Core market map.
	CoreMarketMap mmtypes.MarketMap
	// CoreMarketMapJSON is the JSON representation of the Core MarketMap that can be used
	// to initialize for a genesis state or used by the sidecar as as static market map.
	CoreMarketMapJSON = `
	{
		"markets": {
		  "AAVE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AAVE",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "AAVEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "aaveusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "AAVEUSD"
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "AAVE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "AAVE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "binance_ws",
				"off_chain_ticker": "AAVEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "AAVE-USD"
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "AAVE_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ADA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ADA",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "ADAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ADAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ADA-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ADA_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "adausdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "ADAUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ADA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ADAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ADA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "ADA_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "AEVO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AEVO",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "AEVOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "AEVOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "AEVO_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "AEVOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "AEVO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ALGO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ALGO",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "ALGOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ALGOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ALGO-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "ALGOUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ALGO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ALGOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ALGO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "APE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "APE",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "APEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "APE-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "APE_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "APEUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "APE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "APEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "APE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "APE_USD"
			  }
			]
		  },
		  "APT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "APT",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "APTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "APTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "APT-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "APT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "aptusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "APT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "APTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "APT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "APT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ARB/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ARB",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "ARBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ARBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ARB-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ARB_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "arbusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ARB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ARBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ARB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "ARB_USD"
			  }
			]
		  },
		  "ARKM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ARKM",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "ARKMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ARKMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ARKM_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ARKM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ARKMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ASTR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ASTR",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "ASTRUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ASTR_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "ASTRUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ASTR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ASTRUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ASTR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ATOM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ATOM",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "ATOMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ATOMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ATOM-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ATOM_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "ATOMUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ATOM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ATOMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ATOM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "ATOM_USD"
			  }
			]
		  },
		  "AVAX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AVAX",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "AVAXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "AVAXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "AVAX-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "AVAX_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "avaxusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "AVAXUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "AVAX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "AVAX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "AVAX_USD"
			  }
			]
		  },
		  "AXL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AXL",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "AXLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "AXLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "AXL-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "WAXL_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "WAXLUSD"
			  }
			]
		  },
		  "BCH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BCH",
				"Quote": "USD"
			  },
			  "decimals": 7,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "BCHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "BCHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "BCH-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "BCH_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "bchusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "BCHUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BCH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "BCHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "BCH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "BCH_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "BLUR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BLUR",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "BLUR-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "BLUR_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "BLURUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BLUR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "BLURUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "BLUR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "BNB/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BNB",
				"Quote": "USD"
			  },
			  "decimals": 7,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "okx_ws",
				"off_chain_ticker": "BNB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BNB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "BNBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "binance_ws",
				"off_chain_ticker": "BNBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "BNBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "BNB_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "BONK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BONK",
				"Quote": "USD"
			  },
			  "decimals": 14,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "BONKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "BONKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "BONK-USD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BONK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "BONK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "BONKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "BONK_USD"
			  }
			]
		  },
		  "BTC/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BTC",
				"Quote": "USD"
			  },
			  "decimals": 5,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "BTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "BTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "BTC-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "btcusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XXBTZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BTC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "BTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "BTC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "BTC_USD"
			  }
			]
		  },
		  "COMP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "COMP",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "COMPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "COMP-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "COMP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "COMPUSD"
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "COMPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "COMP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "CRV/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CRV",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "CRVUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "CRV-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "CRV_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "CRVUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "CRV-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "CRVUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "CRV-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "DOGE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DOGE",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "DOGEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "DOGEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "DOGE-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "DOGE_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "dogeusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XDGUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "DOGE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "DOGEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "DOGE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "DOGE_USD"
			  }
			]
		  },
		  "DOT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DOT",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "DOTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "DOTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "DOT-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "DOT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "DOTUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "DOT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "DOTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "DOT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "DOT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "DYDX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DYDX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "DYDXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "DYDXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "DYDX_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "DYDX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "DYDXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "DYDX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "DYM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DYM",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "DYMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "DYMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "DYM_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "DYM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "DYMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "EOS/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "EOS",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "EOSUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "EOSUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "EOS-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "EOS_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "EOSUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "EOS-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "EOS-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "EOSUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ETC/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ETC",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "ETCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ETC-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ETC_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "etcusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ETC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ETCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ETC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ETH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ETH",
				"Quote": "USD"
			  },
			  "decimals": 6,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "ETHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ETHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ETH-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "ethusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XETHZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ETH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ETHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ETH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "ETH_USD"
			  }
			]
		  },
		  "FET/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "FET",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "FETUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "FET-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "FETUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "FET-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "FET-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "FETUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "FET_USD"
			  }
			]
		  },
		  "FIL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "FIL",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "FILUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "FIL-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "FIL_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "filusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "FILUSD"
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "FILUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "FIL-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "GRT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "GRT",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "GRTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "GRTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "GRT-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "GRT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "GRTUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "GRT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "GRTUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "GRT-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "GRT_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "HBAR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "HBAR",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "HBARUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bitstamp_api",
				"off_chain_ticker": "hbarusd"
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "HBARUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "HBAR-USD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "HBAR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "HBARUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "HBAR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "HBAR_USD"
			  }
			]
		  },
		  "ICP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ICP",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "ICPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ICPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "ICP-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "ICPUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ICP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ICP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ICPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "ICP_USD"
			  }
			]
		  },
		  "IMX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "IMX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "IMXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "IMX-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "IMXUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "IMX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "IMXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "IMX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "INJ/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "INJ",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "INJUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "INJUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "INJ-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "INJUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "INJ-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "INJ-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "INJUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "INJ_USD"
			  }
			]
		  },
		  "JTO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "JTO",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "JTO-USD"
			  },
			  {
				"name": "binance_ws",
				"off_chain_ticker": "JTOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "JTOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "JTOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "JTO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "JTO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "JUP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "JUP",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "JUP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "JUP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "binance_ws",
				"off_chain_ticker": "JUPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "JUPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "JUP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "JUPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "JUP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "LDO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "LDO",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "LDOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "LDO-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "LDOUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "LDO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "LDOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "LDO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "LDO_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "LINK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "LINK",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "LINKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "LINKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "LINK-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "LINKUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "LINK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "LINKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "LINK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "LINK_USD"
			  }
			]
		  },
		  "LTC/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "LTC",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "LTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "LTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "LTC-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "ltcusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XLTCZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "LTC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "LTCUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "LTC-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "LTC_USD"
			  }
			]
		  },
		  "MANA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MANA",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "MANAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "MANA-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "MANA_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "MANAUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "MANA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "MANAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "MANA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "POL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "POL",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "POLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "POLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "POL-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "POL_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "polusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "POLUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "POL-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "POLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "POL-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "POL_USD"
			  }
			]
		  },
		  "MKR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MKR",
				"Quote": "USD"
			  },
			  "decimals": 6,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "MKRUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "MKR-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "MKRUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "MKR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "MKRUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "MKR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "NEAR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "NEAR",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "NEARUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "NEAR-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "NEAR_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "nearusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "NEAR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "NEARUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "NEAR-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "NEAR_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "NTRN/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "NTRN",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 2,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "NTRNUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "NTRN_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "NTRN-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "NTRN-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "OP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "OP",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "OPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "OP-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "OP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "OP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "OPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "OP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "OP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "ORDI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ORDI",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "ORDIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "ORDIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "ORDI_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "ordiusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "ORDI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "ORDI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "ORDIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "PEPE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PEPE",
				"Quote": "USD"
			  },
			  "decimals": 16,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "PEPEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "PEPEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "PEPE_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "PEPEUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "PEPE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "PEPEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "PEPE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "PEPE_USD"
			  }
			]
		  },
		  "PYTH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PYTH",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "PYTHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "PYTHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "PYTH_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "PYTH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "PYTH-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "PYTHUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "RUNE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "RUNE",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "RUNEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "RUNE_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "RUNEUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "RUNE-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "RUNEUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "SEI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SEI",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "SEIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "SEIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "SEI-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "SEI_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "seiusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "SEI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "SEIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "SHIB/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SHIB",
				"Quote": "USD"
			  },
			  "decimals": 15,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "SHIBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "SHIBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "SHIB-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "SHIB_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "SHIBUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "SHIB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "SHIBUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "SHIB-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "SHIB_USD"
			  }
			]
		  },
		  "SNX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SNX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "SNXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "SNXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "SNX-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "SNXUSD"
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "SNXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "SNX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "SOL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SOL",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "SOLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "SOLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "SOL-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "solusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "SOLUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "SOL-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "SOLUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "SOL-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "SOL_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "STRK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "STRK",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "STRKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "STRKUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "STRKUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "STRK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "STRK-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "STRK_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "STX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "STX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "STXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "STXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "STX-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "STX_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "STXUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "STX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "STX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "STXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "STX_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "SUI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SUI",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "SUIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "SUIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "SUI-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "SUI_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "suiusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "SUI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "SUIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "SUI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "SUI_USD"
			  }
			]
		  },
		  "TIA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TIA",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "TIAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "TIAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "TIA-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "tiausdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "TIAUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "TIA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "TIAUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "TIA-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "TIA_USD"
			  }
			]
		  },
		  "TRX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TRX",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "TRXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "TRXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "TRX_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "trxusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "TRXUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "TRX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "TRXUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "TRX-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "UNI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "UNI",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "UNIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "UNIUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "UNI-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "UNI_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "UNIUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "UNI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "UNI-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "UNI_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "USDT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "USDT",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "USDCUSDT",
				"invert": true
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "USDCUSDT",
				"invert": true
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "USDT-USD"
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "ethusdt",
				"normalize_by_pair": {
				  "Base": "ETH",
				  "Quote": "USD"
				},
				"invert": true
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "USDTZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "BTC-USDT",
				"normalize_by_pair": {
				  "Base": "BTC",
				  "Quote": "USD"
				},
				"invert": true
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "USDC-USDT",
				"invert": true
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "USDT_USD"
			  }
			]
		  },
		  "WLD/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "WLD",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "WLDUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "WLDUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "WLD_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "wldusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "WLD-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "WLDUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "WLD-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "WLD_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "WOO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "WOO",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "WOOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "WOO_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "WOO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "WOO-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "WOOUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "XLM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "XLM",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "XLMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "XLMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "XLM-USD"
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XXLMZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "XLM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "XLMUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "XLM-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "XLM_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  }
			]
		  },
		  "XRP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "XRP",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 3,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "binance_ws",
				"off_chain_ticker": "XRPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "bybit_ws",
				"off_chain_ticker": "XRPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "coinbase_ws",
				"off_chain_ticker": "XRP-USD"
			  },
			  {
				"name": "gate_ws",
				"off_chain_ticker": "XRP_USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "huobi_ws",
				"off_chain_ticker": "xrpusdt",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "kraken_api",
				"off_chain_ticker": "XXRPZUSD"
			  },
			  {
				"name": "kucoin_ws",
				"off_chain_ticker": "XRP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "mexc_ws",
				"off_chain_ticker": "XRPUSDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "okx_ws",
				"off_chain_ticker": "XRP-USDT",
				"normalize_by_pair": {
				  "Base": "USDT",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "crypto_dot_com_ws",
				"off_chain_ticker": "XRP_USD"
			  }
			]
		  }
		}
	}
	  `

	// UniswapV3BaseMarketMap is used to initialize the Base market map. This only includes
	// the markets that are supported by uniswapv3 on Base.
	UniswapV3BaseMarketMap mmtypes.MarketMap

	// UniswapV3BaseMarketMapJSON is the JSON representation of UniswapV3BaseMarketMap.
	UniswapV3BaseMarketMapJSON = `
	{
		"markets": {
		  "BRETT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BRETT",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "uniswapv3_api-base",
				"off_chain_ticker": "BRETT/ETH",
				"metadata_JSON": "{\"address\":\"0xBA3F945812a83471d709BCe9C3CA699A19FB46f7\",\"base_decimals\":18,\"quote_decimals\":18,\"invert\":true}",
				"normalize_by_pair": {
					"Base": "ETH",
					"Quote": "USD"
				}
			  }
			]
		  },
		  "ETH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ETH",
				"Quote": "USD"
			  },
			  "decimals": 18,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "uniswapv3_api-base",
				"off_chain_ticker": "ETH/USDT",
				"metadata_JSON": "{\"address\":\"0xd92E0767473D1E3FF11Ac036f2b1DB90aD0aE55F\",\"base_decimals\":18,\"quote_decimals\":6,\"invert\":false}",
				"normalize_by_pair": {
					"Base": "USDT",
					"Quote": "USD"
				}
			  }
			]
		  },
		  "USDT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "USDT",
				"Quote": "USD"
			  },
			  "decimals": 6,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "uniswapv3_api-base",
				"off_chain_ticker": "USDT/ETH",
				"metadata_JSON": "{\"address\":\"0xd92E0767473D1E3FF11Ac036f2b1DB90aD0aE55F\",\"base_decimals\":6,\"quote_decimals\":18,\"invert\":true}",
				"normalize_by_pair": {
				  "Base": "ETH",
				  "Quote": "USD"
				}
			  },
			  {
				"name": "uniswapv3_api-base",
				"off_chain_ticker": "USDT/USDC",
				"metadata_JSON": "{\"address\":\"0xD56da2B74bA826f19015E6B7Dd9Dae1903E85DA1\",\"base_decimals\":6,\"quote_decimals\":6,\"invert\":true}"
			  }
			]
		  }
	    }
	}
	`

	// CoinGeckoMarketMap is used to initialize the CoinGecko market map. This only includes
	// the markets that are supported by CoinGecko & are included in the Core market map.
	CoinGeckoMarketMap mmtypes.MarketMap

	// CoinGeckoMarketMapJSON is the JSON representation of CoinGeckoMarketMap.
	CoinGeckoMarketMapJSON = `
	{
		"markets": {
		  "KHAI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "KHAI",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "kitten-haimer/usd"
			  }
			]
		  },
		  "WAFFLES/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "WAFFLES",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "waffles/usd"
			  }
			]
		  },
		  "HEGE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "HEGE",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "hege/usd"
			  }
			]
		  },
		  "WUF/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "WUF",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "wuffi/usd"
			  }
			]
		  },
		  "CHAT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CHAT",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "solchat/usd"
			  }
			]
		  },
		  "BEER/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BEER",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "beercoin-2/usd"
			  }
			]
		  },
		  "MANEKI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MANEKI",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "maneki/usd"
			  }
			]
		  },
		  "SLERF/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SLERF",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "slerf/usd"
			  }
			]
		  },
		  "MYRO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MYRO",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "myro/usd"
			  }
			]
		  },
		  "RAY/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "RAY",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "raydium/usd"
			  }
			]
		  },
		  "WIF/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "WIF",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "dogwifcoin/usd"
			  }
			]
		  },
		  "MICHI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MICHI",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "michicoin/usd"
			  }
			]
		  },
		  "MEW/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MEW",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "cat-in-a-dogs-world/usd"
			  }
			]
		  },
		  "PONKE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PONKE",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "ponke/usd"
			  }
			]
		  },
		  "USA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "USA",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "american-coin/usd"
			  }
			]
		  },
		  "BOME/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BOME",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "book-of-meme/usd"
			  }
			]
		  },
		  "GME/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "GME",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "gme/usd"
			  }
			]
		  },
		  "DJT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DJT",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "trumpcoin/usd"
			  }
			]
		  },
		  "POPCAT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "POPCAT",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "popcat/usd"
			  }
			]
		  },
		  "$RETIRE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "$RETIRE",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "retire-on-sol/usd"
			  }
			]
		  },
		  "AAVE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AAVE",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "aave/usd"
			  }
			]
		  },
		  "ADA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ADA",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "cardano/usd"
			  }
			]
		  },
		  "AEVO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AEVO",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "aevo-exchange/usd"
			  }
			]
		  },
		  "ALGO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ALGO",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "algorand/usd"
			  }
			]
		  },
		  "APE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "APE",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "apecoin/usd"
			  }
			]
		  },
		  "APT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "APT",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "aptos/usd"
			  }
			]
		  },
		  "ARB/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ARB",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "arbitrum/usd"
			  }
			]
		  },
		  "ARKM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ARKM",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "arkham/usd"
			  }
			]
		  },
		  "ASTR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ASTR",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "astar/usd"
			  }
			]
		  },
		  "ATOM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ATOM",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "cosmos/usd"
			  }
			]
		  },
		  "AVAX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AVAX",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "avalanche-2/usd"
			  }
			]
		  },
		  "AXL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "AXL",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "axelar/usd"
			  }
			]
		  },
		  "BAZINGA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BAZINGA",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "bazinga-2/usd"
			  }
			]
		  },
		  "BCH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BCH",
				"Quote": "USD"
			  },
			  "decimals": 7,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "bitcoin-cash/usd"
			  }
			]
		  },
		  "BENDOG/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BENDOG",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "ben-the-dog/usd"
			  }
			]
		  },
		  "BLUR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BLUR",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "blur/usd"
			  }
			]
		  },
		  "BNB/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BNB",
				"Quote": "USD"
			  },
			  "decimals": 7,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "binancecoin/usd"
			  }
			]
		  },
		  "BONK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BONK",
				"Quote": "USD"
			  },
			  "decimals": 14,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "bonk/usd"
			  }
			]
		  },
		  "BTC/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BTC",
				"Quote": "USD"
			  },
			  "decimals": 5,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "bitcoin/usd"
			  }
			]
		  },
		  "BUBBA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "BUBBA",
				"Quote": "USD"
			  },
			  "decimals": 12,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "bubba/usd"
			  }
			]
		  },
		  "CATGPT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CATGPT",
				"Quote": "USD"
			  },
			  "decimals": 12,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "catgpt/usd"
			  }
			]
		  },
		  "CHEESE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CHEESE",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "cheese-2/usd"
			  }
			]
		  },
		  "CHUD/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CHUD",
				"Quote": "USD"
			  },
			  "decimals": 12,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "chudjak/usd"
			  }
			]
		  },
		  "COK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "COK",
				"Quote": "USD"
			  },
			  "decimals": 14,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "catownkimono/usd"
			  }
			]
		  },
		  "COST/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "COST",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "costco-hot-dog/usd"
			  }
			]
		  },
		  "COMP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "COMP",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "compound-governance-token/usd"
			  }
			]
		  },
		  "CRV/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "CRV",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "curve-dao-token/usd"
			  }
			]
		  },
		  "DOGE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DOGE",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "dogecoin/usd"
			  }
			]
		  },
		  "DOT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DOT",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "polkadot/usd"
			  }
			]
		  },
		  "DEVIN/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DEVIN",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "devin-on-solana/usd"
			  }
			]
		  },
		  "DUKO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DUKO",
				"Quote": "USD"
			  },
			  "decimals": 12,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "duko/usd"
			  }
			]
		  },
		  "DYDX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DYDX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "dydx/usd"
			  }
			]
		  },
		  "DYM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "DYM",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "dymension/usd"
			  }
			]
		  },
		  "EGG/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "EGG",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "justanegg-2/usd"
			  }
			]
		  },
		  "FALX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "FALX",
				"Quote": "USD"
			  },
			  "decimals": 12,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "falx/usd"
			  }
			]
		  },
		  "GOL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "GOL",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "golazo-world/usd"
			  }
			]
		  },
		  "GUMMY/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "GUMMY",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "gummy/usd"
			  }
			]
		  },
		  "HABIBI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "HABIBI",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "habibi-sol/usd"
			  }
			]
		  },
		  "HAMMY/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "HAMMY",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "sad-hamster/usd"
			  }
			]
		  },
		  "HARAMBE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "HARAMBE",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "harambe-2/usd"
			  }
			]
		  },
		  "KITTY/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "KITTY",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "roaring-kitty-solana/usd"
			  }
			]
		  },
		  "MUMU/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MUMU",
				"Quote": "USD"
			  },
			  "decimals": 14,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "mumu-the-bull-3/usd"
			  }
			]
		  },
		  "NUB/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "NUB",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "sillynubcat/usd"
			  }
			]
		  },
		  "PAJAMAS/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PAJAMAS",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "pajamas-cat/usd"
			  }
			]
		  },
		  "PENG/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PENG",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "peng/usd"
			  }
			]
		  },
		  "PUNDU/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PUNDU",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "pundu/usd"
			  }
			]
		  },
		  "RETARDIO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "RETARDIO",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "retardio/usd"
			  }
			]
		  },
		  "SELFIE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SELFIE",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "selfiedogcoin/usd"
			  }
			]
		  },
		  "SLOTH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SLOTH",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "slothana/usd"
			  }
			]
		  },
		  "SPEED/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SPEED",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "real-fast/usd"
			  }
			]
		  },
		  "SPIKE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SPIKE",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "spike/usd"
			  }
			]
		  },
		  "STACKS/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "STACKS",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "slap-city/usd"
			  }
			]
		  },
		  "TREMP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TREMP",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "donald-tremp/usd"
			  }
			]
		  },
		  "MOG/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MOG",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "mog-coin/usd"
			  }
			]
		  },
		  "TRUMP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TRUMP",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "maga/usd"
			  }
			]
		  },
		  "MOTHER/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MOTHER",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "mother-iggy/usd"
			  }
			]
		  },
		  "EOS/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "EOS",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "eos/usd"
			  }
			]
		  },
		  "ETC/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ETC",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "ethereum-classic/usd"
			  }
			]
		  },
		  "ETH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ETH",
				"Quote": "USD"
			  },
			  "decimals": 6,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "ethereum/usd"
			  }
			]
		  },
		  "FET/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "FET",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "fetch-ai/usd"
			  }
			]
		  },
		  "FIL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "FIL",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "filecoin/usd"
			  }
			]
		  },
		  "GRT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "GRT",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "the-graph/usd"
			  }
			]
		  },
		  "HBAR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "HBAR",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "hedera-hashgraph/usd"
			  }
			]
		  },
		  "ICP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ICP",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "internet-computer/usd"
			  }
			]
		  },
		  "IMX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "IMX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "immutable-x/usd"
			  }
			]
		  },
		  "INJ/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "INJ",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "injective-protocol/usd"
			  }
			]
		  },
		  "JTO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "JTO",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "jito-governance-token/usd"
			  }
			]
		  },
		  "JUP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "JUP",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "jupiter-exchange-solana/usd"
			  }
			]
		  },
		  "LDO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "LDO",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "lido-dao/usd"
			  }
			]
		  },
		  "LINK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "LINK",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "chainlink/usd"
			  }
			]
		  },
		  "LTC/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "LTC",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "litecoin/usd"
			  }
			]
		  },
		  "MANA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MANA",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "decentraland/usd"
			  }
			]
		  },
		  "POL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "POL",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "pol-network/usd"
			  }
			]
		  },
		  "MKR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "MKR",
				"Quote": "USD"
			  },
			  "decimals": 6,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "maker/usd"
			  }
			]
		  },
		  "NEAR/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "NEAR",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "near/usd"
			  }
			]
		  },
		  "NTRN/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "NTRN",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "neutron-3/usd"
			  }
			]
		  },
		  "OP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "OP",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "optimism/usd"
			  }
			]
		  },
		  "ORDI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "ORDI",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "ordinals/usd"
			  }
			]
		  },
		  "PEPE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PEPE",
				"Quote": "USD"
			  },
			  "decimals": 16,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "pepe/usd"
			  }
			]
		  },
		  "PYTH/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "PYTH",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "pyth-network/usd"
			  }
			]
		  },
		  "RUNE/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "RUNE",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "thorchain/usd"
			  }
			]
		  },
		  "SEI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SEI",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "sei-network/usd"
			  }
			]
		  },
		  "SHIB/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SHIB",
				"Quote": "USD"
			  },
			  "decimals": 15,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "shiba-inu/usd"
			  }
			]
		  },
		  "SNX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SNX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "havven/usd"
			  }
			]
		  },
		  "SOL/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SOL",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "solana/usd"
			  }
			]
		  },
		  "STRK/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "STRK",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "starknet/usd"
			  }
			]
		  },
		  "STX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "STX",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "blockstack/usd"
			  }
			]
		  },
		  "SUI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "SUI",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "sui/usd"
			  }
			]
		  },
		  "TIA/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TIA",
				"Quote": "USD"
			  },
			  "decimals": 8,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "celestia/usd"
			  }
			]
		  },
		  "TRX/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "TRX",
				"Quote": "USD"
			  },
			  "decimals": 11,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "tron/usd"
			  }
			]
		  },
		  "UNI/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "UNI",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "uniswap/usd"
			  }
			]
		  },
		  "USDT/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "USDT",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "tether/usd"
			  }
			]
		  },
		  "WLD/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "WLD",
				"Quote": "USD"
			  },
			  "decimals": 9,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "worldcoin-wld/usd"
			  }
			]
		  },
		  "WOO/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "WOO",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "woo-network/usd"
			  }
			]
		  },
		  "XLM/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "XLM",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "stellar/usd"
			  }
			]
		  },
		  "XRP/USD": {
			"ticker": {
			  "currency_pair": {
				"Base": "XRP",
				"Quote": "USD"
			  },
			  "decimals": 10,
			  "min_provider_count": 1,
			  "enabled": true
			},
			"provider_configs": [
			  {
				"name": "coingecko_api",
				"off_chain_ticker": "ripple/usd"
			  }
			]
		  }
		}
	  } 
	`

	// OsmosisMarketMap is used to initialize the osmosis market map. This only includes
	// the markets that are supported by osmosis.
	OsmosisMarketMap mmtypes.MarketMap

	// OsmosisMarketMapJSON is the JSON representation of OsmosisMarketMap.
	OsmosisMarketMapJSON = `
{
    "markets": {
        "STARS/USD": {
            "ticker": {
                "currency_pair": {
                    "Base": "STARS",
                    "Quote": "USD"
                },
                "decimals": 18,
                "min_provider_count": 1,
                "enabled": true,
                "metadata_JSON": "{\"reference_price\":1,\"liquidity\":0,\"aggregate_ids\":[]}"
            },
            "provider_configs": [
                {
                    "name": "osmosis_api",
                    "off_chain_ticker": "STARS/OSMO",
                    "metadata_JSON": "{\"pool_id\":1096,\"base_token_denom\":\"ibc/987C17B11ABC2B20019178ACE62929FE9840202CE79498E29FE8E5CB02B7C0A4\",\"quote_token_denom\":\"uosmo\"}",
                    "normalize_by_pair": {
                        "Base": "OSMO",
                        "Quote": "USD"
                    }
                }
            ]
        },
        "USDT/USD": {
            "ticker": {
                "currency_pair": {
                    "Base": "USDT",
                    "Quote": "USD"
                },
                "decimals": 9,
                "min_provider_count": 1,
                "enabled": true
            },
            "provider_configs": [
                {
                    "name": "binance_ws",
                    "off_chain_ticker": "USDCUSDT",
                    "invert": true
                },
                {
                    "name": "bybit_ws",
                    "off_chain_ticker": "USDCUSDT",
                    "invert": true
                },
                {
                    "name": "coinbase_ws",
                    "off_chain_ticker": "USDT-USD"
                },
                {
                    "name": "kraken_api",
                    "off_chain_ticker": "USDTZUSD"
                },
                {
                    "name": "okx_ws",
                    "off_chain_ticker": "USDC-USDT",
                    "invert": true
                },
                {
                    "name": "crypto_dot_com_ws",
                    "off_chain_ticker": "USDT_USD"
                }
            ]
        },
        "OSMO/USD": {
            "ticker": {
                "currency_pair": {
                    "Base": "OSMO",
                    "Quote": "USD"
                },
                "decimals": 8,
                "min_provider_count": 1,
                "enabled": true,
                "metadata_JSON": "{\"reference_price\":1,\"liquidity\":0,\"aggregate_ids\":[]}"
            },
            "provider_configs": [
                {
                    "name": "coinbase_ws",
                    "off_chain_ticker": "OSMO-USD"
                },
                {
                    "name": "huobi_ws",
                    "off_chain_ticker": "osmousdt",
                    "normalize_by_pair": {
                        "Base": "USDT",
                        "Quote": "USD"
                    }
                },
                {
                    "name": "binance_api",
                    "off_chain_ticker": "OSMOUSDT",
                    "normalize_by_pair": {
                        "Base": "USDT",
                        "Quote": "USD"
                    }
                }
            ]
        }
    }
}
	`

	// PolymarketMarketMap is used to initialize the Polymarket market map. This only includes one prediction market
	// with one outcome token.
	PolymarketMarketMap mmtypes.MarketMap

	// PolymarketMarketMapJSON is the JSON representation of PolymarketMarketMap.
	PolymarketMarketMapJSON = ` 
{
   "markets":{
      "WILL_BERNIE_SANDERS_WIN_THE_2024_US_PRESIDENTIAL_ELECTION?YES/USD":{
         "ticker":{
            "currency_pair":{
               "Base":"WILL_BERNIE_SANDERS_WIN_THE_2024_US_PRESIDENTIAL_ELECTION?YES",
               "Quote":"USD"
            },
            "decimals":4,
            "min_provider_count":1,
            "enabled":true
         },
         "provider_configs":[
            {
               "name":"polymarket_api",
               "off_chain_ticker":"0x08f5fe8d0d29c08a96f0bc3dfb52f50e0caf470d94d133d95d38fa6c847e0925/95128817762909535143571435260705470642391662537976312011260538371392879420759"
            }
         ]
      },
      "WILL_INSIDE_OUT_2_GROSS_MOST_IN_2024?YES/USD":{
         "ticker":{
            "currency_pair":{
               "Base":"WILL_INSIDE_OUT_2_GROSS_MOST_IN_2024?YES",
               "Quote":"USD"
            },
            "decimals":4,
            "min_provider_count":1,
            "enabled":true
         },
         "provider_configs":[
            {
               "name":"polymarket_api",
               "off_chain_ticker":"0x1ab07117f9f698f28490f57754d6fe5309374230c95867a7eba572892a11d710/50107902083284751016545440401692219408556171231461347396738260657226842527986"
            }
         ]
      }
   }
}`

	// ForexMarketMap is used to initialize the forex market map. This only includes
	// forex markets quoted in usdt.
	ForexMarketMap mmtypes.MarketMap

	// ForexMarketMapJSON is the JSON representation of ForexMarketMap.
	ForexMarketMapJSON = `
{
    "markets": {
      "TRY/USDT": {
        "ticker": {
          "currency_pair": {
            "Base": "TRY",
            "Quote": "USDT"
          },
          "decimals": 11,
          "min_provider_count": 1,
          "enabled": false,
          "metadata_JSON": "{\"reference_price\":2935133548,\"liquidity\":1504939,\"aggregate_ids\":[{\"venue\":\"coinmarketcap\",\"ID\":\"2810\"}]}"
        },
        "provider_configs": [
          {
            "name": "binance_ws",
            "off_chain_ticker": "USDTTRY",
            "invert": true,
            "metadata_JSON": ""
          },
          {
            "name": "okx_ws",
            "off_chain_ticker": "TRY-USDT",
            "invert": false,
            "metadata_JSON": ""
          }
        ]
      },
      "EUR/USDT": {
        "ticker": {
          "currency_pair": {
            "Base": "EUR",
            "Quote": "USDT"
          },
          "decimals": 9,
          "min_provider_count": 1,
          "enabled": false,
          "metadata_JSON": "{\"reference_price\":1100800000,\"liquidity\":3843298,\"aggregate_ids\":[{\"venue\":\"coinmarketcap\",\"ID\":\"2790\"}]}"
        },
        "provider_configs": [
          {
            "name": "binance_ws",
            "off_chain_ticker": "EURUSDT",
            "invert": false,
            "metadata_JSON": ""
          },
          {
            "name": "okx_ws",
            "off_chain_ticker": "EUR-USDT",
            "invert": false,
            "metadata_JSON": ""
          }
        ]
      },
      "BRL/USDT": {
        "ticker": {
          "currency_pair": {
            "Base": "BRL",
            "Quote": "USDT"
          },
          "decimals": 10,
          "min_provider_count": 1,
          "enabled": false,
          "metadata_JSON": "{\"reference_price\":1760563380,\"liquidity\":2479974,\"aggregate_ids\":[{\"venue\":\"coinmarketcap\",\"ID\":\"2783\"}]}"
        },
        "provider_configs": [
          {
            "name": "binance_ws",
            "off_chain_ticker": "USDTBRL",
            "invert": true,
            "metadata_JSON": ""
          },
          {
            "name": "okx_ws",
            "off_chain_ticker": "BRL-USDT",
            "invert": false,
            "metadata_JSON": ""
          }
        ]
      }
    }
  }`
)

func init() {
	err := errors.Join(
		unmarshalValidate("CoinMarketCap", CoinMarketCapMarketMapJSON, &CoinMarketCapMarketMap),
		unmarshalValidate("Raydium", RaydiumMarketMapJSON, &RaydiumMarketMap),
		unmarshalValidate("Core", CoreMarketMapJSON, &CoreMarketMap),
		unmarshalValidate("UniswapV3Base", UniswapV3BaseMarketMapJSON, &UniswapV3BaseMarketMap),
		unmarshalValidate("CoinGecko", CoinGeckoMarketMapJSON, &CoinGeckoMarketMap),
		unmarshalValidate("Osmosis", OsmosisMarketMapJSON, &OsmosisMarketMap),
		unmarshalValidate("Polymarket", PolymarketMarketMapJSON, &PolymarketMarketMap),
		unmarshalValidate("Forex", ForexMarketMapJSON, &ForexMarketMap),
	)
	if err != nil {
		panic(err)
	}
}

// unmarshalValidate unmarshalls data into mm and then calls ValidateBasic.
func unmarshalValidate(name, data string, mm *mmtypes.MarketMap) error {
	if err := json.Unmarshal([]byte(data), mm); err != nil {
		return fmt.Errorf("failed to unmarshal %sMarketMap: %w", name, err)
	}
	if err := mm.ValidateBasic(); err != nil {
		return fmt.Errorf("%sMarketMap failed validation: %w", name, err)
	}
	return nil
}
