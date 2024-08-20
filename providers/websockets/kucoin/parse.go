package kucoin

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	providertypes "github.com/skip-mev/connect/v2/providers/types"

	"github.com/skip-mev/connect/v2/oracle/types"
	"github.com/skip-mev/connect/v2/pkg/math"
)

// parseTickerResponseMessage is used to parse a ticker response message.
func (h *WebSocketHandler) parseTickerResponseMessage(
	msg TickerResponseMessage,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// The response must be from a subscription to the ticker channel.
	if subject := SubjectType(msg.Subject); subject != TickerSubject {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("received unsupported channel %s", subject)
	}

	// Retrieve the ticker data from the message.
	tickerData := strings.Split(msg.Topic, string(TickerTopic))
	if len(tickerData) != ExpectedTopicLength {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("invalid ticker data %s", tickerData)
	}

	// Parse the currency pair from the ticker data.
	offChainTicker := tickerData[TickerIndex]
	ticker, ok := h.cache.FromOffChainTicker(offChainTicker)
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("market not found for ticker %s", offChainTicker)
	}

	// Check if the sequence number is valid.
	sequence, err := strconv.ParseInt(msg.Data.Sequence, 10, 64)
	if err != nil {
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorInvalidResponse),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	seenSequence, ok := h.sequences[ticker]
	switch {
	case !ok || seenSequence < sequence:
		// If the sequence number is not found, then this is the first message
		// received for this currency pair. Set the sequence number to the
		// sequence number received. Additionally, if the sequence number is
		// greater than the sequence number currently stored, then this message
		// was received in order.
		h.sequences[ticker] = sequence
	default:
		// If the sequence number is greater than the sequence number received,
		// then this message was received out of order. Ignore the message.
		err := fmt.Errorf("received out of order ticker response message")
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorInvalidResponse),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	// Parse the price from the message.
	price, err := math.Float64StringToBigFloat(msg.Data.Price)
	if err != nil {
		wErr := fmt.Errorf("failed to parse price %w", err)
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(wErr, providertypes.ErrorFailedToParsePrice),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unResolved), nil
}
