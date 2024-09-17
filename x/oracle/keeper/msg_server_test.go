package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/oracle/keeper"
	"github.com/skip-mev/connect/v2/x/oracle/types"
)

func (s *KeeperTestSuite) TestMsgAddCurrencyPairs() {
	s.SetupWithNoMMKeeper()

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
			"if the authority is not the authority of the module - fail",
			&types.MsgAddCurrencyPairs{
				Authority: sdk.AccAddress("not-authority").String(),
				CurrencyPairs: []connecttypes.CurrencyPair{
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
				Authority: sdk.AccAddress(moduleAuth).String(),
				CurrencyPairs: []connecttypes.CurrencyPair{
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
				Authority: sdk.AccAddress(moduleAuth).String(),
				CurrencyPairs: []connecttypes.CurrencyPair{
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

	initCP := connecttypes.CurrencyPair{
		Base:  "E",
		Quote: "F",
	}

	// set genesis quote price for E/F
	gs := types.GenesisState{
		CurrencyPairGenesis: []types.CurrencyPairGenesis{
			{
				CurrencyPair: initCP,
				CurrencyPairPrice: &types.QuotePrice{
					Price: sdkmath.NewInt(100),
				},
				Nonce: 100,
				Id:    1,
			},
		},
		NextId: 101,
	}
	s.oracleKeeper.InitGenesis(s.ctx, gs)

	// construct message server + wrap context
	ms := keeper.NewMsgServer(s.oracleKeeper)
	for _, tc := range tcs {
		s.Run(tc.name, func() {
			// execute message
			_, err := ms.AddCurrencyPairs(s.ctx, tc.req)

			// expect failure if necessary
			if !tc.expectPass {
				s.Require().NotNil(err)
				return
			}

			// otherwise, check that insertions executed faithfully
			s.Require().Nil(err)

			// check all currency pairs were inserted
			for _, cp := range tc.req.CurrencyPairs {
				// get nonce for cpg.CurrencyPair
				nonce, err := s.oracleKeeper.GetNonceForCurrencyPair(s.ctx, cp)
				s.Require().Nil(err)

				// check the nonce is correct (if the cp had already existed in state, check that it was not overwritten)
				if cp.String() == "E/F" {
					s.Require().Equal(nonce, uint64(100))
				} else {
					s.Require().Equal(nonce, uint64(0))
				}
			}
		})
	}
}

func (s *KeeperTestSuite) TestMsgRemoveCurrencyPairs() {
	s.SetupWithNoMMKeeper()

	// insert CurrencyPairs that will be deleted in the test-cases
	cp1 := connecttypes.CurrencyPair{
		Base:  "AA",
		Quote: "BB",
	}
	cp2 := connecttypes.CurrencyPair{
		Base:  "CC",
		Quote: "DD",
	}
	gs := types.GenesisState{
		CurrencyPairGenesis: []types.CurrencyPairGenesis{
			{
				CurrencyPair: cp1,
				CurrencyPairPrice: &types.QuotePrice{
					Price: sdkmath.NewInt(100),
				},
				Nonce: 100,
				Id:    1,
			},
			{
				CurrencyPair: cp2,
				CurrencyPairPrice: &types.QuotePrice{
					Price: sdkmath.NewInt(100),
				},
				Nonce: 101,
				Id:    2,
			},
		},
		NextId: 102,
	}
	// init genesis
	s.oracleKeeper.InitGenesis(s.ctx, gs)

	// sanity check, assert existence of cps
	// cp1
	qpn, err := s.oracleKeeper.GetPriceWithNonceForCurrencyPair(s.ctx, cp1)
	require.Nil(s.T(), err)
	require.Equal(s.T(), qpn.Nonce(), uint64(100))
	require.Equal(s.T(), qpn.Price.Int64(), int64(100))

	// cp2
	qpn, err = s.oracleKeeper.GetPriceWithNonceForCurrencyPair(s.ctx, cp2)
	require.Nil(s.T(), err)
	require.Equal(s.T(), qpn.Nonce(), uint64(101))
	require.Equal(s.T(), qpn.Price.Int64(), int64(100))

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
			"if the authority is correct + formatted, and the currency pairs are valid - pass",
			&types.MsgRemoveCurrencyPairs{
				Authority: sdk.AccAddress(moduleAuth).String(),
				CurrencyPairIds: []string{
					"AA/BB", "CC/DD",
				},
			},
			true,
		},
		{
			"errors if the currency pair does not exist in state",
			&types.MsgRemoveCurrencyPairs{
				Authority: sdk.AccAddress(moduleAuth).String(),
				CurrencyPairIds: []string{
					"MOG/USD",
				},
			},
			false,
		},
	}

	ms := keeper.NewMsgServer(s.oracleKeeper)
	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			// execute message
			_, err := ms.RemoveCurrencyPairs(s.ctx, tc.req)

			if !tc.expectPass {
				require.NotNil(s.T(), err)
				return
			}

			// otherwise, assert no error
			require.Nil(s.T(), err)

			// check that all currency-pairs were removed
			for _, cps := range tc.req.CurrencyPairIds {
				// get currency pair from request
				cp, err := connecttypes.CurrencyPairFromString(cps)
				require.Nil(s.T(), err)

				// assert that currency-pair was removed
				require.False(t, s.oracleKeeper.HasCurrencyPair(s.ctx, cp))
			}
		})
	}
}
