package bitfinex

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const (
	// indexChannelID is the index of a data stream's channel ID.
	indexChannelID = 0
	// indexPayload is the index of a data stream's payload.
	indexPayload = 1
)

// parseSubscribedMessage updates the channel map for a subscribed message.
func (h *WebsocketDataHandler) parseSubscribedMessage(
	msg SubscribedMessage,
) error {
	return h.UpdateChannelMap(msg.ChannelID, msg.Pair)
}

// parseErrorMessage returns the proper error code from an error message.
func (h *WebsocketDataHandler) parseErrorMessage(
	msg ErrorMessage,
) ([]handlers.WebsocketEncodedMessage, error) {
	e := ErrorCode(msg.Code)
	return nil, e.Error()
}

// handleStream handles a data stream sent from the peer.
//
// Data streams always start with the channelID associated, but
// can have a different payload depending on the type of message.
// This handler handles:
// 1. Heartbeat messages.  These take the following form:
//
//	[ CHANNEL_ID, "hb" ]
//
// 2. Ticker update streams.  These take the following form:
//
// [
//
//	CHANNEL_ID,
//	[
//	  BID,
//	  BID_SIZE,
//	  ASK,
//	  ASK_SIZE,
//	  DAILY_CHANGE,
//	  DAILY_CHANGE_RELATIVE,
//	  LAST_PRICE,
//	  VOLUME,
//	  HIGH,
//	  LOW
//	]
//
// ]
//
// ref: https://docs.bitfinex.com/reference/ws-public-ticker
func (h *WebsocketDataHandler) handleStream(
	message []byte,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], error) {
	var (
		baseStream []interface{}
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unResolved = make(map[oracletypes.CurrencyPair]error)
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseStream); err != nil {
		h.logger.Debug("unable to unmarshal message into base struct", zap.Error(err))
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), err
	}

	if len(baseStream) != ExpectedBaseStreamLength {
		h.logger.Error("invalid length of stream data received. must be 2", zap.Any("data", baseStream), zap.Int("len", len(baseStream)))
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), fmt.Errorf("invalid length of stream data received. must be %d.  stream: %v. len: %d",
			ExpectedBaseStreamLength,
			baseStream,
			len(baseStream),
		)
	}

	// first element is always channel id
	channelID := int(baseStream[indexChannelID].(float64))
	market, ok := h.channelMap[channelID]
	if !ok {
		h.logger.Error("received stream for unknown channel id", zap.Int("channel_id", channelID))
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), fmt.Errorf("received stream for unknown channel id %v", channelID)
	}

	cp := market.CurrencyPair
	h.logger.Debug("received stream", zap.Int("channel_id", channelID), zap.String("market", cp.Ticker()))

	// check if it is a heartbeat
	hbID, ok := baseStream[indexPayload].(string)
	if ok && hbID == IDHeartbeat {
		h.logger.Debug("received heartbeat", zap.Int("channel_id", channelID), zap.String("pair", market.Ticker))
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), nil

	}

	// if it is not a string, it is a stream update
	dataArr, ok := baseStream[indexPayload].([]interface{})
	if !ok || len(dataArr) != ExpectedStreamPayloadLength {
		err := fmt.Errorf("unknown data: %v, len: %d", baseStream[1], len(dataArr))
		unResolved[cp] = err
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), err
	}

	lastPrice := dataArr[6]
	// Convert the price to a big int.
	price := math.Float64ToBigInt(lastPrice.(float64), cp.Decimals)
	resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())

	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unResolved), nil
}
