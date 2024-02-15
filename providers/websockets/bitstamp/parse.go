package bitstamp

import (
	"fmt"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// parseTickerMessage parses a ticker message received from the Bitstamp websocket API.
// All price updates must be made from the live trades channel.
func (h *WebSocketHandler) parseTickerMessage(
	msg TickerResponseMessage,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// Ensure that the price feeds are coming from the live trading channel.
	if !strings.HasPrefix(msg.Channel, string(TickerChannel)) {
		return providertypes.NewGetResponse(resolved, unResolved),
			fmt.Errorf("invalid ticker message %s", msg.Channel)
	}

	tickerSplit := strings.Split(msg.Channel, string(TickerChannel))
	if len(tickerSplit) != ExpectedTickerLength {
		return providertypes.NewGetResponse(resolved, unResolved),
			fmt.Errorf("invalid ticker message length %s", msg.Channel)
	}

	// Get the ticker from the message and market.
	ticker := tickerSplit[TickerCurrencyPairIndex]

	inverted := h.market.Invert()
	market, ok := inverted[ticker]
	if !ok {
		return providertypes.NewGetResponse(resolved, unResolved),
			fmt.Errorf("received unsupported ticker %s", ticker)
	}

	// Get the price from the message.
	price, err := math.Float64StringToBigInt(msg.Data.PriceStr, market.Ticker.Decimals)
	if err != nil {
		unResolved[market.Ticker] = err
		return providertypes.NewGetResponse(resolved, unResolved), err
	}

	resolved[market.Ticker] = providertypes.NewResult(price, time.Now().UTC())
	return providertypes.NewGetResponse(resolved, unResolved), nil
}
