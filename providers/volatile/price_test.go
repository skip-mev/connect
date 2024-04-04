package volatile_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/providers/volatile"
)

const dailySeconds = 24 * 60 * 60

func TestGetVolatilePrice(t *testing.T) {
	testCases := []struct {
		name          string
		tp            volatile.TimeProvider
		amplitude     float64
		offset        float64
		frequency     float64
		expectedPrice *big.Int
	}{
		{
			name:          "test cosinePhase 0",
			tp:            func() time.Time { return time.Unix(0, 0) },
			amplitude:     float64(0.95),
			offset:        float64(100),
			frequency:     float64(1),
			expectedPrice: big.NewInt(195),
		},
		{
			name:          "test cosinePhase .25",
			tp:            func() time.Time { return time.Unix(25*dailySeconds/100, 0) },
			amplitude:     float64(0.95),
			offset:        float64(100),
			frequency:     float64(1),
			expectedPrice: big.NewInt(5),
		},
		{
			name:          "test cosinePhase .26",
			tp:            func() time.Time { return time.Unix(dailySeconds/100, 0) },
			amplitude:     float64(0.95),
			offset:        float64(100),
			frequency:     float64(1),
			expectedPrice: big.NewInt(194),
		},
		{
			name:          "test cosinePhase .5",
			tp:            func() time.Time { return time.Unix(50*dailySeconds/100, 0) },
			amplitude:     float64(0.95),
			offset:        float64(100),
			frequency:     float64(1),
			expectedPrice: big.NewInt(195),
		},
		{
			name:          "test cosinePhase .51",
			tp:            func() time.Time { return time.Unix(51*dailySeconds/100, 0) },
			amplitude:     float64(0.95),
			offset:        float64(100),
			frequency:     float64(1),
			expectedPrice: big.NewInt(5),
		},
		{
			name:          "test cosinePhase .99",
			tp:            func() time.Time { return time.Unix(99*dailySeconds/100, 0) },
			amplitude:     float64(0.95),
			offset:        float64(100),
			frequency:     float64(1),
			expectedPrice: big.NewInt(5),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expectedPrice, volatile.GetVolatilePrice(tc.tp, tc.amplitude, tc.offset, tc.frequency))
		})
	}
}
