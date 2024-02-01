package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/oracle/keeper"
	"github.com/skip-mev/slinky/x/oracle/types"
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
						CurrencyPair: types.CurrencyPair{
							Base:     "AA",
							Quote:    "BB",
							Decimals: types.DefaultDecimals,
						},
					},
					{
						// invalid CurrencyPairGenesis
						CurrencyPair: types.CurrencyPair{
							Base:     "BB",
							Decimals: types.DefaultDecimals,
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
						CurrencyPair: types.CurrencyPair{
							Base:     "AA",
							Quote:    "BB",
							Decimals: types.DefaultDecimals,
						},
						Id: 0,
					},
					{
						CurrencyPair: types.CurrencyPair{
							Base:     "BB",
							Quote:    "CC",
							Decimals: types.DefaultDecimals,
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
		s.T().Run(tc.name, func(t *testing.T) {
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
						require.Nil(s.T(), err)
						checkQuotePriceEqual(s.T(), qp, *cpg.CurrencyPairPrice)
					} else {
						// assert that no price exists for the currency-pair
						require.NotNil(s.T(), err)
					}

					// get nonce, and check equality
					nonce, err := s.oracleKeeper.GetNonceForCurrencyPair(s.ctx, cpg.CurrencyPair)
					require.Nil(s.T(), err)

					// check equality of nonces
					require.Equal(s.T(), nonce, cpg.Nonce)

					// check equality of ids
					id, ok := s.oracleKeeper.GetIDForCurrencyPair(s.ctx, cpg.CurrencyPair)
					require.True(s.T(), ok)

					require.Equal(s.T(), id, cpg.Id)
				}
			}
		})
	}
}

func catchPanic(t *testing.T, k keeper.Keeper, ctx sdk.Context, gs types.GenesisState) {
	defer func() {
		err := recover()
		require.NotNil(t, err)
	}()
	// call init-genesis
	k.InitGenesis(ctx, gs)
}

func (s *KeeperTestSuite) TestExportGenesis() {
	s.T().Run("ExportGenesis with all valid QuotePrices", func(t *testing.T) {
		// insert multiple currency pairs
		cp1 := types.CurrencyPair{
			Base:     "AA",
			Quote:    "BB",
			Decimals: types.DefaultDecimals,
		}
		qp1 := types.QuotePrice{
			Price: sdkmath.NewInt(100),
		}
		cp2 := types.CurrencyPair{
			Base:     "CC",
			Quote:    "DD",
			Decimals: types.DefaultDecimals,
		}
		qp2 := types.QuotePrice{
			Price: sdkmath.NewInt(120),
		}

		// insert
		require.Nil(s.T(), s.oracleKeeper.CreateCurrencyPair(s.ctx, cp1))
		require.Nil(s.T(), s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp1, qp1))

		require.Nil(s.T(), s.oracleKeeper.CreateCurrencyPair(s.ctx, cp2))
		require.Nil(s.T(), s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp2, qp2))

		// insert
		require.Nil(s.T(), s.oracleKeeper.SetPriceForCurrencyPair(s.ctx, cp1, qp1))

		// export genesis
		gs := s.oracleKeeper.ExportGenesis(s.ctx)
		require.Equal(s.T(), len(gs.CurrencyPairGenesis), 2)
		expectedCurrencyPairs := map[string]types.QuotePrice{"AA/BB/8": qp1, "CC/DD/8": qp2}
		expectedNonces := map[string]uint64{"AA/BB/8": 2, "CC/DD/8": 1}

		for _, cpg := range gs.CurrencyPairGenesis {
			qp, ok := expectedCurrencyPairs[cpg.CurrencyPair.String()]
			require.True(s.T(), ok)
			// check equality for quote-prices
			checkQuotePriceEqual(s.T(), qp, *cpg.CurrencyPairPrice)
			// check equality of nonces
			nonce, ok := expectedNonces[cpg.CurrencyPair.String()]
			require.True(s.T(), ok)
			require.Equal(s.T(), nonce, cpg.Nonce)
		}
	})

	s.T().Run("ExportGenesis with some un-price-initialized CurrencyPairs", func(t *testing.T) {
		// initialize genesis w/ price-data
		gs := types.GenesisState{
			CurrencyPairGenesis: []types.CurrencyPairGenesis{
				{
					CurrencyPair: types.CurrencyPair{
						Base:     "AA",
						Quote:    "BB",
						Decimals: types.DefaultDecimals,
					},
					CurrencyPairPrice: &types.QuotePrice{
						Price: sdkmath.NewInt(100),
					},
					Nonce: 100,
					Id:    0,
				},
				{
					CurrencyPair: types.CurrencyPair{
						Base:     "CC",
						Quote:    "DD",
						Decimals: types.DefaultDecimals,
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

		// add un-initialized CurrencyPairs
		ms := keeper.NewMsgServer(s.oracleKeeper)
		_, err := ms.AddCurrencyPairs(s.ctx, &types.MsgAddCurrencyPairs{
			Authority: moduleAuthAddr.String(),
			CurrencyPairs: []types.CurrencyPair{
				{
					Base:     "EE",
					Quote:    "FF",
					Decimals: types.DefaultDecimals,
				},
				{
					Base:     "GG",
					Quote:    "HH",
					Decimals: types.DefaultDecimals,
				},
			},
		})
		require.Nil(s.T(), err)

		// setup expected values
		expectedCurrencyPairs := map[string]struct{}{"AA/BB/8": {}, "CC/DD/8": {}, "EE/FF/8": {}, "GG/HH/8": {}}
		expectedQuotePrices := map[string]types.QuotePrice{
			"AA/BB/8": {
				Price: sdkmath.NewInt(100),
			},
			"CC/DD/8": {
				Price: sdkmath.NewInt(101),
			},
		}
		expectedNonces := map[string]uint64{"AA/BB/8": 100, "CC/DD/8": 101}

		// ExportGenesis
		egs := s.oracleKeeper.ExportGenesis(s.ctx)

		// iterate over CurrencyPairGeneses in egs
		for _, cpg := range egs.CurrencyPairGenesis {
			// expect that all currency-pairs in gen-state are expected
			cps := cpg.CurrencyPair.String()
			_, ok := expectedCurrencyPairs[cps]
			require.True(s.T(), ok)

			// expect that if a CurrencyPrice exists, that it is expected
			if cpg.CurrencyPairPrice != nil {
				qp, ok := expectedQuotePrices[cps]
				require.True(s.T(), ok)

				// assert equality of QuotePrice
				checkQuotePriceEqual(s.T(), qp, *cpg.CurrencyPairPrice)

				nonce, ok := expectedNonces[cps]
				require.True(s.T(), ok)
				// assert equality of Nonce
				require.Equal(s.T(), cpg.Nonce, nonce)
			} else {
				require.Equal(s.T(), cpg.Nonce, uint64(0))
			}

			// check IDs
			id, ok := s.oracleKeeper.GetIDForCurrencyPair(s.ctx, cpg.CurrencyPair)
			require.True(s.T(), ok)
			require.Equal(s.T(), id, cpg.Id)
		}
	})
}
