package metrics

const (
	TickerLabel           = "ticker"
	InclusionLabel        = "included"
	AppNamespace          = "app"
	ProviderLabel         = "provider"
	StatusLabel           = "status"
	ABCIMethodLabel       = "abci_method"
	ChainIDLabel          = "chain_id"
	ABCIMethodStatusLabel = "abci_method_status"
)

type Status int

const (
	StatusFailure Status = iota
	StatusSuccess
)

func (s Status) String() string {
	switch s {
	case StatusFailure:
		return "failure"
	case StatusSuccess:
		return "success"
	default:
		return "unknown"
	}
}

func StatusFromError(err error) Status {
	if err == nil {
		return StatusSuccess
	}
	return StatusFailure
}

type ABCIMethod int

const (
	PrepareProposal ABCIMethod = iota
	ProcessProposal
	ExtendVote
	VerifyVoteExtension
	PreBlock
)

func (a ABCIMethod) String() string {
	switch a {
	case PrepareProposal:
		return "prepare_proposal"
	case ProcessProposal:
		return "process_proposal"
	case ExtendVote:
		return "extend_vote"
	case VerifyVoteExtension:
		return "verify_vote_extension"
	case PreBlock:
		return "pre_blocker"
	default:
		return "not_implemented"
	}
}

type Labeller interface {
	Label() string
}

type Success struct{}

func (s Success) Label() string {
	return "Success"
}
