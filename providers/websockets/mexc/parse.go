package mexc

import (
	"fmt"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
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

	inverted := h.market.Invert()
	market, ok := inverted[msg.Data.Symbol]
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("unknown ticker %s", msg.Data.Symbol)
	}

	// Ensure that the channel received is the ticker channel.
	if !strings.HasPrefix(msg.Channel, string(MiniTickerChannel)) {
		err := fmt.Errorf("invalid channel %s", msg.Channel)
		unResolved[market.Ticker] = err
		return types.NewPriceResponse(resolved, unResolved), err
	}

	// Convert the price.
	price, err := math.Float64StringToBigInt(msg.Data.Price, market.Ticker.Decimals)
	if err != nil {
		unResolved[market.Ticker] = err
		return types.NewPriceResponse(resolved, unResolved), err
	}

	resolved[market.Ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unResolved), nil
}
