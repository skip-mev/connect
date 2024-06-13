package integration

import (
	"context"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type SlinkyOracleValidatorIntegrationSuite struct {
	*SlinkyIntegrationSuite
}

func NewSlinkyOracleValidatorIntegrationSuite(suite *SlinkyIntegrationSuite) *SlinkyOracleValidatorIntegrationSuite {
	return &SlinkyOracleValidatorIntegrationSuite{
		SlinkyIntegrationSuite: suite,
	}
}

func (s *SlinkyOracleValidatorIntegrationSuite) TestUnbonding() {
	ctx := context.Background()
	vals, err := s.chain.StakingQueryValidators(ctx, stakingtypes.Bonded.String())
	s.Require().NoError(err)
	val := vals[0].OperatorAddress

	wasErr := true
	for _, node := range s.chain.Validators {
		err = node.StakingUnbond(ctx, validatorKey, val, vals[0].BondedTokens().String()+s.denom)
		if err == nil {
			wasErr = false
		}
	}
	s.Require().False(wasErr)
}
