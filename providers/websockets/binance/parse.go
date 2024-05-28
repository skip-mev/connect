package binance

import (
	"fmt"
	"time"

	providertypes "github.com/skip-mev/slinky/providers/types"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
)

// parseAggregateTradeMessage parses an aggregate trade message from the Binance websocket feed.
// This message is sent when a trade is executed on the Binance exchange.
func (h *WebSocketHandler) parseAggregateTradeMessage(
	msg AggregatedTradeMessageResponse,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// Determine if the ticker is valid.
	aggMsg := msg.Data
	ticker, ok := h.cache.FromOffChainTicker(aggMsg.Ticker)
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("got response for an unsupported market %s", aggMsg.Ticker)
	}

	// Convert the price to a big Float.
	price, err := math.Float64StringToBigFloat(aggMsg.Price)
	if err != nil {
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToParsePrice),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unResolved), nil
}

// parseTickerMessage parses a ticker message from the Binance websocket feed.
// This message is broadcasted every 1000ms and contains the latest price of a ticker.
func (h *WebSocketHandler) parseTickerMessage(
	msg TickerMessageResponse,
) (types.PriceResponse, error) {
	var (
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// Determine if the ticker is valid.
	tickerMsg := msg.Data
	ticker, ok := h.cache.FromOffChainTicker(tickerMsg.Ticker)
	if !ok {
		return types.NewPriceResponse(resolved, unResolved),
			fmt.Errorf("got response for an unsupported market %s", tickerMsg.Ticker)
	}

	// Convert the price to a big Float.
	price, err := math.Float64StringToBigFloat(tickerMsg.LastPrice)
	if err != nil {
		unResolved[ticker] = providertypes.UnresolvedResult{
			ErrorWithCode: providertypes.NewErrorWithCode(err, providertypes.ErrorFailedToParsePrice),
		}
		return types.NewPriceResponse(resolved, unResolved), err
	}

	resolved[ticker] = types.NewPriceResult(price, time.Now().UTC())
	return types.NewPriceResponse(resolved, unResolved), nil
}
