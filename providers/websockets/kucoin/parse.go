package kucoin

import (
	"fmt"
	"math/big"
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

	// Parse the price from the message.
	cp := market.CurrencyPair
	price, err := math.Float64StringToBigInt(msg.Data.Price, cp.Decimals())
	if err != nil {
		err = fmt.Errorf("failed to parse price %s", err)
		unResolved[cp] = err
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), err
	}

	resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now())
	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), nil
}
