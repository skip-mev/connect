package types

import (
	"fmt"

	"cosmossdk.io/math"
)

// NewPriceFeedSLA returns a new PriceFeedSLA instance.
func NewPriceFeedSLA(
	id string,
	maximumViableWindow uint64,
	expectedUptime math.LegacyDec,
	slashConstant math.LegacyDec,
	minimumBlockUpdates uint64,
	frequency uint64,
) PriceFeedSLA {
	return PriceFeedSLA{
		ID:                  id,
		MaximumViableWindow: maximumViableWindow,
		ExpectedUptime:      expectedUptime,
		SlashConstant:       slashConstant,
		MinimumBlockUpdates: minimumBlockUpdates,
		Frequency:           frequency,
	}
}

// Qualifies determines whether the inputted price feed qualifies for
// an SLA check. A price feed qualifies to be checked if the following
// conditions are met:
//  1. Price feed has the same ID / time window as the SLA.
//  2. The price feed has met the threshold for minimum block updates within
//     the maximum viable window.
func (sla *PriceFeedSLA) Qualifies(priceFeed PriceFeed) (bool, error) {
	if priceFeed.ID != sla.ID {
		return false, nil
	}

	if priceFeed.MaximumViableWindow != sla.MaximumViableWindow {
		return false, fmt.Errorf("price feed %s has a different maximum viable window than the sla with same id", priceFeed.ID)
	}

	// Ensure that the price feed has an acceptable minimum block updates.
	numVotes, err := priceFeed.GetNumVotesWithWindow(uint(sla.MaximumViableWindow))
	if err != nil {
		return false, err
	}

	// Ensure that the price feed has seen enough votes.
	if uint(sla.MinimumBlockUpdates) > numVotes {
		return false, nil
	}

	return true, nil
}

// GetUptimeFromPriceFeed returns the uptime for the given SLA. The calculation for uptime is
// down below:
//
//	uptime = (number of price updates / number of blocks voted on)
//
// This is all done in the context of the maximum viable window.
func (sla *PriceFeedSLA) GetUptimeFromPriceFeed(priceFeed PriceFeed) (math.LegacyDec, error) {
	numUpdates, err := priceFeed.GetNumPriceUpdatesWithWindow(uint(sla.MaximumViableWindow))
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	numVotes, err := priceFeed.GetNumVotesWithWindow(uint(sla.MaximumViableWindow))
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	updates := math.NewIntFromUint64(uint64(numUpdates))
	votes := math.NewIntFromUint64(uint64(numVotes))

	if votes.IsZero() {
		return math.LegacyOneDec(), nil
	}

	// uptime = number of price updates / number of blocks voted on
	uptime := math.LegacyNewDecFromInt(updates).Quo(math.LegacyNewDecFromInt(votes))

	return uptime, nil
}

// ValidateBasic performs basic validation of the PriceFeedSLA returning an
// error for any failed validation criteria.
func (sla *PriceFeedSLA) ValidateBasic() error {
	if len(sla.ID) == 0 {
		return fmt.Errorf("sla id cannot be empty")
	}

	if sla.MaximumViableWindow == 0 {
		return fmt.Errorf("sla %s must have a non-zero maximum viable window", sla.ID)
	}

	if sla.ExpectedUptime.LTE(math.LegacyZeroDec()) {
		return fmt.Errorf("sla %s must have a non-negative expected uptime", sla.ID)
	}

	if sla.SlashConstant.LTE(math.LegacyZeroDec()) {
		return fmt.Errorf("sla %s must have a non-negative slashing constant k", sla.ID)
	}

	if sla.MinimumBlockUpdates == 0 {
		return fmt.Errorf("sla %s must have a non-zero minimum block updates", sla.ID)
	}

	if sla.MinimumBlockUpdates > sla.MaximumViableWindow {
		return fmt.Errorf("sla %s must have a minimum block updates less than the maximum viable window", sla.ID)
	}

	if sla.Frequency == 0 {
		return fmt.Errorf("sla %s must have a non-zero frequency", sla.ID)
	}

	if sla.Frequency > sla.MaximumViableWindow {
		return fmt.Errorf("sla %s must have a frequency less than the maximum viable window", sla.ID)
	}

	return nil
}
