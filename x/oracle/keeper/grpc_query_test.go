package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/stretchr/testify/mock"

	connecttypes "github.com/skip-mev/connect/v2/pkg/types"
	marketmaptypes "github.com/skip-mev/connect/v2/x/marketmap/types"
	"github.com/skip-mev/connect/v2/x/oracle/keeper"
	"github.com/skip-mev/connect/v2/x/oracle/types"
)

func (s *KeeperTestSuite) TestGetAllCurrencyPairs() {
	qs := keeper.NewQueryServer(s.oracleKeeper)

	// test that an error is returned if no CurrencyPairs have been registered in the module
	s.Run("an error is returned if no CurrencyPairs have been registered in the module", func() {
		// execute query
		_, err := qs.GetAllCurrencyPairs(s.ctx, nil)
		s.Require().NotNil(s.T(), err)
	})

	// test that after CurrencyPairs are registered, all of them are returned from the query
	s.Run("after CurrencyPairs are registered, all of them are returned from the query", func() {
		// insert multiple currency Pairs
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{
			Base:  "AA",
			Quote: "BB",
		}))
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{
			Base:  "CC",
			Quote: "DD",
		}))
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{
			Base:  "EE",
			Quote: "FF",
		}))

		// manually insert a new CurrencyPair as well
		s.Require().NoError(s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, connecttypes.CurrencyPair{
			Base:  "GG",
			Quote: "HH",
		}, types.QuotePrice{Price: sdkmath.NewInt(100)}))

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
			CurrencyPair: connecttypes.CurrencyPair{
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
			CurrencyPair: connecttypes.CurrencyPair{
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
		setup      func()
		req        *types.GetPriceRequest
		res        *types.GetPriceResponse
		expectPass bool
	}{
		{
			"if the request is nil, expect failure - fail",
			func() {},
			nil,
			nil,
			false,
		},
		{
			"if the currency pair is empty, expect failure - fail",
			func() {},
			&types.GetPriceRequest{
				CurrencyPair: "",
			},
			nil,
			false,
		},
		{
			"if the currency pair is malformed, expect failure - fail",
			func() {},
			&types.GetPriceRequest{
				CurrencyPair: "AABB",
			},
			nil,
			false,
		},

		{
			"if the query is for a currency pair that does not exist fail - fail",
			func() {},
			&types.GetPriceRequest{
				CurrencyPair: "DD/EE",
			},
			nil,
			false,
		},
		{
			"if the query is for a currency-pair with no price, only the nonce (0) is returned - pass",
			func() {
				s.mockMarketMapKeeper.On("GetMarket", mock.Anything, mock.Anything).Return(marketmaptypes.Market{
					Ticker: marketmaptypes.Ticker{
						CurrencyPair:     connecttypes.CurrencyPair{Base: "CC", Quote: "BB"},
						Decimals:         8,
						MinProviderCount: 3,
						Metadata_JSON:    "",
					},
				}, nil).Once()
			},
			&types.GetPriceRequest{
				CurrencyPair: "CC/BB",
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
			func() {
				s.mockMarketMapKeeper.On("GetMarket", mock.Anything, mock.Anything).Return(marketmaptypes.Market{
					Ticker: marketmaptypes.Ticker{
						CurrencyPair:     connecttypes.CurrencyPair{Base: "AA", Quote: "ETHEREUM"},
						Decimals:         18,
						MinProviderCount: 3,
						Metadata_JSON:    "",
					},
				}, nil).Once()
			},
			&types.GetPriceRequest{
				CurrencyPair: "AA/ETHEREUM",
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
			tc.setup()

			// get the response + error from the query
			res, err := qs.GetPrice(s.ctx, tc.req)
			if !tc.expectPass {
				s.Require().NotNil(err)
				return
			}

			// otherwise, assert no error, and check response
			s.Require().Nil(err)

			// check response
			s.Require().Equal(res.Nonce, tc.res.Nonce)

			// check price if possible
			if tc.res.Price != nil {
				checkQuotePriceEqual(s.T(), *tc.res.Price, *res.Price)
			}

			// check decimals
			s.Require().Equal(tc.res.Decimals, res.Decimals)

			// check id
			s.Require().Equal(tc.res.Id, res.Id)
		})
	}
}

func (s *KeeperTestSuite) TestGetCurrencyPairMappingGRPC() {
	qs := keeper.NewQueryServer(s.oracleKeeper)
	// test that after CurrencyPairs are registered, all of them are returned from the query
	s.Run("after CurrencyPairs are registered, all of them are returned from the query", func() {
		currencyPairs := []connecttypes.CurrencyPair{
			{Base: "TEST", Quote: "COIN1"},
			{Base: "TEST", Quote: "COIN2"},
			{Base: "FOO", Quote: "COIN3"},
		}
		for _, cp := range currencyPairs {
			s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, cp))
		}

		// manually insert a new CurrencyPair as well
		s.Require().NoError(s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, connecttypes.CurrencyPair{
			Base:  "TEST",
			Quote: "COIN1",
		}, types.QuotePrice{Price: sdkmath.NewInt(100)}))

		// query for pairs
		res, err := qs.GetCurrencyPairMapping(s.ctx, nil)
		s.Require().Nil(err)
		for idx, cp := range currencyPairs {
			s.Require().Equal(cp, res.CurrencyPairMapping[uint64(idx)]) //nolint:gosec
		}
	})
}
