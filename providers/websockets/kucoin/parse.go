package kucoin

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// parseTickerResponseMessage is used to parse a ticker response message.
func (h *WebSocketHandler) parseTickerResponseMessage(
	msg TickerResponseMessage,
) (providertypes.GetResponse[mmtypes.Ticker, *big.Int], error) {
	var (
		resolved   = make(map[mmtypes.Ticker]providertypes.Result[*big.Int])
		unResolved = make(map[mmtypes.Ticker]error)
	)

	// The response must be from a subscription to the ticker channel.
	if subject := SubjectType(msg.Subject); subject != TickerSubject {
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unResolved),
			fmt.Errorf("received unsupported channel %s", subject)
	}

	// Retrieve the ticker data from the message.
	tickerData := strings.Split(msg.Topic, string(TickerTopic))
	if len(tickerData) != ExpectedTopicLength {
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unResolved),
			fmt.Errorf("invalid ticker data %s", tickerData)
	}

	// Parse the currency pair from the ticker data.
	ticker := tickerData[TickerIndex]
	inverted := h.market.Invert()
	market, ok := inverted[ticker]
	if !ok {
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unResolved),
			fmt.Errorf("market not found for ticker %s", ticker)
	}

	// Check if the sequence number is valid.
	sequence, err := strconv.ParseInt(msg.Data.Sequence, 10, 64)
	if err != nil {
		unResolved[market.Ticker] = err
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unResolved), err
	}

	seenSequence, ok := h.sequences[market.Ticker]
	switch {
	case !ok || seenSequence < sequence:
		// If the sequence number is not found, then this is the first message
		// received for this currency pair. Set the sequence number to the
		// sequence number received. Additionally, if the sequence number is
		// greater than the sequence number currently stored, then this message
		// was received in order.
		h.sequences[market.Ticker] = sequence
	default:
		// If the sequence number is greater than the sequence number received,
		// then this message was received out of order. Ignore the message.
		err := fmt.Errorf("received out of order ticker response message")
		unResolved[market.Ticker] = err
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unResolved), err
	}

	// Parse the price from the message.
	price, err := math.Float64StringToBigInt(msg.Data.Price, market.Ticker.Decimals)
	if err != nil {
		err = fmt.Errorf("failed to parse price %w", err)
		unResolved[market.Ticker] = err
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unResolved), err
	}

	resolved[market.Ticker] = providertypes.NewResult[*big.Int](price, time.Now())
	return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unResolved), nil
}
