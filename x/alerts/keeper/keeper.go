package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/skip-mev/connect/v2/x/alerts/types"
	"github.com/skip-mev/connect/v2/x/alerts/types/strategies"
)

type Condition func(types.AlertWithStatus) bool

type Keeper struct {
	storeService store.KVStoreService
	cdc          codec.BinaryCodec
	// Expected keepers
	bankKeeper      types.BankKeeper
	oracleKeeper    types.OracleKeeper
	incentiveKeeper types.IncentiveKeeper

	// module authority
	authority sdk.AccAddress

	// ValidatorIncentiveHandler
	validatorIncentiveHandler strategies.ValidatorIncentiveHandler

	// State
	schema collections.Schema
	// alerts are stored under (height, currency-pair) -> Alert
	alerts collections.Map[collections.Pair[uint64, string], types.AlertWithStatus]
	params collections.Item[types.Params]
}

func NewKeeper(
	storeService store.KVStoreService,
	cdc codec.BinaryCodec,
	ok types.OracleKeeper,
	bk types.BankKeeper,
	ik types.IncentiveKeeper,
	vih strategies.ValidatorIncentiveHandler,
	authority sdk.AccAddress,
) *Keeper {
	// create schema builder
	sb := collections.NewSchemaBuilder(storeService)

	k := &Keeper{
		storeService:              storeService,
		cdc:                       cdc,
		bankKeeper:                bk,
		oracleKeeper:              ok,
		incentiveKeeper:           ik,
		validatorIncentiveHandler: vih,
		alerts:                    collections.NewMap(sb, types.AlertStoreKeyPrefix, "alerts", collections.PairKeyCodec(collections.Uint64Key, collections.StringKey), codec.CollValue[types.AlertWithStatus](cdc)),
		params:                    collections.NewItem(sb, types.ParamsStoreKeyPrefix, "params", codec.CollValue[types.Params](cdc)),
		authority:                 authority,
	}

	// build schema
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.schema = schema
	return k
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger()
}

// GetAlert returns the alert for the given UID. This method returns false if no alert exists, and true
// if an alert exists.
func (k *Keeper) GetAlert(ctx sdk.Context, alert types.Alert) (types.AlertWithStatus, bool) {
	alertWithStatus, err := k.alerts.Get(ctx, collections.Join(alert.Height, alert.CurrencyPair.String()))
	if err != nil {
		return types.AlertWithStatus{}, false
	}
	return alertWithStatus, true
}

// SetAlert sets the alert to state, under the (height, currency-pair) key.
func (k *Keeper) SetAlert(ctx sdk.Context, alert types.AlertWithStatus) error {
	return k.alerts.Set(ctx, collections.Join(alert.Alert.Height, alert.Alert.CurrencyPair.String()), alert)
}

// RemoveAlert removes the alert from state, under the (height, currency-pair) key.
func (k *Keeper) RemoveAlert(ctx sdk.Context, alert types.Alert) error {
	return k.alerts.Remove(ctx, collections.Join(alert.Height, alert.CurrencyPair.String()))
}

// GetAllAlerts returns all alerts in state, it does so via an iterator over the alerts table.
func (k *Keeper) GetAllAlerts(ctx sdk.Context) ([]types.AlertWithStatus, error) {
	return k.GetAllAlertsWithCondition(ctx, func(_ types.AlertWithStatus) bool { return true })
}

// GetAllAlertsWithCondition returns all alerts for which the Condition evaluates to true.
func (k *Keeper) GetAllAlertsWithCondition(ctx sdk.Context, c Condition) ([]types.AlertWithStatus, error) {
	if c == nil {
		return nil, fmt.Errorf("condition cannot be nil")
	}

	var alerts []types.AlertWithStatus
	iter, err := k.alerts.Iterate(ctx, nil)
	if err != nil {
		return nil, err
	}

	// expect to close the iterator
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		alert, err := iter.Value()
		if err != nil {
			return nil, err
		}
		if c(alert) {
			alerts = append(alerts, alert)
		}
	}

	return alerts, nil
}

// SetParams sets the params to state.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	return k.params.Set(ctx, params)
}

// GetParams returns the params from state.
func (k *Keeper) GetParams(ctx sdk.Context) types.Params {
	params, err := k.params.Get(ctx)
	if err != nil {
		return types.Params{}
	}
	return params
}
