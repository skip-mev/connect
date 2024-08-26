package types

import "time"

// BlockAgeChecker is a utility type to check if incoming block heights are validly updating.
// If the block heights are not increasing and the time since the last update has exceeded
// a configurable duration, this type will report that the updates are invalid.
//
//go:generate mockery --name BlockAgeChecker --output ./mocks/ --case underscore
type BlockAgeChecker interface {
	IsHeightValid(newHeight uint64) bool
}

var _ BlockAgeChecker = &BlockAgeCheckerImpl{}

// BlockAgeCheckerImpl is an implementation of the BlockAgeChecker interface.
type BlockAgeCheckerImpl struct {
	lastHeight    uint64
	lastTimeStamp time.Time
	maxAge        time.Duration
}

// NewBlockAgeChecker returns a zeroed BlockAgeChecker using the provided maxAge.
func NewBlockAgeChecker(maxAge time.Duration) BlockAgeChecker {
	return &BlockAgeCheckerImpl{
		lastHeight:    0,
		lastTimeStamp: time.Now(),
		maxAge:        maxAge,
	}
}

// IsHeightValid returns true if:
// - the new height is greater than the last height OR
// - the time past the last block height update is less than the configured max age
// returns false if:
// - the time is past the configured max age
func (bc *BlockAgeCheckerImpl) IsHeightValid(newHeight uint64) bool {
	now := time.Now()

	if newHeight > bc.lastHeight {
		bc.lastHeight = newHeight
		bc.lastTimeStamp = now
		return true
	}

	if now.Sub(bc.lastTimeStamp) > bc.maxAge {
		return false
	}

	return true
}
