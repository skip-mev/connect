package types

import (
	"fmt"
	"strings"
)

const (
	ethereum            = "ETHEREUM"
	MaxCPFieldLength    = 256
	fieldSeparator      = ","
	expectedSplitLength = 3
)

// NewCurrencyPair returns a new CurrencyPair with the given base and quote strings.
func NewCurrencyPair(base, quote string) CurrencyPair {
	return CurrencyPair{
		Base:  base,
		Quote: quote,
	}
}

// IsLegacyAssetString returns true if the asset string is of the following format:
// - contains no instances of fieldSeparator.
func IsLegacyAssetString(asset string) bool {
	return !strings.Contains(asset, fieldSeparator)
}

// coreValidate performs checks that are universal across any ticker format style, namely:
// - check that base and quote are not empty.
// - check that the length of the base and quote fields do not exceed MaxCPFieldLength.
func (cp *CurrencyPair) coreValidate() error {
	if cp.Base == "" {
		return fmt.Errorf("base asset cannot be empty")
	}

	if cp.Quote == "" {
		return fmt.Errorf("quote asset cannot be empty")
	}

	if len(cp.Base) > MaxCPFieldLength {
		return fmt.Errorf("base asset exceeds max length of %d", MaxCPFieldLength)
	}

	if len(cp.Quote) > MaxCPFieldLength {
		return fmt.Errorf("quote asset exceeds max length of %d", MaxCPFieldLength)
	}

	return nil
}

// ValidateBasic checks that the Base / Quote strings in the CurrencyPair are formatted correctly, i.e.
// For base and quote asset:
// check if the asset is formatted in the legacy validation form:
//   - if so, check that fields are not empty and all upper case
//   - else, check that the format is in the following form: tokenName,tokenAddress,chainID
func (cp *CurrencyPair) ValidateBasic() error {
	if err := cp.coreValidate(); err != nil {
		return err
	}

	if IsLegacyAssetString(cp.Base) {
		err := ValidateLegacyAssetString(cp.Base)
		if err != nil {
			return fmt.Errorf("base asset %q is invalid: %w", cp.Base, err)
		}
	} else {
		err := ValidateDefiAssetString(cp.Base)
		if err != nil {
			return fmt.Errorf("base defi asset %q is invalid: %w", cp.Base, err)
		}
	}

	// check quote asset
	if IsLegacyAssetString(cp.Quote) {
		err := ValidateLegacyAssetString(cp.Quote)
		if err != nil {
			return fmt.Errorf("quote asset %q is invalid: %w", cp.Quote, err)
		}
	} else {
		err := ValidateDefiAssetString(cp.Quote)
		if err != nil {
			return fmt.Errorf("quote defi asset %q is invalid: %w", cp.Quote, err)
		}
	}

	return nil
}

// ValidateLegacyAssetString checks if the asset string is formatted correctly, i.e.
// - asset string is fully uppercase
// - asset string does not contain the `fieldSeparator`
//
// NOTE: this function assumes that coreValidate() has already been run.
func ValidateLegacyAssetString(asset string) error {
	// check formatting of asset
	if strings.ToUpper(asset) != asset {
		return fmt.Errorf("incorrectly formatted asset string, expected: %q got: %q", strings.ToUpper(asset), asset)
	}

	if !IsLegacyAssetString(asset) {
		return fmt.Errorf("incorrectly formatted asset string, asset %q should not contain the %q character", asset,
			fieldSeparator)
	}

	return nil
}

// ValidateDefiAssetString checks that the asset string is formatted properly as a defi asset (tokenName,tokenAddress,chainID)
// - check that the length of fields separated by fieldSeparator is expectedSplitLength
// - check that the first split (tokenName) is formatted properly as a LegacyAssetString.
//
// NOTE: this function assumes that coreValidate() has already been run.
func ValidateDefiAssetString(asset string) error {
	token, _, _, err := SplitDefiAssetString(asset)
	if err != nil {
		return err
	}

	// first element is a ticker, so we require it to pass legacy asset validation:
	if err := ValidateLegacyAssetString(token); err != nil {
		return fmt.Errorf("token field %q is invalid: %w", token, err)
	}

	if strings.ToUpper(asset) != asset {
		return fmt.Errorf("incorrectly formatted asset string, expected: %q got: %q", strings.ToUpper(asset), asset)
	}

	return nil
}

// SplitDefiAssetString splits a defi asset by the fieldSeparator and checks that it is the proper length.
// returns the split string as (token, address, chainID).
func SplitDefiAssetString(defiString string) (token, address, chainID string, err error) {
	split := strings.Split(defiString, fieldSeparator)
	if len(split) != expectedSplitLength {
		return "", "", "", fmt.Errorf("asset fields have wrong length, expected: %d got: %d", expectedSplitLength, len(split))
	}
	return split[0], split[1], split[2], nil
}

// LegacyValidateBasic checks that the Base / Quote strings in the CurrencyPair are formatted correctly, i.e.
// Base + Quote are non-empty, and are in upper-case.
func (cp *CurrencyPair) LegacyValidateBasic() error {
	// strings must be valid
	if cp.Base == "" || cp.Quote == "" {
		return fmt.Errorf("empty quote or base string")
	}
	// check formatting of base / quote
	if strings.ToUpper(cp.Base) != cp.Base {
		return fmt.Errorf("incorrectly formatted base string, expected: %q got: %q", strings.ToUpper(cp.Base), cp.Base)
	}
	if strings.ToUpper(cp.Quote) != cp.Quote {
		return fmt.Errorf("incorrectly formatted quote string, expected: %q got: %q", strings.ToUpper(cp.Quote),
			cp.Quote)
	}

	if len(cp.Base) > MaxCPFieldLength || len(cp.Quote) > MaxCPFieldLength {
		return fmt.Errorf("string field exceeds maximum length of %d", MaxCPFieldLength)
	}

	return nil
}

// Invert returns an inverted version of cp (where the Base and Quote are swapped).
func (cp *CurrencyPair) Invert() CurrencyPair {
	return CurrencyPair{
		Base:  cp.Quote,
		Quote: cp.Base,
	}
}

// String returns a string representation of the CurrencyPair, in the following form "ETH/BTC".
func (cp CurrencyPair) String() string {
	return fmt.Sprintf("%s/%s", cp.Base, cp.Quote)
}

// CurrencyPairString constructs and returns the string representation of a currency pair.
func CurrencyPairString(base, quote string) string {
	cp := NewCurrencyPair(base, quote)
	return cp.String()
}

// CurrencyPairFromString creates a currency pair from a string. Non-capitalized inputs are sanitized and the resulting
// currency pair is validated.
func CurrencyPairFromString(s string) (CurrencyPair, error) {
	split := strings.Split(s, "/")
	if len(split) != 2 {
		return CurrencyPair{}, fmt.Errorf("incorrectly formatted CurrencyPair: %q", s)
	}

	base, err := sanitizeAssetString(split[0])
	if err != nil {
		return CurrencyPair{}, err
	}

	quote, err := sanitizeAssetString(split[1])
	if err != nil {
		return CurrencyPair{}, err
	}

	cp := CurrencyPair{
		Base:  base,
		Quote: quote,
	}

	return cp, cp.ValidateBasic()
}

func sanitizeAssetString(s string) (string, error) {
	return strings.ToUpper(s), nil
}

// LegacyDecimals returns the number of decimals that the quote will be reported to. If the quote is Ethereum, then
// the number of decimals is 18. Otherwise, the decimals will be reported as 8.
func (cp *CurrencyPair) LegacyDecimals() int {
	if strings.ToUpper(cp.Quote) == ethereum {
		return 18
	}
	return 8
}

// Equal returns true iff the CurrencyPair is equal to the given CurrencyPair.
func (cp *CurrencyPair) Equal(other CurrencyPair) bool {
	return cp.Base == other.Base && cp.Quote == other.Quote
}
