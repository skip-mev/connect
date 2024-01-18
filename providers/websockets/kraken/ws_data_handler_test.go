package kraken_test

import (
	"fmt"
	"testing"

	"github.com/skip-mev/slinky/providers/websockets/kraken"
	"github.com/stretchr/testify/require"
)

func TestDecodeTickerResponseMessage(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		expected kraken.TickerResponseMessage
		expErr   bool
	}{
		{
			name:     "valid response",
			response: `[340,{"a":["42694.60000",31,"31.27308189"],"b":["42694.50000",1,"1.01355072"],"c":["42694.60000","0.00455773"],"v":["2068.49653432","2075.61202911"],"p":["42596.41907","42598.31137"],"t":[21771,22049],"l":["42190.20000","42190.20000"],"h":["43165.00000","43165.00000"],"o":["43134.70000","43159.20000"]},"ticker","XBT/USD"]`,
			expected: kraken.TickerResponseMessage{
				ChannelID: 340,
				TickerData: kraken.TickerData{
					VolumeWeightedAveragePrice: []string{"42596.41907", "42598.31137"},
				},
				ChannelName: "ticker",
				Pair:        "XBT/USD",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := kraken.DecodeTickerResponseMessage([]byte(tc.response))
			fmt.Println(actual, err)
			if tc.expErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, actual)
		})
	}
}
