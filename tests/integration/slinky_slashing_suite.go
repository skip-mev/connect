package integration

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/math"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/strangelove-ventures/interchaintest/v8"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"

	slinkytypes "github.com/skip-mev/slinky/pkg/types"
	alerttypes "github.com/skip-mev/slinky/x/alerts/types"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

const validatorStake = 5000000000000

type SlinkySlashingIntegrationSuite struct {
	*SlinkyIntegrationSuite

	multiSigUser1 cosmos.User
	multiSigUser2 cosmos.User

	multiSigPk1 cryptotypes.PrivKey
	multiSigPk2 cryptotypes.PrivKey
}

func NewSlinkySlashingIntegrationSuite(ss *SlinkyIntegrationSuite) *SlinkySlashingIntegrationSuite {
	return &SlinkySlashingIntegrationSuite{
		SlinkyIntegrationSuite: ss,
	}
}

func (s *SlinkySlashingIntegrationSuite) SetupSuite() {
	s.SlinkyIntegrationSuite.TearDownSuite()

	s.SlinkyIntegrationSuite.SetupSuite()

	// initialize multiSigUsers
	users := interchaintest.GetAndFundTestUsers(s.T(), context.Background(), s.T().Name(), math.NewInt(genesisAmount), s.chain, s.chain)
	s.multiSigUser1 = users[0]
	s.multiSigUser2 = users[1]

	s.multiSigPk1 = CreateKey(Secp256k1)
	s.multiSigPk2 = CreateKey(Ed25519)
}

func (s *SlinkySlashingIntegrationSuite) SetupTest() {
	// get the validators after the slashing
	validators, err := QueryValidators(s.chain)
	s.Require().NoError(err)

	// for any validators that do not have correct stake, delegate to them
	for _, validator := range validators {
		if validator.Tokens.Int64() < validatorStake {
			toStake := sdk.NewCoin(s.denom, math.NewInt(validatorStake-validator.Tokens.Int64()))

			// delegate to the validator
			resp, err := s.Delegate(s.multiSigUser1, validator.OperatorAddress, toStake)
			s.Require().NoError(err)

			s.Require().Equal(uint32(0), resp.CheckTx.Code)
			s.Require().Equal(uint32(0), resp.TxResult.Code)
		}
	}

	// get the validators after the slashing
	validators, err = QueryValidators(s.chain)
	s.Require().NoError(err)

	// ensure all validators have the correct stake
	for _, validator := range validators {
		s.Require().Equal(validator.Tokens.Int64(), int64(validatorStake))
	}
}

func (s *SlinkySlashingIntegrationSuite) Height(ctx context.Context) (uint64, error) {
	height, err := s.chain.Height(ctx)
	if err != nil {
		return 0, err
	}

	return uint64(height), nil
}

func (s *SlinkySlashingIntegrationSuite) TestAlerts() {
	// test getting / setting params
	s.Run("test get params", func() {
		cc, close, err := GetChainGRPC(s.chain)
		s.Require().NoError(err)
		defer close()

		expectedParams := alerttypes.DefaultParams("stake", nil)

		// set default params for the module
		_, err = UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, expectedParams)
		s.Require().NoError(err)

		// create an alerts client
		alertsClient := alerttypes.NewQueryClient(cc)

		// get the params
		params, err := alertsClient.Params(context.Background(), &alerttypes.ParamsRequest{})
		s.Require().NoError(err)

		s.Require().Equal(expectedParams, params.Params)
	})

	s.Run("test setting new params", func() {
		cc, close, err := GetChainGRPC(s.chain)
		s.Require().NoError(err)
		defer close()

		// create an alerts client
		alertsClient := alerttypes.NewQueryClient(cc)

		s.Run("test setting params", func() {
			// update the alert Params
			expect := alerttypes.Params{
				PruningParams: alerttypes.PruningParams{
					Enabled:       true,
					BlocksToPrune: 100,
				},
			}

			// update the module's alert-params
			_, err := UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, expect)
			s.Require().NoError(err)

			// get the params
			got, err := alertsClient.Params(context.Background(), &alerttypes.ParamsRequest{})
			s.Require().NoError(err)

			// expect the params to be updated
			s.Require().True(compareParams(expect, got.Params))
		})
	})
}

func compareParams(a, b alerttypes.Params) bool {
	return a.PruningParams == b.PruningParams &&
		a.AlertParams.Enabled == b.AlertParams.Enabled &&
		a.AlertParams.MaxBlockAge == b.AlertParams.MaxBlockAge
}

// TestSubmittingAlerts will test the alert-submission process, specifically
//   - submitting an alert when alerts are disabled, and expect that the alert is not submitted
//   - submitting an alert when alerts are enabled, but the alert is too old, and expect that the alert is not submitted
//   - submitting an alert when alerts are enabled, and the alert references a non-existent currency-pair, and expect that the alert is not submitted
//   - submitting an alert when alerts are enabled, and the alert is valid, and expect that the alert is submitted. On submission check that the
//
// alert's bond is escrowed at the alert module
//   - submitting an alert when alerts are enabled, and the alert is valid, but the alert has already been submitted, and expect that the alert is not submitted
func (s *SlinkySlashingIntegrationSuite) TestSubmittingAlerts() {
	s.Run("test submitting an alert when alerts are disabled - fail", func() {
		// set params to have alerts disabled
		_, err := UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, alerttypes.Params{
			AlertParams: alerttypes.AlertParams{
				Enabled: false,
			},
		})
		s.Require().NoError(err)

		alertSubmitter, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
		s.Require().NoError(err)

		// submit an alert
		resp, err := s.SubmitAlert(
			s.multiSigUser1,
			alerttypes.NewAlert(
				1,
				alertSubmitter,
				slinkytypes.NewCurrencyPair("BTC", "USD"),
			),
		)
		s.Require().NoError(err)

		// check the response from the chain
		s.Require().Equal(resp.TxResult.Code, uint32(1))
		s.Require().True(strings.Contains(resp.TxResult.Log, "alerts are not enabled"))
	})

	s.Run("test submitting an alert with a block-age > max-block-age - fail", func() {
		// update the max-block-age so that it's feasible for us to submit an alert with a block-age (ctx.Height() - alert.Height) > max-block-age
		_, err := UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, alerttypes.Params{
			AlertParams: alerttypes.AlertParams{
				MaxBlockAge: 1,
				Enabled:     true,
				BondAmount:  sdk.NewCoin(s.denom, alerttypes.DefaultBondAmount),
			},
		})
		s.Require().NoError(err)

		alertSubmitter, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
		s.Require().NoError(err)

		// get the current height, and submit an alert with a height that is too old
		height, err := s.Height(context.Background())
		s.Require().NoError(err)

		// submit an alert
		resp, err := s.SubmitAlert(
			s.multiSigUser1,
			alerttypes.NewAlert(
				uint64(height-5),
				alertSubmitter,
				slinkytypes.NewCurrencyPair("BTC", "USD"),
			),
		)
		s.Require().NoError(err)

		// check the response from the chain
		s.Require().Equal(resp.TxResult.Code, uint32(1))
		s.Require().True(strings.Contains(resp.TxResult.Log, "alert is too old"))
	})

	s.Run("test submitting an alert for a non-existent currency-pair - fail", func() {
		cc, close, err := GetChainGRPC(s.chain)
		s.Require().NoError(err)
		defer close()

		// update the max-block-age
		_, err = UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, alerttypes.Params{
			AlertParams: alerttypes.AlertParams{
				MaxBlockAge: 20,
				Enabled:     true,
				BondAmount:  sdk.NewCoin(s.denom, alerttypes.DefaultBondAmount),
			},
		})
		s.Require().NoError(err)

		// check if the BTC/USD currency pair exists
		oraclesClient := oracletypes.NewQueryClient(cc)
		_, err = oraclesClient.GetPrice(context.Background(), &oracletypes.GetPriceRequest{
			CurrencyPair: slinkytypes.CurrencyPair{
				Base:  "BTC",
				Quote: "USD",
			},
		})
		s.Require().NoError(err)

		// get the current height
		height, err := s.Height(context.Background())
		s.Require().NoError(err)

		alertSubmitter, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
		s.Require().NoError(err)

		// the currency-pair does not exist
		resp, err := s.SubmitAlert(
			s.multiSigUser1,
			alerttypes.NewAlert(
				uint64(height-1),
				alertSubmitter,
				slinkytypes.CurrencyPair{
					Base:  "MOG",
					Quote: "GOM",
				},
			),
		)
		s.Require().NoError(err)

		// check the response from the chain
		s.Require().Equal(resp.TxResult.Code, uint32(1))
		s.Require().True(strings.Contains(resp.TxResult.Log, fmt.Sprint("currency pair MOG/GOM does not exist")), resp.TxResult.Log)
	})

	// submit the alert (the max block-age set previously will suffice)
	height, err := s.Height(context.Background())
	s.Require().NoError(err)

	s.Run("test submitting an alert and check for bond deposit - pass", func() {
		// add the BTC/USD currency pair
		cc, close, err := GetChainGRPC(s.chain)
		s.Require().NoError(err)
		defer close()

		cp := slinkytypes.NewCurrencyPair("BTC", "USD")

		// check if the BTC/USD currency pair exists
		oraclesClient := oracletypes.NewQueryClient(cc)
		_, err = oraclesClient.GetPrice(context.Background(), &oracletypes.GetPriceRequest{
			CurrencyPair: cp,
		})
		if err != nil {
			// add if there was an error
			s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, cp))
		}

		alertSubmitter, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
		s.Require().NoError(err)

		// get the balance of the alert submitter, so that we can check if the bond was deposited after the alert was submitted
		balance, err := s.chain.GetBalance(context.Background(), s.multiSigUser1.FormattedAddress(), s.denom)
		s.Require().NoError(err)

		modAcct := authtypes.NewModuleAddress(alerttypes.ModuleName)

		// do the same for the module account balance
		modAcctBalance, err := s.chain.GetBalance(context.Background(), modAcct.String(), s.denom)
		s.Require().NoError(err)

		// submit the alerts
		alert := alerttypes.NewAlert(
			height-1,
			alertSubmitter,
			cp,
		)

		resp, err := s.SubmitAlert(
			s.multiSigUser1,
			alert,
		)

		// check the response from the chain
		s.Require().NoError(err)
		s.Require().Equal(resp.TxResult.Code, uint32(0))

		alertClient := alerttypes.NewQueryClient(cc)

		// query for the alert
		alertResp, err := alertClient.Alerts(context.Background(), &alerttypes.AlertsRequest{
			Status: alerttypes.AlertStatusID_CONCLUSION_STATUS_UNCONCLUDED,
		})
		s.Require().NoError(err)

		// check the alert
		s.Require().Len(alertResp.Alerts, 1)
		s.Require().Equal(alertResp.Alerts[0], alert)

		// expect submitters balance to be in escrow
		params, err := alertClient.Params(context.Background(), &alerttypes.ParamsRequest{})
		s.Require().NoError(err)

		bondAmt := params.Params.AlertParams.BondAmount.Amount

		// check the balance of the alert submitter
		newBalance, err := s.chain.GetBalance(context.Background(), s.multiSigUser1.FormattedAddress(), s.denom)
		s.Require().NoError(err)

		s.Require().Equal(balance.Sub(bondAmt).Int64()-(resp.TxResult.GasWanted*gasPrice), newBalance.Int64())

		// check the balance of the module account
		newModAcctBalance, err := s.chain.GetBalance(context.Background(), modAcct.String(), s.denom)
		s.Require().NoError(err)

		// diff shld be the bond amount
		s.Require().Equal(modAcctBalance.Add(bondAmt).Int64(), newModAcctBalance.Int64())
	})

	s.Run("test submitting an alert after it has already been submitted - fail", func() {
		// file a duplicate alert
		cp := slinkytypes.NewCurrencyPair("BTC", "USD")

		alertSubmitter, err := sdk.AccAddressFromBech32(s.multiSigUser2.FormattedAddress())
		s.Require().NoError(err)

		alert := alerttypes.NewAlert(
			height-1,
			alertSubmitter,
			cp,
		)

		resp, err := s.SubmitAlert(
			s.multiSigUser2,
			alert,
		)

		// check the response from the chain
		s.Require().NoError(err)
		s.Require().Equal(resp.TxResult.Code, uint32(1))

		s.Require().True(strings.Contains(resp.TxResult.Log, fmt.Sprintf("alert with UID %X already exists", alert.UID())))
	})
}

// TestAlertPruning tests the pruning of alerts, specifically we submit 2 alerts, wait for some period of time between them,
// and expect that they are pruned after BlocksToPrune has passed. Also test that after an alert is concluded, it's BlocksToPrune
// is updated to alert.Height + MaxBlockAge, so that the same alert cannot be submitted + concluded twice.
func (s *SlinkySlashingIntegrationSuite) TestAlertPruning() {
	// check if the BTC/USD currency pair exists
	cc, close, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)

	defer close()

	cp := slinkytypes.NewCurrencyPair("BTC", "USD") // arbitrary

	// expect that the above currency-pair is in state, so we can submit alerts that reference it
	oraclesClient := oracletypes.NewQueryClient(cc)
	_, err = oraclesClient.GetPrice(context.Background(), &oracletypes.GetPriceRequest{
		CurrencyPair: cp,
	})
	if err != nil {
		// remove the currency-pair
		s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, cp))
	}
	// add pruning params with updated max-block-age

	var (
		maxBlockAge   = uint64(20)
		blocksToPrune = uint64(10)
	)

	_, err = UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, alerttypes.Params{
		AlertParams: alerttypes.AlertParams{
			MaxBlockAge: maxBlockAge,
			Enabled:     true,
			BondAmount:  sdk.NewCoin(s.denom, alerttypes.DefaultBondAmount),
		},
		PruningParams: alerttypes.PruningParams{
			Enabled:       true,
			BlocksToPrune: blocksToPrune,
		},
	})
	s.Require().NoError(err)

	s.Run("expect all alerts to have been pruned", func() {
		_, err := ExpectAlerts(s.chain, 3*s.blockTime, nil)
		s.Require().NoError(err)
	})

	s.Run("test that an alert is pruned after no conclusion is submitted, and the alert deposit is returned", func() {
		submitter1, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
		s.Require().NoError(err)

		height, err := s.Height(context.Background())
		s.Require().NoError(err)

		submitter2, err := sdk.AccAddressFromBech32(s.multiSigUser2.FormattedAddress())
		s.Require().NoError(err)

		// submit two alerts for different heights from different accounts
		alert1 := alerttypes.NewAlert(
			height,
			submitter1,
			cp,
		)

		alert2 := alerttypes.NewAlert(
			height+1,
			submitter2,
			cp,
		)

		// get the balances of the accounts before alert submission, so that we can compare after submission whether the bond
		// is in escrow.
		balance1Before, err := s.chain.GetBalance(context.Background(), s.multiSigUser1.FormattedAddress(), s.denom)
		s.Require().NoError(err)

		balance2Before, err := s.chain.GetBalance(context.Background(), s.multiSigUser2.FormattedAddress(), s.denom)
		s.Require().NoError(err)

		var (
			commitHeight, commitHeight2  uint64
			balance2After, balance1After math.Int
		)
		s.Run("submit the first alert, from the first multi-sig address", func() {
			// submit the first alert
			resp, err := s.SubmitAlert(
				s.multiSigUser1,
				alert1,
			)
			s.Require().NoError(err)

			// check the response from the chain
			s.Require().Equal(resp.TxResult.Code, uint32(0))

			// expect the alert to have been committed in state
			commitHeight = uint64(resp.Height)

			// query for the alert,
			alertClient := alerttypes.NewQueryClient(cc)
			alertResp, err := alertClient.Alerts(context.Background(), &alerttypes.AlertsRequest{
				Status: alerttypes.AlertStatusID_CONCLUSION_STATUS_UNCONCLUDED,
			})

			s.Require().NoError(err)

			// check the alerts
			s.Require().Len(alertResp.Alerts, 1)
			s.Require().Equal(alertResp.Alerts[0], alert1)

			// check the bond was escrowed
			balance1After, err = s.chain.GetBalance(context.Background(), s.multiSigUser1.FormattedAddress(), s.denom)
			s.Require().NoError(err)

			// check the balance, expect the bond to be in escrow
			s.Require().Equal(balance1Before.Sub(balance1After).Int64()-resp.TxResult.GasWanted*gasPrice, alerttypes.DefaultBondAmount.Int64())
		})

		s.Run("submit the second alert, from the second multi-sig", func() {
			// submit the second alert
			resp, err := s.SubmitAlert(
				s.multiSigUser2,
				alert2,
			)
			s.Require().NoError(err)

			commitHeight2 = uint64(resp.Height)
			// check the response from the chain
			s.Require().Equal(resp.TxResult.Code, uint32(0))

			// expect the alerts to have been committed in state
			alertClient := alerttypes.NewQueryClient(cc)
			alertResp, err := alertClient.Alerts(context.Background(), &alerttypes.AlertsRequest{
				Status: alerttypes.AlertStatusID_CONCLUSION_STATUS_UNCONCLUDED,
			})

			s.Require().NoError(err)

			// check the alerts
			s.Require().Len(alertResp.Alerts, 2)

			// check the bond was escrowed
			balance2After, err = s.chain.GetBalance(context.Background(), s.multiSigUser2.FormattedAddress(), s.denom)
			s.Require().NoError(err)

			// check the balance
			s.Require().Equal(balance2Before.Sub(balance2After).Int64()-resp.TxResult.GasWanted*gasPrice, alerttypes.DefaultBondAmount.Int64())
		})

		s.Run("wait for 10 blocks after commit height of first alert, and expect it to be pruned", func() {
			// wait for commitheight + 10
			height, err = ExpectAlerts(s.chain, s.blockTime*3, []alerttypes.Alert{alert2})
			s.Require().NoError(err)

			// check that height > commitHeight + 10
			s.Require().True(height >= commitHeight+10)

			// check the bond was returned
			balance1Final, err := s.chain.GetBalance(context.Background(), s.multiSigUser1.FormattedAddress(), s.denom)
			s.Require().NoError(err)

			// check the balance
			s.Require().Equal(balance1Final.Sub(balance1After).Int64(), alerttypes.DefaultBondAmount.Int64())
		})

		s.Run("wait for blocksToPrune blocks after commit height of second alert, and expect it to be pruned", func() {
			// wait for commitheight + 10
			height, err = ExpectAlerts(s.chain, s.blockTime*3, []alerttypes.Alert{})
			s.Require().NoError(err)

			// check that height > commitHeight + 10
			s.Require().True(height >= commitHeight2+10)

			// check the bond was returned
			balance2Final, err := s.chain.GetBalance(context.Background(), s.multiSigUser2.FormattedAddress(), s.denom)
			s.Require().NoError(err)

			// check the balance
			s.Require().Equal(balance2Final.Sub(balance2After).Int64(), alerttypes.DefaultBondAmount.Int64())
		})
	})

	s.Run("test that after an alert is concluded, the alerts pruning height is updated", func() {
		cc, close, err = GetChainGRPC(s.chain)
		s.Require().NoError(err)

		defer close()

		alertClient := alerttypes.NewQueryClient(cc)

		var (
			moduleBalanceBefore math.Int
			alertHeight         uint64
			maxBlockAge         uint64 = 20
			alert               alerttypes.Alert
		)

		s.Run("update the params to have multiSigAddress1 / 2 as signers", func() {
			cvp, err := alerttypes.NewMultiSigVerificationParams(
				[]cryptotypes.PubKey{
					s.multiSigPk1.PubKey(),
					s.multiSigPk2.PubKey(),
				},
			)
			s.Require().NoError(err)

			any, err := codectypes.NewAnyWithValue(cvp)
			s.Require().NoError(err)

			// update the params
			_, err = UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, alerttypes.Params{
				AlertParams: alerttypes.AlertParams{
					MaxBlockAge: maxBlockAge,
					Enabled:     true,
					BondAmount:  sdk.NewCoin(s.denom, alerttypes.DefaultBondAmount),
				},
				PruningParams: alerttypes.PruningParams{
					Enabled:       true,
					BlocksToPrune: blocksToPrune,
				},
				ConclusionVerificationParams: any,
			})
			s.Require().NoError(err)

			// get the params
			params, err := alertClient.Params(context.Background(), &alerttypes.ParamsRequest{})
			s.Require().NoError(err)

			// check the params
			cdc := s.chain.Config().EncodingConfig.Codec

			var expectedCVP alerttypes.ConclusionVerificationParams
			s.Require().NoError(cdc.UnpackAny(params.Params.ConclusionVerificationParams, &expectedCVP))
			cvpGot, ok := expectedCVP.(*alerttypes.MultiSigConclusionVerificationParams)
			s.Require().True(ok)

			for i, pk := range cvpGot.Signers {
				var pk1 cryptotypes.PubKey
				s.Require().NoError(cdc.UnpackAny(pk, &pk1))

				var pk2 cryptotypes.PubKey
				s.Require().NoError(cdc.UnpackAny(cvp.(*alerttypes.MultiSigConclusionVerificationParams).Signers[i], &pk2))
				s.Require().True(pk1.Equals(pk2))
			}
		})

		alertHeight, err = s.Height(context.Background())
		s.Require().NoError(err)

		s.Run("submit an alert from multiSigAddress1", func() {
			submitAddr, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())

			alert = alerttypes.NewAlert(
				alertHeight,
				submitAddr,
				slinkytypes.NewCurrencyPair(
					"BTC",
					"USD",
				),
			)
			s.Require().NoError(err)

			balance1Before, err := s.chain.GetBalance(context.Background(), s.multiSigUser1.FormattedAddress(), s.denom)
			s.Require().NoError(err)

			// get balance of module account
			moduleBalanceBefore, err = s.chain.GetBalance(context.Background(), authtypes.NewModuleAddress(alerttypes.ModuleName).String(), s.denom)
			s.Require().NoError(err)

			// submit the alert
			resp, err := s.SubmitAlert(
				s.multiSigUser1,
				alert,
			)
			s.Require().NoError(err)

			s.T().Log()

			// check the response (it was committed)
			s.Require().Equal(resp.CheckTx.Code, uint32(0))
			s.Require().Equal(resp.TxResult.Code, uint32(0))

			// query for the alert
			alertResp, err := alertClient.Alerts(context.Background(), &alerttypes.AlertsRequest{
				Status: alerttypes.AlertStatusID_CONCLUSION_STATUS_UNCONCLUDED,
			})
			s.Require().NoError(err)

			// check the alerts
			s.Require().Len(alertResp.Alerts, 1)
			s.Require().Equal(alertResp.Alerts[0], alert)

			// expect the bond to be escrowed
			balance1, err := s.chain.GetBalance(context.Background(), s.multiSigUser1.FormattedAddress(), s.denom)
			s.Require().NoError(err)

			// check the balance
			s.Require().Equal(balance1Before.Sub(balance1).Int64(), alerttypes.DefaultBondAmount.Int64()+resp.TxResult.GasWanted*gasPrice)

			// check the balance of the module account after alert submission
			moduleBalanceAfter, err := s.chain.GetBalance(context.Background(), authtypes.NewModuleAddress(alerttypes.ModuleName).String(), s.denom)
			s.Require().NoError(err)

			// check the balance
			s.Require().Equal(moduleBalanceBefore.Add(alerttypes.DefaultBondAmount).Int64(), moduleBalanceAfter.Int64())
		})

		s.Run("submit a conclusion from the multi-sig", func() {
			// get the alert
			alertResp, err := alertClient.Alerts(context.Background(), &alerttypes.AlertsRequest{
				Status: alerttypes.AlertStatusID_CONCLUSION_STATUS_UNCONCLUDED,
			})
			s.Require().NoError(err)
			alert := alertResp.Alerts[0]

			conclusion := alerttypes.MultiSigConclusion{
				Alert:              alert,
				ExtendedCommitInfo: cmtabci.ExtendedCommitInfo{},
				Status:             false,
				PriceBound: alerttypes.PriceBound{
					High: "1",
					Low:  "0",
				},
				Signatures: make([]alerttypes.Signature, 0),
			}

			// sign the conclusion
			sigBytes, err := conclusion.SignBytes()
			s.Require().NoError(err)

			// sign the conclusions
			sig1, err := s.multiSigPk1.Sign(sigBytes)
			s.Require().NoError(err)

			sig2, err := s.multiSigPk2.Sign(sigBytes)
			s.Require().NoError(err)

			conclusion.Signatures = append(conclusion.Signatures, alerttypes.Signature{
				Signer:    sdk.AccAddress(s.multiSigPk1.PubKey().Address()).String(),
				Signature: sig1,
			})
			conclusion.Signatures = append(conclusion.Signatures, alerttypes.Signature{
				Signer:    sdk.AccAddress(s.multiSigPk2.PubKey().Address()).String(),
				Signature: sig2,
			})

			// submit the conclusion
			resp, err := s.SubmitConclusion(s.multiSigUser1, &conclusion)
			s.Require().NoError(err)

			// expect conclusion submission to be successful
			s.Require().Equal(resp.CheckTx.Code, uint32(0))
			s.Require().Equal(resp.TxResult.Code, uint32(0))

			s.Run("expect that the bond is burned", func() {
				// get the final module account balance
				moduleBalanceAfter, err := s.chain.GetBalance(context.Background(), authtypes.NewModuleAddress(alerttypes.ModuleName).String(), s.denom)
				s.Require().NoError(err)

				// check the balance, should be returned to before alert submission (alert bond is burned)
				s.Require().Equal(moduleBalanceBefore.Int64(), moduleBalanceAfter.Int64())
			})
		})

		s.Run("expect that the alert still exists, even after commitHeight + purge-height", func() {
			s.Run("alert exists", func() {
				// query the current height, if less than max-block-age + alert.Height, check for alert's existence, and expect status to be concluded
				currentHeight, err := s.Height(context.Background())
				s.Require().NoError(err)

				// query for the alert (if it still should exist)
				if alertHeight+maxBlockAge > currentHeight {
					s.T().Log("alert should still exist", "alertHeight", alertHeight, "max-block-age", maxBlockAge, "currentHeight", currentHeight)
					alertResp, err := alertClient.Alerts(context.Background(), &alerttypes.AlertsRequest{
						Status: alerttypes.AlertStatusID_CONCLUSION_STATUS_CONCLUDED,
					})
					s.Require().NoError(err)

					// check the alerts
					s.Require().Len(alertResp.Alerts, 1)
					gotAlert := alertResp.Alerts[0]

					// expect equality (this was the only alert available)
					s.Require().Equal(alert, gotAlert)
				}
			})

			s.Run("submission of alert fails", func() {
				// query the current height, if less than max-block-age + alert.Height, attempt to submit an alert for the same height
				currentHeight, err := s.Height(context.Background())
				s.Require().NoError(err)

				// submit the same alert again
				if alertHeight+maxBlockAge > currentHeight {
					s.T().Log("alert should still exist", "alertHeight", alertHeight, "max-block-age", maxBlockAge, "currentHeight", currentHeight)

					alertSigner, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
					s.Require().NoError(err)

					badAlert := alerttypes.NewAlert(
						alert.Height,
						alertSigner,
						alert.CurrencyPair,
					)

					// submit the alert
					resp, err := s.SubmitAlert(s.multiSigUser1, badAlert)
					s.Require().NoError(err)

					// expect the alert submission to fail
					s.Require().Equal(resp.CheckTx.Code, uint32(0))
					s.Require().Equal(resp.TxResult.Code, uint32(1))

					s.Require().True(strings.Contains(resp.TxResult.Log, fmt.Sprintf("alert with UID %X already exists", alert.UID())))
				}
			})

			s.Run("submitting same alert now fails (alert is too old)", func() {
				// wait for original alert to be pruned
				ExpectAlerts(s.chain, s.blockTime*3, nil)

				// wait for the alert's height to pass max block-age
				WaitForHeight(s.chain, alertHeight+maxBlockAge, s.blockTime*10)

				// attempt to resubmit the alert, and expect it to fail
				resp, err := s.SubmitAlert(s.multiSigUser1, alert)
				s.Require().NoError(err)

				// expect the alert submission to fail
				s.Require().Equal(resp.CheckTx.Code, uint32(0))
				s.Require().Equal(resp.TxResult.Code, uint32(1))

				s.Require().True(strings.Contains(resp.TxResult.Log, "alert is too old"))
			})
		})
	})
}

// TestConclusionSubmission tests the conclusion submission process, specifically, this method tests
// the following:
// - conclusion submission fails if alerts are disabled
// - conclusion submission fails if the alert does not exist
// - conclusion submission fails if conclusion verification fails
// - conclusion submission fails if the alert has already been concluded
// - For negative conclusions
//   - conclusion submission passes, and there are no slashing events, and the alert signer is refunded
//
// - For positive conclusions
//   - conclusion submission passes, and there are slashing events, and the alert signer is rewarded the amount slashed
func (s *SlinkySlashingIntegrationSuite) TestConclusionSubmission() {
	// check if the BTC/USD currency pair exists
	cc, close, err := GetChainGRPC(s.chain)
	s.Require().NoError(err)
	defer close()

	cp := slinkytypes.NewCurrencyPair("BTC", "USD")

	// check if the currency pair exists, if not, add it
	oraclesClient := oracletypes.NewQueryClient(cc)
	ctx := context.Background()
	_, err = oraclesClient.GetPrice(ctx, &oracletypes.GetPriceRequest{
		CurrencyPair: cp,
	})
	if err != nil {
		// add the currency-pair
		s.Require().NoError(s.AddCurrencyPairs(s.chain, s.user, 1.1, cp))
	}

	// get the id for the currency-pair
	_, err = getIDForCurrencyPair(ctx, oraclesClient, cp)
	s.Require().NoError(err)

	s.Run("test Conclusion failures", func() {
		s.Run("fails when alerts are disabled", func() {
			// update the params to disable alerts
			_, err := UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, alerttypes.Params{
				ConclusionVerificationParams: nil,
				AlertParams: alerttypes.AlertParams{
					Enabled: false,
				},
			})
			s.Require().NoError(err)

			submitter, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
			s.Require().NoError(err)

			// submit a conclusion
			conclusion := &alerttypes.MultiSigConclusion{
				Alert: alerttypes.NewAlert(1, submitter, slinkytypes.NewCurrencyPair("BASE", "USDC")),
				PriceBound: alerttypes.PriceBound{
					High: "1",
					Low:  "0",
				},
				Signatures: []alerttypes.Signature{
					{
						Signer: submitter.String(),
					},
				},
			}

			resp, err := s.SubmitConclusion(s.multiSigUser1, conclusion)
			s.Require().NoError(err)

			// expect the alert submission to fail
			s.Require().Equal(resp.CheckTx.Code, uint32(0))
			s.Require().Equal(resp.TxResult.Code, uint32(1))

			s.Require().True(strings.Contains(resp.TxResult.Log, "alerts are not enabled"))
		})

		s.Run("fails when the conclusion verification fails", func() {
			cvp, err := alerttypes.NewMultiSigVerificationParams(
				[]cryptotypes.PubKey{
					s.multiSigPk1.PubKey(),
					s.multiSigPk2.PubKey(),
				},
			)
			s.Require().NoError(err)

			params := alerttypes.DefaultParams(s.denom, cvp)

			// update the params to enable alerts
			_, err = UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, params)
			s.Require().NoError(err)

			submitter, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
			s.Require().NoError(err)

			// submit a conclusion
			conclusion := &alerttypes.MultiSigConclusion{
				Alert: alerttypes.NewAlert(1, submitter, slinkytypes.NewCurrencyPair("BASE", "USDC")),
				PriceBound: alerttypes.PriceBound{
					High: "1",
					Low:  "0",
				},
			}

			// set only one signature
			sigBytes, err := conclusion.SignBytes()
			s.Require().NoError(err)

			// sign the bytes with multiSigPk1
			sig, err := s.multiSigPk1.Sign(sigBytes)
			s.Require().NoError(err)

			// set the signature
			conclusion.Signatures = []alerttypes.Signature{
				{
					Signer:    sdk.AccAddress(s.multiSigPk1.PubKey().Address()).String(),
					Signature: sig,
				},
			}

			resp, err := s.SubmitConclusion(s.multiSigUser1, conclusion)
			s.Require().NoError(err)

			// expect the alert submission to fail
			s.Require().Equal(resp.CheckTx.Code, uint32(0))
			s.Require().Equal(resp.TxResult.Code, uint32(1))

			s.Require().True(strings.Contains(resp.TxResult.Log, fmt.Sprintf("no signature provided for signer: %s", sdk.AccAddress(s.multiSigPk2.PubKey().Address()).String())))
		})

		s.Run("fails when the conclusion references a non-existent alert", func() {
			cvp, err := alerttypes.NewMultiSigVerificationParams(
				[]cryptotypes.PubKey{
					s.multiSigPk1.PubKey(),
					s.multiSigPk2.PubKey(),
				},
			)
			s.Require().NoError(err)

			params := alerttypes.DefaultParams(s.denom, cvp)

			// update the params to enable alerts
			_, err = UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, params)
			s.Require().NoError(err)

			submitter, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
			s.Require().NoError(err)

			alert := alerttypes.NewAlert(1, submitter, slinkytypes.NewCurrencyPair("BASE", "USDC"))
			// submit a conclusion
			conclusion := &alerttypes.MultiSigConclusion{
				Alert: alert,
				PriceBound: alerttypes.PriceBound{
					High: "1",
					Low:  "0",
				},
			}

			// set only one signature
			sigBytes, err := conclusion.SignBytes()
			s.Require().NoError(err)

			// sign the bytes with multiSigPk1
			sig, err := s.multiSigPk1.Sign(sigBytes)
			s.Require().NoError(err)

			// set the signature
			conclusion.Signatures = []alerttypes.Signature{
				{
					Signer:    sdk.AccAddress(s.multiSigPk1.PubKey().Address()).String(),
					Signature: sig,
				},
			}

			// set the second signature
			sig, err = s.multiSigPk2.Sign(sigBytes)
			s.Require().NoError(err)
			conclusion.Signatures = append(conclusion.Signatures, alerttypes.Signature{
				Signer:    sdk.AccAddress(s.multiSigPk2.PubKey().Address()).String(),
				Signature: sig,
			})

			resp, err := s.SubmitConclusion(s.multiSigUser1, conclusion)
			s.Require().NoError(err)

			// expect the alert submission to fail
			s.Require().Equal(resp.CheckTx.Code, uint32(0))
			s.Require().Equal(resp.TxResult.Code, uint32(1))

			s.Require().True(strings.Contains(resp.TxResult.Log, fmt.Sprintf("alert not found: %v", alert)))
		})

		s.Run("test negatively concluded alert", func() {
			cvp, err := alerttypes.NewMultiSigVerificationParams(
				[]cryptotypes.PubKey{
					s.multiSigPk1.PubKey(),
					s.multiSigPk2.PubKey(),
				},
			)
			s.Require().NoError(err)

			params := alerttypes.DefaultParams(s.denom, cvp)

			// update the params to enable alerts
			_, err = UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, params)
			s.Require().NoError(err)

			submitter, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
			s.Require().NoError(err)

			alert := alerttypes.NewAlert(1, submitter, slinkytypes.NewCurrencyPair("BTC", "USD"))

			// get the balance of the sender / module to check balance differences for escrow
			senderBalanceBeforeAlert, err := s.chain.GetBalance(context.Background(), s.multiSigUser1.FormattedAddress(), s.denom)
			s.Require().NoError(err)

			moduleBalanceBeforeAlert, err := s.chain.GetBalance(context.Background(), authtypes.NewModuleAddress(alerttypes.ModuleName).String(), s.denom)
			s.Require().NoError(err)

			// submit the alert
			resp, err := s.SubmitAlert(s.multiSigUser1, alert)
			s.Require().NoError(err)

			// expect the alert submission to succeed
			s.Require().Equal(resp.CheckTx.Code, uint32(0))
			s.Require().Equal(resp.TxResult.Code, uint32(0))

			// expect bond to be escrowed
			senderBalanceAfterAlert, err := s.chain.GetBalance(context.Background(), s.multiSigUser1.FormattedAddress(), s.denom)
			s.Require().NoError(err)

			moduleBalanceAfterAlert, err := s.chain.GetBalance(context.Background(), authtypes.NewModuleAddress(alerttypes.ModuleName).String(), s.denom)
			s.Require().NoError(err)

			s.Require().Equal(senderBalanceBeforeAlert.Sub(senderBalanceAfterAlert).Int64(), alerttypes.DefaultBondAmount.Int64()+resp.TxResult.GasWanted*gasPrice)
			s.Require().Equal(moduleBalanceAfterAlert.Sub(moduleBalanceBeforeAlert).Int64(), alerttypes.DefaultBondAmount.Int64())

			validatorsBeforeConclusion, err := QueryValidators(s.chain)
			s.Require().NoError(err)

			conclusion := &alerttypes.MultiSigConclusion{
				Alert: alert,
				PriceBound: alerttypes.PriceBound{
					High: "1",
					Low:  "0",
				},
			}

			// set only one signature
			sigBytes, err := conclusion.SignBytes()
			s.Require().NoError(err)

			// sign the bytes with multiSigPk1
			sig, err := s.multiSigPk1.Sign(sigBytes)
			s.Require().NoError(err)

			// set the signature
			conclusion.Signatures = []alerttypes.Signature{
				{
					Signer:    sdk.AccAddress(s.multiSigPk1.PubKey().Address()).String(),
					Signature: sig,
				},
			}

			// set the second signature
			sig, err = s.multiSigPk2.Sign(sigBytes)
			s.Require().NoError(err)
			conclusion.Signatures = append(conclusion.Signatures, alerttypes.Signature{
				Signer:    sdk.AccAddress(s.multiSigPk2.PubKey().Address()).String(),
				Signature: sig,
			})

			resp, err = s.SubmitConclusion(s.multiSigUser2, conclusion)
			s.Require().NoError(err)

			// expect the conclusion submission to succeed
			s.Require().Equal(resp.CheckTx.Code, uint32(0))
			s.Require().Equal(resp.TxResult.Code, uint32(0))

			// expect bond to have been burned
			senderBalanceFinal, err := s.chain.GetBalance(context.Background(), s.multiSigUser1.FormattedAddress(), s.denom)
			s.Require().NoError(err)

			moduleBalanceFinal, err := s.chain.GetBalance(context.Background(), authtypes.NewModuleAddress(alerttypes.ModuleName).String(), s.denom)
			s.Require().NoError(err)

			// sender balance is the same
			s.Require().Equal(senderBalanceAfterAlert.Uint64(), senderBalanceFinal.Uint64())

			// module balance is decremented by the bond amount
			s.Require().Equal(moduleBalanceFinal.Int64(), moduleBalanceAfterAlert.Sub(alerttypes.DefaultBondAmount).Int64())

			// validators after conclusion
			validatorsAfterConclusion, err := QueryValidators(s.chain)
			s.Require().NoError(err)

			// expect the validator's bond to not have changed
			for i, val := range validatorsBeforeConclusion {
				s.Require().True(val.Tokens.Equal(validatorsAfterConclusion[i].Tokens))
			}
		})

		s.Run("fails when the alert is alr concluded", func() {
			submitter, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
			s.Require().NoError(err)

			alert := alerttypes.NewAlert(1, submitter, slinkytypes.NewCurrencyPair("BTC", "USD"))

			// submit a conclusion
			conclusion := &alerttypes.MultiSigConclusion{
				Alert: alert,
				PriceBound: alerttypes.PriceBound{
					High: "1",
					Low:  "0",
				},
			}

			// set only one signature
			sigBytes, err := conclusion.SignBytes()
			s.Require().NoError(err)

			// sign the bytes with multiSigPk1
			sig, err := s.multiSigPk1.Sign(sigBytes)
			s.Require().NoError(err)

			// set the signature
			conclusion.Signatures = []alerttypes.Signature{
				{
					Signer:    sdk.AccAddress(s.multiSigPk1.PubKey().Address()).String(),
					Signature: sig,
				},
			}

			// set the second signature
			sig, err = s.multiSigPk2.Sign(sigBytes)
			s.Require().NoError(err)
			conclusion.Signatures = append(conclusion.Signatures, alerttypes.Signature{
				Signer:    sdk.AccAddress(s.multiSigPk2.PubKey().Address()).String(),
				Signature: sig,
			})

			resp, err := s.SubmitConclusion(s.multiSigUser1, conclusion)
			s.Require().NoError(err)

			// expect the conclusion submission to succeed
			s.Require().Equal(resp.CheckTx.Code, uint32(0))
			s.Require().Equal(resp.TxResult.Code, uint32(1))

			s.Require().True(strings.Contains(resp.TxResult.Log, "alert already concluded"))
		})
	})

	s.Run("valid conclusion submissions", func() {
		// TODO restore once slashing is more finalized

		/*
			var honestPrice int64 = 150

			s.Run("update validator oracles", func() {
				// update first validator's oracle to submit incorrect Prices
				nodes := s.chain.Nodes()

				btcusdTicker := mmtypes.Ticker{
					CurrencyPair:     cp,
					Decimals:         8,
					MinProviderCount: 1,
					Enabled:          true,
				}

				// update the first node to report incorrect Prices (too high)
				s.Require().NoError(UpdateNodePrices(nodes[0], btcusdTicker, 152))

				// update the second node to report incorrect Prices (too low)
				s.Require().NoError(UpdateNodePrices(nodes[1], btcusdTicker, 148))

				// update the third node to report correct Prices
				s.Require().NoError(UpdateNodePrices(nodes[2], btcusdTicker, float64(honestPrice)))

				// update the fourth node to report correct Prices
				s.Require().NoError(UpdateNodePrices(nodes[3], btcusdTicker, float64(honestPrice)))
			})

			validatorsPreSlash, err := QueryValidators(s.chain)
			s.Require().NoError(err)

			cdc := s.chain.Config().EncodingConfig.Codec
			validatorPreSlashMap := mapValidators(validatorsPreSlash, cdc)

			zero := big.NewInt(0)
			two := big.NewInt(2)
			negativeTwo := big.NewInt(-2)

			zeroBz, err := zero.GobEncode()
			s.Require().NoError(err)

			twoBz, err := two.GobEncode()
			s.Require().NoError(err)

			negativeTwoBz, err := negativeTwo.GobEncode()
			s.Require().NoError(err)

			infractionHeight, err := ExpectVoteExtensions(s.chain, s.blockTime*10, []slinkyabci.OracleVoteExtension{
				{
					Prices: map[uint64][]byte{
						id: negativeTwoBz, // 148
					},
				},
				{
					Prices: map[uint64][]byte{
						id: zeroBz, // 150
					},
				},
				{
					Prices: map[uint64][]byte{
						id: zeroBz, // 150
					},
				},
				{
					Prices: map[uint64][]byte{
						id: twoBz, // 152
					},
				},
			})
			s.Require().NoError(err)

			// get the latest extended commit info
			extendedCommit, err := GetExtendedCommit(s.chain, int64(infractionHeight))
			s.Require().NoError(err)

			valsToOracleReport := make(map[string]int64)

			s.Run("map validators to their oracle responses", func() {
				// map validators to their oracle responses
				for _, vote := range extendedCommit.Votes {
					oracleData, err := GetOracleDataFromVote(vote)
					s.Require().NoError(err)

					key := sdk.ConsAddress(vote.Validator.Address).String()

					// get the price from the oracle data
					priceBz, ok := oracleData.Prices[id]
					s.Require().True(ok)

					// get the big from string value
					var price big.Int
					price.SetBytes(priceBz)

					// convert the big to int64
					valsToOracleReport[key] = int64(price.Uint64())
				}
			})

			alertSigner, err := sdk.AccAddressFromBech32(s.multiSigUser1.FormattedAddress())
			s.Require().NoError(err)

			s.Run("update params to enable alerts + conclusion-verification-params", func() {
				cvp, err := alerttypes.NewMultiSigVerificationParams(
					[]cryptotypes.PubKey{
						s.multiSigPk1.PubKey(),
						s.multiSigPk2.PubKey(),
					},
				)
				s.Require().NoError(err)

				params := alerttypes.DefaultParams(s.denom, cvp)

				// update the params to enable alerts
				_, err = UpdateAlertParams(s.chain, s.authority.String(), s.denom, deposit, 2*s.blockTime, s.multiSigUser1, params)
				s.Require().NoError(err)
			})

			var alertSignerBalance math.Int

			s.Run("submit alert + conclusion", func() {
				// submit an alert for the infraction height
				alert := alerttypes.NewAlert(infractionHeight, alertSigner, cp)

				// submit the alert
				resp, err := s.SubmitAlert(s.multiSigUser1, alert)
				s.Require().NoError(err)

				// expect the alert submission to succeed
				s.Require().Equal(resp.CheckTx.Code, uint32(0))
				s.Require().Equal(resp.TxResult.Code, uint32(0))

				alertSignerBalance, err = s.chain.GetBalance(context.Background(), alertSigner.String(), s.denom)
				s.Require().NoError(err)

				// create the conclusion
				conclusion := &alerttypes.MultiSigConclusion{
					Alert:              alert,
					ExtendedCommitInfo: extendedCommit,
					PriceBound: alerttypes.PriceBound{
						High: "151",
						Low:  "149",
					},
					Status:         true,
					CurrencyPairID: id,
				}

				sigBytes, err := conclusion.SignBytes()
				s.Require().NoError(err)

				// sign from the first multi-sig key
				sig, err := s.multiSigPk1.Sign(sigBytes)
				s.Require().NoError(err)

				conclusion.Signatures = []alerttypes.Signature{
					{
						Signer:    sdk.AccAddress(s.multiSigPk1.PubKey().Address()).String(),
						Signature: sig,
					},
				}

				// sign from second multi-sig key
				sig, err = s.multiSigPk2.Sign(sigBytes)
				s.Require().NoError(err)

				conclusion.Signatures = append(conclusion.Signatures, alerttypes.Signature{
					Signer:    sdk.AccAddress(s.multiSigPk2.PubKey().Address()).String(),
					Signature: sig,
				})

				// submit the conclusion
				resp, err = s.SubmitConclusion(s.multiSigUser2, conclusion)
				s.Require().NoError(err)

				// expect the conclusion submission to succeed
				s.Require().Equal(resp.CheckTx.Code, uint32(0))
				s.Require().Equal(resp.TxResult.Code, uint32(0))

				// wait for a block for the incentive to be executed
				WaitForHeight(s.chain, uint64(resp.Height)+2, 4*s.blockTime)
			})

			s.Run("expect that slashing / rewarding is executed", func() {
				// get the validators after the slashing
				validatorsPostSlash, err := QueryValidators(s.chain)
				s.Require().NoError(err)

				cdc := s.chain.Config().EncodingConfig.Codec
				validatorPostSlashMap := mapValidators(validatorsPostSlash, cdc)

				reward := math.NewInt(0)
				// check that the validators are slashed / rewarded correctly
				for consAddr, priceReport := range valsToOracleReport {
					// if the validator's report was honest, expect no slash
					if priceReport == honestPrice {
						preSlashValidator, ok := validatorPreSlashMap[consAddr]
						s.Require().True(ok)

						postSlashValidator, ok := validatorPostSlashMap[consAddr]
						s.Require().True(ok)

						s.Require().True(preSlashValidator.Tokens.Equal(postSlashValidator.Tokens))
						continue
					}

					// otherwise expect a slash
					preSlashValidator, ok := validatorPreSlashMap[consAddr]
					s.Require().True(ok)

					postSlashValidator, ok := validatorPostSlashMap[consAddr]
					s.Require().True(ok)

					s.Require().True(postSlashValidator.Tokens.LT(preSlashValidator.Tokens))
					reward = reward.Add(preSlashValidator.Tokens.Sub(postSlashValidator.Tokens))
				}

				// check that the alert signer is rewarded correctly
				alertSignerBalanceAfter, err := s.chain.GetBalance(context.Background(), alertSigner.String(), s.denom)
				s.Require().NoError(err)

				s.Require().True(alertSignerBalanceAfter.Equal(alertSignerBalance.Add(reward).Add(alerttypes.DefaultBondAmount)))
			})
		*/
	})
}

func mapValidators(vals []stakingtypes.Validator, cdc codec.Codec) map[string]stakingtypes.Validator {
	m := make(map[string]stakingtypes.Validator)
	for _, v := range vals {
		key, err := pkToKey(v.ConsensusPubkey, cdc)
		if err != nil {
			continue
		}
		m[key] = v
	}
	return m
}

func pkToKey(pkAny *codectypes.Any, cdc codec.Codec) (string, error) {
	protoCodec, ok := cdc.(*codec.ProtoCodec)
	if !ok {
		return "", fmt.Errorf("expected codec to be a proto codec")
	}

	var pk cryptotypes.PubKey

	if err := protoCodec.UnpackAny(pkAny, &pk); err != nil {
		return "", err
	}

	return sdk.ConsAddress(pk.Address()).String(), nil
}
