package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/suite"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	"github.com/skip-mev/slinky/x/marketmap/keeper"
	"github.com/skip-mev/slinky/x/marketmap/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	// Keeper variables
	authority sdk.AccAddress
	keeper    keeper.Keeper
}

func (s *KeeperTestSuite) initKeeper() keeper.Keeper {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	ss := runtime.NewKVStoreService(key)
	encCfg := moduletestutil.MakeTestEncodingConfig()
	s.authority = sdk.AccAddress("authority")
	s.ctx = testutil.DefaultContext(key, storetypes.NewTransientStoreKey("transient_key")).WithBlockHeight(10)
	return keeper.NewKeeper(ss, encCfg.Codec, s.authority)
}

func (s *KeeperTestSuite) SetupTest() {
	s.keeper = s.initKeeper()
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

var (
	btcusdt = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "BITCOIN",
			Quote: "USDT",
		},
		Decimals:         8,
		MinProviderCount: 1,
		Paths: []types.Path{
			{
				Operations: []types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:  "BITCOIN",
							Quote: "USDT",
						},
					},
				},
			},
		},
	}

	usdtusd = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "USDT",
			Quote: "USD",
		},
		Decimals:         8,
		MinProviderCount: 1,
		Paths: []types.Path{
			{
				Operations: []types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:  "USDT",
							Quote: "USD",
						},
					},
				},
			},
		},
	}

	usdcusd = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "USDC",
			Quote: "USD",
		},
		Decimals:         8,
		MinProviderCount: 1,
		Paths: []types.Path{
			{
				Operations: []types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:  "USDC",
							Quote: "USD",
						},
					},
				},
			},
		},
	}

	ethusdt = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:  "ETHEREUM",
			Quote: "USDT",
		},
		Decimals:         8,
		MinProviderCount: 1,
		Paths: []types.Path{
			{
				Operations: []types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:  "ETHEREUM",
							Quote: "USDT",
						},
					},
				},
			},
		},
	}

	tickers = []types.Ticker{
		btcusdt,
		usdcusd,
		usdtusd,
		ethusdt,
	}
)

func (s *KeeperTestSuite) TestTickers() {
	s.Run("get no tickers", func() {
		got, err := s.keeper.GetAllTickers(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal([]types.Ticker(nil), got)
	})

	s.Run("setup initial markets", func() {
		for _, ticker := range tickers {
			s.Require().NoError(s.keeper.CreateTicker(s.ctx, ticker))
		}

		s.Run("unable to set markets again", func() {
			for _, ticker := range tickers {
				s.Require().ErrorIs(s.keeper.CreateTicker(s.ctx, ticker), types.NewTickerAlreadyExistsError(types.TickerString(ticker.String())))
			}
		})
	})

	s.Run("get all tickers", func() {
		got, err := s.keeper.GetAllTickers(s.ctx)
		s.Require().NoError(err)

		s.Require().Equal(len(tickers), len(got))
		s.Require().True(unorderedEqual(tickers, got))
	})
}

func unorderedEqual(first, second []types.Ticker) bool {
	if len(first) != len(second) {
		return false
	}
	exists := make(map[string]bool)
	for _, value := range first {
		exists[value.String()] = true
	}
	for _, value := range second {
		if !exists[value.String()] {
			return false
		}
	}
	return true
}
