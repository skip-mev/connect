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
	// Config is the Coinbase config.
	Config
}

// NewCoinBaseAPIHandler returns a new Coinbase APIDataHandler.
func NewCoinBaseAPIHandler(
	providerCfg config.ProviderConfig,
) (*CoinBaseAPIHandler, error) {
	if providerCfg.Name != Name {
		return nil, fmt.Errorf("expected provider config name %s, got %s", Name, providerCfg.Name)
	}

	cfg, err := ReadCoinbaseConfigFromFile(providerCfg.Path)
	if err != nil {
		return nil, err
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
	base, ok := h.SymbolMap[cp.Base]
	if !ok {
		return "", fmt.Errorf("unknown base currency %s", cp.Base)
	}

	quote, ok := h.SymbolMap[cp.Quote]
	if !ok {
		return "", fmt.Errorf("unknown quote currency %s", cp.Quote)
	}

	return fmt.Sprintf(BaseURL, base, quote), nil
}

// Atomic returns true as this API handler must request data from the Coinbase API
// sequentially for each currency pair.
func (h *CoinBaseAPIHandler) Atomic() bool {
	return false
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

	// If the response quote currency does not match the requested quote currency, return an error.
	cp := cps[0]
	quote, ok := h.SymbolMap[cp.Quote]
	if !ok {
		return providertypes.NewGetResponseWithErr[oracletypes.CurrencyPair, *big.Int](
			cps,
			fmt.Errorf("unknown quote currency %s", cp.Quote),
		)
	}

	// Parse the response into a CoinBaseResponse.
	var result CoinBaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return providertypes.NewGetResponseWithErr[oracletypes.CurrencyPair, *big.Int](cps, err)
	}

	if quote != result.Data.Currency {
		return providertypes.NewGetResponseWithErr[oracletypes.CurrencyPair, *big.Int](
			cps,
			fmt.Errorf("expected quote currency %s, got %s", cp.Quote, result.Data.Currency),
		)
	}

	// Convert the float64 price into a big.Int.
	floatAmount := result.Data.Amount
	price, err := math.Float64StringToBigInt(floatAmount, cp.Decimals())
	if err != nil {
		return providertypes.NewGetResponseWithErr[oracletypes.CurrencyPair, *big.Int](cps, err)
	}

	resolved := map[oracletypes.CurrencyPair]providertypes.Result[*big.Int]{
		cp: providertypes.NewResult[*big.Int](price, time.Now()),
	}

	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, nil)
}

// Name returns the name of the handler.
func (h *CoinBaseAPIHandler) Name() string {
	return Name
}
