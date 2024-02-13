package bitstamp

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// parseTickerMessage parses a ticker message received from the Bitstamp websocket API.
// All price updates must be made from the live trades channel.
func (h *WebSocketDataHandler) parseTickerMessage(
	msg TickerResponseMessage,
) (providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int], error) {
	var (
		resolved   = make(map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int])
		unResolved = make(map[slinkytypes.CurrencyPair]error)
	)

	// Ensure that the price feeds are coming from the live trading channel.
	if !strings.HasPrefix(msg.Channel, string(TickerChannel)) {
		return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unResolved),
			fmt.Errorf("invalid ticker message %s", msg.Channel)
	}

	tickerSplit := strings.Split(msg.Channel, string(TickerChannel))
	if len(tickerSplit) != ExpectedTickerLength {
		return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unResolved),
			fmt.Errorf("invalid ticker message length %s", msg.Channel)
	}

	// Get the ticker from the message and market.
	ticker := tickerSplit[TickerCurrencyPairIndex]
	market, ok := h.cfg.Market.TickerToMarketConfigs[ticker]
	if !ok {
		return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unResolved),
			fmt.Errorf("received unsupported ticker %s", ticker)
	}

	// Get the price from the message.
	cp := market.CurrencyPair
	price, err := math.Float64StringToBigInt(msg.Data.PriceStr, cp.Decimals())
	if err != nil {
		unResolved[cp] = err
		return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unResolved), err
	}

	resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
	return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unResolved), nil
}
