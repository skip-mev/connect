package abci

import (
	"fmt"
	"time"

	"cosmossdk.io/log"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/holiman/uint256"
	abcitypes "github.com/skip-mev/slinky/abci/types"
	"github.com/skip-mev/slinky/service"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

type VoteExtensionHandler struct {
	logger log.Logger

	// oracle is the oracle service that is responsible for fetching prices
	oracle service.OracleService

	// timeout is the maximum amount of time to wait for the oracle to respond
	// to a price request.
	timeout time.Duration
}

func NewVoteExtensionHandler(logger log.Logger, oracle service.OracleService, timeout time.Duration) *VoteExtensionHandler {
	return &VoteExtensionHandler{
		logger:  logger.With("module", "VoteExtensionHandler"),
		oracle:  oracle,
		timeout: timeout,
	}
}

// ExtendVoteHandler returns a handler that extends a vote with the oracle's
// current price feed. In the case where oracle data is unable to be fetched
// or correctly marshalled, the handler will return an empty vote extension to
// ensure liveliness.
func (h *VoteExtensionHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *abci.RequestExtendVote) (resp *abci.ResponseExtendVote, err error) {
		// Catch any panic that occurs in the oracle request.
		defer func() {
			if r := recover(); r != nil {
				h.logger.Error(
					"recovered from panic in ExtendVoteHandler",
					"err", r,
				)

				resp, err = &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
			}
		}()

		// To ensure liveliness, we return a vote even if the oracle is not running
		// or if the oracle returns a bad response.
		oracleResp, err := h.oracle.Prices(ctx, &service.QueryPricesRequest{})
		if err != nil {
			h.logger.Error(
				"failed to retrieve oracle prices for vote extension",
				"height", req.Height,
				"err", err,
			)

			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
		}

		// If we get no response, we return an empty vote extension.
		if oracleResp == nil {
			h.logger.Error(
				"oracle returned nil prices for vote extension",
				"height", req.Height,
			)

			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
		}

		h.logger.Info("extending vote with oracle prices", "num_prices", len(oracleResp.Prices))

		voteExt := &abcitypes.OracleVoteExtension{
			Height:    req.Height,
			Prices:    oracleResp.Prices,
			Timestamp: oracleResp.Timestamp,
		}

		bz, err := voteExt.Marshal()
		if err != nil {
			h.logger.Error(
				"failed to marshal vote extension",
				"height", req.Height,
				"err", err,
			)

			return &abci.ResponseExtendVote{VoteExtension: []byte{}}, nil
		}

		return &abci.ResponseExtendVote{VoteExtension: bz}, nil
	}
}

// VerifyVoteExtensionHandler returns a handler that verifies the vote extension provided by
// a validator is valid. In the case when the vote extension is empty, we return ACCEPT. This means
// that the validator was unable to fetch prices from the oracle and is voting an empty vote extension.
// We reject any vote extensions that are not empty and fail to unmarshal or contain invalid prices.
func (h *VoteExtensionHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
		voteExtension := req.VoteExtension

		// If we get an empty vote extension, we return ACCEPT. This means that the validator was unable
		// to fetch prices from the oracle and is voting on an empty vote extension.
		if len(voteExtension) == 0 {
			return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
		}

		voteExt := &abcitypes.OracleVoteExtension{}
		if err := voteExt.Unmarshal(voteExtension); err != nil {
			return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_REJECT},
				fmt.Errorf("failed to unmarshal vote extension: %w", err)
		}

		// The height of the vote extension must match the height of the request.
		if voteExt.Height != req.Height {
			return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_REJECT},
				fmt.Errorf("vote extension height does not match request height; expected: %d, got: %d", req.Height, voteExt.Height)
		}

		// Verify tickers and prices are valid.
		for currencyPair, price := range voteExt.Prices {
			if _, err := oracletypes.CurrencyPairFromString(currencyPair); err != nil {
				return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_REJECT},
					fmt.Errorf("invalid ticker in oracle vote extension %s: %w", currencyPair, err)
			}

			if _, err := uint256.FromHex(price); err != nil {
				return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_REJECT},
					fmt.Errorf("invalid price in oracle vote extension %s: %w", currencyPair, err)
			}
		}

		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}
