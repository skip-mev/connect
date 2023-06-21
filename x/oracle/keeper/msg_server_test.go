package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/x/oracle/keeper"
	"github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/assert"
)

func (s *KeeperTestSuite) TestMsgAddCurrencyPairs() {
	tcs := []struct {
		name       string
		req        *types.MsgAddCurrencyPairs
		expectPass bool
	}{
		{
			"if the request is empty - fail",
			nil,
			false,
		},
		{
			"if the message is incorrectly formatted (authority) - fail",
			&types.MsgAddCurrencyPairs{
				Authority: "abc",
			},
			false,
		},
		{
			"if the message is incorrectly formatted (CurrencyPairs) - fail",
			&types.MsgAddCurrencyPairs{
				Authority: sdk.AccAddress([]byte("abc")).String(),
				CurrencyPairs: []types.CurrencyPair{
					// incorrectly formatted currency-pair
					{
						Base:  "A",
						Quote: "b",
					},
				},
			},
			false,
		},
		{
			"if the authority is not the authority of the module - fail",
			&types.MsgAddCurrencyPairs{
				Authority: sdk.AccAddress([]byte("not-authority")).String(),
				CurrencyPairs: []types.CurrencyPair{
					{
						Base:  "A",
						Quote: "B",
					},
				},
			},
			false,
		},
		{
			"if the authority is correct + formatted, and the currency pairs are valid - pass",
			&types.MsgAddCurrencyPairs{
				Authority: sdk.AccAddress([]byte(moduleAuth)).String(),
				CurrencyPairs: []types.CurrencyPair{
					{
						Base:  "A",
						Quote: "B",
					},
					{
						Base:  "C",
						Quote: "D",
					},
				},
			},
			true,
		},
		{
			"if there is a CurrencyPair that already exists in module, it is not overwritten",
			&types.MsgAddCurrencyPairs{
				Authority: sdk.AccAddress([]byte(moduleAuth)).String(),
				CurrencyPairs: []types.CurrencyPair{
					{
						Base:  "A",
						Quote: "B",
					},
					{
						Base:  "C",
						Quote: "D",
					},
					{
						Base:  "E",
						Quote: "F",
					},
				},
			},
			true,
		},
	}

	initCP := types.CurrencyPair{
		Base:  "E",
		Quote: "F",
	}

	// set genesis quote price for E/F
	gs := types.GenesisState{
		CurrencyPairGenesis: []types.CurrencyPairGenesis{
			{
				CurrencyPair: initCP,
				CurrencyPairPrice: &types.QuotePrice{
					Price: sdk.NewInt(100),
				},
				Nonce: 100,
			},
		},
	}
	s.oracleKeeper.InitGenesis(s.ctx, gs)

	// construct message server + wrap context
	ms := keeper.NewMsgServer(s.oracleKeeper)
	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			// execute message
			_, err := ms.AddCurrencyPairs(s.ctx, tc.req)

			// expect failure if necessary
			if !tc.expectPass {
				assert.NotNil(s.T(), err)
				return
			}

			// otherwise, check that insertions executed faithfully
			assert.Nil(s.T(), err)

			// check all currency pairs were inserted
			for _, cp := range tc.req.CurrencyPairs {
				// get nonce for cpg.CurrencyPair
				nonce, err := s.oracleKeeper.GetNonceForCurrencyPair(s.ctx, cp)
				assert.Nil(s.T(), err)

				// check the nonce is correct (if the cp had already existed in state, check that it was not overwritten)
				if cp.ToString() == "E/F" {
					assert.Equal(s.T(), nonce, uint64(100))
				} else {
					assert.Equal(s.T(), nonce, uint64(0))
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgRemoveCurrencyPairs() {
	// insert CurrencyPairs that will be deleted in the test-cases
	cp1 := types.CurrencyPair{
		Base:  "AA",
		Quote: "BB",
	}
	cp2 := types.CurrencyPair{
		Base:  "CC",
		Quote: "DD",
	}
	gs := types.GenesisState{
		CurrencyPairGenesis: []types.CurrencyPairGenesis{
			{
				CurrencyPair: cp1,
				CurrencyPairPrice: &types.QuotePrice{
					Price: sdk.NewInt(100),
				},
				Nonce: 100,
			},
			{
				CurrencyPair: cp2,
				CurrencyPairPrice: &types.QuotePrice{
					Price: sdk.NewInt(100),
				},
				Nonce: 101,
			},
		},
	}
	// init genesis
	s.oracleKeeper.InitGenesis(s.ctx, gs)

	// sanity check, assert existence of cps
	// cp1
	qpn, err := s.oracleKeeper.GetPriceWithNonceForCurrencyPair(s.ctx, cp1)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), qpn.Nonce(), uint64(100))
	assert.Equal(s.T(), qpn.Price.Int64(), int64(100))

	// cp2
	qpn, err = s.oracleKeeper.GetPriceWithNonceForCurrencyPair(s.ctx, cp2)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), qpn.Nonce(), uint64(101))
	assert.Equal(s.T(), qpn.Price.Int64(), int64(100))

	// define test-cases
	tcs := []struct {
		name       string
		req        *types.MsgRemoveCurrencyPairs
		expectPass bool
	}{
		{
			"if the request is empty - fail",
			nil,
			false,
		},
		{
			"if the message is incorrectly formatted (authority) - fail",
			&types.MsgRemoveCurrencyPairs{
				Authority: "abc",
			},
			false,
		},
		{
			"if the message is incorrectly formatted (CurrencyPairIDs) - fail",
			&types.MsgRemoveCurrencyPairs{
				Authority: sdk.AccAddress([]byte("abc")).String(),
				CurrencyPairIds: []string{
					// incorrectly formatted currency-pair
					"abc", "AA/BB",
				},
			},
			false,
		},
		{
			"if the authority is correct + formatted, and the currency pairs are valid - pass",
			&types.MsgRemoveCurrencyPairs{
				Authority: sdk.AccAddress([]byte(moduleAuth)).String(),
				CurrencyPairIds: []string{
					"AA/BB", "CC/DD",
				},
			},
			true,
		},
	}

	ms := keeper.NewMsgServer(s.oracleKeeper)
	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			// execute message
			_, err := ms.RemoveCurrencyPairs(s.ctx, tc.req)

			if !tc.expectPass {
				assert.NotNil(s.T(), err)
				return
			}

			// otherwise, assert no error
			assert.Nil(s.T(), err)

			// check that all currency-pairs were removed
			for _, cps := range tc.req.CurrencyPairIds {
				// get currency pair from request
				cp, err := types.CurrencyPairFromString(cps)
				assert.Nil(s.T(), err)

				// assert that currency-pair was removed
				_, err = s.oracleKeeper.GetNonceForCurrencyPair(s.ctx, cp)
				_, ok := err.(*types.CurrencyPairNotExistError)
				assert.True(s.T(), ok)
			}
		})
	}
}
