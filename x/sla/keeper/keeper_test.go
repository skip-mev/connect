package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/suite"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/sla/keeper"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
	"github.com/skip-mev/connect/v2/x/sla/types/mocks"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	// Keeper variables
	authority      sdk.AccAddress
	stakingKeeper  *mocks.StakingKeeper
	slashingKeeper *mocks.SlashingKeeper
	keeper         *keeper.Keeper

	// Message server variables
	msgServer slatypes.MsgServer

	// Query server variables
	queryServer slatypes.QueryServer
}

func (s *KeeperTestSuite) SetupTest() {
	s.keeper = s.initKeeper()
	s.msgServer = keeper.NewMsgServer(*s.keeper)
	s.queryServer = keeper.NewQueryServer(*s.keeper)
}

func (s *KeeperTestSuite) SetupSubTest() {
	s.keeper = s.initKeeper()
	s.msgServer = keeper.NewMsgServer(*s.keeper)
	s.queryServer = keeper.NewQueryServer(*s.keeper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) TestSetParams() {
	params := slatypes.DefaultParams()

	s.Run("can set and get params", func() {
		err := s.keeper.SetParams(s.ctx, params)
		s.Require().NoError(err)

		params2, err := s.keeper.GetParams(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(params, params2)
	})
}

func (s *KeeperTestSuite) TestSetCurrencyPairs() {
	cp1 := slinkytypes.NewCurrencyPair("btc", "usd")
	cp2 := slinkytypes.NewCurrencyPair("eth", "usd")

	testCPs := map[slinkytypes.CurrencyPair]struct{}{
		cp1: {},
		cp2: {},
	}

	s.Run("can set and get currency pairs", func() {
		err := s.keeper.SetCurrencyPairs(s.ctx, testCPs)
		s.Require().NoError(err)

		cps, err := s.keeper.GetCurrencyPairs(s.ctx)
		s.Require().NoError(err)
		s.Require().Equal(testCPs, cps)
	})
}

func (s *KeeperTestSuite) initKeeper() *keeper.Keeper {
	// Set up context
	key := storetypes.NewKVStoreKey(slatypes.StoreKey)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx

	// Set up store and encoding configs
	storeService := runtime.NewKVStoreService(key)
	encodingConfig := moduletestutil.MakeTestEncodingConfig()

	slatypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	s.stakingKeeper = mocks.NewStakingKeeper(s.T())
	s.slashingKeeper = mocks.NewSlashingKeeper(s.T())
	s.authority = sdk.AccAddress("authority")

	// Set up keeper
	k := keeper.NewKeeper(
		storeService,
		encodingConfig.Codec,
		s.authority,
		s.stakingKeeper,
		s.slashingKeeper,
	)

	s.Require().NoError(k.SetParams(s.ctx, slatypes.DefaultParams()))

	return k
}
