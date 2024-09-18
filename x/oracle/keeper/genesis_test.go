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

func (s *KeeperTestSuite) TestInitGenesis() {
	tcs := []struct {
		name       string
		gs         types.GenesisState
		expectPass bool
	}{
		{
			"if the genesis-state is incorrectly formatted - fail",
			types.GenesisState{
				CurrencyPairGenesis: []types.CurrencyPairGenesis{
					{
						CurrencyPair: connecttypes.CurrencyPair{
							Base:  "AA",
							Quote: "BB",
						},
					},
					{
						// invalid CurrencyPairGenesis
						CurrencyPair: connecttypes.CurrencyPair{
							Base: "BB",
						},
					},
				},
			},
			false,
		},
		{
			"if the genesis-state is correctly formatted - pass",
			types.GenesisState{
				CurrencyPairGenesis: []types.CurrencyPairGenesis{
					{
						CurrencyPair: connecttypes.CurrencyPair{
							Base:  "AA",
							Quote: "BB",
						},
						Id: 0,
					},
					{
						CurrencyPair: connecttypes.CurrencyPair{
							Base:  "BB",
							Quote: "CC",
						},
						CurrencyPairPrice: &types.QuotePrice{
							Price: sdkmath.NewInt(100),
						},
						Nonce: 12,
						Id:    1,
					},
				},
				NextId: 2,
			},
			true,
		},
	}

	for _, tc := range tcs {
		s.Run(tc.name, func() {
			if !tc.expectPass {
				// call init-genesis, and catch the panic
				catchPanic(s.T(), s.oracleKeeper, s.ctx, tc.gs)
			} else {
				// call init-genesis
				s.oracleKeeper.InitGenesis(s.ctx, tc.gs)

				// expect all the currency-pairs to be stored in state
				for _, cpg := range tc.gs.CurrencyPairGenesis {
					// get the quote-price
					qp, err := s.oracleKeeper.GetPriceForCurrencyPair(s.ctx, cpg.CurrencyPair)

					// check equality of quote-price if one is given
					if cpg.CurrencyPairPrice != nil {
						// check equality
						s.Require().Nil(err)
						checkQuotePriceEqual(s.T(), qp, *cpg.CurrencyPairPrice)
					} else {
						// assert that no price exists for the currency-pair
						s.Require().NotNil(err)
					}

					// get nonce, and check equality
					nonce, err := s.oracleKeeper.GetNonceForCurrencyPair(s.ctx, cpg.CurrencyPair)
					s.Require().Nil(err)

					// check equality of nonces
					s.Require().Equal(nonce, cpg.Nonce)

					// check equality of ids
					id, ok := s.oracleKeeper.GetIDForCurrencyPair(s.ctx, cpg.CurrencyPair)
					s.Require().True(ok)

					s.Require().Equal(id, cpg.Id)
				}
			}
		})
	}
}

func catchPanic(t *testing.T, k keeper.Keeper, ctx sdk.Context, gs types.GenesisState) {
	t.Helper()

	defer func() {
		err := recover()
		require.NotNil(t, err)
	}()
	// call init-genesis
	k.InitGenesis(ctx, gs)
}

func (s *KeeperTestSuite) TestExportGenesis() {
	s.Run("ExportGenesis with all valid QuotePrices", func() {
		// insert multiple currency pairs
		cp1 := connecttypes.CurrencyPair{
			Base:  "AA",
			Quote: "BB",
		}
		qp1 := types.QuotePrice{
			Price: sdkmath.NewInt(100),
		}
		cp2 := connecttypes.CurrencyPair{
			Base:  "CC",
			Quote: "DD",
		}
		qp2 := types.QuotePrice{
			Price: sdkmath.NewInt(120),
		}

		// insert
		s.Require().Nil(s.oracleKeeper.CreateCurrencyPair(s.ctx, cp1))
		s.Require().Nil(s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp1, qp1))

		s.Require().Nil(s.oracleKeeper.CreateCurrencyPair(s.ctx, cp2))
		s.Require().Nil(s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp2, qp2))

		// insert
		s.Require().Nil(s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp1, qp1))

		// export genesis
		gs := s.oracleKeeper.ExportGenesis(s.ctx)
		s.Require().Equal(len(gs.CurrencyPairGenesis), 2)
		expectedCurrencyPairs := map[string]types.QuotePrice{"AA/BB": qp1, "CC/DD": qp2}
		expectedNonces := map[string]uint64{"AA/BB": 2, "CC/DD": 1}

		for _, cpg := range gs.CurrencyPairGenesis {
			qp, ok := expectedCurrencyPairs[cpg.CurrencyPair.String()]
			s.Require().True(ok)
			// check equality for quote-prices
			checkQuotePriceEqual(s.T(), qp, *cpg.CurrencyPairPrice)
			// check equality of nonces
			nonce, ok := expectedNonces[cpg.CurrencyPair.String()]
			s.Require().True(ok)
			s.Require().Equal(nonce, cpg.Nonce)
		}
	})

	s.Run("ExportGenesis with some un-price-initialized CurrencyPairs", func() {
		// initialize genesis w/ price-data
		gs := types.GenesisState{
			CurrencyPairGenesis: []types.CurrencyPairGenesis{
				{
					CurrencyPair: connecttypes.CurrencyPair{
						Base:  "AA",
						Quote: "BB",
					},
					CurrencyPairPrice: &types.QuotePrice{
						Price: sdkmath.NewInt(100),
					},
					Nonce: 100,
					Id:    0,
				},
				{
					CurrencyPair: connecttypes.CurrencyPair{
						Base:  "CC",
						Quote: "DD",
					},
					CurrencyPairPrice: &types.QuotePrice{
						Price: sdkmath.NewInt(101),
					},
					Nonce: 101,
					Id:    1,
				},
			},
			NextId: 2,
		}
		// init genesis
		s.oracleKeeper.InitGenesis(s.ctx, gs)

		// add currency pairs
		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{
			Base:  "EE",
			Quote: "FF",
		}))

		s.Require().NoError(s.oracleKeeper.CreateCurrencyPair(s.ctx, connecttypes.CurrencyPair{
			Base:  "GG",
			Quote: "HH",
		}))

		// setup expected values
		expectedCurrencyPairs := map[string]struct{}{"AA/BB": {}, "CC/DD": {}, "EE/FF": {}, "GG/HH": {}}
		expectedQuotePrices := map[string]types.QuotePrice{
			"AA/BB": {
				Price: sdkmath.NewInt(100),
			},
			"CC/DD": {
				Price: sdkmath.NewInt(101),
			},
		}
		expectedNonces := map[string]uint64{"AA/BB": 100, "CC/DD": 101}

		// ExportGenesis
		egs := s.oracleKeeper.ExportGenesis(s.ctx)

		// iterate over CurrencyPairGeneses in egs
		for _, cpg := range egs.CurrencyPairGenesis {
			// expect that all currency-pairs in gen-state are expected
			cps := cpg.CurrencyPair.String()
			_, ok := expectedCurrencyPairs[cps]
			s.Require().True(ok)

			// expect that if a CurrencyPrice exists, that it is expected
			if cpg.CurrencyPairPrice != nil {
				qp, ok := expectedQuotePrices[cps]
				s.Require().True(ok)

				// assert equality of QuotePrice
				checkQuotePriceEqual(s.T(), qp, *cpg.CurrencyPairPrice)

				nonce, ok := expectedNonces[cps]
				s.Require().True(ok)
				// assert equality of Nonce
				s.Require().Equal(cpg.Nonce, nonce)
			} else {
				s.Require().Equal(cpg.Nonce, uint64(0))
			}

			// check IDs
			id, ok := s.oracleKeeper.GetIDForCurrencyPair(s.ctx, cpg.CurrencyPair)

			s.Require().True(ok)

			s.Require().Equal(id, cpg.Id)
		}
	})
}
