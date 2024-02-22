package bitstamp

import (
	"fmt"
	"strings"
	"time"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
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
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("invalid ticker message %s", msg.Channel)
	}

	tickerSplit := strings.Split(msg.Channel, string(TickerChannel))
	if len(tickerSplit) != ExpectedTickerLength {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("invalid ticker message length %s", msg.Channel)
	}

	// Get the ticker from the message and market.
	offChainTicker := tickerSplit[TickerCurrencyPairIndex]
	ticker, ok := h.market.OffChainMap[offChainTicker]
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("received unsupported ticker %s", ticker)
	}

	// Get the price from the message.
	price, err := math.Float64StringToBigInt(msg.Data.PriceStr, ticker.Decimals)
	if err != nil {
		unResolved[ticker] = err
		return types.NewPriceResponse(resolved, unResolved), err
	}

	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unResolved), nil
}
