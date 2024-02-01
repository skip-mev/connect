package petri_test

import (
	"github.com/skip-mev/slinky/tests/petri"
	"github.com/stretchr/testify/suite"
	"testing"
)

// TestSlinkyIntegration runs all of the petri slinky testapp tests
func TestSlinkyIntegration(t *testing.T) {
	chainCfg := petri.GetChainConfig()
	suite.Run(t, petri.NewSlinkyIntegrationSuite(&chainCfg))
}
