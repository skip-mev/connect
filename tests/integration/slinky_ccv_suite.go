package integration

// Type SlinkyCCVSuite is a testing-suite for testing slinky's integration with ics consumer chains
type SlinkyCCVSuite struct {
	*SlinkyIntegrationSuite
}

func (s *SlinkyCCVSuite) TestCCVAggregation() {
	// test that prices are reported as expected when stake-weight is the same across validators
	
	// test that when provider stake-weight changes, the price changes accordingly
}
