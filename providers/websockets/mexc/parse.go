package mexc

import (
	"fmt"
	"strings"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
)

// parseTickerResponseMessage parses a price update received from the MEXC websocket
// and returns a GetResponse.
func (h *WebSocketHandler) parseTickerResponseMessage(
	msg TickerResponseMessage,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	ticker, ok := h.cache.FromOffChainTicker(msg.Data.Symbol)
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("unknown ticker %s", msg.Data.Symbol)
	}

	// Ensure that the channel received is the ticker channel.
	if !strings.HasPrefix(msg.Channel, string(MiniTickerChannel)) {
		err := fmt.Errorf("invalid channel %s", msg.Channel)
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorInvalidWebSocketTopic),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	// Convert the price.
	price, err := math.Float64StringToBigFloat(msg.Data.Price)
	if err != nil {
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToParsePrice),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unResolved), nil
}
