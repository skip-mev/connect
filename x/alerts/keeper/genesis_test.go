package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/alerts/types"
)

func (s *KeeperTestSuite) TestInitGenesisInvalidGenesis() {
	s.Run("test that init genesis with invalid genesis fails", func() {
		// create a fake genesis state
		gs := types.GenesisState{
			Params: types.NewParams(types.AlertParams{
				Enabled:     false,
				BondAmount:  sdk.NewCoin("test", math.NewInt(100)),
				MaxBlockAge: 1,
			}, nil, types.PruningParams{}),
		}

		// assert that InitGenesis panics
		s.Require().Panics(func() {
			s.alertKeeper.InitGenesis(s.ctx, gs)
		})
	})
}

func (s *KeeperTestSuite) TestInitGenesisValidGenesis() {
	// create genesis
	params := types.NewParams(types.AlertParams{
		Enabled:     true,
		BondAmount:  sdk.NewCoin("test", math.NewInt(100)),
		MaxBlockAge: 1,
	}, nil, types.PruningParams{})
	alert2 := types.NewAlertWithStatus(
		types.NewAlert(2, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
		types.NewAlertStatus(1, 2, time.Now(), 1),
	)

	alert1 := types.NewAlertWithStatus(
		types.NewAlert(1, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
		types.NewAlertStatus(1, 2, time.Now(), 1),
	)

	alert3 := types.NewAlertWithStatus(
		types.NewAlert(3, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
		types.NewAlertStatus(1, 2, time.Now(), 1),
	)

	gs := types.GenesisState{
		Params: params,
		Alerts: []types.AlertWithStatus{alert1, alert2, alert3},
	}

	s.Run("initialize with valid genesis", func() {
		// assert that InitGenesis does not panic
		s.Require().NotPanics(func() {
			s.alertKeeper.InitGenesis(s.ctx, gs)
		})
	})

	// check that all alerts are added
	s.Run("check that all alerts are added", func() {
		// check that alert1 was added to state
		alertInState, ok := s.alertKeeper.GetAlert(s.ctx, alert1.Alert)
		s.Require().True(ok)
		s.Require().Equal(alert1, alertInState)

		// check that alert2 was added to state
		alertInState, ok = s.alertKeeper.GetAlert(s.ctx, alert2.Alert)
		s.Require().True(ok)
		s.Require().Equal(alert2, alertInState)

		// check that alert3 was added to state
		alertInState, ok = s.alertKeeper.GetAlert(s.ctx, alert3.Alert)
		s.Require().True(ok)
		s.Require().Equal(alert3, alertInState)
	})

	s.Run("check that the params are set correctly", func() {
		// check that the params are set correctly
		s.Require().Equal(params, s.alertKeeper.GetParams(s.ctx))
	})
}

func (s *KeeperTestSuite) TestExportGenesis() {
	// create values + genesis-state
	params := types.NewParams(types.AlertParams{
		Enabled:     false,
		BondAmount:  sdk.NewCoin("test", math.NewInt(0)),
		MaxBlockAge: 0,
	}, nil, types.PruningParams{})
	alert2 := types.NewAlertWithStatus(
		types.NewAlert(2, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
		types.NewAlertStatus(1, 2, time.Now(), 1),
	)

	alert1 := types.NewAlertWithStatus(
		types.NewAlert(1, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
		types.NewAlertStatus(1, 2, time.Now(), 1),
	)

	alert3 := types.NewAlertWithStatus(
		types.NewAlert(3, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
		types.NewAlertStatus(1, 2, time.Now(), 1),
	)

	alert4 := types.NewAlertWithStatus(
		types.NewAlert(4, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
		types.NewAlertStatus(1, 2, time.Now(), 1),
	)

	gs := types.GenesisState{
		Params: params,
		Alerts: []types.AlertWithStatus{alert1, alert2, alert3},
	}
	var exportedGenesis types.GenesisState
	s.Run("test that init-genesis is successful", func() {
		// assert that InitGenesis does not panic
		s.Require().NotPanics(func() {
			s.alertKeeper.InitGenesis(s.ctx, gs)
		})
	})

	s.Run("test that additional alerts can be added to state", func() {
		s.Require().Nil(s.alertKeeper.SetAlert(s.ctx, alert4))
	})

	s.Run("test that genesis is exported as expected", func() {
		// assert that ExportGenesis does not panic
		s.Require().NotPanics(func() {
			gsTemp := s.alertKeeper.ExportGenesis(s.ctx)
			exportedGenesis = *gsTemp
		})

		s.Run("check that params is correct", func() {
			s.Require().Equal(params, exportedGenesis.Params)
		})

		s.Run("check that alerts are correct", func() {
			expectedAlerts := map[string]struct{}{
				string(alert1.Alert.UID()): {},
				string(alert2.Alert.UID()): {},
				string(alert3.Alert.UID()): {},
				string(alert4.Alert.UID()): {},
			}

			for _, alert := range exportedGenesis.Alerts {
				_, ok := expectedAlerts[string(alert.Alert.UID())]
				s.Require().True(ok)
			}
		})
	})

	s.Run("test that genesis exported is valid", func() {
		// validate the exported genesis
		s.Require().Nil(exportedGenesis.ValidateBasic())
	})
}
