package kucoin

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
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
	ticker, ok := h.market.OffChainMap[offChainTicker]
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("market not found for ticker %s", offChainTicker)
	}

	// Check if the sequence number is valid.
	sequence, err := strconv.ParseInt(msg.Data.Sequence, 10, 64)
	if err != nil {
		unResolved[ticker] = err
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
		unResolved[ticker] = err
		return types.NewPriceResponse(resolved, unResolved), err
	}

	// Parse the price from the message.
	price, err := math.Float64StringToBigInt(msg.Data.Price, ticker.Decimals)
	if err != nil {
		err = fmt.Errorf("failed to parse price %w", err)
		unResolved[ticker] = err
		return types.NewPriceResponse(resolved, unResolved), err
	}

	resolved[ticker] = types.NewPriceResult(price, time.Now())
	return types.NewPriceResponse(resolved, unResolved), nil
}
