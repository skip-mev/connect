package cryptodotcom

import (
	"fmt"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"go.uber.org/zap"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
)

// parseInstrumentMessage is used to parse an instrument message received from the Crypto.com
// websocket API. This message contains the latest price data for a set of instruments.
func (h *WebSocketHandler) parseInstrumentMessage(
	msg InstrumentResponseMessage,
) (types.PriceResponse, error) {
	var (
		resolved    = make(types.ResolvedPrices)
		unresolved  = make(types.UnResolvedPrices)
		instruments = msg.Result.Data
	)

	// If the response contained no instrument data, return an error.
	if len(instruments) == 0 {
		return types.NewPriceResponse(resolved, unresolved),
			fmt.Errorf("no instrument data was returned")
	}

	// Iterate through each market and attempt to parse the price.
	for _, instrument := range instruments {
		// If we don't have a mapping for the instrument, return an error. This is likely a configuration
		// error.
		ticker, ok := h.cache.FromOffChainTicker(instrument.Name)
		if !ok {
			h.logger.Debug("failed to find currency pair for instrument", zap.String("instrument", instrument.Name))
			continue
		}

		// Attempt to parse the price.
		if price, err := math.Float64StringToBigFloat(instrument.LatestTradePrice); err != nil {
			wErr := fmt.Errorf("failed to parse price %s:"+" %w", instrument.LatestTradePrice, err)
			unresolved[ticker] = providertypes.UnresolvedResult{
				ErrorWithCode: providertypes.NewErrorWithCode(wErr, providertypes.ErrorFailedToParsePrice),
			}
		} else {
			resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
		}

	}

	return types.NewPriceResponse(resolved, unresolved), nil
}
