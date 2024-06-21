package tickermetadatajson

// DyDx is the Ticker.Metadata_JSON published to every Ticker in the x/marketmap module on dYdX.
type DyDx struct {
	// ReferencePrice gives a spot price for that Ticker at the point in time when the ReferencePrice was updated.
	// You should _not_ use this for up-to-date/instantaneous spot pricing data since it is updated infrequently.
	// The price is scaled by Ticker.Decimals.
	ReferencePrice uint64 `json:"reference_price"`
	// Liquidity gives a _rough_ estimate of the amount of liquidity in the Providers for a given Market.
	// It is _not_ updated in coordination with spot prices and only gives rough order of magnitude accuracy at the time
	// which the update for it is published.
	// The liquidity value stored here is USD denominated.
	Liquidity uint64 `json:"liquidity"`
}
