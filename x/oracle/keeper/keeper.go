package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/x/oracle/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
}

func NewKeeper(sk storetypes.StoreKey) Keeper {
	return Keeper{
		storeKey: sk,
	}
}

// GetPriceForCurrencyPair retrieves the QuotePrice for a given CurrencyPair. if a QuotePrice does not
// exist for the given CurrencyPair, this function errors and returns an empty QuotePrice
func (k Keeper) GetPriceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) (types.QuotePrice, error) {
	store := ctx.KVStore(k.storeKey)

	// get QuotePrice for CurrencyPair (if any is stored)
	bz := store.Get(cp.GetStoreKeyForCurrencyPair())

	if len(bz) == 0 {
		return types.QuotePrice{}, fmt.Errorf("no QuotePrice price found for CurrencyPair: %s", cp)
	}

	// unmarshal
	tp := types.QuotePrice{}
	if err := tp.Unmarshal(bz); err != nil {
		return types.QuotePrice{}, err
	}

	return tp, nil
}

// SetPriceForCurrencyPair sets the given QuotePrice for a given CurrencyPair. Note, no validation is performed on
// either the CurrencyPair or the QuotePrice (it is expected the caller performs this validation).
func (k Keeper) SetPriceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair, qp types.QuotePrice) error {
	store := ctx.KVStore(k.storeKey)

	// marshal QuotePrice
	bz, err := qp.Marshal()
	if err != nil {
		return fmt.Errorf("error marshalling QuotePrice: %v", err)
	}

	// set the marshalled QuotePrice to state under the CurrencyPair's store-key
	store.Set(cp.GetStoreKeyForCurrencyPair(), bz)

	return nil
}

// GetAllTickers returns all tickers that have currently been stored to state.
func (k Keeper) GetAllTickers(ctx sdk.Context) []types.CurrencyPair {
	store := ctx.KVStore(k.storeKey)

	// iterate over all keys in store
	it := storetypes.KVStorePrefixIterator(store, types.KeyPrefixCurrencyPair)
	cps := make([]types.CurrencyPair, 0)

	// iterate over all keys
	for ; it.Valid(); it.Next() {
		cp, err := types.GetCurrencyPairFromKey(it.Key())
		if err != nil {
			continue
		}
		cps = append(cps, cp)
	}

	return cps
}
