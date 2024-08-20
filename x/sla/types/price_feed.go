package types

import (
	"fmt"

	"github.com/bits-and-blooms/bitset"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
)

const (
	// NoVote indicates that the validator did not vote on the current block.
	NoVote UpdateStatus = iota

	// VoteWithPrice indicates that the validator voted with a price update. i.e. the
	// vote extension included any price update for the given currency pair.
	VoteWithPrice

	// VoteWithoutPrice indicates that the validator voted without a price update. i.e. the
	// vote extension did not include any price update for the given currency pair.
	VoteWithoutPrice
)

// UpdateStatus is an enum that represents the status of a price update.
type UpdateStatus int

// NewPriceFeed returns a new price feed for the given parameters. This is meant to be
// called every time a new SLA is created, new currency pair is added, or a new validator
// is added to the network.
func NewPriceFeed(
	maximumViableWindow uint,
	validator sdk.ConsAddress,
	currencyPair slinkytypes.CurrencyPair,
	id string,
) (PriceFeed, error) {
	updateMap := bitset.New(maximumViableWindow)
	updateMapBz, err := updateMap.MarshalBinary()
	if err != nil {
		return PriceFeed{}, err
	}

	inclusionMap := bitset.New(maximumViableWindow)
	inclusionMapBz, err := inclusionMap.MarshalBinary()
	if err != nil {
		return PriceFeed{}, err
	}

	return PriceFeed{
		MaximumViableWindow: uint64(maximumViableWindow),
		Validator:           validator.Bytes(),
		CurrencyPair:        currencyPair,
		UpdateMap:           updateMapBz,
		InclusionMap:        inclusionMapBz,
		ID:                  id,
	}, nil
}

// SetUpdate updates the state of the SLA given the status of the update.
func (feed *PriceFeed) SetUpdate(status UpdateStatus) error {
	index := uint(feed.Index)

	inclusionMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := inclusionMap.UnmarshalBinary(feed.InclusionMap); err != nil {
		return err
	}

	updateMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := updateMap.UnmarshalBinary(feed.UpdateMap); err != nil {
		return err
	}

	switch status {
	case VoteWithPrice:
		// Vote + price update
		inclusionMap.Set(index)
		updateMap.Set(index)
	case VoteWithoutPrice:
		// Vote without price update
		inclusionMap.Set(index)
		updateMap.Clear(index)
	default:
		// Vote was not included in previous block
		inclusionMap.Clear(index)
		updateMap.Clear(index)
	}

	feed.Index = (feed.Index + 1) % feed.MaximumViableWindow

	// Marshal the bitsets.
	var err error
	feed.UpdateMap, err = updateMap.MarshalBinary()
	if err != nil {
		return err
	}

	feed.InclusionMap, err = inclusionMap.MarshalBinary()
	if err != nil {
		return err
	}

	return nil
}

// GetInclusionBit returns the bit at the given index in the inclusion map.
func (feed *PriceFeed) GetInclusionBit(index uint) (bool, error) {
	inclusionMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := inclusionMap.UnmarshalBinary(feed.InclusionMap); err != nil {
		return false, err
	}

	return inclusionMap.Test(index), nil
}

// GetInclusionCount returns the total number of votes in the maximum viable window.
func (feed *PriceFeed) GetInclusionCount() (uint, error) {
	inclusionMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := inclusionMap.UnmarshalBinary(feed.InclusionMap); err != nil {
		return 0, err
	}

	return inclusionMap.Count(), nil
}

// GetUpdateBit returns the bit at the given index in the update map.
func (feed *PriceFeed) GetUpdateBit(index uint) (bool, error) {
	updateMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := updateMap.UnmarshalBinary(feed.UpdateMap); err != nil {
		return false, err
	}

	return updateMap.Test(index), nil
}

// GetUpdateCount returns the number of updates in the maximum viable window.
func (feed *PriceFeed) GetUpdateCount() (uint, error) {
	updateMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := updateMap.UnmarshalBinary(feed.UpdateMap); err != nil {
		return 0, err
	}

	return updateMap.Count(), nil
}

// GetNumPriceUpdatesWithWindow returns the number of price updates in the moving window. This
// corresponds to the number of blocks that the validator has voted on and included
// a price update in the previous n blocks.
func (feed *PriceFeed) GetNumPriceUpdatesWithWindow(n uint) (uint, error) {
	bitRange, err := feed.getBitRange(n)
	if err != nil {
		return 0, err
	}

	updateMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := updateMap.UnmarshalBinary(feed.UpdateMap); err != nil {
		return 0, err
	}

	return updateMap.Intersection(bitRange).Count(), nil
}

// GetNumVotesWithWindow returns the number of blocks the validator has voted on
// in the previous n blocks.
func (feed *PriceFeed) GetNumVotesWithWindow(n uint) (uint, error) {
	bitRange, err := feed.getBitRange(n)
	if err != nil {
		return 0, err
	}

	inclusionMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := inclusionMap.UnmarshalBinary(feed.InclusionMap); err != nil {
		return 0, err
	}

	return inclusionMap.Intersection(bitRange).Count(), nil
}

// Stringify returns a string representation of the price feed. Primarily used for
// debugging purposes.
func (feed *PriceFeed) Stringify() string {
	inclusionMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := inclusionMap.UnmarshalBinary(feed.InclusionMap); err != nil {
		panic(err)
	}

	updateMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := updateMap.UnmarshalBinary(feed.UpdateMap); err != nil {
		panic(err)
	}

	return fmt.Sprintf(`Price Feed:
	Maximum Viable Window: %d
	Validator: %s
	Currency Pair: %s
	Update Map: %s
	Inclusion Map: %s
	Index: %d
	ID: %s`,
		feed.MaximumViableWindow,
		feed.Validator,
		feed.CurrencyPair,
		updateMap.DumpAsBits(),
		inclusionMap.DumpAsBits(),
		feed.Index,
		feed.ID,
	)
}

// ValidateBasic performs basic validation of the price feed.
func (feed *PriceFeed) ValidateBasic() error {
	if len(feed.ID) == 0 {
		return fmt.Errorf("sla id cannot be empty")
	}

	if feed.MaximumViableWindow == 0 {
		return fmt.Errorf("sla %s must have a non-zero maximum viable window", feed.ID)
	}

	inclusionMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := inclusionMap.UnmarshalBinary(feed.InclusionMap); err != nil {
		return err
	}

	updateMap := bitset.New(uint(feed.MaximumViableWindow))
	if err := updateMap.UnmarshalBinary(feed.UpdateMap); err != nil {
		return err
	}

	if feed.Validator == nil {
		return fmt.Errorf("validator cannot be nil")
	}

	return feed.CurrencyPair.ValidateBasic()
}

// getBitRange returns a bitset that represents the range of bits in the moving
// window that we are interested in. This returns a bitset with all bits set to
// 1 if the range is valid.
func (feed *PriceFeed) getBitRange(previousNumBlocks uint) (*bitset.BitSet, error) {
	maximumViableWindow := uint(feed.MaximumViableWindow)
	if previousNumBlocks > maximumViableWindow {
		return nil, fmt.Errorf(
			"previousNumBlocks cannot be greater than the maximum viable window; got %d, expected %d",
			previousNumBlocks,
			feed.MaximumViableWindow,
		)
	}

	bitRange := bitset.New(maximumViableWindow)
	index := uint(feed.Index)
	if previousNumBlocks < index {
		lowerBound := index - previousNumBlocks
		bitRange = bitRange.FlipRange(lowerBound, index)
	} else {
		lowerBound := maximumViableWindow - (previousNumBlocks - index)
		bitRange = bitRange.FlipRange(lowerBound, maximumViableWindow)
		bitRange = bitRange.FlipRange(0, index)
	}

	return bitRange, nil
}
