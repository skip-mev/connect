package keeper_test

import (
	"testing"
	"time"

	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttime "github.com/cometbft/cometbft/types/time"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/suite"

	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/alerts/keeper"
	alerttypes "github.com/skip-mev/connect/v2/x/alerts/types"
	"github.com/skip-mev/connect/v2/x/alerts/types/mocks"
	"github.com/skip-mev/connect/v2/x/alerts/types/strategies"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context

	// alert keeper
	alertKeeper *keeper.Keeper
	// bank-keeper
	bk *mocks.BankKeeper
	// oracle-keeper
	ok *mocks.OracleKeeper
	// incentive-keeper
	ik *mocks.IncentiveKeeper
	// private-key
	privateKey cryptotypes.PrivKey
	// authority
	authority sdk.AccAddress
}

func (s *KeeperTestSuite) SetupTest() {
	key := storetypes.NewKVStoreKey(alerttypes.StoreKey)
	ss := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(s.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	s.ctx = testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: cmttime.Now()})
	encCfg := moduletestutil.MakeTestEncodingConfig()

	// register strategies interfaces to the encoding config
	strategies.RegisterInterfaces(encCfg.InterfaceRegistry)
	alerttypes.RegisterInterfaces(encCfg.InterfaceRegistry)

	s.bk = mocks.NewBankKeeper(s.T())
	s.ok = mocks.NewOracleKeeper(s.T())
	s.ik = mocks.NewIncentiveKeeper(s.T())

	s.authority = sdk.AccAddress("authority")
	s.alertKeeper = keeper.NewKeeper(ss, encCfg.Codec, s.ok, s.bk, s.ik, strategies.DefaultHandleValidatorIncentive(), s.authority)

	// create a private key
	s.privateKey = secp256k1.GenPrivKey()
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// test that we can set / remove / get alerts from the keeper.
func (s *KeeperTestSuite) TestAlerts() {
	alert := alerttypes.NewAlertWithStatus(
		alerttypes.NewAlert(1, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BTC", "USD")),
		alerttypes.NewAlertStatus(1, 2, time.Now(), alerttypes.Unconcluded),
	)
	alert2 := alerttypes.NewAlertWithStatus(
		alerttypes.NewAlert(2, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BTC", "USD")),
		alerttypes.NewAlertStatus(2, 3, time.Now(), alerttypes.Unconcluded),
	)
	// set an alert in the keeper
	s.Run("set alerts", func() {
		s.Require().NoError(s.alertKeeper.SetAlert(s.ctx, alert))
		s.Require().NoError(s.alertKeeper.SetAlert(s.ctx, alert2))

		// check that both alerts are in the keeper
		alertInState, ok := s.alertKeeper.GetAlert(s.ctx, alert.Alert)
		s.Require().True(ok)
		s.Require().Equal(alert, alertInState)

		alertInState, ok = s.alertKeeper.GetAlert(s.ctx, alert2.Alert)
		s.Require().True(ok)
		s.Require().Equal(alert2, alertInState)
	})

	// remove alert from the keeper
	s.Run("remove alert", func() {
		// remove a single alert, and check that alert2 exists, but alert1 does not
		s.Require().NoError(s.alertKeeper.RemoveAlert(s.ctx, alert.Alert))

		_, ok := s.alertKeeper.GetAlert(s.ctx, alert.Alert)

		s.Require().False(ok)

		alertInState, ok := s.alertKeeper.GetAlert(s.ctx, alert2.Alert)
		s.Require().Equal(alert2, alertInState)
		s.Require().True(ok)
	})

	// remove all alerts from the keeper
	s.Run("remove all alerts", func() {
		s.Require().NoError(s.alertKeeper.RemoveAlert(s.ctx, alert2.Alert))
		_, ok := s.alertKeeper.GetAlert(s.ctx, alert.Alert)
		s.Require().False(ok)

		_, ok = s.alertKeeper.GetAlert(s.ctx, alert2.Alert)
		s.Require().False(ok)
	})
}

func (s *KeeperTestSuite) TestGetAllAlerts() {
	// set some alerts in the keeper
	alert := alerttypes.NewAlertWithStatus(
		alerttypes.NewAlert(1, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BTC", "USD")),
		alerttypes.NewAlertStatus(1, 2, time.Now(), alerttypes.Unconcluded),
	)
	alert2 := alerttypes.NewAlertWithStatus(
		alerttypes.NewAlert(2, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BTC", "USD")),
		alerttypes.NewAlertStatus(2, 3, time.Now(), alerttypes.Unconcluded),
	)
	alert3 := alerttypes.NewAlertWithStatus(
		alerttypes.NewAlert(2, sdk.AccAddress("test"), slinkytypes.NewCurrencyPair("BTC", "USD")),
		alerttypes.NewAlertStatus(2, 3, time.Now(), alerttypes.Unconcluded),
	)

	// set alerts
	s.Require().NoError(s.alertKeeper.SetAlert(s.ctx, alert))
	s.Require().NoError(s.alertKeeper.SetAlert(s.ctx, alert2))
	s.Require().NoError(s.alertKeeper.SetAlert(s.ctx, alert3))

	// get all alerts
	expectedAlerts := make(map[string]struct{})

	for _, alert := range []alerttypes.AlertWithStatus{alert, alert2, alert3} {
		expectedAlerts[string(alert.Alert.UID())] = struct{}{}
	}

	alerts, err := s.alertKeeper.GetAllAlerts(s.ctx)
	s.Require().NoError(err)

	for _, alert := range alerts {
		_, ok := expectedAlerts[string(alert.Alert.UID())]
		s.Require().True(ok)
		delete(expectedAlerts, string(alert.Alert.UID()))
	}
}
