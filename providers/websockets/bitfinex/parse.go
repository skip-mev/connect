package bitfinex

import (
	"github.com/skip-mev/slinky/providers/base/websocket/handlers"
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
