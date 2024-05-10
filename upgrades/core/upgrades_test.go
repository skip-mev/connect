package core_test

import (
	"fmt"
	"os"
	"testing"

	"cosmossdk.io/log"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/slinky/tests/simapp"
	"github.com/skip-mev/slinky/upgrades"
	"github.com/skip-mev/slinky/upgrades/core"
	marketmaptypes "github.com/skip-mev/slinky/x/marketmap/types"
)

type UpgradeTestSuite struct {
	suite.Suite

	ctx sdk.Context
	app *simapp.SimApp
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}

func (s *UpgradeTestSuite) SetupTest() {
	dir, err := os.MkdirTemp("", "simapp")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)

	s.app = simapp.NewSimApp(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(dir))

	s.ctx = s.app.NewContext(true)
}

func (s *UpgradeTestSuite) TestSlinkyCoreUpgrade() {
	app := s.app
	ctx := s.ctx

	markets, err := marketmaptypes.ReadMarketsFromFile("markets.json")
	s.Require().NoError(err)
	marketMap := markets.ToMarketMap()

	initializeUpgrade := core.NewDefaultInitializeUpgrade(marketmaptypes.DefaultParams())

	app.UpgradeKeeper.SetUpgradeHandler(core.InitializationName, initializeUpgrade.CreateUpgradeHandler(
		app.ModuleManager,
		app.Configurator(),
		app.OracleKeeper,
		app.MarketMapKeeper,
		app.AppCodec(),
		upgrades.EmptyUpgrade,
	))

	upgrade := upgradetypes.Plan{
		Name:   core.InitializationName,
		Info:   "some text here",
		Height: 100,
	}
	require.NoError(s.T(), app.UpgradeKeeper.ApplyUpgrade(ctx, upgrade))

	params, err := app.MarketMapKeeper.GetParams(ctx)
	s.Require().NoError(err)
	s.Require().Equal(params.MarketAuthorities[0], "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn")
	s.Require().Equal(params.Admin, "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn")

	// check that the market map was properly set
	mm, err := app.MarketMapKeeper.GetAllMarkets(ctx)
	gotMM := marketmaptypes.MarketMap{Markets: mm}
	s.Require().NoError(err)
	s.Require().True(marketMap.Equal(gotMM))

	numCps, err := app.OracleKeeper.GetNumCurrencyPairs(ctx)
	s.Require().NoError(err)
	s.Require().Equal(numCps, uint64(len(markets)))

	// check that all currency pairs have been initialized in the oracle module
	for _, market := range markets {
		decimals, err := app.OracleKeeper.GetDecimalsForCurrencyPair(ctx, market.Ticker.CurrencyPair)
		s.Require().NoError(err)
		s.Require().Equal(market.Ticker.Decimals, decimals)

		price, err := app.OracleKeeper.GetPriceWithNonceForCurrencyPair(ctx, market.Ticker.CurrencyPair)
		s.Require().NoError(err)
		s.Require().Equal(uint64(0), price.Nonce())     // no nonce because no updates yet
		s.Require().Equal(uint64(0), price.BlockHeight) // no block height because no price written yet

		s.Require().True(market.Ticker.Enabled)
	}
}
