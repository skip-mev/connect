package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TickerPrice defines price information for a given symbol given
// by a price provider.
type TickerPrice struct {
	Price     sdk.Dec   // last trade price
	Timestamp time.Time // timestamp
}

func NewTickerPrice(lastPrice string, timestamp time.Time) (TickerPrice, error) {
	price, err := sdk.NewDecFromStr(lastPrice)
	if err != nil {
		return TickerPrice{}, fmt.Errorf("failed to parse %s price (%s)", lastPrice, err)
	}

	return TickerPrice{Price: price, Timestamp: timestamp}, nil
}
