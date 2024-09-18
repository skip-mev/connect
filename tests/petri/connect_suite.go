package petri

import (
	"context"

	petritypes "github.com/skip-mev/petri/types/v2"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

// SlinkyIntegrationSuite is the testify suite for building and running slinky testapp networks using petri
type SlinkyIntegrationSuite struct {
	suite.Suite

	logger *zap.Logger

	spec *petritypes.ChainConfig

	chain petritypes.ChainI
}

func NewSlinkyIntegrationSuite(spec *petritypes.ChainConfig) *SlinkyIntegrationSuite {
	return &SlinkyIntegrationSuite{
		spec: spec,
	}
}

func (s *SlinkyIntegrationSuite) SetupSuite() {
	// create the logger
	var err error
	s.logger, err = zap.NewDevelopment()
	s.Require().NoError(err)

	// create the chain
	s.chain, err = GetChain(context.Background(), s.logger)
	s.Require().NoError(err)

	// initialize the chain
	err = s.chain.Init(context.Background())
	s.Require().NoError(err)
}

func (s *SlinkyIntegrationSuite) TearDownSuite() {
	err := s.chain.Teardown(context.Background())
	s.Require().NoError(err)
}

// TestSlinkyIntegration waits for the chain to reach height 5
func (s *SlinkyIntegrationSuite) TestSlinkyIntegration() {
	err := s.chain.WaitForHeight(context.Background(), 5)
	s.Require().NoError(err)
}
