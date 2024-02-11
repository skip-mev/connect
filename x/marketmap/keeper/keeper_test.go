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
	"github.com/skip-mev/slinky/x/marketmap/types/mocks"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	// Keeper variables
	authority    sdk.AccAddress
	oracleKeeper *mocks.OracleKeeper
	keeper       keeper.Keeper
}

func (s *KeeperTestSuite) initKeeper() keeper.Keeper {
	key := storetypes.NewKVStoreKey(types.StoreKey)
	ss := runtime.NewKVStoreService(key)
	encCfg := moduletestutil.MakeTestEncodingConfig()
	s.oracleKeeper = mocks.NewOracleKeeper(s.T())
	s.authority = sdk.AccAddress("authority")
	s.ctx = testutil.DefaultContext(key, storetypes.NewTransientStoreKey("transient_key"))
	return keeper.NewKeeper(ss, encCfg.Codec, s.oracleKeeper, s.authority)
}

func (s *KeeperTestSuite) SetupTest() {
	s.keeper = s.initKeeper()
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) TestMarketConfigs() {
	btcEthTickerConfig := types.TickerConfig{
		Ticker: types.Ticker{
			Base:             "BTC",
			Quote:            "ETH",
			Decimals:         8,
			MinProviderCount: 1,
		},
		OffChainTicker: "BTC-ETH",
	}
	atomUsdcTickerConfig := types.TickerConfig{
		Ticker: types.Ticker{
			Base:             "BTC",
			Quote:            "ETH",
			Decimals:         8,
			MinProviderCount: 1,
		},
		OffChainTicker: "BTC-ETH",
	}
	marketCfg1 := types.MarketConfig{
		Name: "provider1",
		TickerConfigs: map[string]types.TickerConfig{
			"BTC/ETH": btcEthTickerConfig,
		},
	}
	marketCfg2 := types.MarketConfig{
		Name: "provider2",
		TickerConfigs: map[string]types.TickerConfig{
			"BTC/ETH":   btcEthTickerConfig,
			"ATOM/USDC": atomUsdcTickerConfig,
		},
	}
	s.Run("creating valid market configs passes", func() {
		s.Require().NoError(s.keeper.CreateMarketConfig(s.ctx, marketCfg1))
		s.Require().NoError(s.keeper.CreateMarketConfig(s.ctx, marketCfg2))
	})
	s.Run("creating market config for existing provider fails", func() {
		s.Require().ErrorIs(s.keeper.CreateMarketConfig(s.ctx, marketCfg1), types.NewMarketConfigAlreadyExistsError(marketCfg1.Name))
	})
	s.Run("fetching all market configs returns all initialized market configs", func() {
		marketCfgs, err := s.keeper.GetAllMarketConfigs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(2, len(marketCfgs))
	})
}

func (s *KeeperTestSuite) TestAggregationConfigs() {
	cp1 := slinkytypes.CurrencyPair{Base: "BTC", Quote: "ETH"}
	aggCfg1 := types.PathsConfig{
		Ticker: types.Ticker{
			Base:             cp1.Base,
			Quote:            cp1.Quote,
			Decimals:         0,
			MinProviderCount: 0,
		},
		Paths: []types.Path{
			{Operations: []types.Operation{{Ticker: types.Ticker{Base: cp1.Base, Quote: cp1.Quote}}}},
		},
	}
	cp2 := slinkytypes.CurrencyPair{Base: "ATOM", Quote: "USDC"}
	aggCfg2 := types.PathsConfig{
		Ticker: types.Ticker{
			Base:             cp2.Base,
			Quote:            cp2.Quote,
			Decimals:         0,
			MinProviderCount: 0,
		},
		Paths: []types.Path{
			{Operations: []types.Operation{{Ticker: types.Ticker{Base: cp2.Base, Quote: cp2.Quote}}}},
		},
	}
	s.Run("creating valid agg configs passes", func() {
		s.Require().NoError(s.keeper.CreateAggregationConfig(s.ctx, aggCfg1))
		s.Require().NoError(s.keeper.CreateAggregationConfig(s.ctx, aggCfg2))
	})
	s.Run("creating agg config for existing ticker fails", func() {
		s.Require().ErrorIs(s.keeper.CreateAggregationConfig(s.ctx, aggCfg1), types.NewAggregationConfigAlreadyExistsError(cp1.String()))
	})
	s.Run("fetching all agg configs returns all initialized agg configs", func() {
		aggCfgs, err := s.keeper.GetAllAggregationConfigs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(2, len(aggCfgs))
	})
}

func (s *KeeperTestSuite) TestMarketMap() {
	cp1 := slinkytypes.CurrencyPair{Base: "BTC", Quote: "ETH"}
	aggCfg1 := types.PathsConfig{
		Ticker: types.Ticker{
			Base:             cp1.Base,
			Quote:            cp1.Quote,
			Decimals:         0,
			MinProviderCount: 0,
		},
		Paths: []types.Path{
			{Operations: []types.Operation{{Ticker: types.Ticker{Base: cp1.Base, Quote: cp1.Quote}}}},
		},
	}
	btcEthTickerConfig := types.TickerConfig{
		Ticker: types.Ticker{
			Base:             "BTC",
			Quote:            "ETH",
			Decimals:         8,
			MinProviderCount: 1,
		},
		OffChainTicker: "BTC-ETH",
	}
	marketCfg1 := types.MarketConfig{
		Name: "provider1",
		TickerConfigs: map[string]types.TickerConfig{
			"BTC/ETH": btcEthTickerConfig,
		},
	}
	s.Run("market map returns the full set of data in the module", func() {
		s.Require().NoError(s.keeper.CreateMarketConfig(s.ctx, marketCfg1))
		s.Require().NoError(s.keeper.CreateAggregationConfig(s.ctx, aggCfg1))
		marketMap, err := s.keeper.GetMarketMap(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(1, len(marketMap.MarketConfigs))
		s.Require().Equal(1, len(marketMap.TickerConfigs))
		marketCfg, ok := marketMap.MarketConfigs[marketCfg1.Name]
		s.Require().True(ok)
		s.Require().Equal(marketCfg1.String(), marketCfg.String())
		aggCfg, ok := marketMap.TickerConfigs[cp1.String()]
		s.Require().True(ok)
		s.Require().Equal(aggCfg1.String(), aggCfg.String())
	})
}
