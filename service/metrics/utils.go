package metrics

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
