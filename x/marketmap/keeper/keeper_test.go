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
	k := keeper.NewKeeper(ss, encCfg.Codec, s.authority)
	s.Require().NoError(k.SetLastUpdated(s.ctx))
	return k
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

	tickers = map[string]types.Ticker{
		btcusdt.String(): btcusdt,
		usdcusd.String(): usdcusd,
		usdtusd.String(): usdtusd,
		ethusdt.String(): ethusdt,
	}

	paths = map[string]types.Paths{
		btcusdt.String(): {Paths: btcusdt.Paths},
		usdcusd.String(): {Paths: usdcusd.Paths},
		usdtusd.String(): {Paths: usdtusd.Paths},
		ethusdt.String(): {Paths: ethusdt.Paths},
	}

	providers = map[string]types.Providers{
		btcusdt.String(): {Providers: btcusdt.Providers},
		usdcusd.String(): {Providers: usdcusd.Providers},
		usdtusd.String(): {Providers: usdtusd.Providers},
		ethusdt.String(): {Providers: ethusdt.Providers},
	}
)

func (s *KeeperTestSuite) TestGets() {
	s.Run("get no tickers", func() {
		got, err := s.keeper.GetAllTickers(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal([]types.Ticker(nil), got)
	})

	s.Run("setup initial markets", func() {
		for _, ticker := range tickers {
			s.Require().NoError(s.keeper.CreateMarket(s.ctx, ticker, types.Paths{Paths: ticker.Paths}, types.Providers{Providers: ticker.Providers}))
		}

		s.Run("unable to set markets again", func() {
			for _, ticker := range tickers {
				s.Require().ErrorIs(s.keeper.CreateTicker(s.ctx, ticker), types.NewTickerAlreadyExistsError(types.TickerString(ticker.String())))
			}
		})
	})

	s.Run("get all tickers", func() {
		got, err := s.keeper.GetAllTickersMap(s.ctx)
		s.Require().NoError(err)

		s.Require().Equal(len(tickers), len(got))
		s.Require().Equal(tickers, got)
	})

	s.Run("get all paths", func() {
		got, err := s.keeper.GetAllPathsMap(s.ctx)
		s.Require().NoError(err)

		s.Require().Equal(len(paths), len(got))
		s.Require().Equal(paths, got)
	})

	s.Run("get all providers", func() {
		got, err := s.keeper.GetAllProvidersMap(s.ctx)
		s.Require().NoError(err)

		s.Require().Equal(len(providers), len(got))
		s.Require().Equal(providers, got)
	})
}
