package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Candle defines price, volume, and time information for an exchange rate.
//
// XXX: Consider replacing sdk.Dec with another decimal type.
type Candle struct {
	Price     sdk.Dec // last trade price
	Volume    sdk.Dec // volume
	Timestamp int64   // timestamp
}

func NewCandle(provider, symbol, lastPrice, volume string, timestamp int64) (Candle, error) {
	price, err := sdk.NewDecFromStr(lastPrice)
	if err != nil {
		return Candle{}, fmt.Errorf("failed to parse %s price (%s) for %s: %w", provider, lastPrice, symbol, err)
	}

	volumeDec, err := sdk.NewDecFromStr(volume)
	if err != nil {
		return Candle{}, fmt.Errorf("failed to parse %s volume (%s) for %s: %w", provider, volume, symbol, err)
	}

	return Candle{Price: price, Volume: volumeDec, Timestamp: timestamp}, nil
}
