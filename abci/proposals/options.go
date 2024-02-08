package proposals

// Option is a function that enables optional configuration of the ProposalHandler.
type Option func(*ProposalHandler)

// RetainOracleDataInWrappedProposalHandler returns an Option that configures the
// ProposalHandler to pass the injected extend-commit-info to the wrapped proposal handler.
func RetainOracleDataInWrappedProposalHandler() Option {
	return func(p *ProposalHandler) {
		p.retainOracleDataInWrappedHandler = true
	}
}
