package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	slatypes "github.com/skip-mev/connect/v2/x/sla/types"
)

func (s *KeeperTestSuite) TestExecSLA() {
	// Strategy inputs
	id := "testID"
	maximumViableWindow := 20
	expectedUptime := math.LegacyMustNewDecFromStr("0.8")
	k := math.LegacyMustNewDecFromStr("1.0")
	frequency := 10
	minimumBlockUpdates := 10

	sla := slatypes.NewPriceFeedSLA(
		id,
		uint64(maximumViableWindow),
		expectedUptime,
		k,
		uint64(minimumBlockUpdates),
		uint64(frequency),
	)

	// Price feed parameters
	validator := sdk.ConsAddress([]byte("validator"))
	cp := slinkytypes.CurrencyPair{
		Base:  "BTC",
		Quote: "ETH",
	}

	s.Run("returns incentive when sla should not be checked", func() {
		s.ctx = s.ctx.WithBlockHeight(1)

		priceFeed, err := slatypes.NewPriceFeed(uint(maximumViableWindow), validator, cp, id)
		s.Require().NoError(err)
		err = s.keeper.SetPriceFeed(s.ctx, priceFeed)
		s.Require().NoError(err)

		s.keeper.SetSLA(s.ctx, sla)
		err = s.keeper.ExecSLA(s.ctx, sla)
		s.Require().NoError(err)
	})

	s.Run("returns when the validator does not exist", func() {
		s.ctx = s.ctx.WithBlockHeight(int64(frequency))

		priceFeed, err := slatypes.NewPriceFeed(uint(maximumViableWindow), validator, cp, id)
		s.Require().NoError(err)
		for i := 0; i < 20; i++ {
			priceFeed.SetUpdate(slatypes.VoteWithPrice)
		}
		err = s.keeper.SetPriceFeed(s.ctx, priceFeed)
		s.Require().NoError(err)

		s.stakingKeeper.On("GetLastValidatorPower", s.ctx, sdk.ValAddress(validator.Bytes())).Return(int64(0), fmt.Errorf("validator does not exist"))

		s.keeper.SetSLA(s.ctx, sla)
		err = s.keeper.ExecSLA(s.ctx, sla)
		s.Require().NoError(err)

		// Check that the price feed was deleted
		contains, err := s.keeper.ContainsPriceFeed(s.ctx, id, cp, validator)
		s.Require().NoError(err)
		s.Require().False(contains)
	})

	s.Run("validator did not vote on enough blocks to be considered", func() {
		s.ctx = s.ctx.WithBlockHeight(int64(frequency))

		priceFeed, err := slatypes.NewPriceFeed(uint(maximumViableWindow), validator, cp, id)
		s.Require().NoError(err)
		err = s.keeper.SetPriceFeed(s.ctx, priceFeed)
		s.Require().NoError(err)

		s.keeper.SetSLA(s.ctx, sla)
		err = s.keeper.ExecSLA(s.ctx, sla)
		s.Require().NoError(err)
	})

	s.Run("price feed with same size window as sla has not met SLA (0 uptime)", func() {
		s.ctx = s.ctx.WithBlockHeight(int64(frequency))

		priceFeed, err := slatypes.NewPriceFeed(uint(maximumViableWindow), validator, cp, id)
		s.Require().NoError(err)
		for i := 0; i < 20; i++ {
			priceFeed.SetUpdate(slatypes.VoteWithoutPrice)
		}
		err = s.keeper.SetPriceFeed(s.ctx, priceFeed)
		s.Require().NoError(err)

		s.stakingKeeper.On("GetLastValidatorPower", s.ctx, sdk.ValAddress(validator)).Return(int64(100), nil)

		expectedSlashFactor := k
		s.slashingKeeper.On(
			"Slash",
			s.ctx,
			validator,
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			int64(100),
			expectedSlashFactor,
		).Return(math.NewInt(10), nil)

		s.keeper.SetSLA(s.ctx, sla)
		err = s.keeper.ExecSLA(s.ctx, sla)
		s.Require().NoError(err)
	})

	s.Run("price feed with same size window as sla and has met SLA", func() {
		s.ctx = s.ctx.WithBlockHeight(int64(frequency))

		priceFeed, err := slatypes.NewPriceFeed(uint(maximumViableWindow), validator, cp, id)
		s.Require().NoError(err)
		for i := 0; i < 16; i++ {
			priceFeed.SetUpdate(slatypes.VoteWithPrice)
		}
		for i := 0; i < 4; i++ {
			priceFeed.SetUpdate(slatypes.VoteWithoutPrice)
		}
		err = s.keeper.SetPriceFeed(s.ctx, priceFeed)
		s.Require().NoError(err)

		s.stakingKeeper.On("GetLastValidatorPower", mock.Anything, sdk.ValAddress(validator)).Return(int64(100), nil)

		s.keeper.SetSLA(s.ctx, sla)
		err = s.keeper.ExecSLA(s.ctx, sla)
		s.Require().NoError(err)
	})

	s.Run("price feed with same size window as sla and has exceeded SLA", func() {
		s.ctx = s.ctx.WithBlockHeight(int64(frequency))

		priceFeed, err := slatypes.NewPriceFeed(uint(maximumViableWindow), validator, cp, id)
		s.Require().NoError(err)
		for i := 0; i < 20; i++ {
			priceFeed.SetUpdate(slatypes.VoteWithPrice)
		}
		err = s.keeper.SetPriceFeed(s.ctx, priceFeed)
		s.Require().NoError(err)

		s.stakingKeeper.On("GetLastValidatorPower", mock.Anything, sdk.ValAddress(validator)).Return(int64(100), nil)

		s.keeper.SetSLA(s.ctx, sla)
		err = s.keeper.ExecSLA(s.ctx, sla)
		s.Require().NoError(err)
	})

	s.Run("price feed with same size window as sla has not met SLA (half uptime)", func() {
		s.ctx = s.ctx.WithBlockHeight(int64(frequency))

		priceFeed, err := slatypes.NewPriceFeed(uint(maximumViableWindow), validator, cp, id)
		s.Require().NoError(err)
		for i := 0; i < 5; i++ {
			priceFeed.SetUpdate(slatypes.VoteWithoutPrice)
		}
		for i := 0; i < 5; i++ {
			priceFeed.SetUpdate(slatypes.VoteWithPrice)
		}
		err = s.keeper.SetPriceFeed(s.ctx, priceFeed)
		s.Require().NoError(err)

		s.stakingKeeper.On("GetLastValidatorPower", mock.Anything, sdk.ValAddress(validator)).Return(int64(100), nil)

		expectedDevation := (expectedUptime.Sub(math.LegacyMustNewDecFromStr("0.5"))).Quo(expectedUptime)
		expectedSlashFactor := k.Mul(expectedDevation)
		s.slashingKeeper.On(
			"Slash",
			mock.Anything,
			validator,
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			int64(100),
			expectedSlashFactor,
		).Return(math.NewInt(10), nil)

		s.keeper.SetSLA(s.ctx, sla)
		err = s.keeper.ExecSLA(s.ctx, sla)
		s.Require().NoError(err)
	})

	s.Run("price feed with same size window as sla has not met SLA (half uptime)", func() {
		s.ctx = s.ctx.WithBlockHeight(int64(frequency))

		priceFeed, err := slatypes.NewPriceFeed(uint(maximumViableWindow), validator, cp, id)
		s.Require().NoError(err)
		for i := 0; i < 10; i++ {
			priceFeed.SetUpdate(slatypes.VoteWithoutPrice)
		}
		for i := 0; i < 10; i++ {
			priceFeed.SetUpdate(slatypes.VoteWithPrice)
		}
		err = s.keeper.SetPriceFeed(s.ctx, priceFeed)
		s.Require().NoError(err)

		s.stakingKeeper.On("GetLastValidatorPower", mock.Anything, sdk.ValAddress(validator)).Return(int64(100), nil)

		expectedDevation := (expectedUptime.Sub(math.LegacyMustNewDecFromStr("0.5"))).Quo(expectedUptime)
		expectedSlashFactor := k.Mul(expectedDevation)
		s.slashingKeeper.On(
			"Slash",
			mock.Anything,
			validator,
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			int64(100),
			expectedSlashFactor,
		).Return(math.NewInt(10), nil)

		s.keeper.SetSLA(s.ctx, sla)
		err = s.keeper.ExecSLA(s.ctx, sla)
		s.Require().NoError(err)
	})
}

func (s *KeeperTestSuite) TestEnforceSLA() {
	id := "testID"
	expectedUptime := math.LegacyMustNewDecFromStr("0.8")
	slashConstant := math.LegacyMustNewDecFromStr("0.25")
	sla := slatypes.NewPriceFeedSLA(
		id,
		uint64(20),
		expectedUptime,
		slashConstant,
		uint64(10),
		uint64(10),
	)

	consAddress := sdk.ConsAddress([]byte("validator"))
	cp := slinkytypes.NewCurrencyPair("mog", "usd")
	feed, err := slatypes.NewPriceFeed(
		uint(20),
		consAddress,
		cp,
		id,
	)
	s.Require().NoError(err)

	s.Run("returns when the validator does not exist", func() {
		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(0), fmt.Errorf("validator does not exist"))

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)
	})

	s.Run("will remove price feed when validator does not exist", func() {
		err := s.keeper.SetPriceFeed(s.ctx, feed)
		s.Require().NoError(err)

		contains, err := s.keeper.ContainsPriceFeed(s.ctx, id, cp, consAddress)
		s.Require().NoError(err)
		s.Require().True(contains)

		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(0), fmt.Errorf("validator does not exist"))

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)

		contains, err = s.keeper.ContainsPriceFeed(s.ctx, id, cp, consAddress)
		s.Require().NoError(err)
		s.Require().False(contains)
	})

	s.Run("does not slash on price feed with no updates", func() {
		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(100), nil)

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)
	})

	s.Run("100% uptime with minimum number of blocks", func() {
		for i := 0; i < 10; i++ {
			feed.SetUpdate(slatypes.VoteWithPrice)
		}

		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(100), nil)

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)
	})

	s.Run("100% uptime with maximum number of blocks", func() {
		for i := 0; i < 20; i++ {
			feed.SetUpdate(slatypes.VoteWithPrice)
		}

		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(100), nil)

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)
	})

	s.Run("50% uptime with minimum number of blocks", func() {
		feed, err := slatypes.NewPriceFeed(
			uint(20),
			consAddress,
			cp,
			id,
		)
		s.Require().NoError(err)

		for i := 0; i < 5; i++ {
			feed.SetUpdate(slatypes.VoteWithPrice)
		}

		for i := 0; i < 5; i++ {
			feed.SetUpdate(slatypes.VoteWithoutPrice)
		}

		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(100), nil)

		expectedDevation := (expectedUptime.Sub(math.LegacyMustNewDecFromStr("0.5"))).Quo(expectedUptime)
		expectedSlashFactor := slashConstant.Mul(expectedDevation)

		s.slashingKeeper.On(
			"Slash",
			mock.Anything,
			consAddress,
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			int64(100),
			expectedSlashFactor,
		).Return(math.NewInt(10), nil)

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)
	})

	s.Run("50% uptime with maximum number of blocks", func() {
		feed, err := slatypes.NewPriceFeed(
			uint(20),
			consAddress,
			cp,
			id,
		)
		s.Require().NoError(err)

		for i := 0; i < 10; i++ {
			feed.SetUpdate(slatypes.VoteWithPrice)
		}

		for i := 0; i < 10; i++ {
			feed.SetUpdate(slatypes.VoteWithoutPrice)
		}

		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(100), nil)

		expectedDevation := (expectedUptime.Sub(math.LegacyMustNewDecFromStr("0.5"))).Quo(expectedUptime)
		expectedSlashFactor := slashConstant.Mul(expectedDevation)

		s.slashingKeeper.On(
			"Slash",
			mock.Anything,
			consAddress,
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			int64(100),
			expectedSlashFactor,
		).Return(math.NewInt(10), nil)

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)
	})

	s.Run("0% uptime with minimum number of blocks", func() {
		feed, err := slatypes.NewPriceFeed(
			uint(20),
			consAddress,
			cp,
			id,
		)
		s.Require().NoError(err)

		for i := 0; i < 10; i++ {
			feed.SetUpdate(slatypes.VoteWithoutPrice)
		}

		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(100), nil)

		expectedSlashFactor := slashConstant

		s.slashingKeeper.On(
			"Slash",
			mock.Anything,
			consAddress,
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			int64(100),
			expectedSlashFactor,
		).Return(math.NewInt(10), nil)

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)
	})

	s.Run("0% uptime with maximum number of blocks", func() {
		feed, err := slatypes.NewPriceFeed(
			uint(20),
			consAddress,
			cp,
			id,
		)
		s.Require().NoError(err)

		for i := 0; i < 20; i++ {
			feed.SetUpdate(slatypes.VoteWithoutPrice)
		}

		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(100), nil)

		expectedSlashFactor := slashConstant

		s.slashingKeeper.On(
			"Slash",
			mock.Anything,
			consAddress,
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			int64(100),
			expectedSlashFactor,
		).Return(math.NewInt(10), nil)

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)
	})

	s.Run("75% uptime", func() {
		feed, err := slatypes.NewPriceFeed(
			uint(20),
			consAddress,
			cp,
			id,
		)
		s.Require().NoError(err)

		for i := 0; i < 15; i++ {
			feed.SetUpdate(slatypes.VoteWithPrice)
		}

		for i := 0; i < 5; i++ {
			feed.SetUpdate(slatypes.VoteWithoutPrice)
		}

		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(100), nil)

		expectedDevation := (expectedUptime.Sub(math.LegacyMustNewDecFromStr("0.75"))).Quo(expectedUptime)
		expectedSlashFactor := slashConstant.Mul(expectedDevation)

		s.slashingKeeper.On(
			"Slash",
			mock.Anything,
			consAddress,
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			int64(100),
			expectedSlashFactor,
		).Return(math.NewInt(10), nil)

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)
	})

	s.Run("75% uptime with wrap around", func() {
		feed, err := slatypes.NewPriceFeed(
			uint(20),
			consAddress,
			cp,
			id,
		)
		s.Require().NoError(err)

		for i := 0; i < 20; i++ {
			feed.SetUpdate(slatypes.VoteWithPrice)
		}

		for i := 0; i < 5; i++ {
			feed.SetUpdate(slatypes.VoteWithoutPrice)
		}

		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(100), nil)

		expectedDevation := (expectedUptime.Sub(math.LegacyMustNewDecFromStr("0.75"))).Quo(expectedUptime)
		expectedSlashFactor := slashConstant.Mul(expectedDevation)

		s.slashingKeeper.On(
			"Slash",
			mock.Anything,
			consAddress,
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			int64(100),
			expectedSlashFactor,
		).Return(math.NewInt(10), nil)

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)
	})

	s.Run("75% uptime with many wrap arounds", func() {
		feed, err := slatypes.NewPriceFeed(
			uint(20),
			consAddress,
			cp,
			id,
		)
		s.Require().NoError(err)

		for i := 0; i < 50; i++ {
			feed.SetUpdate(slatypes.VoteWithPrice)
		}

		for i := 0; i < 5; i++ {
			feed.SetUpdate(slatypes.VoteWithoutPrice)
		}

		s.stakingKeeper.On(
			"GetLastValidatorPower",
			mock.Anything,
			sdk.ValAddress(consAddress.Bytes()),
		).Return(int64(100), nil)

		expectedDevation := (expectedUptime.Sub(math.LegacyMustNewDecFromStr("0.75"))).Quo(expectedUptime)
		expectedSlashFactor := slashConstant.Mul(expectedDevation)

		s.slashingKeeper.On(
			"Slash",
			mock.Anything,
			consAddress,
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			int64(100),
			expectedSlashFactor,
		).Return(math.NewInt(10), nil)

		err = s.keeper.EnforceSLA(s.ctx, sla, feed)
		s.Require().NoError(err)
	})
}

func (s *KeeperTestSuite) TestSlash() {
	validator := sdk.ValAddress([]byte("validator"))
	power := int64(100)
	slashFactor := math.LegacyMustNewDecFromStr("0.1")

	s.Run("slash validator", func() {
		s.slashingKeeper.On(
			"Slash",
			mock.Anything,
			sdk.ConsAddress(validator.Bytes()),
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			power,
			slashFactor,
		).Return(math.NewInt(10), nil)

		err := s.keeper.Slash(s.ctx, validator, power, slashFactor)
		s.Require().NoError(err)
	})

	s.Run("returns error when slashing fails", func() {
		s.slashingKeeper.On(
			"Slash",
			mock.Anything,
			sdk.ConsAddress(validator.Bytes()),
			s.ctx.BlockHeight()-sdk.ValidatorUpdateDelay,
			power,
			slashFactor,
		).Return(math.NewInt(0), fmt.Errorf("failed to slash validator"))

		err := s.keeper.Slash(s.ctx, validator, power, slashFactor)
		s.Require().Error(err)
	})
}
