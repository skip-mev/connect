package strategies_test

import (
	"fmt"
	"math/big"
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"

	slinkyabci "github.com/skip-mev/connect/v2/abci/ve/types"
	slinkytypes "github.com/skip-mev/connect/v2/pkg/types"
	"github.com/skip-mev/connect/v2/x/alerts/types"
	"github.com/skip-mev/connect/v2/x/alerts/types/mocks"
	"github.com/skip-mev/connect/v2/x/alerts/types/strategies"
	incentivetypes "github.com/skip-mev/connect/v2/x/incentives/types"
	"github.com/skip-mev/connect/v2/x/incentives/types/examples/goodprice"
)

func TestValidatorAlertIncentive(t *testing.T) {
	t.Run("test validate basic", func(t *testing.T) {
		cases := []struct {
			name               string
			validatorIncentive incentivetypes.Incentive
			valid              bool
		}{
			{
				"nil address",
				strategies.NewValidatorAlertIncentive(cmtabci.Validator{
					Address: nil,
					Power:   1,
				}, 1, sdk.AccAddress("test")),
				false,
			},
			{
				"negative power",
				strategies.NewValidatorAlertIncentive(cmtabci.Validator{
					Address: []byte("test"),
					Power:   -1,
				}, 1, sdk.AccAddress("test")),
				false,
			},
			{
				"zero power",
				strategies.NewValidatorAlertIncentive(cmtabci.Validator{
					Address: []byte("test"),
					Power:   0,
				}, 1, sdk.AccAddress("test")),
				false,
			},
			{
				"invalid acc-address string - fail",
				&strategies.ValidatorAlertIncentive{
					Validator: cmtabci.Validator{
						Address: []byte("test"),
						Power:   1,
					},
					AlertHeight: 1,
					AlertSigner: "",
				},
				false,
			},
			{
				"valid",
				strategies.NewValidatorAlertIncentive(cmtabci.Validator{
					Address: []byte("test"),
					Power:   1,
				}, 1, sdk.AccAddress("test")),
				true,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.validatorIncentive.ValidateBasic()

				if tc.valid && err != nil {
					t.Errorf("expected no error, got %s", err)
				}

				if !tc.valid && err == nil {
					t.Errorf("expected error, got nil")
				}
			})
		}
	})

	t.Run("test type", func(t *testing.T) {
		ic := strategies.NewValidatorAlertIncentive(cmtabci.Validator{
			Address: []byte("test"),
			Power:   1,
		}, 1, sdk.AccAddress("test"))

		require.Equal(t, ic.Type(), strategies.ValidatorAlertIncentiveType)
	})

	t.Run("test copy", func(t *testing.T) {
		ic := strategies.NewValidatorAlertIncentive(cmtabci.Validator{
			Address: []byte("test"),
			Power:   1,
		}, 1, sdk.AccAddress("test"))

		icCopy := ic.Copy()

		require.Equal(t, ic, icCopy)
		require.False(t, &ic == &icCopy)

		// assert addresses are diff
		addr1 := ic.(*strategies.ValidatorAlertIncentive).Validator.Address
		addr2 := icCopy.(*strategies.ValidatorAlertIncentive).Validator.Address

		addr1[0] = 1

		require.NotEqual(t, addr1, addr2)
	})
}

func TestStrategy(t *testing.T) {
	var (
		mockStakingKeeper *mocks.StakingKeeper
		mockBankKeeper    *mocks.BankKeeper
	)

	ctx := sdk.Context{}.WithLogger(log.NewNopLogger())

	slashFraction := math.LegacyNewDecFromIntWithPrec(math.NewInt(5), 1)

	cases := []struct {
		name               string
		validatorIncentive incentivetypes.Incentive
		setup              func()
		expectedErr        error
	}{
		{
			"incorrect incentive type",
			&goodprice.GoodPriceIncentive{},
			func() {
				mockStakingKeeper = mocks.NewStakingKeeper(t)
				mockBankKeeper = mocks.NewBankKeeper(t)
			},
			fmt.Errorf("incentive must be of type ValidatorAlertIncentive, got %T", &goodprice.GoodPriceIncentive{}),
		},
		{
			"validator not found error",
			strategies.NewValidatorAlertIncentive(cmtabci.Validator{
				Address: []byte("test"),
				Power:   1,
			}, 1, sdk.AccAddress("test")),
			func() {
				mockStakingKeeper = mocks.NewStakingKeeper(t)
				mockBankKeeper = mocks.NewBankKeeper(t)
				mockStakingKeeper.On("GetValidatorByConsAddr", ctx, sdk.ConsAddress("test")).Return(stakingtypes.Validator{}, fmt.Errorf("error")).Once()
			},
			fmt.Errorf("validator with address %s does not exist", sdk.ConsAddress("test")),
		},
		{
			"slash error",
			strategies.NewValidatorAlertIncentive(cmtabci.Validator{
				Address: []byte("test"),
				Power:   1,
			}, 1, sdk.AccAddress("test")),
			func() {
				mockStakingKeeper = mocks.NewStakingKeeper(t)
				mockBankKeeper = mocks.NewBankKeeper(t)
				mockStakingKeeper.On("GetValidatorByConsAddr", ctx, sdk.ConsAddress("test")).Return(stakingtypes.Validator{}, nil).Once()
				mockStakingKeeper.On("Slash", ctx, sdk.ConsAddress("test"), 1-sdk.ValidatorUpdateDelay, int64(1), slashFraction).Return(math.Int{}, fmt.Errorf("slash error")).Once()
			},
			fmt.Errorf("failed to slash validator: slash error"),
		},
		{
			"bond denom error",
			strategies.NewValidatorAlertIncentive(cmtabci.Validator{
				Address: []byte("test"),
				Power:   1,
			}, 1, sdk.AccAddress("test")),
			func() {
				mockStakingKeeper = mocks.NewStakingKeeper(t)
				mockBankKeeper = mocks.NewBankKeeper(t)
				mockStakingKeeper.On("GetValidatorByConsAddr", ctx, sdk.ConsAddress("test")).Return(stakingtypes.Validator{}, nil).Once()
				mockStakingKeeper.On("Slash", ctx, sdk.ConsAddress("test"), 1-sdk.ValidatorUpdateDelay, int64(1), slashFraction).Return(math.NewInt(1), nil).Once()
				mockStakingKeeper.On("BondDenom", ctx).Return("", fmt.Errorf("error"))
			},
			fmt.Errorf("failed to get bond denom: error"),
		},
		{
			"mint coins error",
			strategies.NewValidatorAlertIncentive(cmtabci.Validator{
				Address: []byte("test"),
				Power:   1,
			}, 1, sdk.AccAddress("test")),
			func() {
				mockStakingKeeper = mocks.NewStakingKeeper(t)
				mockBankKeeper = mocks.NewBankKeeper(t)
				mockStakingKeeper.On("GetValidatorByConsAddr", ctx, sdk.ConsAddress("test")).Return(stakingtypes.Validator{}, nil).Once()
				mockStakingKeeper.On("Slash", ctx, sdk.ConsAddress("test"), 1-sdk.ValidatorUpdateDelay, int64(1), slashFraction).Return(math.NewInt(1), nil).Once()
				mockStakingKeeper.On("BondDenom", ctx).Return("stake", nil).Once()
				mockBankKeeper.On("MintCoins", ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1)))).Return(fmt.Errorf("mint error")).Once()
			},
			fmt.Errorf("failed to mint coins: mint error"),
		},
		{
			"send coins error",
			strategies.NewValidatorAlertIncentive(cmtabci.Validator{
				Address: []byte("test"),
				Power:   1,
			}, 1, sdk.AccAddress("signer")),
			func() {
				mockStakingKeeper = mocks.NewStakingKeeper(t)
				mockBankKeeper = mocks.NewBankKeeper(t)
				mockStakingKeeper.On("GetValidatorByConsAddr", ctx, sdk.ConsAddress("test")).Return(stakingtypes.Validator{}, nil).Once()
				mockStakingKeeper.On("Slash", ctx, sdk.ConsAddress("test"), 1-sdk.ValidatorUpdateDelay, int64(1), slashFraction).Return(math.NewInt(1), nil).Once()
				mockStakingKeeper.On("BondDenom", ctx).Return("stake", nil).Once()
				mockBankKeeper.On("MintCoins", ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1)))).Return(nil).Once()
				mockBankKeeper.On("SendCoinsFromModuleToAccount", ctx, types.ModuleName, sdk.AccAddress("signer"), sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1)))).Return(fmt.Errorf("send error")).Once()
			},
			fmt.Errorf("failed to send coins: send error"),
		},
		{
			"success",
			strategies.NewValidatorAlertIncentive(cmtabci.Validator{
				Address: []byte("test"),
				Power:   1,
			}, 1, sdk.AccAddress("signer")),
			func() {
				mockStakingKeeper = mocks.NewStakingKeeper(t)
				mockBankKeeper = mocks.NewBankKeeper(t)
				mockStakingKeeper.On("GetValidatorByConsAddr", ctx, sdk.ConsAddress("test")).Return(stakingtypes.Validator{}, nil).Once()
				mockStakingKeeper.On("Slash", ctx, sdk.ConsAddress("test"), 1-sdk.ValidatorUpdateDelay, int64(1), slashFraction).Return(math.NewInt(1), nil).Once()
				mockStakingKeeper.On("BondDenom", ctx).Return("stake", nil).Once()
				mockBankKeeper.On("MintCoins", ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1)))).Return(nil).Once()
				mockBankKeeper.On("SendCoinsFromModuleToAccount", ctx, types.ModuleName, sdk.AccAddress("signer"), sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1)))).Return(nil).Once()
			},
			nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			strategy := strategies.DefaultValidatorAlertIncentiveStrategy(mockStakingKeeper, mockBankKeeper)

			_, err := strategy(ctx, tc.validatorIncentive)
			if err == nil {
				if tc.expectedErr != nil {
					t.Errorf("expected error %s, got nil", tc.expectedErr)
				}
				return
			}

			if err.Error() != tc.expectedErr.Error() {
				t.Errorf("expected error %s, got %s", tc.expectedErr, err)
			}
		})
	}
}

func TestDefaultHandler(t *testing.T) {
	cases := []struct {
		name  string
		od    slinkyabci.OracleVoteExtension
		pb    types.PriceBound
		a     types.Alert
		v     cmtabci.Validator
		ic    *strategies.ValidatorAlertIncentive
		id    uint64
		valid bool
	}{
		{
			"invalid alert",
			slinkyabci.OracleVoteExtension{
				Prices: map[uint64][]byte{},
			},
			types.PriceBound{
				High: "1",
				Low:  "2",
			},
			types.Alert{
				Signer: sdk.AccAddress("test").String(),
				Height: 1,
			},
			cmtabci.Validator{
				Address: []byte("test"),
			},
			nil,
			0,
			false,
		},
		{
			"invalid price-bound",
			slinkyabci.OracleVoteExtension{
				Prices: map[uint64][]byte{},
			},
			types.PriceBound{
				High: "1",
				Low:  "1",
			},
			types.Alert{
				Signer: sdk.AccAddress("test").String(),
				Height: 1,
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "A",
					Quote: "B",
				},
			},
			cmtabci.Validator{
				Address: []byte("test"),
			},
			nil,
			0,
			false,
		},
		{
			"no price report, nil incentive",
			slinkyabci.OracleVoteExtension{
				Prices: map[uint64][]byte{},
			},
			types.PriceBound{
				High: "2",
				Low:  "1",
			},
			types.Alert{
				Signer: sdk.AccAddress("test").String(),
				Height: 1,
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "A",
					Quote: "B",
				},
			},
			cmtabci.Validator{
				Address: []byte("test"),
			},
			nil,
			0,
			true,
		},
		{
			"if price is higher than high bound, incentive is non-nil",
			slinkyabci.OracleVoteExtension{
				Prices: map[uint64][]byte{
					0: big.NewInt(3).Bytes(),
				},
			},
			types.PriceBound{
				High: "2",
				Low:  "1",
			},
			types.Alert{
				Signer: sdk.AccAddress("signer").String(),
				Height: 1,
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "A",
					Quote: "B",
				},
			},
			cmtabci.Validator{
				Address: []byte("test"),
			},
			&strategies.ValidatorAlertIncentive{
				Validator: cmtabci.Validator{
					Address: []byte("test"),
				},
				AlertSigner: sdk.AccAddress("signer").String(),
				AlertHeight: 1,
			},
			0,
			true,
		},
		{
			"if price is lower than low bound, incentive is non-nil",
			slinkyabci.OracleVoteExtension{
				Prices: map[uint64][]byte{
					0: big.NewInt(0).Bytes(),
				},
			},
			types.PriceBound{
				High: "2",
				Low:  "1",
			},
			types.Alert{
				Signer: sdk.AccAddress("signer").String(),
				Height: 1,
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "A",
					Quote: "B",
				},
			},
			cmtabci.Validator{
				Address: []byte("test"),
			},
			&strategies.ValidatorAlertIncentive{
				Validator: cmtabci.Validator{
					Address: []byte("test"),
				},
				AlertSigner: sdk.AccAddress("signer").String(),
				AlertHeight: 1,
			},
			0,
			true,
		},
		{
			"if price is within bounds, incentive is nil",
			slinkyabci.OracleVoteExtension{
				Prices: map[uint64][]byte{
					0: big.NewInt(1).Bytes(),
				},
			},
			types.PriceBound{
				High: "2",
				Low:  "1",
			},
			types.Alert{
				Signer: sdk.AccAddress("signer").String(),
				Height: 1,
				CurrencyPair: slinkytypes.CurrencyPair{
					Base:  "A",
					Quote: "B",
				},
			},
			cmtabci.Validator{
				Address: []byte("test"),
			},
			nil,
			0,
			true,
		},
	}

	handler := strategies.DefaultHandleValidatorIncentive()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bz, err := tc.od.Marshal()
			require.NoError(t, err)

			// create extendedVoteInfo
			incentive, err := handler(
				cmtabci.ExtendedVoteInfo{
					VoteExtension: bz,
					Validator:     tc.v,
				},
				tc.pb,
				tc.a,
				tc.id,
			)

			// check for expected errors
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				return
			}

			// if error is not expected, check for expected incentive
			if tc.ic != nil {
				gotIc, ok := incentive.(*strategies.ValidatorAlertIncentive)
				require.True(t, ok)

				require.Equal(t, tc.ic, gotIc)
			} else {
				require.Nil(t, incentive)
			}
		})
	}
}
