package keeper_test

import (
	"fmt"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"

	"github.com/skip-mev/slinky/x/alerts/keeper"
	"github.com/skip-mev/slinky/x/alerts/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

// matcher for sdk.Coins
type coinsMatcher struct {
	coins sdk.Coins
}

func (c *coinsMatcher) Matches(x interface{}) bool {
	coins, ok := x.(sdk.Coins)
	if !ok {
		return false
	}

	if len(coins) != len(c.coins) {
		return false
	}

	for i, coin := range coins {
		if !coin.IsEqual(c.coins[i]) {
			return false
		}
	}

	return true
}

func (c *coinsMatcher) String() string {
	return fmt.Sprintf("is equal to %v", c.coins)
}

func CoinsMatcher(coins sdk.Coins) gomock.Matcher {
	return &coinsMatcher{coins: coins}
}

func (s *KeeperTestSuite) TestConcludeAlert() {
	// set the params
	s.alertKeeper.SetParams(s.ctx, types.Params{
		AlertParams: types.AlertParams{
			Enabled:     true,
			MaxBlockAge: 10,
			BondAmount: sdk.NewCoin(
				"stake",
				math.NewInt(100),
			),
		},
	})

	cases := []struct {
		name          string
		alert         types.Alert
		status        keeper.ConclusionStatus
		setup         func(ctx sdk.Context)
		err           error
		expectedAlert types.AlertWithStatus
	}{
		{
			"invalid alert - fail",
			types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("base", "")),
			keeper.Negative,
			func(ctx sdk.Context) {},
			fmt.Errorf("invalid alert: empty quote or base string"),
			types.AlertWithStatus{},
		},
		{
			"alert not found - fail",
			types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
			keeper.Negative,
			func(ctx sdk.Context) {},
			fmt.Errorf("alert not found: %v", types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE"))),
			types.AlertWithStatus{},
		},
		{
			"alert already concluded",
			types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
			keeper.Negative,
			func(ctx sdk.Context) {
				// set the alert with concluded AlertStatus
				s.alertKeeper.SetAlert(
					ctx,
					types.NewAlertWithStatus(
						types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
						types.NewAlertStatus(1, 1, time.Now(), types.Concluded),
					),
				)
			},
			fmt.Errorf("alert already concluded"),
			types.AlertWithStatus{},
		},
		{
			"alert status unknown",
			types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
			keeper.ConclusionStatus(3),
			func(ctx sdk.Context) {
				// set the alert with concluded AlertStatus
				s.alertKeeper.SetAlert(
					ctx,
					types.NewAlertWithStatus(
						types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
						types.NewAlertStatus(1, 1, time.Now(), types.Unconcluded),
					),
				)
			},
			fmt.Errorf("invalid status: 3"),
			types.AlertWithStatus{},
		},
		{
			"negative alert - bond is burned",
			types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
			keeper.Negative,
			func(ctx sdk.Context) {
				alert := types.NewAlertWithStatus(
					types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
					types.NewAlertStatus(10, 11, time.Time{}, types.Unconcluded),
				)
				// set the unconcluded alert
				s.alertKeeper.SetAlert(
					ctx,
					alert,
				)

				s.bk.EXPECT().BurnCoins(
					gomock.Any(),
					types.ModuleName,
					CoinsMatcher(sdk.NewCoins(s.alertKeeper.GetParams(s.ctx).AlertParams.BondAmount)),
				).Return(nil)
			},
			nil,
			types.NewAlertWithStatus(
				types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
				types.NewAlertStatus(10, 11, time.Time{}, types.Concluded),
			),
		},
		{
			"positive alert - bond is returned",
			types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
			keeper.Positive,
			func(ctx sdk.Context) {
				alert := types.NewAlertWithStatus(
					types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
					types.NewAlertStatus(10, 11, time.Time{}, types.Unconcluded),
				)
				// set the unconcluded alert
				s.alertKeeper.SetAlert(
					ctx,
					alert,
				)

				s.bk.EXPECT().SendCoinsFromModuleToAccount(
					gomock.Any(),
					types.ModuleName,
					sdk.AccAddress("abc"),
					CoinsMatcher(sdk.NewCoins(s.alertKeeper.GetParams(s.ctx).AlertParams.BondAmount)),
				).Return(nil)
			},
			nil,
			types.NewAlertWithStatus(
				types.NewAlert(1, sdk.AccAddress("abc"), oracletypes.NewCurrencyPair("BASE", "QUOTE")),
				types.NewAlertStatus(10, 11, time.Time{}, types.Concluded),
			),
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			// setup
			tc.setup(s.ctx)

			// when
			err := s.alertKeeper.ConcludeAlert(s.ctx, tc.alert, tc.status)

			// then
			if tc.err != nil {
				s.Require().Error(err)
				s.Require().Equal(tc.err.Error(), err.Error())
				return
			}

			s.Require().NoError(err)

			// assert equality of saved alerts
			alert, found := s.alertKeeper.GetAlert(s.ctx, tc.alert)
			s.Require().True(found)
			s.Require().Equal(tc.expectedAlert, alert)
		})
	}
}
