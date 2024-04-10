package volatile

import (
	"math"
	"math/big"
	"time"
)

type TimeProvider func() time.Time

var (
	// dailySeconds is the number of seconds in a day.
	dailySeconds = float64(24 * 60 * 60)
	// normalizedPhaseSize is the radians in our repeating function (4π) before adjusting for frequency and time.
	normalizedPhaseSize = float64(4)
)

// GetVolatilePrice generates a time-based price value. The value follows a cosine wave
// function, but that includes jumps from the lowest value to the highest value (and vice versa)
// once per period. The general formula is written below.
// - price = offset * (1 + amplitude * cosVal)
// - cosVal = math.Cos(radians)
// - radians = (cosinePhase <= 0.5 ? cosinePhase * 4 : cosinePhase * 4 - 1) * π
// - cosinePhase = (frequency * unix_time(in seconds) / dailySeconds) % 1.
func GetVolatilePrice(tp TimeProvider, amplitude float64, offset float64, frequency float64) *big.Float {
	// The phase is the location of the final price within our repeating price function.
	// The resulting value is taken mod(1) i.e. it is between 0 and 1 inclusive
	cosinePhase := math.Mod(
		frequency*float64(tp().Unix())/dailySeconds,
		1,
	)
	// To achieve our price "jump", we implement a piecewise function at 0.5
	radians := cosinePhase * normalizedPhaseSize
	if cosinePhase > 0.5 {
		radians -= float64(1)
	}
	radians *= math.Pi
	cosVal := math.Cos(radians)
	return big.NewFloat(offset * (1 + amplitude*cosVal))
}
