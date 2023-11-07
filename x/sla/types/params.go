package types

// DefaultParams returns default incentive parameters.
func DefaultParams() Params {
	return Params{
		Enabled: true,
	}
}

// NewParams returns a new Params instance.
func NewParams(enabled bool) Params {
	return Params{
		Enabled: enabled,
	}
}
