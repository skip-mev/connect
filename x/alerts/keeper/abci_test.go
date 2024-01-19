package keeper_test

import (
	"time"

	"github.com/stretchr/testify/mock"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/slinky/x/alerts/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func (s *KeeperTestSuite) TestEndBlocker() {
	// set context
	s.ctx = s.ctx.WithBlockHeight(10)

	// set three alerts (this shld be purged first)
	alert1 := types.NewAlertWithStatus(
		types.NewAlert(1, sdk.AccAddress("abc1"), oracletypes.NewCurrencyPair("BTC", "USD")),
		types.NewAlertStatus(10, 10, time.Time{}, types.Concluded),
	)

	// this will be purged next
	alert2 := types.NewAlertWithStatus(
		types.NewAlert(2, sdk.AccAddress("abc2"), oracletypes.NewCurrencyPair("BTC", "USD")),
		types.NewAlertStatus(10, 11, time.Time{}, types.Concluded),
	)

	// this will be purged last
	alert3 := types.NewAlertWithStatus(
		types.NewAlert(3, sdk.AccAddress("abc3"), oracletypes.NewCurrencyPair("BTC", "USD")),
		types.NewAlertStatus(10, 12, time.Time{}, types.Unconcluded),
	)

	// set alerts in the store
	s.Require().NoError(s.alertKeeper.SetAlert(s.ctx, alert1))
	s.Require().NoError(s.alertKeeper.SetAlert(s.ctx, alert2))
	s.Require().NoError(s.alertKeeper.SetAlert(s.ctx, alert3))

	s.Run("expect no alerts are pruned at endblock if pruning is disabled in end-block", func() {
		err := s.alertKeeper.SetParams(
			s.ctx,
			types.NewParams(
				types.AlertParams{},
				nil,
				types.PruningParams{
					Enabled: false,
				},
			),
		)
		s.Require().NoError(err)

		// run endblocker
		updates, err := s.alertKeeper.EndBlocker(s.ctx)
		s.Require().NoError(err)
		s.Require().Nil(updates)

		// assert that all alerts are still in the store
		alerts, err := s.alertKeeper.GetAllAlerts(s.ctx)
		s.Require().NoError(err)
		s.Require().Len(alerts, 3)
	})

	// enable pruning
	err := s.alertKeeper.SetParams(
		s.ctx,
		types.NewParams(
			types.AlertParams{
				Enabled:     true,
				BondAmount:  sdk.NewCoin("test", math.NewInt(100)),
				MaxBlockAge: 1,
			},
			nil,
			types.PruningParams{
				Enabled: true,
			},
		),
	)
	s.Require().NoError(err)

	s.Run("expect first alert is pruned at the end of endblock", func() {
		updates, err := s.alertKeeper.EndBlocker(s.ctx)
		s.Require().NoError(err)
		s.Require().Nil(updates)

		// assert that the first alert is pruned
		alerts, err := s.alertKeeper.GetAllAlerts(s.ctx)
		s.Require().NoError(err)

		s.Require().Len(alerts, 2)

		// query the first alert
		_, ok := s.alertKeeper.GetAlert(s.ctx, alert1.Alert)
		s.Require().False(ok)

		// query the second alert
		_, ok = s.alertKeeper.GetAlert(s.ctx, alert2.Alert)
		s.Require().True(ok)

		// query the third alert
		_, ok = s.alertKeeper.GetAlert(s.ctx, alert3.Alert)
		s.Require().True(ok)
	})

	// increment block height
	s.ctx = s.ctx.WithBlockHeight(11)
	s.Run("expect second alert is pruned at the end of endblock", func() {
		updates, err := s.alertKeeper.EndBlocker(s.ctx)
		s.Require().NoError(err)
		s.Require().Nil(updates)

		// assert that the second alert is pruned
		_, ok := s.alertKeeper.GetAlert(s.ctx, alert2.Alert)
		s.Require().False(ok)

		// assert that the third alert is still in the store
		_, ok = s.alertKeeper.GetAlert(s.ctx, alert3.Alert)
		s.Require().True(ok)
	})

	// increment block height
	s.ctx = s.ctx.WithBlockHeight(12)
	s.Run("expect third alert is pruned at the end of endblock", func() {
		s.bk.On("SendCoinsFromModuleToAccount",
			mock.Anything,
			types.ModuleName,
			sdk.AccAddress("abc3"),
			sdk.NewCoins(s.alertKeeper.GetParams(s.ctx).AlertParams.BondAmount),
		).Return(nil)

		updates, err := s.alertKeeper.EndBlocker(s.ctx)
		s.Require().NoError(err)
		s.Require().Nil(updates)

		// assert that the third alert is pruned
		_, ok := s.alertKeeper.GetAlert(s.ctx, alert3.Alert)
		s.Require().False(ok)
	})
}
