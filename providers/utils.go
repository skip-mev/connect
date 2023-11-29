package providers

import (
	"context"
	"io"
	"math/big"
	"net/http"
	"strconv"
)

type (
	// ReadFn is a convenience type for reading from a HTTP response body
	ReadFn func([]byte) error

	// ReqFn is a convenience type for adding headers, etc. to an HTTP request header
	ReqFn func(*http.Request)
)

// Float64StringToBigInt converts a float64 string to a big.Int.
func Float64StringToBigInt(s string, decimals int) (*big.Int, error) {
	floatNum, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, err
	}

	return Float64ToBigInt(floatNum, decimals), nil
}

// Float64ToBigInt converts a float64 to a big.Int.
//
// TODO: Is there a better approach to this?
func Float64ToBigInt(val float64, decimals int) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)

	coin := new(big.Float)
	factor := big.NewInt(1).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	coin.SetInt(factor)

	bigval.Mul(bigval, coin)

	result := new(big.Int)
	bigval.Int(result) // store converted number in result

	return result
}

// GetWithContext provides logic for making an http get request, whose duration is bounded / controlled by a given context.
func GetWithContext(ctx context.Context, url string, reader ReadFn) error {
	return GetWithContextAndHeader(ctx, url, reader, nil)
}

// GetWithContextAndHeader provides logic for making an http get request, whose duration is bounded / controlled by a given context, and also updating
// fields in the header of the request
func GetWithContextAndHeader(ctx context.Context, url string, reader ReadFn, reqfn ReqFn) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	if reqfn != nil {
		reqfn(req)
	}

	// execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return reader(body)
}
