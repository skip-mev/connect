package integration

import (
	"context"
	"testing"
	"time"

	"cosmossdk.io/math"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/strangelove-ventures/interchaintest/v8/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v8/testutil"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/skip-mev/slinky/abci/strategies/codec"
	slinkyabci "github.com/skip-mev/slinky/abci/ve/types"
	alerttypes "github.com/skip-mev/slinky/x/alerts/types"
)

const gasPrice = 100

// UpdateAlertParams creates + submits the proposal to update the alert params, votes for the prop w/ all nodes,
// and waits for the proposal to pass.
func UpdateAlertParams(chain *cosmos.CosmosChain, authority, denom string, deposit int64, timeout time.Duration, user cosmos.User, params alerttypes.Params) (string, error) {
	propId, err := SubmitProposal(chain, sdk.NewCoin(denom, math.NewInt(deposit)), user.KeyName(), []sdk.Msg{&alerttypes.MsgUpdateParams{
		Authority: authority,
		Params:    params,
	}}...)
	if err != nil {
		return "", err
	}

	return propId, PassProposal(chain, propId, timeout)
}

// SubmitAlert submits an alert to the chain, submitted by the given address
func (s *SlinkySlashingIntegrationSuite) SubmitAlert(user cosmos.User, alert alerttypes.Alert) (*coretypes.ResultBroadcastTxCommit, error) {
	tx := CreateTx(s.T(), s.chain, user, gasPrice, alerttypes.NewMsgAlert(
		alert,
	))

	// get an rpc endpoint for the chain
	client := s.chain.Nodes()[0].Client

	// broadcast the tx
	return client.BroadcastTxCommit(context.Background(), tx)
}

// SubmitConclusion submits the provided conclusion to the chain
func (s *SlinkySlashingIntegrationSuite) SubmitConclusion(user cosmos.User, conclusion alerttypes.Conclusion) (*coretypes.ResultBroadcastTxCommit, error) {
	addr, err := sdk.AccAddressFromBech32(user.FormattedAddress())
	s.Require().NoError(err)

	tx := CreateTx(s.T(), s.chain, user, gasPrice, alerttypes.NewMsgConclusion(
		conclusion,
		addr,
	))

	// get an rpc endpoint for the chain
	client := s.chain.Nodes()[0].Client
	return client.BroadcastTxCommit(context.Background(), tx)
}

// Delegate sends a delegation tx from the given user to the given validator
func (s *SlinkySlashingIntegrationSuite) Delegate(user cosmos.User, validatorOperator string, tokens sdk.Coin) (*coretypes.ResultBroadcastTxCommit, error) {
	msg := stakingtypes.NewMsgDelegate(user.FormattedAddress(), validatorOperator, tokens)

	tx := CreateTx(s.T(), s.chain, user, gasPrice, msg)

	// get an rpc endpoint for the chain
	client := s.chain.Nodes()[0].Client
	return client.BroadcastTxCommit(context.Background(), tx)
}

// CreateTx creates a new transaction to be signed by the given user, including a provided set of messages
func CreateTx(t *testing.T, chain *cosmos.CosmosChain, user cosmos.User, GasPrice int64, msgs ...sdk.Msg) []byte {
	bc := cosmos.NewBroadcaster(t, chain)

	ctx := context.Background()
	// create tx factory + Client Context
	txf, err := bc.GetFactory(ctx, user)
	require.NoError(t, err)

	cc, err := bc.GetClientContext(ctx, user)
	require.NoError(t, err)

	txf = txf.WithSimulateAndExecute(true)

	txf, err = txf.Prepare(cc)
	require.NoError(t, err)

	// get gas for tx
	txf.WithGas(25000000)

	// update sequence number
	txf = txf.WithSequence(txf.Sequence())
	txf = txf.WithGasPrices(sdk.NewDecCoins(sdk.NewDecCoin(chain.Config().Denom, math.NewInt(GasPrice))).String())

	// sign the tx
	txBuilder, err := txf.BuildUnsignedTx(msgs...)
	require.NoError(t, err)

	require.NoError(t, tx.Sign(cc.CmdContext, txf, cc.GetFromName(), txBuilder, true))

	// encode and return
	bz, err := cc.TxConfig.TxEncoder()(txBuilder.GetTx())
	require.NoError(t, err)
	return bz
}

type PrivKeyType uint64

const (
	Secp256k1 PrivKeyType = iota
	Ed25519
)

// CreateKey creates a private key wrt. the given private-key type
func CreateKey(typ PrivKeyType) cryptotypes.PrivKey {
	switch typ {
	case Secp256k1:
		return secp256k1.GenPrivKey()
	case Ed25519:
		return ed25519.GenPrivKey()
	default:
		panic("unknown private key type")
	}
}

// QueryValidators queries for all network's validators
func QueryValidators(chain *cosmos.CosmosChain) ([]stakingtypes.Validator, error) {
	// get grpc client of the node
	grpcAddr := chain.GetHostGRPCAddress()
	cc, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := stakingtypes.NewQueryClient(cc)

	// query validators
	resp, err := client.Validators(context.Background(), &stakingtypes.QueryValidatorsRequest{})
	if err != nil {
		return nil, err
	}

	return resp.Validators, nil
}

// GetExtendedCommit gets the extended commit from a block for a designated height
func GetExtendedCommit(chain *cosmos.CosmosChain, height int64) (cmtabci.ExtendedCommitInfo, error) {
	// get a client from the chain
	client := chain.Nodes()[0].Client

	// query the block
	block, err := client.Block(context.Background(), &height)
	if err != nil {
		return cmtabci.ExtendedCommitInfo{}, err
	}

	// get the extended commit
	eci, err := extCommitCodec.Decode(block.Block.Txs[0])
	if err != nil {
		return cmtabci.ExtendedCommitInfo{}, err
	}

	// unmarshal votes
	voteEncoder := codec.NewDefaultVoteExtensionCodec()
	for i, vote := range eci.Votes {
		// unmarshal compressed ve
		voteInfo, err := veCodec.Decode(vote.VoteExtension)
		if err != nil {
			return cmtabci.ExtendedCommitInfo{}, err
		}

		eci.Votes[i].VoteExtension, err = voteEncoder.Encode(voteInfo)
		if err != nil {
			return cmtabci.ExtendedCommitInfo{}, err
		}
	}

	return eci, nil
}

// GetOracleDataFromVote gets the oracle-data from a vote included in a last commit info
func GetOracleDataFromVote(vote cmtabci.ExtendedVoteInfo) (slinkyabci.OracleVoteExtension, error) {
	var ve slinkyabci.OracleVoteExtension
	if err := ve.Unmarshal(vote.VoteExtension); err != nil {
		return slinkyabci.OracleVoteExtension{}, err
	}

	return ve, nil
}

// ExpectAlerts waits until the provided alerts are in module state or until timeout. This method returns an error if it times-out. Otherwise,
// it returns the height for which the condition was satisfied.
//
// Notice: the height returned is safe for querying
func ExpectAlerts(chain *cosmos.CosmosChain, timeout time.Duration, alerts []alerttypes.Alert) (uint64, error) {
	cc, close, err := GetChainGRPC(chain)
	if err != nil {
		return 0, err
	}
	defer close()

	alertsClient := alerttypes.NewQueryClient(cc)

	var height int64

	if err := testutil.WaitForCondition(timeout, 100*time.Millisecond, func() (bool, error) {
		height, err = chain.Height(context.Background())

		resp, err := alertsClient.Alerts(context.Background(), &alerttypes.AlertsRequest{})
		if err != nil {
			return false, err
		}

		if len(resp.Alerts) != len(alerts) {
			return false, nil
		}

		expectedAlerts := mapAlerts(alerts)
		for _, alert := range alerts {
			if _, ok := expectedAlerts[alertKey(alert)]; !ok {
				return false, nil
			}
		}

		return true, nil
	}); err != nil {
		return 0, err
	}

	return uint64(height), WaitForHeight(chain, uint64(height+1), timeout)
}

func mapAlerts(alerts []alerttypes.Alert) map[string]struct{} {
	m := make(map[string]struct{})

	for _, alert := range alerts {
		m[alertKey(alert)] = struct{}{}
	}

	return m
}

func alertKey(alert alerttypes.Alert) string {
	return string(alert.UID())
}
