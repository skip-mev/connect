package metrics

const (
	// Metric namespace
	AppNamespace = "app"
	// Metrics labels
	TickerLabel           = "ticker"
	InclusionLabel        = "included"
	ProviderLabel         = "provider"
	StatusLabel           = "status"
	ABCIMethodLabel       = "abci_method"
	ChainIDLabel          = "chain_id"
	ABCIMethodStatusLabel = "abci_method_status"
	MessageTypeLabel      = "message_type"
)

// StatusFromError returns a Labeller that can be used to label metrics based on the error. This
// is used to label metrics based on the error returned from oracle client requests.
func StatusFromError(err error) Labeller {
	if err == nil {
		return Success{}
	}
	return Failure{}
}

// MessageType is an identifier used to represent the different types of data that is transmitted between validators in Slinky.
// This ID is used to paginate metrics corresponding to these messages
type MessageType int

const (
	ExtendedCommit MessageType = iota
	VoteExtension
)

func (m MessageType) String() string {
	switch m {
	case ExtendedCommit:
		return "extended_commit"
	case VoteExtension:
		return "vote_extension"
	default:
		return "not_implemented"
	}
}

// ABCIMethod is an identifier for ABCI methods, this is used to paginate latencies / responses in prometheus
// metrics.
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

// Labeller is an interface that can be implemented by errors to provide a label for prometheus metrics.
type Labeller interface {
	Label() string
}

type Success struct{}

func (s Success) Label() string {
	return "Success"
}

type Failure struct{}

func (f Failure) Label() string {
	return "Failure"
}
