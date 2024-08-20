package keeper_test

import (
	"fmt"
	"math/big"
	"time"

	"github.com/stretchr/testify/mock"

	"cosmossdk.io/math"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	slinkyabci "github.com/skip-mev/connect/v2/abci/ve/types"
	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/alerts/keeper"
	"github.com/skip-mev/connect/v2/x/alerts/types"
	"github.com/skip-mev/connect/v2/x/alerts/types/strategies"
	incentivetypes "github.com/skip-mev/connect/v2/x/incentives/types"
)

func (s *KeeperTestSuite) TestMsgAlert() {
	type testCase struct {
		setup func(sdk.Context)
		name  string
		msg   *types.MsgAlert
		valid bool
	}

	validAlert := types.NewAlert(8, sdk.AccAddress("abc"), slinkytypes.NewCurrencyPair("BTC", "USD"))

	s.ctx = s.ctx.WithBlockHeight(10)
	s.ctx = s.ctx.WithBlockTime(time.Now())

	testCases := []testCase{
		{
			name:  "nil message - fail",
			setup: func(_ sdk.Context) {},
			msg:   nil,
			valid: false,
		},
		{
			name:  "invalid message - fail",
			setup: func(_ sdk.Context) {},
			msg: &types.MsgAlert{
				Alert: types.Alert{
					Height:       1,
					Signer:       "",
					CurrencyPair: slinkytypes.NewCurrencyPair("base", "quote"),
				},
			},
			valid: false,
		},
		{
			name: "alerts disabled - fail",
			setup: func(ctx sdk.Context) {
				err := s.alertKeeper.SetParams(ctx, types.Params{
					AlertParams: types.AlertParams{
						Enabled: false,
					},
				})
				s.Require().NoError(err)
			},
			msg: &types.MsgAlert{
				Alert: types.NewAlert(1, sdk.AccAddress("abc"), slinkytypes.NewCurrencyPair("base", "quote")),
			},
			valid: false,
		},
		{
			name: "alert is too old",
			setup: func(ctx sdk.Context) {
				// ensure that alerts are enabled
				s.Require().NoError(s.alertKeeper.SetParams(ctx, types.Params{
					AlertParams: types.AlertParams{
						Enabled:     true,
						BondAmount:  sdk.NewInt64Coin("test", 100),
						MaxBlockAge: 2,
					},
				}))
			},
			msg: &types.MsgAlert{
				Alert: types.NewAlert(7, sdk.AccAddress("abc"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
			},
			valid: false,
		},
		{
			name: "alert already exists - fail",
			setup: func(ctx sdk.Context) {
				// ensure that alerts are enabled
				s.Require().NoError(s.alertKeeper.SetParams(ctx, types.Params{
					AlertParams: types.AlertParams{
						Enabled:     true,
						BondAmount:  sdk.NewInt64Coin("test", 100),
						MaxBlockAge: 2,
					},
				}))

				// set the alert to state
				alert := types.NewAlertWithStatus(
					types.NewAlert(9, sdk.AccAddress("abc1"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
					types.NewAlertStatus(9, 11, s.ctx.BlockTime(), types.Unconcluded),
				)
				s.Require().NoError(s.alertKeeper.SetAlert(ctx, alert))
			},
			msg: &types.MsgAlert{
				Alert: types.NewAlert(9, sdk.AccAddress("abc"), slinkytypes.NewCurrencyPair("BASE", "QUOTE")),
			},
			valid: false,
		},
		{
			name: "currency pair does not exist - fail",
			setup: func(ctx sdk.Context) {
				// ensure that alerts are enabled
				s.Require().NoError(s.alertKeeper.SetParams(ctx, types.Params{
					AlertParams: types.AlertParams{
						Enabled:     true,
						BondAmount:  sdk.NewInt64Coin("test", 100),
						MaxBlockAge: 2,
					},
				}))

				// expect a failed response from the oracle keeper (no currency pair)
				s.ok.On("HasCurrencyPair", mock.Anything, slinkytypes.NewCurrencyPair("BTC", "USD")).Return(false).Once()
			},
			msg: &types.MsgAlert{
				Alert: types.NewAlert(8, sdk.AccAddress("abc"), slinkytypes.NewCurrencyPair("BTC", "USD")),
			},
			valid: false,
		},
		{
			name: "bond amount cannot be escrowed - fail",
			setup: func(ctx sdk.Context) {
				// ensure that alerts are enabled
				s.Require().NoError(s.alertKeeper.SetParams(ctx, types.Params{
					AlertParams: types.AlertParams{
						Enabled:     true,
						BondAmount:  sdk.NewInt64Coin("test", 100),
						MaxBlockAge: 2,
					},
				}))

				// expect a correct response from the oracle keeper
				s.ok.On("HasCurrencyPair", mock.Anything, slinkytypes.NewCurrencyPair("BTC", "USD")).Return(true).Once()

				// expect a failed response from the bank keeper
				s.bk.On("SendCoinsFromAccountToModule",
					mock.Anything,
					sdk.AccAddress("abc"),
					types.ModuleName,
					sdk.NewCoins(s.alertKeeper.GetParams(s.ctx).AlertParams.BondAmount),
				).Return(fmt.Errorf("bank error")).Once()
			},
			msg: &types.MsgAlert{
				Alert: types.NewAlert(8, sdk.AccAddress("abc"), slinkytypes.NewCurrencyPair("BTC", "USD")),
			},
			valid: false,
		},
		{
			name: "valid message - success",
			setup: func(ctx sdk.Context) {
				// ensure that alerts are enabled
				s.Require().NoError(s.alertKeeper.SetParams(ctx, types.Params{
					AlertParams: types.AlertParams{
						Enabled:     true,
						BondAmount:  sdk.NewInt64Coin("test", 100),
						MaxBlockAge: 2,
					},
					PruningParams: types.PruningParams{
						BlocksToPrune: 10,
					},
				}))

				// expect a correct response from the oracle keeper
				s.ok.On("HasCurrencyPair",
					mock.Anything,
					slinkytypes.NewCurrencyPair("BTC", "USD"),
				).Return(true).Once()

				// expect a correct response from the bank keeper
				s.bk.On("SendCoinsFromAccountToModule",
					mock.Anything,
					sdk.AccAddress("abc"),
					types.ModuleName,
					sdk.NewCoins(s.alertKeeper.GetParams(s.ctx).AlertParams.BondAmount),
				).Return(nil).Once()
			},
			msg: &types.MsgAlert{
				Alert: validAlert,
			},
			valid: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// perform setup for test-case
			tc.setup(s.ctx)
			ms := keeper.NewMsgServer(*s.alertKeeper)

			// run the message server
			_, err := ms.Alert(s.ctx, tc.msg)
			if tc.valid {
				s.Require().NoError(err)

				// check that the alert was added to the state
				alert, ok := s.alertKeeper.GetAlert(s.ctx, tc.msg.Alert)
				s.Require().Equal(tc.msg.Alert, alert.Alert)
				s.Require().Equal(alert.Status, types.AlertStatus{
					PurgeHeight:         uint64(20),
					SubmissionHeight:    10,
					ConclusionStatus:    uint64(types.Unconcluded),
					SubmissionTimestamp: uint64(time.Now().UTC().Unix()),
				})
				s.Require().True(ok)
			} else {
				s.Require().Error(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestConclusion() {
	msgServer := keeper.NewMsgServer(*s.alertKeeper)
	ctx := s.ctx

	alert := types.Alert{
		Height:       1,
		Signer:       sdk.AccAddress("cosmos1").String(),
		CurrencyPair: slinkytypes.NewCurrencyPair("BTC", "USD"),
	}

	conclusion := &types.MultiSigConclusion{
		ExtendedCommitInfo: cmtabci.ExtendedCommitInfo{},
		Alert:              alert,
		PriceBound: types.PriceBound{
			High: big.NewInt(1).String(),
			Low:  big.NewInt(0).String(),
		},
		Signatures: make([]types.Signature, 0),
		Status:     false,
	}

	// sign the conclusion
	signBytes, err := conclusion.SignBytes()
	s.Require().NoError(err)

	// set the signature
	sig, err := s.privateKey.Sign(signBytes)
	conclusion.Signatures = append(conclusion.Signatures, types.Signature{
		Signer:    sdk.AccAddress(s.privateKey.PubKey().Address()).String(),
		Signature: sig,
	})
	s.Require().NoError(err)

	validConclusionAny, err := codectypes.NewAnyWithValue(conclusion)
	s.Require().NoError(err)

	conclusionVerificationParams, err := types.NewMultiSigVerificationParams([]cryptotypes.PubKey{s.privateKey.PubKey()})
	s.Require().NoError(err)

	// set as any
	verificationParams, err := codectypes.NewAnyWithValue(conclusionVerificationParams)
	s.Require().NoError(err)

	s.Run("if the msg is nil - Conclusion fails", func() {
		_, err := msgServer.Conclusion(ctx, nil)
		s.Require().Error(err)
		s.Require().Equal("message cannot be empty", err.Error())
	})

	s.Run("if the conclusion fails validate basic", func() {
		msg := &types.MsgConclusion{
			Signer: "",
		}
		_, err := msgServer.Conclusion(ctx, msg)

		s.Require().Error(err)
		s.Require().Equal(fmt.Errorf("message validation failed: %w", msg.ValidateBasic()).Error(), err.Error())
	})

	s.Run("if alerts are not enabled", func() {
		// set Alerts as disabled in Params
		s.Require().NoError(s.alertKeeper.SetParams(ctx, types.Params{
			AlertParams: types.AlertParams{
				Enabled: false,
			},
		}))

		// msg should pass validate basic
		msg := &types.MsgConclusion{
			Signer:     sdk.AccAddress("cosmos1").String(),
			Conclusion: validConclusionAny,
		}

		// Conclusion should fail
		_, err := msgServer.Conclusion(ctx, msg)
		s.Require().Error(err)
		s.Require().Equal("alerts are not enabled", err.Error())
	})

	s.Run("if the conclusion fails in verification", func() {
		pk := secp256k1.GenPrivKey().PubKey()
		pkany, err := codectypes.NewAnyWithValue(pk)
		s.Require().NoError(err)

		invalidVerificationParams, err := codectypes.NewAnyWithValue(&types.MultiSigConclusionVerificationParams{
			Signers: []*codectypes.Any{pkany},
		})
		s.Require().NoError(err)

		s.Require().NoError(s.alertKeeper.SetParams(ctx, types.Params{
			AlertParams: types.AlertParams{
				Enabled: true,
			},
			ConclusionVerificationParams: invalidVerificationParams,
		}))

		// msg should pass validate basic
		msg := &types.MsgConclusion{
			Signer:     sdk.AccAddress("cosmos1").String(),
			Conclusion: validConclusionAny,
		}

		// Conclusion should fail
		_, err = msgServer.Conclusion(ctx, msg)
		s.Require().Error(err)

		s.Require().Equal(fmt.Errorf("failed to verify conclusion: %w", conclusion.Verify(&types.MultiSigConclusionVerificationParams{
			Signers: []*codectypes.Any{pkany},
		})).Error(), err.Error())
	})

	s.Run("if the alert is not in state - fail", func() {
		s.Require().NoError(s.alertKeeper.SetParams(ctx, types.Params{
			AlertParams: types.AlertParams{
				Enabled: true,
			},
			ConclusionVerificationParams: verificationParams,
		}))

		// msg should pass validate basic
		msg := &types.MsgConclusion{
			Signer:     sdk.AccAddress("cosmos1").String(),
			Conclusion: validConclusionAny,
		}

		// Conclusion should fail
		_, err = msgServer.Conclusion(ctx, msg)
		s.Require().Error(err)

		s.Require().Equal(fmt.Errorf("failed to conclude alert: alert not found: %v", alert).Error(), err.Error())
	})

	s.Run("if the alert has already been concluded", func() {
		s.Require().NoError(s.alertKeeper.SetAlert(
			ctx,
			types.NewAlertWithStatus(
				alert,
				types.AlertStatus{
					PurgeHeight:      uint64(20),
					SubmissionHeight: 10,
					ConclusionStatus: uint64(types.Concluded),
				},
			),
		))

		// msg should pass validate basic
		msg := &types.MsgConclusion{
			Signer:     sdk.AccAddress("cosmos1").String(),
			Conclusion: validConclusionAny,
		}

		// Conclusion should fail
		_, err = msgServer.Conclusion(ctx, msg)
		s.Require().Error(err)

		s.Require().Equal(
			fmt.Errorf("failed to conclude alert: alert already concluded").Error(),
			err.Error(),
		)
	})

	s.Run("if the alert is concluded negatively - expect the bond to be burned", func() {
		s.Require().NoError(s.alertKeeper.SetParams(ctx, types.Params{
			AlertParams: types.AlertParams{
				Enabled:     true,
				MaxBlockAge: 10,
				BondAmount: sdk.NewCoin(
					sdk.DefaultBondDenom,
					math.NewInt(100),
				),
			},
			ConclusionVerificationParams: verificationParams,
		}))

		s.Require().NoError(s.alertKeeper.SetAlert(
			ctx,
			types.NewAlertWithStatus(
				alert,
				types.AlertStatus{
					PurgeHeight:      uint64(11),
					SubmissionHeight: 10,
					ConclusionStatus: uint64(types.Unconcluded),
				},
			),
		))

		s.bk.On("BurnCoins",
			mock.Anything,
			types.ModuleName,
			sdk.NewCoins(s.alertKeeper.GetParams(s.ctx).AlertParams.BondAmount),
		).Return(nil)

		// msg should pass validate basic
		msg := &types.MsgConclusion{
			Signer:     sdk.AccAddress("cosmos1").String(),
			Conclusion: validConclusionAny,
		}

		// Conclusion should fail
		_, err = msgServer.Conclusion(ctx, msg)
		s.Require().NoError(err)

		// get the alert
		alert, found := s.alertKeeper.GetAlert(ctx, alert)
		s.Require().True(found)

		// check the alert status
		s.Require().Equal(
			types.AlertStatus{
				PurgeHeight:      uint64(11),
				SubmissionHeight: 10,
				ConclusionStatus: uint64(types.Concluded),
			},
			alert.Status,
		)
	})

	s.Run("if the alert is concluded positively - expect the bond to be returned, and issue incentives", func() {
		s.Require().NoError(s.alertKeeper.SetAlert(
			ctx,
			types.NewAlertWithStatus(
				alert,
				types.AlertStatus{
					PurgeHeight:      uint64(11),
					SubmissionHeight: 10,
					ConclusionStatus: uint64(types.Unconcluded),
				},
			),
		))

		pb := types.PriceBound{
			High: big.NewInt(100).String(),
			Low:  big.NewInt(90).String(),
		}

		// create 3 validators
		val1 := cmtabci.Validator{
			Address: sdk.ConsAddress("val1"),
			Power:   10,
		}
		val2 := cmtabci.Validator{
			Address: sdk.ConsAddress("val2"),
			Power:   10,
		}
		val3 := cmtabci.Validator{
			Address: sdk.ConsAddress("val3"),
			Power:   10,
		}

		// val1 is not within the price-bound
		val1VE := slinkyabci.OracleVoteExtension{
			Prices: map[uint64][]byte{
				0: big.NewInt(101).Bytes(),
			},
		}
		val1VEbz, err := val1VE.Marshal()
		s.Require().NoError(err)

		// val2 is within the price-bound
		val2VE := slinkyabci.OracleVoteExtension{
			Prices: map[uint64][]byte{
				0: big.NewInt(99).Bytes(),
			},
		}
		val2VEbz, err := val2VE.Marshal()
		s.Require().NoError(err)

		// val3 is not within the price-bound
		val3VE := slinkyabci.OracleVoteExtension{
			Prices: map[uint64][]byte{
				0: big.NewInt(89).Bytes(),
			},
		}
		val3VEbz, err := val3VE.Marshal()
		s.Require().NoError(err)

		// construct extended commit
		commit := cmtabci.ExtendedCommitInfo{
			Votes: []cmtabci.ExtendedVoteInfo{
				{
					Validator:     val1,
					VoteExtension: val1VEbz,
				},
				{
					Validator:     val2,
					VoteExtension: val2VEbz,
				},
				{
					Validator:     val3,
					VoteExtension: val3VEbz,
				},
			},
		}

		// create conclusion
		conclusion := types.MultiSigConclusion{
			ExtendedCommitInfo: commit,
			PriceBound:         pb,
			Alert:              alert,
			Signatures:         make([]types.Signature, 0),
			Status:             true,
			CurrencyPairID:     0,
		}
		signBz, err := conclusion.SignBytes()
		s.Require().NoError(err)

		// sign the conclusion
		signature, err := s.privateKey.Sign(signBz)
		s.Require().NoError(err)

		// add the signature
		conclusion.Signatures = append(conclusion.Signatures, types.Signature{
			Signer:    sdk.AccAddress(s.privateKey.PubKey().Address()).String(),
			Signature: signature,
		})

		conclusionAny, err := codectypes.NewAnyWithValue(&conclusion)
		s.Require().NoError(err)

		s.bk.On("SendCoinsFromModuleToAccount",
			mock.Anything,
			types.ModuleName,
			sdk.AccAddress("cosmos1"),
			sdk.NewCoins(s.alertKeeper.GetParams(s.ctx).AlertParams.BondAmount),
		).Return(nil)

		s.ik.On("AddIncentives",
			mock.Anything,
			[]incentivetypes.Incentive{
				&strategies.ValidatorAlertIncentive{
					Validator:   val1,
					AlertSigner: sdk.AccAddress("cosmos1").String(),
					AlertHeight: uint64(1),
				},
				&strategies.ValidatorAlertIncentive{
					Validator:   val3,
					AlertSigner: sdk.AccAddress("cosmos1").String(),
					AlertHeight: uint64(1),
				},
			},
		).Return(nil)

		// msg should pass validate basic
		msg := &types.MsgConclusion{
			Signer:     sdk.AccAddress("cosmos1").String(),
			Conclusion: conclusionAny,
		}

		// Conclusion should fail
		_, err = msgServer.Conclusion(ctx, msg)
		s.Require().NoError(err)

		// get the alert
		alert, found := s.alertKeeper.GetAlert(ctx, alert)
		s.Require().True(found)

		// check the alert status
		s.Require().Equal(
			types.AlertStatus{
				PurgeHeight:      uint64(11),
				SubmissionHeight: 10,
				ConclusionStatus: uint64(types.Concluded),
			},
			alert.Status,
		)
	})
}

func (s *KeeperTestSuite) TestUpdateParams() {
	invalidParams := types.MsgUpdateParams{
		Authority: "invalid",
	}

	cases := []struct {
		name      string
		msg       *types.MsgUpdateParams
		expectErr error
	}{
		{
			"nil request - fail",
			nil,
			fmt.Errorf("message cannot be empty"),
		},
		{
			"invalid message - fail",
			&types.MsgUpdateParams{
				Authority: "invalid",
			},
			fmt.Errorf("message validation failed: %w", invalidParams.ValidateBasic()),
		},
		{
			"signer is not the authority - fail",
			&types.MsgUpdateParams{
				Authority: sdk.AccAddress("cosmos1").String(),
				Params:    types.DefaultParams("new_denom", nil),
			},
			fmt.Errorf("signer is not the authority of this module: signer %v, authority %v", sdk.AccAddress("cosmos1").String(), s.authority.String()),
		},
		{
			"valid message - success",
			&types.MsgUpdateParams{
				Authority: s.authority.String(),
				Params: types.NewParams(
					types.AlertParams{
						Enabled:     true,
						BondAmount:  sdk.NewCoin("denom", math.NewInt(100)),
						MaxBlockAge: 10,
					},
					nil,
					types.PruningParams{
						Enabled:       true,
						BlocksToPrune: 10,
					},
				),
			},
			nil,
		},
	}

	msgServer := keeper.NewMsgServer(*s.alertKeeper)

	for _, tc := range cases {
		params := types.DefaultParams("denom", nil)
		s.Require().NoError(s.alertKeeper.SetParams(s.ctx, params))
		_, err := msgServer.UpdateParams(s.ctx, tc.msg)
		if tc.expectErr == nil {
			s.Require().NoError(err)

			// check params in module
			params := s.alertKeeper.GetParams(s.ctx)
			s.Require().Equal(tc.msg.Params, params)
		} else {
			s.Require().Error(err)
			s.Require().Equal(tc.expectErr.Error(), err.Error())

			// check params in module
			moduleParams := s.alertKeeper.GetParams(s.ctx)
			s.Require().Equal(params, moduleParams)
		}
	}
}
