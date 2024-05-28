package binance

import (
	"fmt"
	"strings"
	"time"

	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
)

// parseAggregateTradeMessage parses an aggregate trade message from the Binance websocket feed.
// This message is sent when a trade is executed on the Binance exchange.
func (h *WebSocketHandler) parseAggregateTradeMessage(
	msg StreamMessageResponse,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// Determine if the stream type is valid.
	if !strings.HasSuffix(msg.Stream, string(AggregateTradeStream)) {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("invalid stream type %s; expected %s", msg.Stream, string(AggregateTradeStream))
	}

	// Determine if the ticker is valid.
	aggMsg := msg.Data
	ticker, ok := h.cache.FromOffChainTicker(aggMsg.Ticker)
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("got response for an unsupported market %s", aggMsg.Ticker)
	}

	// Convert the price to a big Float.
	price, err := math.Float64StringToBigFloat(aggMsg.Price)
	if err != nil {
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToParsePrice),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unResolved), nil
}
