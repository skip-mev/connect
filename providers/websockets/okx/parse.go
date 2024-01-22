package okx

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
	"go.uber.org/zap"
)

const (
	// ExpectedErrorPrefix is the prefix of an error message that is returned by the OKX API.
	// Specifically, this is the prefix of the error message that is returned when the user
	// attempts to subscribe to a channel but could not be subscribed.
	ExpectedErrorPrefix = "Invalid request: "

	// ExpectedErrorElements is the number of elements that are expected in the error message.
	ExpectedErrorElements = 2
)

// parseSubscribeResponseMessage parses a subscribe response message. The format of the message
// is defined in the messages.go file. There are two cases that are handled:
//
// 1. Successfully subscribed to the channel. In this case, no further action is required.
// 2. Error message. In this case, we attempt to re-subscribe to the channel.
func (h *WebsocketDataHandler) parseSubscribeResponseMessage(resp SubscribeResponseMessage) ([]byte, error) {
	// A response with an event type of subscribe means that we have successfully subscribed to the channel.
	if t := EventType(resp.Event); t == EventSubscribe {
		h.logger.Info("successfully subscribed to channel", zap.String("instrument", resp.Arguments.InstrumentID))
		return nil, nil
	}

	// Attempt to re-subscribe to the channel.
	// Format of the message is:
	//  ...
	//	"msg": "Invalid request: {\"op\": \"subscribe\", \"argss\":[{ \"channel\" : \"index-tickers\", \"instId\" : \"BTC-USDT\"}]}",
	//  ...
	//
	// The message is an exact copy of the request message, so we can just unmarshal it and re-subscribe.
	h.logger.Error("received error message", zap.String("message", resp.Message), zap.String("code", resp.Code))
	jsonString := strings.Split(resp.Message, ExpectedErrorPrefix)
	if len(jsonString) != ExpectedErrorElements {
		return nil, fmt.Errorf("unable to parse subscription message from message: %s", resp.Message)
	}

	// Attempt to unmarshal the request.
	var request SubscribeRequestMessage
	if err := json.Unmarshal([]byte(jsonString[1]), &request); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request: %s", err)
	}

	// Re-subscribe to the channel.
	h.logger.Debug("re-subscribing to channel", zap.Any("instrument", request.Arguments))
	return json.Marshal(request)
}

// parseTickerResponseMessage parses a ticker response message. The format of the message is defined
// in the messages.go file. This message contains the latest price data for a set of instruments.
func (h *WebsocketDataHandler) parseTickerResponseMessage(
	resp IndexTickersResponseMessage,
) (providertypes.GetResponse[oracletypes.CurrencyPair, *big.Int], error) {
	var (
		resolved   = make(map[oracletypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved = make(map[oracletypes.CurrencyPair]error)
	)

	// The channel must be the index tickers channel.
	if Channel(resp.Arguments.Channel) != IndexTickersChannel {
		return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved),
			fmt.Errorf("invalid channel %s", resp.Arguments.Channel)
	}

	// Iterate through all of the tickers and add them to the response.
	for _, ticker := range resp.Data {
		market, ok := h.invertedMarketCfg.MarketToCurrencyPairConfigs[ticker.InstrumentID]
		if !ok {
			h.logger.Debug("currency pair not found for instrument ID", zap.String("instrument_id", ticker.InstrumentID))
			continue
		}

		// Convert the price to a big.Int.
		cp := market.CurrencyPair
		price, err := math.Float64StringToBigInt(ticker.IndexPrice, cp.Decimals())
		if err != nil {
			h.logger.Error("failed to convert price to big.Int", zap.Error(err))
			unresolved[cp] = fmt.Errorf("failed to convert price to big.Int: %s", err)
			continue
		}

		resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
	}

	return providertypes.NewGetResponse[oracletypes.CurrencyPair, *big.Int](resolved, unresolved), nil
}
