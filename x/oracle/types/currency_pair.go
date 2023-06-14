package types

import (
	"strings"
	"fmt"
)

// ValidateBasic checks that the Base / Quote strings in the CurrencyPair are formatted correctly
func (cp CurrencyPair) ValidateBasic() error {
	// strings must be valid
	if cp.Base == "" || cp.Quote == "" {
		return fmt.Errorf("empty quote or base string")
	}
	// check formatting of base / quote
	if strings.ToUpper(cp.Base) != cp.Base {
		return fmt.Errorf("incorrectly formatted base string, expected: %s got: %s", strings.ToUpper(cp.Base), strings.ToUpper(cp.Base))
	}
	if strings.ToUpper(cp.Quote) != cp.Quote {
		return fmt.Errorf("incorrectly formatted quote string, expected: %s got: %s", strings.ToUpper(cp.Quote), strings.ToUpper(cp.Quote))
	}
	return nil
}

func (cp CurrencyPair) ToString() string {
	return fmt.Sprintf("%s/%s", cp.Base, cp.Quote)
}

func CurrencyPairFromString(s string) CurrencyPair {
	split := strings.Split(s, "/")
	if len(split) != 2 {
		return CurrencyPair{}
	}
	return CurrencyPair{
		Base: split[0],
		Quote: split[1],
	}
}
