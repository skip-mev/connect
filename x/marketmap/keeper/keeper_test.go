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
	oraclekeeper "github.com/skip-mev/slinky/x/oracle/keeper"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	// Keeper variables
	authority    sdk.AccAddress
	keeper       *keeper.Keeper
	oracleKeeper oraclekeeper.Keeper

	hooks types.MarketMapHooks
}

func (s *KeeperTestSuite) initKeeper() *keeper.Keeper {
	mmKey := storetypes.NewKVStoreKey(types.StoreKey)
	oracleKey := storetypes.NewKVStoreKey(oracletypes.StoreKey)
	mmSS := runtime.NewKVStoreService(mmKey)
	oracleSS := runtime.NewKVStoreService(oracleKey)
	encCfg := moduletestutil.MakeTestEncodingConfig()

	keys := map[string]*storetypes.KVStoreKey{
		types.StoreKey:       mmKey,
		oracletypes.StoreKey: oracleKey,
	}

	transientKeys := map[string]*storetypes.TransientStoreKey{
		types.StoreKey:       storetypes.NewTransientStoreKey("transient_mm"),
		oracletypes.StoreKey: storetypes.NewTransientStoreKey("transient_oracle"),
	}

	s.authority = sdk.AccAddress("authority")
	s.ctx = testutil.DefaultContextWithKeys(keys, transientKeys, nil).WithBlockHeight(10)

	k := keeper.NewKeeper(mmSS, encCfg.Codec, s.authority)
	s.Require().NoError(k.SetLastUpdated(s.ctx, uint64(s.ctx.BlockHeight())))

	params := types.NewParams(s.authority.String(), 10)
	s.Require().NoError(k.SetParams(s.ctx, params))

	s.oracleKeeper = oraclekeeper.NewKeeper(oracleSS, encCfg.Codec, k, s.authority)
	s.hooks = types.MultiMarketMapHooks{
		s.oracleKeeper.Hooks(),
	}
	k.SetHooks(s.hooks)

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
			Base:      "BITCOIN",
			Quote:     "USDT",
			Delimiter: slinkytypes.DefaultDelimiter,
		},
		Decimals:         8,
		MinProviderCount: 1,
	}

	btcusdtPaths = types.Paths{
		Paths: []types.Path{
			{
				Operations: []types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:      "BITCOIN",
							Quote:     "USDT",
							Delimiter: slinkytypes.DefaultDelimiter,
						},
					},
				},
			},
		},
	}

	btcusdtProviders = types.Providers{
		Providers: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "btc-usdt",
			},
		},
	}

	usdtusd = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:      "USDT",
			Quote:     "USD",
			Delimiter: slinkytypes.DefaultDelimiter,
		},
		Decimals:         8,
		MinProviderCount: 1,
	}

	usdtusdPaths = types.Paths{
		Paths: []types.Path{
			{
				Operations: []types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:      "USDT",
							Quote:     "USD",
							Delimiter: slinkytypes.DefaultDelimiter,
						},
					},
				},
			},
		},
	}

	usdtusdProviders = types.Providers{
		Providers: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "usdt-usd",
			},
		},
	}

	usdcusd = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:      "USDC",
			Quote:     "USD",
			Delimiter: slinkytypes.DefaultDelimiter,
		},
		Decimals:         8,
		MinProviderCount: 1,
	}

	usdcusdPaths = types.Paths{
		Paths: []types.Path{
			{
				Operations: []types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:      "USDC",
							Quote:     "USD",
							Delimiter: slinkytypes.DefaultDelimiter,
						},
					},
				},
			},
		},
	}

	usdcusdProviders = types.Providers{
		Providers: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "usdc-usd",
			},
		},
	}

	ethusdt = types.Ticker{
		CurrencyPair: slinkytypes.CurrencyPair{
			Base:      "ETHEREUM",
			Quote:     "USDT",
			Delimiter: slinkytypes.DefaultDelimiter,
		},
		Decimals:         8,
		MinProviderCount: 1,
	}

	ethusdtPaths = types.Paths{
		Paths: []types.Path{
			{
				Operations: []types.Operation{
					{
						CurrencyPair: slinkytypes.CurrencyPair{
							Base:      "ETHEREUM",
							Quote:     "USDT",
							Delimiter: slinkytypes.DefaultDelimiter,
						},
					},
				},
			},
		},
	}

	ethusdtProviders = types.Providers{
		Providers: []types.ProviderConfig{
			{
				Name:           "kucoin",
				OffChainTicker: "eth-usdt",
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
		btcusdt.String(): btcusdtPaths,
		usdcusd.String(): usdcusdPaths,
		usdtusd.String(): usdtusdPaths,
		ethusdt.String(): ethusdtPaths,
	}

	providers = map[string]types.Providers{
		btcusdt.String(): btcusdtProviders,
		usdcusd.String(): usdcusdProviders,
		usdtusd.String(): usdtusdProviders,
		ethusdt.String(): ethusdtProviders,
	}

	markets = struct {
		tickers   map[string]types.Ticker
		paths     map[string]types.Paths
		providers map[string]types.Providers
	}{
		tickers:   tickers,
		paths:     paths,
		providers: providers,
	}
)

func (s *KeeperTestSuite) TestGets() {
	s.Run("get no tickers", func() {
		got, err := s.keeper.GetAllTickers(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal([]types.Ticker(nil), got)
	})

	s.Run("setup initial markets", func() {
		for _, ticker := range markets.tickers {
			marketPaths, ok := markets.paths[ticker.String()]
			s.Require().True(ok)
			marketProviders, ok := markets.providers[ticker.String()]
			s.Require().True(ok)
			s.Require().NoError(s.keeper.CreateMarket(s.ctx, ticker, marketPaths, marketProviders))
		}

		s.Run("unable to set markets again", func() {
			for _, ticker := range markets.tickers {
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

func (s *KeeperTestSuite) TestSetParams() {
	params := types.DefaultParams()

	s.Run("can set and get params", func() {
		err := s.keeper.SetParams(s.ctx, params)
		s.Require().NoError(err)

		params2, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(params, params2)
	})
}
