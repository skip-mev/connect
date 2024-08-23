package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/alerts/keeper"
	"github.com/skip-mev/connect/v2/x/alerts/types"
)

func (s *KeeperTestSuite) TestQueryServer() {
	s.Run("Alerts", func() {
		// add a concluded alert to state
		concludedAlert := types.NewAlertWithStatus(
			types.NewAlert(1, sdk.AccAddress("abc1"), slinkytypes.NewCurrencyPair("BTC", "USD")),
			types.NewAlertStatus(1, 1, s.ctx.BlockTime(), types.Concluded),
		)

		// add an unconcluded alert to state
		unconcludedAlert := types.NewAlertWithStatus(
			types.NewAlert(2, sdk.AccAddress("abc2"), slinkytypes.NewCurrencyPair("BTC", "USD")),
			types.NewAlertStatus(1, 1, s.ctx.BlockTime(), types.Unconcluded),
		)

		// add alerts to state
		s.Require().NoError(s.alertKeeper.SetAlert(s.ctx, concludedAlert))
		s.Require().NoError(s.alertKeeper.SetAlert(s.ctx, unconcludedAlert))

		qs := keeper.NewQueryServer(*s.alertKeeper)

		s.Run("nil request - fail", func() {
			_, err := qs.Alerts(s.ctx, nil)
			s.Require().Error(err)
		})

		s.Run("AlertStatusID Concluded - pass", func() {
			res, err := qs.Alerts(s.ctx, &types.AlertsRequest{
				Status: types.AlertStatusID_CONCLUSION_STATUS_CONCLUDED,
			})
			s.Require().NoError(err)

			s.Require().Len(res.Alerts, 1)
			s.Require().Equal(concludedAlert.Alert, res.Alerts[0])
		})

		s.Run("AlertStatusID Unconcluded - pass", func() {
			res, err := qs.Alerts(s.ctx, &types.AlertsRequest{
				Status: types.AlertStatusID_CONCLUSION_STATUS_UNCONCLUDED,
			})
			s.Require().NoError(err)

			s.Require().Len(res.Alerts, 1)
			s.Require().Equal(unconcludedAlert.Alert, res.Alerts[0])
		})

		s.Run("AlertStatusID All - pass", func() {
			res, err := qs.Alerts(s.ctx, &types.AlertsRequest{
				Status: types.AlertStatusID_CONCLUSION_STATUS_UNSPECIFIED,
			})
			s.Require().NoError(err)

			s.Require().Len(res.Alerts, 2)
			expectedAlerts := map[string]struct{}{
				string(concludedAlert.Alert.UID()):   {},
				string(unconcludedAlert.Alert.UID()): {},
			}

			for _, a := range res.Alerts {
				_, ok := expectedAlerts[string(a.UID())]
				s.Require().True(ok)
			}
		})
	})
}

func (s *KeeperTestSuite) TestParams() {
	params := s.alertKeeper.GetParams(s.ctx)

	qs := keeper.NewQueryServer(*s.alertKeeper)
	s.Run("nil request", func() {
		_, err := qs.Params(s.ctx, nil)
		s.Require().Error(err)
	})

	s.Run("valid request", func() {
		res, err := qs.Params(s.ctx, &types.ParamsRequest{})
		s.Require().NoError(err)
		s.Require().Equal(params, res.Params)
	})
}
