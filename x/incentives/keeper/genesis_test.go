package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/incentives/types"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/badprice"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/goodprice"
)

func (s *KeeperTestSuite) TestInitGenesis() {
	s.Run("can initialize genesis with no incentives", func() {
		genesis := types.NewDefaultGenesisState()
		s.incentivesKeeper.InitGenesis(s.ctx, *genesis)
	})

	s.Run("can initialize genesis with a single incentive", func() {
		validator := sdk.ValAddress([]byte("validator"))
		amount := math.NewInt(100)
		badPrice := badprice.NewBadPriceIncentive(validator, amount)

		// Create the genesis state.
		bz, err := badPrice.Marshal()
		s.Require().NoError(err)

		badPriceIncentives := types.NewIncentives(badprice.BadPriceIncentiveType, [][]byte{bz})
		genesis := types.NewGenesisState([]types.IncentivesByType{badPriceIncentives})

		// Initialize the genesis state.
		s.incentivesKeeper.InitGenesis(s.ctx, *genesis)

		// Check that the incentive was added to the store.
		incentives, err := s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		// Check that the incentive is the same as the one we added.
		i, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator.String(), i.Validator)
		s.Require().Equal(amount.String(), i.Amount)
	})

	s.Run("can initialize genesis with multiple incentives", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))

		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice1 := badprice.NewBadPriceIncentive(validator1, amount1)
		badPrice2 := badprice.NewBadPriceIncentive(validator2, amount2)

		// Create the genesis state.
		bz1, err := badPrice1.Marshal()
		s.Require().NoError(err)

		bz2, err := badPrice2.Marshal()
		s.Require().NoError(err)

		badPriceIncentives := types.NewIncentives(badprice.BadPriceIncentiveType, [][]byte{bz1, bz2})
		genesis := types.NewGenesisState([]types.IncentivesByType{badPriceIncentives})

		// Initialize the genesis state.
		s.incentivesKeeper.InitGenesis(s.ctx, *genesis)

		// Check that the incentives were added to the store.
		incentives, err := s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 2)

		// Check that the incentives are the same as the ones we added.
		i1, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator1.String(), i1.Validator)
		s.Require().Equal(amount1.String(), i1.Amount)

		i2, ok := incentives[1].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator2.String(), i2.Validator)
		s.Require().Equal(amount2.String(), i2.Amount)
	})

	s.Run("can initialize genesis with multiple incentive types", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))

		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice := badprice.NewBadPriceIncentive(validator1, amount1)
		goodPrice := goodprice.NewGoodPriceIncentive(validator2, amount2)

		// Create the genesis state.
		bz1, err := badPrice.Marshal()
		s.Require().NoError(err)

		bz2, err := goodPrice.Marshal()
		s.Require().NoError(err)

		badPriceIncentives := types.NewIncentives(badprice.BadPriceIncentiveType, [][]byte{bz1})
		goodPriceIncentives := types.NewIncentives(goodprice.GoodPriceIncentiveType, [][]byte{bz2})
		genesis := types.NewGenesisState([]types.IncentivesByType{badPriceIncentives, goodPriceIncentives})

		// Initialize the genesis state.
		s.incentivesKeeper.InitGenesis(s.ctx, *genesis)

		// Check that the incentives were added to the store.
		incentives, err := s.incentivesKeeper.GetIncentivesByType(s.ctx, &badprice.BadPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i1, ok := incentives[0].(*badprice.BadPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator1.String(), i1.Validator)
		s.Require().Equal(amount1.String(), i1.Amount)

		incentives, err = s.incentivesKeeper.GetIncentivesByType(s.ctx, &goodprice.GoodPriceIncentive{})
		s.Require().NoError(err)
		s.Require().Len(incentives, 1)

		i2, ok := incentives[0].(*goodprice.GoodPriceIncentive)
		s.Require().True(ok)
		s.Require().Equal(validator2.String(), i2.Validator)
		s.Require().Equal(amount2.String(), i2.Amount)
	})

	s.Run("errors when initializing genesis with unsupported incentive type", func() {
		unsupportedIncentive := types.NewIncentives("unsupported", [][]byte{[]byte("unsupported")})
		genesis := types.NewGenesisState([]types.IncentivesByType{unsupportedIncentive})

		// catch and check that initgenesis panics
		defer func() {
			r := recover()
			s.Require().NotNil(r)
		}()

		s.incentivesKeeper.InitGenesis(s.ctx, *genesis)
	})

	s.Run("errors when initializing genesis with invalid incentive", func() {
		validator := sdk.ValAddress([]byte("validator"))
		amount := math.NewInt(100)
		badPrice := badprice.NewBadPriceIncentive(validator, amount)

		// marshal the incentive to a byte slice
		bz, err := badPrice.Marshal()
		s.Require().NoError(err)

		// modify the byte slice to make it invalid
		bz[0] = 0x00

		// create the genesis state
		badPriceIncentives := types.NewIncentives(badprice.BadPriceIncentiveType, [][]byte{bz})
		genesis := types.NewGenesisState([]types.IncentivesByType{badPriceIncentives})

		// catch and check that initgenesis panics
		defer func() {
			r := recover()
			s.Require().NotNil(r)
		}()

		s.incentivesKeeper.InitGenesis(s.ctx, *genesis)
	})
}

func (s *KeeperTestSuite) TestExportGenesis() {
	s.Run("can export genesis with no incentives", func() {
		genesis := s.incentivesKeeper.ExportGenesis(s.ctx)
		s.Require().NotNil(genesis)

		// Check that the genesis state is valid.
		err := genesis.ValidateBasic()
		s.Require().NoError(err)

		// Check that the genesis state is empty.
		s.Require().Len(genesis.Registry, 0)
	})

	s.Run("can export genesis with a single incentive", func() {
		validator := sdk.ValAddress([]byte("validator"))
		amount := math.NewInt(100)
		badPrice := badprice.NewBadPriceIncentive(validator, amount)

		// Add the incentive to the store.
		err := s.incentivesKeeper.AddIncentives(s.ctx, []types.Incentive{badPrice})
		s.Require().NoError(err)

		// Export the genesis state.
		genesis := s.incentivesKeeper.ExportGenesis(s.ctx)
		s.Require().NotNil(genesis)

		// Check that the genesis state is valid.
		err = genesis.ValidateBasic()
		s.Require().NoError(err)

		// Check that the genesis state contains the incentive.
		s.Require().Len(genesis.Registry, 1)
		s.Require().Equal(badprice.BadPriceIncentiveType, genesis.Registry[0].IncentiveType)
		s.Require().Len(genesis.Registry[0].Entries, 1)

		bz, err := badPrice.Marshal()
		s.Require().NoError(err)
		s.Require().Equal(bz, genesis.Registry[0].Entries[0])
	})

	s.Run("can export genesis with multiple incentives", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))

		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice1 := badprice.NewBadPriceIncentive(validator1, amount1)
		badPrice2 := badprice.NewBadPriceIncentive(validator2, amount2)

		// Add the incentives to the store.
		err := s.incentivesKeeper.AddIncentives(s.ctx, []types.Incentive{badPrice1, badPrice2})
		s.Require().NoError(err)

		// Export the genesis state.
		genesis := s.incentivesKeeper.ExportGenesis(s.ctx)
		s.Require().NotNil(genesis)

		// Check that the genesis state is valid.
		err = genesis.ValidateBasic()
		s.Require().NoError(err)

		// Check that the genesis state contains the incentives.
		s.Require().Len(genesis.Registry, 1)
		s.Require().Equal(badprice.BadPriceIncentiveType, genesis.Registry[0].IncentiveType)
		s.Require().Len(genesis.Registry[0].Entries, 2)

		bz1, err := badPrice1.Marshal()
		s.Require().NoError(err)

		bz2, err := badPrice2.Marshal()
		s.Require().NoError(err)

		s.Require().Equal(bz1, genesis.Registry[0].Entries[0])
		s.Require().Equal(bz2, genesis.Registry[0].Entries[1])
	})

	s.Run("can export genesis with multiple incentive types", func() {
		validator1 := sdk.ValAddress([]byte("validator1"))
		validator2 := sdk.ValAddress([]byte("validator2"))

		amount1 := math.NewInt(100)
		amount2 := math.NewInt(200)

		badPrice := badprice.NewBadPriceIncentive(validator1, amount1)
		goodPrice := goodprice.NewGoodPriceIncentive(validator2, amount2)

		// Add the incentives to the store.
		err := s.incentivesKeeper.AddIncentives(s.ctx, []types.Incentive{badPrice, goodPrice})
		s.Require().NoError(err)

		// Export the genesis state.
		genesis := s.incentivesKeeper.ExportGenesis(s.ctx)
		s.Require().NotNil(genesis)

		// Check that the genesis state is valid.
		err = genesis.ValidateBasic()
		s.Require().NoError(err)

		// Check that the genesis state contains the incentives.
		s.Require().Len(genesis.Registry, 2)
		s.Require().Equal(badprice.BadPriceIncentiveType, genesis.Registry[0].IncentiveType)
		s.Require().Len(genesis.Registry[0].Entries, 1)
		s.Require().Equal(goodprice.GoodPriceIncentiveType, genesis.Registry[1].IncentiveType)
		s.Require().Len(genesis.Registry[1].Entries, 1)

		bz1, err := badPrice.Marshal()
		s.Require().NoError(err)

		bz2, err := goodPrice.Marshal()
		s.Require().NoError(err)

		s.Require().Equal(bz1, genesis.Registry[0].Entries[0])
		s.Require().Equal(bz2, genesis.Registry[1].Entries[0])
	})
}
