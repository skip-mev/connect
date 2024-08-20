package keeper

import (
	"context"
	"fmt"

	"github.com/skip-mev/connect/v2/x/alerts/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	incentivetypes "github.com/skip-mev/connect/v2/x/incentives/types"
)

// msgServer is the default implementation of the x/alerts MsgService.
type msgServer struct {
	k Keeper
}

// NewMsgServer returns an implementation of the x/alerts MsgServer interface
// for the provided Keeper.
func NewMsgServer(k Keeper) types.MsgServer {
	return &msgServer{k}
}

var _ types.MsgServer = msgServer{}

// Alert implements the MsgServer.Alerts method, which is used to create a new alert. This method
// will check that the referenced alert does not already exist, that the currency-pair referenced in the alert
// exists, that the alert's age is less than max-block-age, and that the alert itself is valid. If any of these
// checks fail, the method will return an error. If the alert is valid, and Alerts are enabled, then params.BondAmount
// will be escrowed at the module account, and the alert will be added to the module's state.
func (m msgServer) Alert(goCtx context.Context, req *types.MsgAlert) (*types.MsgAlertResponse, error) {
	// request should not be nil
	if req == nil {
		return nil, fmt.Errorf("message cannot be empty")
	}

	// check that the message is valid
	if err := req.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("message validation failed: %w", err)
	}

	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check that alerts are enabled
	params := m.k.GetParams(ctx)
	if !params.AlertParams.Enabled {
		return nil, fmt.Errorf("alerts are not enabled")
	}

	// check that the alert's age < MaxBlockAge
	height := uint64(ctx.BlockHeight())
	if alertBlockAge := height - req.Alert.Height; alertBlockAge > params.AlertParams.MaxBlockAge {
		return nil, fmt.Errorf("alert is too old: %d > %d", alertBlockAge, params.AlertParams.MaxBlockAge)
	}

	// check that the referenced alert does not already exist
	if alert, ok := m.k.GetAlert(ctx, req.Alert); ok {
		return nil, fmt.Errorf("alert with UID %X already exists: %v", req.Alert.UID(), alert)
	}

	// check that the referenced currency-pair exists
	if !m.k.oracleKeeper.HasCurrencyPair(ctx, req.Alert.CurrencyPair) {
		return nil, fmt.Errorf("currency pair %s does not exist", req.Alert.CurrencyPair)
	}

	// escrow the bond amount
	if err := m.k.escrowBondAmount(ctx, req.Alert.Signer, params.AlertParams.BondAmount); err != nil {
		return nil, fmt.Errorf("failed to escrow bond amount: %w", err)
	}

	// add the alert + alert-status to the module's state
	if err := m.k.SetAlert(ctx, types.NewAlertWithStatus(
		req.Alert,
		types.NewAlertStatus(
			height,
			height+params.PruningParams.BlocksToPrune, // keep the alert in state until height + blocksToPrune
			ctx.BlockTime(), // this alert can be safely concluded until ctx.BlockTime() + unbondingTime
			types.Unconcluded,
		),
	)); err != nil {
		return nil, fmt.Errorf("failed to set alert: %w", err)
	}

	// return the response
	return &types.MsgAlertResponse{}, nil
}

// escrowBondAmount is a helper function that will escrow the bond amount at the module address.
func (k *Keeper) escrowBondAmount(ctx sdk.Context, signer string, bondAmount sdk.Coin) error {
	// get the sdk address for the signer
	addr, err := sdk.AccAddressFromBech32(signer)
	if err != nil {
		return fmt.Errorf("failed to get sdk address for signer: %w", err)
	}

	// send coins to the module account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, sdk.NewCoins(bondAmount)); err != nil {
		return fmt.Errorf("failed to escrow bond amount: %w", err)
	}

	return nil
}

// Conclusion is a no-op if alerts are not enabled, or if the conclusion is not valid (according to MsgConclusion.ValidateBasic()). The conclusion
// must be verifiable in accordance w/ the registered Conclusions / ConclusionVerificationParams for this module. If the above criteria
// are met, then depending on the status of the conclusion, incentives will be issued to the parties deemed at fault, and the referenced
// alert will be marked as concluded.
func (m msgServer) Conclusion(goCtx context.Context, req *types.MsgConclusion) (*types.MsgConclusionResponse, error) {
	// check if the msg is nil
	if req == nil {
		return nil, fmt.Errorf("message cannot be empty")
	}

	// check that the message is valid
	if err := req.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("message validation failed: %w", err)
	}

	// unwrap the context
	ctx := sdk.UnwrapSDKContext(goCtx)

	// check that alerts are enabled
	params := m.k.GetParams(ctx)
	if !params.AlertParams.Enabled {
		return nil, fmt.Errorf("alerts are not enabled")
	}

	// unmarshal the conclusion
	conclusion, ok := req.Conclusion.GetCachedValue().(types.Conclusion)
	if !ok {
		return nil, fmt.Errorf("failed to unmarshal conclusion")
	}

	// unmarshal the conclusion verification params, this should never error
	var verificationParams types.ConclusionVerificationParams
	if err := m.k.cdc.UnpackAny(params.ConclusionVerificationParams, &verificationParams); err != nil {
		return nil, fmt.Errorf("failed to unmarshal conclusion verification params: %w", err)
	}

	// verify the conclusion in accordance with the ConclusionVerificationParams
	if err := conclusion.Verify(verificationParams); err != nil {
		return nil, fmt.Errorf("failed to verify conclusion: %w", err)
	}

	m.k.Logger(ctx).Info("conclusion verified", "conclusion", conclusion.String(), "params", verificationParams.String())

	// conclusion has been verified, mark the alert as concluded
	if err := m.k.ConcludeAlert(ctx, conclusion.GetAlert(), boolToConclusionStatus(conclusion.GetStatus())); err != nil {
		return nil, fmt.Errorf("failed to conclude alert: %w", err)
	}

	// finally, if the conclusion was positive, issue incentives to all validators referenced in the conclusion
	if conclusion.GetStatus() {
		extCommit := conclusion.GetExtendedCommitInfo()

		incentives := make([]incentivetypes.Incentive, 0)

		// determine whether to issue an incentive to each validator who signed a vote in the Commit referenced
		for _, vote := range extCommit.Votes {
			alert := conclusion.GetAlert()
			m.k.Logger(ctx).Info("issuing incentive to validator", "validator", sdk.ConsAddress(vote.Validator.Address).String(), "alert", fmt.Sprintf("%X", alert.UID()))

			// execute the ValidatorIncentiveHandler to determine if validator should be issued an incentive
			incentive, err := m.k.validatorIncentiveHandler(vote, conclusion.GetPriceBound(), conclusion.GetAlert(), conclusion.GetCurrencyPairID())
			if err != nil {
				return nil, fmt.Errorf("failed to determine incentive: %w", err)
			}

			// if the incentive is non-nil, then add it to the list of incentives to issue
			if incentive != nil {
				getAlert := conclusion.GetAlert()
				m.k.Logger(ctx).Info("incentive issued to validator", "validator", vote.Validator.Address, "incentive", incentive.String(), "alert", fmt.Sprintf("%X", getAlert.UID()))
				incentives = append(incentives, incentive)
			}
		}

		// finally, issue the incentives
		if err := m.k.incentiveKeeper.AddIncentives(ctx, incentives); err != nil {
			return nil, fmt.Errorf("failed to issue incentives: %w", err)
		}
	}

	return nil, nil
}

func boolToConclusionStatus(status bool) ConclusionStatus {
	if status {
		return Positive
	}
	return Negative
}

// UpdateParams is the handler for the UpdateParams RPC. This method expects a MsgUpdateParams message. This method fails if the msg fails validation
// or if the provided signer is not the authority of this module. Otherwise, the given params are set to state.
func (m msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	// fail for nil requests
	if req == nil {
		return nil, fmt.Errorf("message cannot be empty")
	}

	// validate the message
	if err := req.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("message validation failed: %w", err)
	}

	// check that the signer (authority) of the message is the authority of this module
	if authority, err := sdk.AccAddressFromBech32(req.Authority); err == nil {
		// get the sdk address for the signer
		if !authority.Equals(m.k.authority) {
			return nil, fmt.Errorf("signer is not the authority of this module: signer %v, authority %v", req.Authority, m.k.authority.String())
		}
	} else {
		return nil, fmt.Errorf("failed to get sdk address for authority: %w", err)
	}

	// signer is the authority of the module, update params
	if err := m.k.SetParams(sdk.UnwrapSDKContext(goCtx), req.Params); err != nil {
		return nil, fmt.Errorf("failed to set params: %w", err)
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
