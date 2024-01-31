package gate

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

// parseSubscribeResponse attempts to parse a SubscribeResponse to see if it was successful.
func (h *WebsocketDataHandler) parseSubscribeResponse(
	msg SubscribeResponse,
) ([]handlers.WebsocketEncodedMessage, error) {
	if msg.Error.Message != "" {
		errCode := ErrorCode(msg.Error.Code)
		h.logger.Error("found error in subscribe response", zap.Error(errCode.Error()))
		return nil, errCode.Error()
	}

	if Status(msg.Result.Status) != StatusSuccess {
		h.logger.Error("subscription was not successful", zap.String("status", msg.Result.Status), zap.String("pair", msg.Error.Message))
		return nil, fmt.Errorf("subscription was not successful: %s", msg.Result.Status)
	}

	h.logger.Debug("successfully subscribed", zap.Int("id", msg.ID))
	return nil, nil
}

// parseTickerStream attempts to parse a TickerStream and translate it to the corresponding
// CurrencyPair update.
func (h *WebsocketDataHandler) parseTickerStream(
	stream TickerStream,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], error) {
	var (
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved = make(map[oracletypes.CurrencyPair]error)
	)

	// The channel must be the tickers channel.
	if Channel(stream.Channel) != ChannelTickers {
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved),
			fmt.Errorf("invalid channel %s", stream.Channel)
	}

	// Get the currency pair from the ticker.
	h.logger.Debug("received price update", zap.String("symbol", stream.Result.CurrencyPair))
	market, ok := h.cfg.Market.TickerToMarketConfigs[stream.Result.CurrencyPair]
	if !ok {
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved),
			fmt.Errorf("no currency pair found for symbol %s", stream.Result.CurrencyPair)
	}

	// Parse the price update.
	cp := market.CurrencyPair
	priceStr := stream.Result.Last
	price, err := math.Float64StringToBigInt(priceStr, cp.Decimals)
	if err != nil {
		unresolved[cp] = fmt.Errorf("failed to parse price %s: %s", priceStr, err)
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved), unresolved[cp]
	}

	resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
	return providertypes.NewGetResponse(resolved, unresolved), nil
}
