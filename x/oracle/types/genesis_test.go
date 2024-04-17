package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"

	"github.com/skip-mev/slinky/x/oracle/types"
)

func TestGenesisValidation(t *testing.T) {
	tcs := []struct {
		name       string
		cpgs       []types.CurrencyPairGenesis
		expectPass bool
	}{
		{
			"if any of the currency-pair geneses are invalid - fail",
			[]types.CurrencyPairGenesis{
				{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "AA",
						Quote: "BB",
					},
				},
				{
					// invalid CurrencyPairGenesis
					CurrencyPair: slinkytypes.CurrencyPair{
						Base: "BB",
					},
				},
			},
			false,
		},
		{
			"if the CurrencyPairPrice is nil, but the nonce is non-zero - fail",
			[]types.CurrencyPairGenesis{
				{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "AA",
						Quote: "BB",
					},
					Nonce: 10,
				},
			},
			false,
		},
		{
			"if all of the currency-pair geneses are valid - pass",
			[]types.CurrencyPairGenesis{
				{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "AA",
						Quote: "BB",
					},
				},
				{
					// invalid CurrencyPairGenesis
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "BB",
						Quote: "CC",
					},
				},
			},
			true,
		},
		{
			"if any of the CurrencyPairGenesis ID's are duplicated - fail",
			[]types.CurrencyPairGenesis{
				{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "AA",
						Quote: "BB",
					},
				},
				{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "AA",
						Quote: "BB",
					},
				},
			},
			false,
		},
		{
			"if any of the CurrencyPairs are repeated - fail",
			[]types.CurrencyPairGenesis{
				{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "AA",
						Quote: "BB",
					},
				},
				{
					CurrencyPair: slinkytypes.CurrencyPair{
						Base:  "AA",
						Quote: "BB",
					},
				},
			},
			false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			gs := types.NewGenesisState(tc.cpgs)
			err := gs.Validate()

			if tc.expectPass {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}
		})
	}
}
