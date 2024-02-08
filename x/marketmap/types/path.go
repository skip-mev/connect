package types

//

// NewOperation returns a new Operation instance.
func NewOperation(ticker Ticker, invert bool) (Operation, error) {
	o := Operation{
		Ticker: ticker,
		Invert: invert,
	}

	if err := o.ValidateBasic(); err != nil {
		return Operation{}, err
	}

	return o, nil
}

// ValidateBasic performs basic validation on the Operation.
func (o *Operation) ValidateBasic() error {
	if err := o.Ticker.ValidateBasic(); err != nil {
		return err
	}

	return nil
}
