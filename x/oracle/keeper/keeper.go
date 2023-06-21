package keeper

import (
	"encoding/binary"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/skip-mev/slinky/x/oracle/types"
)

// Keeper is the base keeper for the x/oracle module
type Keeper struct {
	storeKey  storetypes.StoreKey
	authority sdk.AccAddress
}

// NewKeeper constructs a new keeper from a store-key + authority account address
func NewKeeper(sk storetypes.StoreKey, authority sdk.AccAddress) Keeper {
	return Keeper{
		storeKey:  sk,
		authority: authority,
	}
}

// RemoveCurrencyPair removes a given CurrencyPair from state, i.e removes its nonce + QuotePrice from the module's store.
func (k Keeper) RemoveCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) {
	// remove quote-price
	k.removeQuotePriceForCurrencyPair(ctx, cp)

	// remove nonce
	k.removeNonceForCurrencyPair(ctx, cp)
}

// removeQuotePriceForCurrencyPair removes the QuotePrice for a given CurrencyPair from the module's state.
func (k Keeper) removeQuotePriceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) {
	store := ctx.KVStore(k.storeKey)

	// remove entry for currency-pair from QuotePrice store
	store.Delete(cp.GetStoreKeyForQuotePrice())
}

// removeNonceForCurrencyPair removes the nonce for a given CurrencyPair from the module's state.
func (k Keeper) removeNonceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) {
	store := ctx.KVStore(k.storeKey)

	// remove entry for currency-pair from nonce store
	store.Delete(cp.GetStoreKeyForNonce())
}

// GetPriceWithNonceForCurrencyPair returns a QuotePriceWithNonce for a given CurrencyPair. The nonce for the QuotePrice represents
// the number of times that a given QuotePrice has been updated. Notice: prefer GetPriceWithNonceForCurrencyPair over GetPriceForCurrencyPair.
func (k Keeper) GetPriceWithNonceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) (types.QuotePriceWithNonce, error) {
	// get the QuotePrice for the currency pair
	qp, err := k.GetPriceForCurrencyPair(ctx, cp)
	if err != nil {
		// only fail if the Price Query failed for a reason other than there being no QuotePrice for cp
		if _, ok := err.(QuotePriceNotExistError); !ok {
			return types.QuotePriceWithNonce{}, err
		}
	}

	// get the nonce
	nonce, err := k.GetNonceForCurrencyPair(ctx, cp)
	if err != nil {
		return types.QuotePriceWithNonce{}, err
	}

	return types.NewQuotePriceWithNonce(qp, nonce), nil
}

// GetNonceForCurrency Pair returns the nonce for a given CurrencyPair. If one has not been stored, return an error.
func (k Keeper) GetNonceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) (uint64, error) {
	store := ctx.KVStore(k.storeKey)

	key := cp.GetStoreKeyForNonce()
	// get current nonce for cp from store
	bz := store.Get(key)

	// return the nonce, if one has not been stored yet alert caller
	if len(bz) == 0 {
		return 0, types.NewCurrencyPairNotExistError(cp.ToString())
	}

	// set the nonce to whatever exists + 1
	return binary.BigEndian.Uint64(bz), nil
}

type QuotePriceNotExistError struct {
	cp string
}

func (e QuotePriceNotExistError) Error() string {
	return fmt.Sprintf("no price updates for CurrencyPair: %s", e.cp)
}

// GetPriceForCurrencyPair retrieves the QuotePrice for a given CurrencyPair. if a QuotePrice does not
// exist for the given CurrencyPair, this function errors and returns an empty QuotePrice
func (k Keeper) GetPriceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) (types.QuotePrice, error) {
	store := ctx.KVStore(k.storeKey)

	// get QuotePrice for CurrencyPair (if any is stored)
	bz := store.Get(cp.GetStoreKeyForQuotePrice())

	if len(bz) == 0 {
		return types.QuotePrice{}, QuotePriceNotExistError{cp.ToString()}
	}

	// unmarshal
	qp := types.QuotePrice{}
	if err := qp.Unmarshal(bz); err != nil {
		return types.QuotePrice{}, err
	}

	return qp, nil
}

// SetPriceForCurrencyPair sets the given QuotePrice for a given CurrencyPair, and updates the CurrencyPair's nonce. Note, no validation is performed on
// either the CurrencyPair or the QuotePrice (it is expected the caller performs this validation).
func (k Keeper) SetPriceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair, qp types.QuotePrice) error {
	store := ctx.KVStore(k.storeKey)

	// marshal QuotePrice
	bz, err := qp.Marshal()
	if err != nil {
		return fmt.Errorf("error marshalling QuotePrice: %v", err)
	}

	// update the nonce for the currency-pair
	k.incrementNonceForCurrencyPair(ctx, cp)

	// set the marshalled QuotePrice to state under the CurrencyPair's store-key
	store.Set(cp.GetStoreKeyForQuotePrice(), bz)

	return nil
}

// Increment the nonce for a given currency pair. This should be called each time that a CurrencyPair
// has a QuotePrice stored for it. This method should only be called when we set a new QuotePrice for a CurrencyPair (i.e SetPriceForCurrencyPair)
func (k Keeper) incrementNonceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) error {
	// get the nonce
	nonce, err := k.GetNonceForCurrencyPair(ctx, cp)
	if err != nil {
		// return err only if the error is not from the CurrencyPair failing to exist in state
		if _, ok := err.(*types.CurrencyPairNotExistError); !ok {
			return err
		}
	} else {
		// if the nonce exists in state, increment
		nonce++
	}

	// set the nonce in state
	k.setNonceForCurrencyPair(ctx, cp, nonce)
	return nil
}

// Set the given nonce in state for the given currency pair
func (k Keeper) setNonceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair, nonce uint64) {
	store := ctx.KVStore(k.storeKey)

	// set the nonce in state
	bz := make([]byte, binary.MaxVarintLen64)
	binary.BigEndian.PutUint64(bz, nonce)

	store.Set(cp.GetStoreKeyForNonce(), bz)
}

// GetAllCurrencyPairs returns all CurrencyPairs that have currently been stored to state.
func (k Keeper) GetAllCurrencyPairs(ctx sdk.Context) []types.CurrencyPair {
	cps := make([]types.CurrencyPair, 0)

	// aggregate CurrencyPairs stored under KeyPrefixNonce
	k.IterateNonces(ctx, func(cp types.CurrencyPair) {
		cps = append(cps, cp)
	})

	return cps
}

// IterateQuotePrices iterates over all CurrencyPairs stored under the QuotePrice key in state, and executes a call-back w/ parameters CurrencyPair + QuotePrice.
// This method errors if there are any errors in the process of unmarshalling CurrencyPairs, or QuotePrices
func (k Keeper) IterateQuotePrices(ctx sdk.Context, cb func(cp types.CurrencyPair, qp types.QuotePrice) error) error {
	// construct iterator func
	f := func(it db.Iterator) error {
		// unmarshal key into a CurrencyPair
		cp, err := types.GetCurrencyPairFromPriceKey(it.Key())
		if err != nil {
			return err
		}

		// unmarshal QuotePrice
		qp := types.QuotePrice{}
		if err := qp.Unmarshal(it.Value()); err != nil {
			return err
		}

		// execute call-back on the values
		return cb(cp, qp)
	}

	// iterate over store w/ KeyPrefixQuotePrice
	return k.iteratorFunc(ctx, types.KeyPrefixQuotePrice, f)
}

// IterateNonces iterates over all CurrencyPairs stored under the nonce-key in state, and executes a call-back taking a CurrencyPair as a parameter.
// This method errors if there are any errors encountered in the process of unmarshalling CurrencyPairs
func (k Keeper) IterateNonces(ctx sdk.Context, cb func(cp types.CurrencyPair)) error {
	// construct iterator func
	f := func(it db.Iterator) error {
		// unmarshal key into a CurrencyPair
		cp, err := types.GetCurrencyPairFromNonceKey(it.Key())
		if err != nil {
			return err
		}

		// execute call-back
		cb(cp)

		return nil
	}

	// iterate store w/ KeyPrefixNonce
	return k.iteratorFunc(ctx, types.KeyPrefixNonce, f)
}

// helper method for iterating over a store w/ a call-back
func (k Keeper) iteratorFunc(ctx sdk.Context, prefix []byte, f func(db.Iterator) error) error {
	// get iterator for store w/ prefix
	store := ctx.KVStore(k.storeKey)
	it := storetypes.KVStorePrefixIterator(store, prefix)

	// close the iterator
	defer it.Close()
	for ; it.Valid(); it.Next() {
		// execute call-back, and return error if necessary
		if err := f(it); err != nil {
			return err
		}
	}
	return nil
}
