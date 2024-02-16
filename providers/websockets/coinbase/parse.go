package coinbase

import (
	"fmt"
	"time"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
)

// parseTickerResponseMessage is used to parse a ticker response message. Note
// that each response will include a sequence number. This sequence number
// should be used to determine if any messages were missed. If any previous
// messages were missed, the client should ignore the previous messages if they
// are received at a later time.
func (h *WebSocketHandler) parseTickerResponseMessage(
	msg TickerResponseMessage,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// Determine if the ticker is valid.
	inverted := h.market.Invert()
	market, ok := inverted[msg.Ticker]
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("got response for an unsupported market %s", msg.Ticker)
	}

	// Determine if the sequence number is valid.
	sequence, ok := h.sequence[market.Ticker]
	switch {
	case !ok || sequence < msg.Sequence:
		// If the sequence number is not found, then this is the first message
		// received for this currency pair. Set the sequence number to the
		// sequence number received. Additionally, if the sequence number is
		// greater than the sequence number currently stored, then this message
		// was received in order.
		h.sequence[market.Ticker] = msg.Sequence
	default:
		// If the sequence number is greater than the sequence number received,
		// then this message was received out of order. Ignore the message.
		err := fmt.Errorf("received out of order ticker response message")
		unResolved[market.Ticker] = err
		return types.NewPriceResponse(resolved, unResolved), err
	}

	// Convert the price to a big int.
	price, err := math.Float64StringToBigInt(msg.Price, market.Ticker.Decimals)
	if err != nil {
		unResolved[market.Ticker] = err
		return types.NewPriceResponse(resolved, unResolved), err
	}

	// Convert the time to a time object and resolve the price into the response.
	resolved[market.Ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unResolved), nil
}
