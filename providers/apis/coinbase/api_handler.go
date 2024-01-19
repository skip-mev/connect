package coinbase

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/skip-mev/slinky/oracle/config"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/api/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// Name is the name of the provider.
	Name = "coinbase"
)

var _ handlers.APIDataHandler[oracletypes.CurrencyPair, *big.Int] = (*CoinBaseAPIHandler)(nil)

// CoinBaseAPIHandler implements the APIDataHandler interface for Coinbase, which can be used
// by a base provider. The DataHandler fetches data from the spot price Coinbase API. It is
// atomic in that it must request data from the Coinbase API sequentially for each currency pair.
type CoinBaseAPIHandler struct { //nolint
	cfg config.ProviderConfig
}

// NewCoinBaseAPIHandler returns a new Coinbase APIDataHandler.
func NewCoinBaseAPIHandler(
	cfg config.ProviderConfig,
) (*CoinBaseAPIHandler, error) {
	if err := cfg.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("invalid provider config %s", err)
	}

	if !cfg.API.Enabled {
		return nil, fmt.Errorf("api is not enabled for provider %s", cfg.Name)
	}

	if cfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, cfg.Name)
	}

	return &CoinBaseAPIHandler{
		cfg,
	}, nil
}

// CreateURL returns the URL that is used to fetch data from the Coinbase API for the
// given currency pair. Since the Coinbase API only supports fetching spot prices for
// a single currency pair at a time, this function will return an error if the currency
// pair slice contains more than one currency pair.
func (h *CoinBaseAPIHandler) CreateURL(
	cps []oracletypes.CurrencyPair,
) (string, error) {
	if len(cps) != 1 {
		return "", fmt.Errorf("expected 1 currency pair, got %d", len(cps))
	}

	// Ensure that the base and quote currencies are supported by the Coinbase API and
	// are configured for the handler.
	cp := cps[0]
	market, ok := h.cfg.MarketConfig.CurrencyPairToMarketConfigs[cp.ToString()]
	if !ok {
		return "", fmt.Errorf("unknown currency pair %s", cp)
	}

	return fmt.Sprintf(h.cfg.API.URL, market.Ticker), nil
}

// ParseResponse parses the spot price HTTP response from the Coinbase API and returns
// the resulting price. Note that this can only parse a single currency pair at a time.
func (h *CoinBaseAPIHandler) ParseResponse(
	cps []oracletypes.CurrencyPair,
	resp *http.Response,
) providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int] {
	if len(cps) != 1 {
		return providertypes.NewGetResponseWithErr[oracletypes.CurrencyPair, *big.Int](
			cps,
			fmt.Errorf("expected 1 currency pair, got %d", len(cps)),
		)
	}

	// Check if this currency pair is supported by the Coinbase API.
	cp := cps[0]
	_, ok := h.cfg.MarketConfig.CurrencyPairToMarketConfigs[cp.ToString()]
	if !ok {
		return providertypes.NewGetResponseWithErr[oracletypes.CurrencyPair, *big.Int](
			cps,
			fmt.Errorf("unknown currency pair %s", cp),
		)
	}

	// Parse the response into a CoinBaseResponse.
	var result CoinBaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return providertypes.NewGetResponseWithErr[oracletypes.CurrencyPair, *big.Int](cps, err)
	}

	// Convert the float64 price into a big.Int.
	price, err := math.Float64StringToBigInt(result.Data.Amount, cp.Decimals())
	if err != nil {
		return providertypes.NewGetResponseWithErr[oracletypes.CurrencyPair, *big.Int](cps, err)
	}

	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](
		map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
			cp: providertypes.NewResult[*big.Int](price, time.Now()),
		},
		nil,
	)
}
