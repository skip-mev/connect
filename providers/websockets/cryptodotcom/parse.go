package cryptodotcom

import (
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/skip-mev/slinky/pkg/math"
	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	providertypes "github.com/skip-mev/slinky/providers/types"
)

// parseInstrumentMessage is used to parse an instrument message received from the Crypto.com
// websocket API. This message contains the latest price data for a set of instruments.
func (h *WebSocketDataHandler) parseInstrumentMessage(
	msg InstrumentResponseMessage,
) (providertypes.GetResponse[slinkytypes.CurrencyPair, *big.Int], error) {
	var (
		resolved    = make(map[slinkytypes.CurrencyPair]providertypes.Result[*big.Int])
		unresolved  = make(map[slinkytypes.CurrencyPair]error)
		instruments = msg.Result.Data
	)

	// If the response contained no instrument data, return an error.
	if len(instruments) == 0 {
		return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unresolved),
			fmt.Errorf("no instrument data was returned")
	}

	// Iterate through each market and attempt to parse the price.
	for _, instrument := range instruments {
		// If we don't have a mapping for the instrument, return an error. This is likely a configuration
		// error.
		market, ok := h.cfg.Market.TickerToMarketConfigs[instrument.Name]
		if !ok {
			h.logger.Error("failed to find currency pair for instrument", zap.String("instrument", instrument.Name))
			continue
		}

		// Attempt to parse the price.
		cp := market.CurrencyPair
		if price, err := math.Float64StringToBigInt(instrument.LatestTradePrice, cp.Decimals()); err != nil {
			unresolved[cp] = fmt.Errorf("failed to parse price %s: %w", instrument.LatestTradePrice, err)
		} else {
			resolved[cp] = providertypes.NewResult[*big.Int](price, time.Now().UTC())
		}

	}

	return providertypes.NewGetResponse[slinkytypes.CurrencyPair, *big.Int](resolved, unresolved), nil
}
