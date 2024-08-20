package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skip-mev/connect/v2/x/alerts/types"
)

func TestPriceBound(t *testing.T) {
	t.Run("test ValidateBasic()", func(t *testing.T) {
		cases := []struct {
			name       string
			priceBound types.PriceBound
			valid      bool
		}{
			{
				"valid price-bound",
				types.PriceBound{
					High: "1",
					Low:  "0",
				},
				true,
			},
			{
				"invalid price-bound",
				types.PriceBound{
					High: "0",
					Low:  "1",
				},
				false,
			},
			{
				"invalid price-bound high == low",
				types.PriceBound{
					High: "1",
					Low:  "1",
				},
				false,
			},
			{
				"invalid price-bound invalid strings",
				types.PriceBound{
					High: "",
					Low:  "",
				},
				false,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.priceBound.ValidateBasic()
				if tc.valid && err != nil {
					t.Fatalf("expected price-bound to be valid but got error: %s", err)
				}

				if !tc.valid && err == nil {
					t.Fatal("expected price-bound to be invalid but got no error")
				}
			})
		}
	})

	t.Run("test GetHighInt()", func(t *testing.T) {
		// invalid high value fails
		pb := types.PriceBound{
			High: "x",
			Low:  "2",
		}
		_, err := pb.GetHighInt()
		require.NotNil(t, err)

		// valid high value succeeds
		pb = types.PriceBound{
			High: "1",
			Low:  "0",
		}

		high, err := pb.GetHighInt()
		require.Nil(t, err)
		require.Equal(t, high.Uint64(), uint64(1))
	})

	t.Run("test GetLowInt()", func(t *testing.T) {
		// invalid low value fails
		pb := types.PriceBound{
			High: "1",
			Low:  "x",
		}
		_, err := pb.GetLowInt()
		require.NotNil(t, err)

		// valid low value succeeds
		pb = types.PriceBound{
			High: "2",
			Low:  "1",
		}

		low, err := pb.GetLowInt()
		require.Nil(t, err)
		require.Equal(t, low.Uint64(), uint64(1))
	})
}
