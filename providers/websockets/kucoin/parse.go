package kucoin

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// parseTickerResponseMessage is used to parse a ticker response message.
func (h *WebSocketDataHandler) parseTickerResponseMessage(
	msg TickerResponseMessage,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], error) {
	var (
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unResolved = make(map[oracletypes.CurrencyPair]error)
	)

	// The response must be from a subscription to the ticker channel.
	if subject := SubjectType(msg.Subject); subject != TickerSubject {
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved),
			fmt.Errorf("received unsupported channel %s", subject)
	}

	// Retrieve the ticker data from the message.
	tickerData := strings.Split(msg.Topic, string(TickerTopic))
	if len(tickerData) != ExpectedTopicLength {
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved),
			fmt.Errorf("invalid ticker data %s", tickerData)
	}

	// Parse the currency pair from the ticker data.
	ticker := tickerData[TickerIndex]
	market, ok := h.cfg.Market.TickerToMarketConfigs[ticker]
	if !ok {
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved),
			fmt.Errorf("market not found for ticker %s", ticker)
	}

	// Check if the sequence number is valid.
	cp := market.CurrencyPair
	sequence, err := strconv.ParseInt(msg.Data.Sequence, 10, 64)
	if err != nil {
		unResolved[cp] = err
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), err
	}

	seenSequence, ok := h.sequences[cp]
	switch {
	case !ok || seenSequence < sequence:
		// If the sequence number is not found, then this is the first message
		// received for this currency pair. Set the sequence number to the
		// sequence number received. Additionally, if the sequence number is
		// greater than the sequence number currently stored, then this message
		// was received in order.
		h.sequences[cp] = sequence
	default:
		// If the sequence number is greater than the sequence number received,
		// then this message was received out of order. Ignore the message.
		err := fmt.Errorf("received out of order ticker response message")
		unResolved[cp] = err
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), err
	}

	// Parse the price from the message.
	price, err := math.Float64StringToBigInt(msg.Data.Price, cp.Decimals())
	if err != nil {
		err = fmt.Errorf("failed to parse price %s", err)
		unResolved[cp] = err
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), err
	}

	resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now())
	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), nil
}
