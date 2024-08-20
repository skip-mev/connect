package keeper_test

import (
	"testing"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/connect/v2/x/incentives/keeper"
	"github.com/skip-mev/connect/v2/x/incentives/types"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/badprice"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/goodprice"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/mocks"
)

type KeeperTestSuite struct {
	suite.Suite

	incentivesKeeper keeper.Keeper
	queryServer      keeper.QueryServer
	key              storetypes.StoreKey
	ctx              sdk.Context

	// mock strategies
	stakingKeeper mocks.StakingKeeper
	bankKeeper    mocks.BankKeeper
}

func (s *KeeperTestSuite) SetupTest() {
	s.key = storetypes.NewKVStoreKey(types.StoreKey)
	s.ctx = testutil.DefaultContext(s.key, storetypes.NewTransientStoreKey("transient_key"))
	s.incentivesKeeper = keeper.NewKeeper(s.key, nil)
	s.queryServer = keeper.NewQueryServer(s.incentivesKeeper)
}

func (s *KeeperTestSuite) SetupSubTest() {
	err := s.incentivesKeeper.RemoveIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
	s.Require().NoError(err)

	incentives, err := s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
	s.Require().NoError(err)
	s.Require().Len(incentives, 0)

	err = s.incentivesKeeper.RemoveIncentivesByType(s.ctx, &goodprice.GoodPriceIncentive{})
	s.Require().NoError(err)

	incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &goodprice.GoodPriceIncentive{})
	s.Require().NoError(err)
	s.Require().Len(incentives, 0)

	// Reset the mock strategies.
	s.bankKeeper = *mocks.NewBankKeeper(s.T())
	s.stakingKeeper = *mocks.NewStakingKeeper(s.T())

	badPriceStrategy := badprice.NewBadPriceIncentiveStrategy(&s.stakingKeeper).GetStrategy()
	goodPriceStrategy := goodprice.NewGoodPriceIncentiveStrategy(&s.bankKeeper).GetStrategy()
	strategies := map[types.Incentive]types.Strategy{
		&badprice.BadPriceIncentive{}:   badPriceStrategy,
		&goodprice.GoodPriceIncentive{}: goodPriceStrategy,
	}

	// Reset the keeper with the new strategies.
	s.incentivesKeeper = keeper.NewKeeper(s.key, strategies)
	s.queryServer = keeper.NewQueryServer(s.incentivesKeeper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
