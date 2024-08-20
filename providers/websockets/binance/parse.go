package binance

import (
	"fmt"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
)

// parsePriceUpdateMessage parses a price update message from the Binance websocket feed.
// This is repurposed for ticker and aggregate trade messages.
func (h *WebSocketHandler) parsePriceUpdateMessage(offChainTicker string, price string) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	ticker, ok := h.cache.FromOffChainTicker(offChainTicker)
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("got response for an unsupported market %s", offChainTicker)
	}

	// Convert the price to a big Float.
	priceFloat, err := math.Float64StringToBigFloat(price)
	if err != nil {
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToParsePrice),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	resolved[ticker] = types.NewPriceResult(priceFloat, time.Now().UTC())
	return types.NewPriceResponse(resolved, unResolved), nil
}
