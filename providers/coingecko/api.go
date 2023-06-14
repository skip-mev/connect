package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/oracle/types"
)

// getPriceForPair returns the price of a currency pair.
//
// Response format:
//
//	{
//	  "atom": {
//	    "usd": 11.35
//	  },
//	  "btc": {
//	    "usd": 10000
//	  }
//	}
func (p *Provider) getPrices() (map[types.CurrencyPair]types.TickerPrice, error) {
	url := getPriceEndpoint(strings.Join(p.bases, ","), strings.Join(p.quotes, ","))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respMap := make(map[string]map[string]float64)
	if err := json.Unmarshal(body, &respMap); err != nil {
		return nil, err
	}

	prices := make(map[types.CurrencyPair]types.TickerPrice)

	for base, data := range respMap {
		for quote, price := range data {
			base = strings.ToUpper(base)
			quote = strings.ToUpper(quote)

			if _, ok := p.cache[base]; !ok {
				continue
			}

			if _, ok := p.cache[base][quote]; !ok {
				continue
			}

			cp := types.CurrencyPair{
				Base:  base,
				Quote: quote,
			}

			price := types.TickerPrice{
				Price:     float64ToDec(price),
				Timestamp: time.Now(),
			}

			prices[cp] = price
		}
	}

	return prices, nil
}

// float64ToDec converts a float64 to a sdk.Dec.
func float64ToDec(f float64) sdk.Dec {
	float := strconv.FormatFloat(f, 'g', 10, 64)
	return sdk.MustNewDecFromStr(float)
}

// getPriceEndpoint is the CoinGecko endpoint for getting the price of a
// currency pair.
func getPriceEndpoint(base, quote string) string {
	return fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", base, quote)
}
