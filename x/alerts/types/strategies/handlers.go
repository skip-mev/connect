package strategies

import (
	cmtabci "github.com/cometbft/cometbft/abci/types"

	"github.com/skip-mev/connect/v2/x/alerts/types"
	incentivetypes "github.com/skip-mev/connect/v2/x/incentives/types"
)

// ValidatorIncentiveHandler determines whether a validator's price report deviated significantly from
// what was expected off-chain, and returns the alert to be issued to the incentives keeper if so.
type ValidatorIncentiveHandler func(ve cmtabci.ExtendedVoteInfo, pb types.PriceBound, a types.Alert, cpID uint64) (incentivetypes.Incentive, error)
