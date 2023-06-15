package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/oracle/types"
)

// getPriceForPair returns the price of a currency pair. The price is fetched
// from the CoinGecko API in a single request for all pairs. Since the CoinGecko
// response will match some base denoms to quote denoms that should not be supported,
// we filter out pairs that are not supported by the provider.
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
	url := getPriceEndpoint(p.bases, p.quotes)

	resp, err := http.Get(url) //nolint:all
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

	for _, pair := range p.pairs {
		if _, ok := respMap[pair.Base]; !ok {
			continue
		}

		if _, ok := respMap[pair.Base][pair.Quote]; !ok {
			continue
		}

		price := float64ToDec(respMap[pair.Base][pair.Quote])
		prices[pair] = types.TickerPrice{
			Price:     price,
			Timestamp: time.Now(),
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
