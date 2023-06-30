package e2e

// TestSetUp tests that the e2e suite was spun up correctly.
func (s *IntegrationTestSuite) TestSetUp() {
	// Ensure there are some CPs set on genesis.
	currencyPairs := s.queryAllCurrencyPairs()
	s.Require().Greater(len(currencyPairs), 0)

	// Ensure that the height is incrementing.
	currentHeight := s.queryCurrentHeight()
	s.Require().Greater(currentHeight, uint64(0))
	s.waitForABlock()
	s.Require().Equal(currentHeight+1, s.queryCurrentHeight())

	// Ensure that the validator set is correct.
	validators := s.queryValidators()
	s.Require().Len(validators, numValidators)

	// Able to query the price for the CP.
	resp, err := s.queryPriceForCurrencyPair("BITCOIN", "USD")
	s.Require().NoError(err)
	s.Require().NotNil(resp)
	s.Require().Equal(resp.Decimals, uint64(8))
}
