package types_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/x/alerts/types"
	oracle_types "github.com/skip-mev/slinky/x/oracle/types"
	"github.com/stretchr/testify/assert"
)

func TestAlertUnmarshal(t *testing.T) {
	// create an interface registry
	ir := codectypes.NewInterfaceRegistry()

	// create a codec
	cdc := codec.NewProtoCodec(ir)

	// create an alert
	alert := types.Alert{
		Height: 1,
		Signer: "signer",
		CurrencyPair: oracle_types.CurrencyPair{
			Base:  "base",
			Quote: "quote",
		},
	}

	// marshal the alert
	bz, err := cdc.Marshal(&alert)
	assert.NoError(t, err)

	// unmarshal the alert
	alert2 := types.Alert{}
	assert.NoError(t, alert2.Unmarshal(bz))

	// assert that the two alerts are equal
	assert.Equal(t, alert, alert2)
}

func TestAlertValidateBasic(t *testing.T) {
	type testCase struct {
		name  string
		alert types.Alert
		valid bool
	}

	cases := []testCase{
		{
			"invalid signer - fail",
			types.Alert{
				Height: 1,
				Signer: "",
				CurrencyPair: oracle_types.CurrencyPair{
					Base:  "BASE",
					Quote: "QUOTE",
				},
			},
			false,
		},
		{
			"invalid currency-pair - fail",
			types.Alert{
				Height: 1,
				Signer: sdk.AccAddress("signer").String(),
				CurrencyPair: oracle_types.CurrencyPair{
					Base:  "",
					Quote: "",
				},
			},
			false,
		},
		{
			name: "valid alert - pass",
			alert: types.Alert{
				Height: 1,
				Signer: sdk.AccAddress("signer").String(),
				CurrencyPair: oracle_types.CurrencyPair{
					Base:  "BASE",
					Quote: "QUOTE",
				},
			},
			valid: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.alert.ValidateBasic()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestAlertUID(t *testing.T) {
	// create two alerts
	alert1 := types.Alert{
		Height: 1,
		Signer: "signer",
		CurrencyPair: oracle_types.CurrencyPair{
			Base:  "base",
			Quote: "quote",
		},
	}

	alert2 := types.Alert{
		Height: 2,
		Signer: "signer",
		CurrencyPair: oracle_types.CurrencyPair{
			Base:  "base",
			Quote: "quote",
		},
	}
	t.Run("test Alert UID uniqueness", func(t *testing.T) {
		// check that the alerts have different UIDs
		assert.NotEqual(t, alert1.UID(), alert2.UID())
	})
}

func TestAlertStatus(t *testing.T) {
	t.Run("test String", func(t *testing.T) {
		cs := types.ConclusionStatus(0)
		assert.Equal(t, "Unconcluded", cs.String())

		cs = types.ConclusionStatus(1)
		assert.Equal(t, "Concluded", cs.String())

		cs = types.ConclusionStatus(2)
		assert.Equal(t, "unknown", cs.String())
	})

	t.Run("test validate basic", func(t *testing.T) {
		cases := []struct {
			name        string
			alertStatus types.AlertStatus
			valid       bool
		}{
			{
				"invalid submission height",
				types.AlertStatus{
					SubmissionHeight: 3,
					PurgeHeight:      2,
				},
				false,
			},
			{
				"invalid conclusion status",
				types.AlertStatus{
					SubmissionHeight: 1,
					PurgeHeight:      2,
					ConclusionStatus: 2,
				},
				false,
			},
			{
				"valid alert status",
				types.AlertStatus{
					SubmissionHeight: 1,
					PurgeHeight:      2,
					ConclusionStatus: 1,
				},
				true,
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				err := tc.alertStatus.ValidateBasic()
				if tc.valid {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		}
	})
}

func TestAlertWithStatus(t *testing.T) {
	cases := []struct {
		name  string
		alert types.AlertWithStatus
		valid bool
	}{
		{
			"invalid alert",
			types.AlertWithStatus{
				Alert: types.Alert{
					Height: 1,
					Signer: "",
					CurrencyPair: oracle_types.CurrencyPair{
						Base:  "BASE",
						Quote: "QUOTE",
					},
				},
				Status: types.NewAlertStatus(1, 2, time.Now(), 1),
			},
			false,
		},
		{
			"invalid alert-status",
			types.AlertWithStatus{
				Alert: types.NewAlert(1, sdk.AccAddress("signer"), oracle_types.NewCurrencyPair("BASE", "QUOTE")),
				Status: types.AlertStatus{
					SubmissionHeight: 3,
					PurgeHeight:      2,
				},
			},
			false,
		},
		{
			"valid alert with status",
			types.NewAlertWithStatus(
				types.NewAlert(1, sdk.AccAddress("signer"), oracle_types.NewCurrencyPair("BASE", "QUOTE")),
				types.NewAlertStatus(1, 2, time.Now(), 1),
			),
			true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.alert.ValidateBasic()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
