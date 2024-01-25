package coinbase

import (
	"fmt"
	"math/big"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// parseTickerResponseMessage is used to parse a ticker response message. Note
// that each response will include a sequence number. This sequence number
// should be used to determine if any messages were missed. If any previous
// messages were missed, the client should ignore the previous messages if they
// are received at a later time.
func (h *WebSocketDataHandler) parseTickerResponseMessage(
	msg TickerResponseMessage,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], error) {
	var (
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unResolved = make(map[oracletypes.CurrencyPair]error)
	)

	// Determine if the ticker is valid.
	market, ok := h.cfg.Market.TickerToMarketConfigs[msg.Ticker]
	if !ok {
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved),
			fmt.Errorf("got response for an unsupported market %s", msg.Ticker)
	}

	// Determine if the sequence number is valid.
	cp := market.CurrencyPair
	sequence, ok := h.sequence[market.CurrencyPair]
	switch {
	case !ok || sequence < msg.Sequence:
		// If the sequence number is not found, then this is the first message
		// received for this currency pair. Set the sequence number to the
		// sequence number received. Additionally, if the sequence number is
		// greater than the sequence number currently stored, then this message
		// was received in order.
		h.sequence[cp] = msg.Sequence
	default:
		// If the sequence number is greater than the sequence number received,
		// then this message was received out of order. Ignore the message.
		err := fmt.Errorf("received out of order ticker response message")
		unResolved[cp] = err
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), err
	}

	// Convert the price to a big int.
	price, err := math.Float64StringToBigInt(msg.Price, cp.Decimals())
	if err != nil {
		unResolved[cp] = err
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), err
	}

	// Convert the time to a time object and resolve the price into the response.
	resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())

	h.logger.Debug("successfully parsed ticker response message")
	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), nil
}
