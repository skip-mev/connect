package petri_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/connect/v2/tests/petri"
)

// TestConnectIntegration runs all petri connect testapp tests
func TestConnectIntegration(t *testing.T) {
	chainCfg := petri.GetChainConfig()
	suite.Run(t, petri.NewConnectIntegrationSuite(&chainCfg))
}
