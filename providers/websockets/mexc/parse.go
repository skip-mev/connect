package mexc

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// parseTickerResponseMessage parses a price update received from the MEXC websocket
// and returns a GetResponse.
func (h *WebSocketHandler) parseTickerResponseMessage(
	msg TickerResponseMessage,
) (providertypes.GetResponse[mmtypes.Ticker, *big.Int], error) {
	var (
		resolved   = make(map[mmtypes.Ticker]providertypes.Result[*big.Int])
		unResolved = make(map[mmtypes.Ticker]error)
	)

	inverted := h.market.Invert()
	market, ok := inverted[msg.Data.Symbol]
	if !ok {
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unResolved),
			fmt.Errorf("unknown ticker %s", msg.Data.Symbol)
	}

	// Ensure that the channel received is the ticker channel.
	if !strings.HasPrefix(msg.Channel, string(MiniTickerChannel)) {
		err := fmt.Errorf("invalid channel %s", msg.Channel)
		unResolved[market.Ticker] = err
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unResolved), err
	}

	// Convert the price.
	price, err := math.Float64StringToBigInt(msg.Data.Price, market.Ticker.Decimals)
	if err != nil {
		unResolved[market.Ticker] = err
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unResolved), err
	}

	resolved[market.Ticker] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
	return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unResolved), nil
}
