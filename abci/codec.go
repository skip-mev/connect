package abci

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/abci/types"
)

// NewOracleTxDecoder wraps the BaseApp's TxDecoder so that injected oracle data does not fail in tx decoding.
func NewOracleTxDecoder(td sdk.TxDecoder) sdk.TxDecoder { // TODO(nikhil): Deprecate this once the upstream changes to SDK are merged: https://github.com/cosmos/cosmos-sdk/pull/16700
	return func(txBytes []byte) (sdk.Tx, error) {
		// attempt to unmarshal into oracle-data first, if no error return nil
		if err := (&types.OracleData{}).Unmarshal(txBytes); err == nil {
			return nil, nil
		}

		// otherwise, decode the transaction
		return td(txBytes)
	}
}
