package types_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/x/oracle/types"
)

func TestQuotePrice(t *testing.T) {
	tcs := []struct {
		name       string
		quotePrice types.QuotePrice
		err        error
	}{
		{
			"negative price",
			types.QuotePrice{
				Price:          math.NewInt(-1),
				BlockTimestamp: time.Now().UTC(),
				BlockHeight:    1,
			},
			fmt.Errorf("price cannot be negative: %s", math.NewInt(-1)),
		},
		{
			"zero price",
			types.QuotePrice{
				Price:          math.NewInt(0),
				BlockTimestamp: time.Now().UTC(),
				BlockHeight:    1,
			},
			nil,
		},
		{
			"positive price",
			types.QuotePrice{
				Price:          math.NewInt(1),
				BlockTimestamp: time.Now().UTC(),
				BlockHeight:    1,
			},
			nil,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.quotePrice.ValidateBasic()
			require.Equal(t, tc.err, err)
		})
	}
}
