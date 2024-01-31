package mexc

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// parseTickerResponseMessage parses a price update received from the MEXC websocket
// and returns a GetResponse.
func (h *WebSocketDataHandler) parseTickerResponseMessage(
	msg TickerResponseMessage,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], error) {
	var (
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unResolved = make(map[oracletypes.CurrencyPair]error)
	)

	market, ok := h.cfg.Market.TickerToMarketConfigs[msg.Data.Symbol]
	if !ok {
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved),
			fmt.Errorf("unknown ticker %s", msg.Data.Symbol)
	}

	// Ensure that the channel received is the ticker channel.
	cp := market.CurrencyPair
	if !strings.HasPrefix(msg.Channel, string(MiniTickerChannel)) {
		err := fmt.Errorf("invalid channel %s", msg.Channel)
		unResolved[cp] = err
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), err
	}

	// Convert the price.
	price, err := math.Float64StringToBigInt(msg.Data.Price, cp.Decimals)
	if err != nil {
		unResolved[cp] = err
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), err
	}

	resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), nil
}
