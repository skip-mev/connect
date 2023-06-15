package keeper

import (
	"encoding/binary"
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

// GetPriceWithNonceForCurrencyPair returns a QuotePriceWithNonce for a given CurrencyPair. The nonce for the QuotePrice represents
// the number of times that a given QuotePrice has been updated. Notice: prefer GetPriceWithNonceForCurrencyPair over GetPriceForCurrencyPair.
func (k Keeper) GetPriceWithNonceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) (types.QuotePriceWithNonce, error) {
	// get the QuotePrice for the currency pair
	qp, err := k.GetPriceForCurrencyPair(ctx, cp)
	if err != nil {
		return types.QuotePriceWithNonce{}, err
	}

	// get the nonce
	store := ctx.KVStore(k.storeKey)
	nonce, _ := getNonceForCurrencyPair(store, cp)

	return types.NewQuotePriceWithNonce(qp, nonce), nil
}

// GetPriceForCurrencyPair retrieves the QuotePrice for a given CurrencyPair. if a QuotePrice does not
// exist for the given CurrencyPair, this function errors and returns an empty QuotePrice
func (k Keeper) GetPriceForCurrencyPair(ctx sdk.Context, cp types.CurrencyPair) (types.QuotePrice, error) {
	store := ctx.KVStore(k.storeKey)

	// get QuotePrice for CurrencyPair (if any is stored)
	bz := store.Get(cp.GetStoreKeyForQuotePrice())

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
	incrementNonceForCurrencyPair(store, cp)

	// set the marshalled QuotePrice to state under the CurrencyPair's store-key
	store.Set(cp.GetStoreKeyForQuotePrice(), bz)

	return nil
}

// Increment the nonce for a given currency pair. This should be called each time that a CurrencyPair
// has a QuotePrice stored for it. This method should only be called when we set a new QuotePrice for a CurrencyPair (i.e SetPriceForCurrencyPair)
func incrementNonceForCurrencyPair(store storetypes.KVStore, cp types.CurrencyPair) {
	// get the nonce
	nonce, ok := getNonceForCurrencyPair(store, cp)

	// if the nonce exists in state, increment
	if ok {
		nonce++
	}

	// set the nonce in state
	bz := make([]byte, binary.MaxVarintLen64)
	binary.BigEndian.PutUint64(bz, nonce)

	store.Set(cp.GetStoreKeyForNonce(), bz)
}

// get the nonce for a given currency pair, if a nonce does not exist yet, then return (0, false). Otherwise, return (nonce, true).
func getNonceForCurrencyPair(store storetypes.KVStore, cp types.CurrencyPair) (uint64, bool) {
	key := cp.GetStoreKeyForNonce()
	// get current nonce for cp from store
	bz := store.Get(key)

	// return the nonce, if one has not been stored yet alert caller
	if len(bz) == 0 {
		return 0, false
	}

	// set the nonce to whatever exists + 1
	return binary.BigEndian.Uint64(bz), true
}

// GetAllTickers returns all tickers that have currently been stored to state.
func (k Keeper) GetAllTickers(ctx sdk.Context) []types.CurrencyPair {
	store := ctx.KVStore(k.storeKey)

	// iterate over all keys in store
	it := storetypes.KVStorePrefixIterator(store, types.KeyPrefixQuotePrice)
	defer it.Close()
	cps := make([]types.CurrencyPair, 0)

	// iterate over all keys
	for ; it.Valid(); it.Next() {
		cp, err := types.GetCurrencyPairFromQuotePriceKey(it.Key())
		if err != nil {
			continue
		}
		cps = append(cps, cp)
	}

	return cps
}
