package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
					CurrencyPair: types.CurrencyPair{
						Base:  "AA",
						Quote: "BB",
					},
				},
				{
					// invalid CurrencyPairGenesis
					CurrencyPair: types.CurrencyPair{
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
					CurrencyPair: types.CurrencyPair{
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
					CurrencyPair: types.CurrencyPair{
						Base:  "AA",
						Quote: "BB",
					},
				},
				{
					// invalid CurrencyPairGenesis
					CurrencyPair: types.CurrencyPair{
						Base:  "BB",
						Quote: "CC",
					},
				},
			},
			true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			gs := types.NewGenesisState(tc.cpgs)
			err := gs.Validate()

			if tc.expectPass {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
