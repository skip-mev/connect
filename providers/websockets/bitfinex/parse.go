package bitfinex

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/oracle/types"
	"github.com/skip-mev/slinky/pkg/math"
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

const (
	// indexChannelID is the index of a data stream's channel ID.
	indexChannelID = 0
	// indexPayload is the index of a data stream's payload.
	indexPayload = 1
)

// parseSubscribedMessage updates the channel map for a subscribed message.
func (h *WebSocketHandler) parseSubscribedMessage(
	msg SubscribedMessage,
) error {
	return h.updateChannelMap(msg.ChannelID, msg.Pair)
}

// parseErrorMessage returns the proper error code from an error message.
func (h *WebSocketHandler) parseErrorMessage(
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
func (h *WebSocketHandler) handleStream(
	message []byte,
) (types.PriceResponse, error) {
	var (
		baseStream []interface{}
		resolved   = make(types.ResolvedPrices)
		unResolved = make(types.UnResolvedPrices)
	)

	// Attempt to unmarshal the message into a base message. This is used to determine the type
	// of message that was received.
	if err := json.Unmarshal(message, &baseStream); err != nil {
		return providertypes.NewGetResponse(resolved, unResolved), err
	}

	if len(baseStream) != ExpectedBaseStreamLength {
		return providertypes.NewGetResponse(resolved, unResolved),
			fmt.Errorf("invalid length of stream data received. must be %d.  stream: %v. len: %d",
				ExpectedBaseStreamLength,
				baseStream,
				len(baseStream),
			)
	}

	// first element is always channel id
	channelID := int(baseStream[indexChannelID].(float64))
	market, ok := h.channelMap[channelID]
	if !ok {
		return providertypes.NewGetResponse(resolved, unResolved),
			fmt.Errorf("received stream for unknown channel id %v", channelID)
	}

	ticker := market.Ticker
	h.logger.Debug("received stream", zap.Int("channel_id", channelID), zap.String("ticker", ticker.String()))

	// check if it is a heartbeat
	hbID, ok := baseStream[indexPayload].(string)
	if ok && hbID == IDHeartbeat {
		h.logger.Debug("received heartbeat", zap.Int("channel_id", channelID), zap.String("ticker", ticker.String()))
		return providertypes.NewGetResponse(resolved, unResolved), nil

	}

	// if it is not a string, it is a stream update
	dataArr, ok := baseStream[indexPayload].([]interface{})
	if !ok || len(dataArr) != ExpectedStreamPayloadLength {
		err := fmt.Errorf("unknown data: %v, len: %d", baseStream[1], len(dataArr))
		unResolved[ticker] = err
		return providertypes.NewGetResponse(resolved, unResolved), err
	}

	lastPrice := dataArr[6]
	// Convert the price to a big int.
	price := math.Float64ToBigInt(lastPrice.(float64), ticker.Decimals)
	resolved[ticker] = providertypes.NewResult[*big.Int](price, time.Now().UTC())

	return providertypes.NewGetResponse(resolved, unResolved), nil
}

// updateChannelMap updates the internal map for the given channelID and ticker.
func (h *WebSocketHandler) updateChannelMap(channelID int, ticker string) error {
	inverted := h.market.Invert()
	market, ok := inverted[ticker]
	if !ok {
		return fmt.Errorf("unable to find market for currency pair: %s", ticker)
	}

	h.channelMap[channelID] = market
	return nil
}
