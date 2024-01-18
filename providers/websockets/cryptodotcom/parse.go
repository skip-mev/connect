package cryptodotcom

import (
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// parseInstrumentMessage is used to parse an instrument message received from the Crypto.com
// web socket API. This message contains the latest price data for a set of instruments.
func (h *WebSocketDataHandler) parseInstrumentMessage(
	msg InstrumentResponseMessage,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], error) {
	var (
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved = make(map[oracletypes.CurrencyPair]error)
		markets    = msg.Result.Data
	)

	// If the response contained no instrument data, return an error.
	if len(markets) == 0 {
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved),
			fmt.Errorf("no instrument data was returned")
	}

	// Iterate through each market and attempt to parse the price.
	for _, market := range markets {
		// If we don't have a mapping for the instrument, return an error. This is likely a configuration
		// error.
		cp, ok := h.config.ReverseCache[market.InstrumentName]
		if !ok {
			h.logger.Error("failed to find currency pair for instrument", zap.String("instrument", market.InstrumentName))
			continue
		}

		// Attempt to parse the price.
		if price, err := math.Float64StringToBigInt(market.LatestTradePrice, cp.Decimals()); err != nil {
			unresolved[cp] = fmt.Errorf("failed to parse price %s: %s", market.LatestTradePrice, err)
		} else {
			resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
		}

	}

	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved), nil
}
