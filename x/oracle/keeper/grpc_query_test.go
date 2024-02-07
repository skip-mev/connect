package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/oracle/keeper"
	"github.com/skip-mev/slinky/x/oracle/types"
)

func (s *KeeperTestSuite) TestGetAllCurrencyPairs() {
	qs := keeper.NewQueryServer(s.oracleKeeper)

	// test that an error is returned if no CurrencyPairs have been registered in the module
	s.T().Run("an error is returned if no CurrencyPairs have been registered in the module", func(t *testing.T) {
		// execute query
		_, err := qs.GetAllCurrencyPairs(s.ctx, nil)
		require.Nil(s.T(), err)
	})

	// test that after CurrencyPairs are registered, all of them are returned from the query
	s.T().Run("after CurrencyPairs are registered, all of them are returned from the query", func(t *testing.T) {
		// insert multiple currency Pairs
		cp1 := types.CurrencyPair{
			Base:     "AA",
			Quote:    "BB",
			Decimals: types.DefaultDecimals,
		}
		cp2 := types.CurrencyPair{
			Base:     "CC",
			Quote:    "DD",
			Decimals: types.DefaultDecimals,
		}
		cp3 := types.CurrencyPair{
			Base:     "EE",
			Quote:    "FF",
			Decimals: types.DefaultDecimals,
		}

		// insert into module
		ms := keeper.NewMsgServer(s.oracleKeeper)
		_, err := ms.AddCurrencyPairs(s.ctx, &types.MsgAddCurrencyPairs{
			CurrencyPairs: []types.CurrencyPair{cp1, cp2, cp3},
			Authority:     sdk.AccAddress(moduleAuth).String(),
		})
		require.Nil(s.T(), err)

		// manually insert a new CurrencyPair as well
		require.NoError(s.T(), s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, types.CurrencyPair{
			Base:     "GG",
			Quote:    "HH",
			Decimals: types.DefaultDecimals,
		}, types.QuotePrice{Price: sdkmath.NewInt(100)}))

		expectedCurrencyPairs := map[string]struct{}{"AA/BB": {}, "CC/DD": {}, "EE/FF": {}, "GG/HH": {}}

		// query for pairs
		res, err := qs.GetAllCurrencyPairs(s.ctx, nil)
		require.Nil(s.T(), err)

		// assert that currency-pairs are correctly returned
		for _, cp := range res.CurrencyPairs {
			_, ok := expectedCurrencyPairs[cp.Ticker()]
			require.True(t, ok)
		}
	})
}

func (s *KeeperTestSuite) TestGetPrice() {
	// set CPs on genesis for testing
	cpg := []types.CurrencyPairGenesis{
		{
			CurrencyPair: types.CurrencyPair{
				Base:     "AA",
				Quote:    "ETHEREUM",
				Decimals: types.EthereumDecimals,
			},
			CurrencyPairPrice: &types.QuotePrice{
				Price: sdkmath.NewInt(100),
			},
			Nonce: 12,
			Id:    2,
		},
		{
			CurrencyPair: types.CurrencyPair{
				Base:     "CC",
				Quote:    "BB",
				Decimals: types.DefaultDecimals,
			},
			Id: 1,
		},
	}

	// init genesis
	s.oracleKeeper.InitGenesis(s.ctx, *types.NewGenesisState(cpg, 3))

	testCP := types.CurrencyPair{Base: "AA", Quote: "ETHEREUM", Decimals: types.EthereumDecimals}

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
			nil,
			nil,
			false,
		},
		{
			"if the query is for a currency pair that does not exist fail - fail",
			&types.GetPriceRequest{
				CurrencyPairId: "DD/EE",
			},
			nil,
			false,
		},
		{
			"if the query is for a currency-pair with no price, only the nonce (0) is returned - pass",
			&types.GetPriceRequest{
				CurrencyPairId: "CC/BB",
			},
			&types.GetPriceResponse{
				Nonce:    0,
				Decimals: 8,
				Id:       1,
			},
			true,
		},
		{
			"if the query is for a currency pair that has valid price data, return the price + the nonce - pass",
			&types.GetPriceRequest{
				CurrencyPairId: testCP.Ticker(),
			},
			&types.GetPriceResponse{
				Nonce: 12,
				Price: &types.QuotePrice{
					Price: sdkmath.NewInt(100),
				},
				Decimals: 18,
				Id:       2,
			},
			true,
		},
	}

	qs := keeper.NewQueryServer(s.oracleKeeper)

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			// get the response + error from the query
			res, err := qs.GetPrice(s.ctx, tc.req)
			if !tc.expectPass {
				require.NotNil(s.T(), err)
				return
			}

			// otherwise, assert no error, and check response
			require.Nil(s.T(), err)

			// check response
			require.Equal(s.T(), res.Nonce, tc.res.Nonce)

			// check price if possible
			if tc.res.Price != nil {
				checkQuotePriceEqual(s.T(), *tc.res.Price, *res.Price)
			}

			// check decimals
			require.Equal(s.T(), tc.res.Decimals, res.Decimals)

			// check id
			require.Equal(s.T(), tc.res.Id, res.Id)
		})
	}
}
