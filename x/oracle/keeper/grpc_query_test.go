package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/oracle/keeper"
	"github.com/skip-mev/slinky/x/oracle/types"
)

func (s *KeeperTestSuite) TestGetAllCurrencyPairs() {
	qs := keeper.NewQueryServer(s.oracleKeeper)

	// test that an error is returned if no CurrencyPairs have been registered in the module
	s.Run("an error is returned if no CurrencyPairs have been registered in the module", func() {
		// execute query
		_, err := qs.GetAllCurrencyPairs(s.ctx, nil)
		assert.Nil(s.T(), err)
	})

	// test that after CurrencyPairs are registered, all of them are returned from the query
	s.Run("after CurrencyPairs are registered, all of them are returned from the query", func() {
		// insert multiple currency Pairs
		cp1 := slinkytypes.CurrencyPair{
			Base:  "AA",
			Quote: "BB",
		}
		cp2 := slinkytypes.CurrencyPair{
			Base:  "CC",
			Quote: "DD",
		}
		cp3 := slinkytypes.CurrencyPair{
			Base:  "EE",
			Quote: "FF",
		}

		// insert into module
		ms := keeper.NewMsgServer(s.oracleKeeper)
		_, err := ms.AddCurrencyPairs(s.ctx, &types.MsgAddCurrencyPairs{
			CurrencyPairs: []slinkytypes.CurrencyPair{cp1, cp2, cp3},
			Authority:     sdk.AccAddress([]byte(moduleAuth)).String(),
		})
		s.Require().Nil(err)

		// manually insert a new CurrencyPair as well
		s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, slinkytypes.CurrencyPair{
			Base:  "GG",
			Quote: "HH",
		}, types.QuotePrice{Price: sdkmath.NewInt(100)})

		expectedCurrencyPairs := map[string]struct{}{"AA/BB": {}, "CC/DD": {}, "EE/FF": {}, "GG/HH": {}}

		// query for pairs
		res, err := qs.GetAllCurrencyPairs(s.ctx, nil)
		s.Require().Nil(err)

		// assert that currency-pairs are correctly returned
		for _, cp := range res.CurrencyPairs {
			_, ok := expectedCurrencyPairs[cp.String()]
			s.Require().True(ok)
		}
	})
}

func (s *KeeperTestSuite) TestGetPrice() {
	// set CPs on genesis for testing
	cpg := []types.CurrencyPairGenesis{
		{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "AA",
				Quote: "ETHEREUM",
			},
			CurrencyPairPrice: &types.QuotePrice{
				Price: sdkmath.NewInt(100),
			},
			Nonce: 12,
			Id:    2,
		},
		{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "CC",
				Quote: "BB",
			},
			Id: 1,
		},
	}

	// init genesis
	s.oracleKeeper.InitGenesis(s.ctx, *types.NewGenesisState(cpg, 3))

	tcs := []struct {
		name       string
		req        *types.GetPriceRequest
		res        *types.GetPriceResponse
		expectPass bool
	}{
		{
			"if the request is nil, expect failure - fail",
			nil,
			nil,
			false,
		},
		{
			"if the currency pair selector is nil, expect failure - fail",
			&types.GetPriceRequest{
				CurrencyPairSelector: nil,
			},
			nil,
			false,
		},
		{
			"if the currency pair selector's currency pair is nil, expect failure - fail",
			&types.GetPriceRequest{
				CurrencyPairSelector: &types.GetPriceRequest_CurrencyPair{CurrencyPair: nil},
			},
			nil,
			false,
		},
		{
			"if the query is for a currency pair that does not exist fail - fail",
			&types.GetPriceRequest{
				CurrencyPairSelector: &types.GetPriceRequest_CurrencyPairId{CurrencyPairId: "DD/EE"},
			},
			nil,
			false,
		},
		{
			"if the query is for a currency-pair with no price, only the nonce (0) is returned - pass",
			&types.GetPriceRequest{
				CurrencyPairSelector: &types.GetPriceRequest_CurrencyPairId{CurrencyPairId: "CC/BB"},
			},
			&types.GetPriceResponse{
				Nonce:    0,
				Decimals: uint64(8),
				Id:       1,
			},
			true,
		},
		{
			"if the query is for a currency pair that has valid price data, return the price + the nonce - pass",
			&types.GetPriceRequest{
				CurrencyPairSelector: &types.GetPriceRequest_CurrencyPair{CurrencyPair: &slinkytypes.CurrencyPair{Base: "AA", Quote: "ETHEREUM"}},
			},
			&types.GetPriceResponse{
				Nonce: 12,
				Price: &types.QuotePrice{
					Price: sdkmath.NewInt(100),
				},
				Decimals: uint64(18),
				Id:       2,
			},
			true,
		},
	}

	qs := keeper.NewQueryServer(s.oracleKeeper)

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			// get the response + error from the query
			res, err := qs.GetPrice(s.ctx, tc.req)
			if !tc.expectPass {
				assert.NotNil(s.T(), err)
				return
			}

			// otherwise, assert no error, and check response
			assert.Nil(s.T(), err)

			// check response
			assert.Equal(s.T(), res.Nonce, tc.res.Nonce)

			// check price if possible
			if tc.res.Price != nil {
				checkQuotePriceEqual(s.T(), *tc.res.Price, *res.Price)
			}

			// check decimals
			assert.Equal(s.T(), tc.res.Decimals, res.Decimals)

			// check id
			assert.Equal(s.T(), tc.res.Id, res.Id)
		})
	}
}
