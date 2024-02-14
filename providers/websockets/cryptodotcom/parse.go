package cryptodotcom

import (
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/pkg/math"
	providertypes "github.com/skip-mev/slinky/providers/types"
	mmtypes "github.com/skip-mev/slinky/x/marketmap/types"
)

// parseInstrumentMessage is used to parse an instrument message received from the Crypto.com
// websocket API. This message contains the latest price data for a set of instruments.
func (h *WebSocketHandler) parseInstrumentMessage(
	msg InstrumentResponseMessage,
) (providertypes.GetResponse[mmtypes.Ticker, *big.Int], error) {
	var (
		resolved    = make(map[mmtypes.Ticker]providertypes.Result[*big.Int])
		unresolved  = make(map[mmtypes.Ticker]error)
		instruments = msg.Result.Data
	)

	// If the response contained no instrument data, return an error.
	if len(instruments) == 0 {
		return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unresolved),
			fmt.Errorf("no instrument data was returned")
	}

	// Iterate through each market and attempt to parse the price.
	inverted := h.market.Invert()
	for _, instrument := range instruments {
		// If we don't have a mapping for the instrument, return an error. This is likely a configuration
		// error.
		market, ok := inverted[instrument.Name]
		if !ok {
			h.logger.Error("failed to find currency pair for instrument", zap.String("instrument", instrument.Name))
			continue
		}

		// Attempt to parse the price.
		if price, err := math.Float64StringToBigInt(instrument.LatestTradePrice, market.Ticker.Decimals); err != nil {
			unresolved[market.Ticker] = fmt.Errorf("failed to parse price %s: %w", instrument.LatestTradePrice, err)
		} else {
			resolved[market.Ticker] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
		}

	}

	return providertypes.NewGetResponse[mmtypes.Ticker, *big.Int](resolved, unresolved), nil
}
