package petri

import (
	"context"

	petritypes "github.com/skip-mev/petri/types/v2"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

// ConnectIntegrationSuite is the testify suite for building and running connect testapp networks using petri
type ConnectIntegrationSuite struct {
	suite.Suite

	logger *zap.Logger

	spec *petritypes.ChainConfig

	chain petritypes.ChainI
}

func NewConnectIntegrationSuite(spec *petritypes.ChainConfig) *ConnectIntegrationSuite {
	return &ConnectIntegrationSuite{
		spec: spec,
	}
}

func (s *ConnectIntegrationSuite) SetupSuite() {
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

func (s *ConnectIntegrationSuite) TearDownSuite() {
	err := s.chain.Teardown(context.Background())
	s.Require().NoError(err)
}

// TestConnectIntegration waits for the chain to reach height 5
func (s *ConnectIntegrationSuite) TestConnectIntegration() {
	err := s.chain.WaitForHeight(context.Background(), 5)
	s.Require().NoError(err)
}
