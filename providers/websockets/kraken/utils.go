package kraken

import (
	"encoding/json"
	"fmt"
)

const (
	// URL is the websocket URL for Kraken. You can find the documentation here:
	// https://docs.kraken.com/websockets/. Kraken provides a authenticated and
	// unauthenticated websocket. The URLs defined below are all unauthenticated.

	// ProductionURL is the production websocket URL for Kraken.
	ProductionURL = "wss://ws.kraken.com"

	// BetaURL is the demo websocket URL for Kraken.
	BetaURL = "wss://beta-ws.kraken.com"
)

// DecodeTickerResponseMessage decodes a ticker response message .
func DecodeTickerResponseMessage(message []byte) (TickerResponseMessage, error) {
	var rawResponse []json.RawMessage
	if err := json.Unmarshal(message, &rawResponse); err != nil {
		return TickerResponseMessage{}, err
	}

	if len(rawResponse) != ExpectedTickerResponseMessageLength {
		return TickerResponseMessage{}, fmt.Errorf(
			"invalid ticker response message; expected length %d, got %d", ExpectedTickerResponseMessageLength, len(rawResponse),
		)
	}

	var response TickerResponseMessage
	if err := json.Unmarshal(rawResponse[ChannelIDIndex], &response.ChannelID); err != nil {
		return TickerResponseMessage{}, err
	}

	if err := json.Unmarshal(rawResponse[TickerDataIndex], &response.TickerData); err != nil {
		return TickerResponseMessage{}, err
	}

	if err := json.Unmarshal(rawResponse[ChannelNameIndex], &response.ChannelName); err != nil {
		return TickerResponseMessage{}, err
	}

	if err := json.Unmarshal(rawResponse[PairIndex], &response.Pair); err != nil {
		return TickerResponseMessage{}, err
	}

	return response, nil
}
