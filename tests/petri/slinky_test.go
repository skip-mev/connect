package petri_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/skip-mev/slinky/tests/petri"
)

// TestSlinkyIntegration runs all of the petri slinky testapp tests
func TestSlinkyIntegration(t *testing.T) {
	chainCfg := petri.GetChainConfig()
	suite.Run(t, petri.NewSlinkyIntegrationSuite(&chainCfg))
}
