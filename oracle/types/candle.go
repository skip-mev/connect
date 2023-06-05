package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CandlePrice defines price, volume, and time information for an exchange rate.
//
// XXX: Consider replacing sdk.Dec with another decimal type.
type CandlePrice struct {
	Price     sdk.Dec // last trade price
	Volume    sdk.Dec // volume
	Timestamp int64   // timestamp
}

func NewCandlePrice(provider, symbol, lastPrice, volume string, timestamp int64) (CandlePrice, error) {
	price, err := sdk.NewDecFromStr(lastPrice)
	if err != nil {
		return CandlePrice{}, fmt.Errorf("failed to parse %s price (%s) for %s: %w", provider, lastPrice, symbol, err)
	}

	volumeDec, err := sdk.NewDecFromStr(volume)
	if err != nil {
		return CandlePrice{}, fmt.Errorf("failed to parse %s volume (%s) for %s: %w", provider, volume, symbol, err)
	}

	return CandlePrice{Price: price, Volume: volumeDec, Timestamp: timestamp}, nil
}
