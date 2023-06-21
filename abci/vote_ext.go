package abci

import (
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/holiman/uint256"
	"github.com/skip-mev/slinky/abci/types"
	"github.com/skip-mev/slinky/service"
)

type VoteExtHandler struct {
	logger       log.Logger
	currentBlock int64
	oracle       service.OracleService
}

func NewVoteExtHandler(logger log.Logger, oracle service.OracleService) *VoteExtHandler {
	return &VoteExtHandler{
		logger: logger.With("module", "VoteExtHandler"),
		oracle: oracle,
	}
}

func (h *VoteExtHandler) ExtendVoteHandler() sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
		h.currentBlock = req.Height

		// fetch all prices
		resp, err := h.oracle.Prices(ctx, &service.QueryPricesRequest{})
		if err != nil {
			return nil, err
		}

		h.logger.Info("retrieved oracle prices for vote extension", "height", req.Height, "last_sync_time", resp.Timestamp)

		voteExt := &types.OracleVoteExtension{
			Height:    req.Height,
			Prices:    resp.Prices,
			Timestamp: resp.Timestamp,
		}

		bz, err := voteExt.Marshal()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal oracle prices: %w", err)
		}

		return &abci.ResponseExtendVote{VoteExtension: bz}, nil
	}
}

func (h *VoteExtHandler) VerifyVoteExtensionHandler() sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
		voteExt := &types.OracleVoteExtension{}

		if err := voteExt.Unmarshal(req.VoteExtension); err != nil {
			return nil, fmt.Errorf("failed to unmarshal vote extension: %w", err)
		}

		if voteExt.Height != req.Height {
			return nil, fmt.Errorf("vote extension height does not match request height; expected: %d, got: %d", req.Height, voteExt.Height)
		}

		// verify tickers and prices are valid
		for ticker, price := range voteExt.Prices {
			if _, err := uint256.FromHex(price); err != nil {
				return nil, fmt.Errorf("invalid price in oracle vote extension %s: %w", ticker, err)
			}
		}

		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}
