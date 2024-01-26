package bitfinex

import (
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func (h *WebsocketDataHandler) parseSubscribedMessage(
	msg SubscribedMessage,
) error {
	return h.UpdateChannelMap(msg.ChannelID, msg.Pair)
}

func (h *WebsocketDataHandler) parseErrorMessage(
	msg ErrorMessage,
) ([]handlers.WebsocketEncodedMessage, error) {
	e := ErrorCode(msg.Code)
	return nil, e.Error()
}

func (h *WebsocketDataHandler) parseTickerStream(
	stream TickerStream,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], error) {
	var (
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved = make(map[oracletypes.CurrencyPair]error)
	)

	// handle stream for one of the tickers
	market, ok := h.channelMap[stream.ChannelID]
	if !ok {
		h.logger.Debug("currency pair not found for ticker channel ID", zap.String("channel_id", stream.ChannelID))
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved), nil
	}

	cp := market.CurrencyPair
	price, err := math.Float64StringToBigInt(stream.LastPrice, cp.Decimals())
	if err != nil {
		h.logger.Error("failed to convert price to big.Int", zap.Error(err))
		unresolved[cp] = fmt.Errorf("failed to convert price to big.Int: %s", err)
	}

	resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved), nil
}
