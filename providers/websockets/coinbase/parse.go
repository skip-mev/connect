package coinbase

import (
	"fmt"
	"math/big"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
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
	ticker, ok := h.cache.FromOffChainTicker(msg.Ticker)
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("got response for an unsupported market %s", msg.Ticker)
	}

	// Determine if the sequence number is valid.
	if err := h.checkSequenceNumber(ticker, msg.Sequence); err != nil {
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorInvalidResponse),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	// Convert the price to a big Float.
	price, err := math.Float64StringToBigFloat(msg.Price)
	if err != nil {
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToParsePrice),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	// Update the trade ID.
	h.tradeIDs[ticker] = msg.TradeID

	// Convert the time to a time object and resolve the price into the response.
	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unResolved), nil
}

// parseHeartbeatResponseMessage is used to parse a heartbeat response message. In particular.
// this function checks that the trade ID and sequence number are valid. If the trade ID is the
// same as what is cached, then we know that the price has not changed. If the sequence number is
// out of order, then we know that we missed a message and should ignore the message.
func (h *WebSocketHandler) parseHeartbeatResponseMessage(
	msg HeartbeatResponseMessage,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// Determine if the ticker is valid.
	ticker, ok := h.cache.FromOffChainTicker(msg.Ticker)
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("got response for an unsupported trade ID %d", msg.LastTradeID)
	}

	// Determine if the sequence number is valid.
	if err := h.checkSequenceNumber(ticker, msg.Sequence); err != nil {
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorInvalidResponse),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	currentTradeID, ok := h.tradeIDs[ticker]
	if !ok || currentTradeID != msg.LastTradeID {
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(fmt.Errorf("no price update received"), providertypes.ErrorNoExistingPrice),
		}
		return types.NewPriceResponse(resolved, unResolved), nil
	}

	// If the trade ID is the same as the current trade ID, then the price has not changed.
	resolved[ticker] = types.NewPriceResultWithCode(big.NewFloat(0), time.Now().UTC(), providertypes.ResponseCodeUnchanged)
	return types.NewPriceResponse(resolved, unResolved), nil
}

// checkSequenceNumber is used to check the sequence number of the message. If the sequence number
// is out of order, then the message should be ignored. If the sequence number is in order, then the
// sequence number should be updated.
func (h *WebSocketHandler) checkSequenceNumber(
	ticker types.ProviderTicker,
	sequence int64,
) error {
	observedSequence, ok := h.sequence[ticker]
	switch {
	case !ok || observedSequence <= sequence:
		// If the sequence number is not found, then this is the first message
		// received for this currency pair. Set the sequence number to the
		// sequence number received. Additionally, if the sequence number is
		// greater than the sequence number currently stored, then this message
		// was received in order.
		h.sequence[ticker] = sequence
	default:
		// If the sequence number is greater than the sequence number received,
		// then this message was received out of order. Ignore the message.
		return fmt.Errorf(
			"received out of order ticker response message; is %d, expected %d",
			sequence,
			observedSequence,
		)
	}

	return nil
}
