package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module.
	ModuleName = "incentives"
	// StoreKey is the store key string for incentives.
	StoreKey = ModuleName
)

const (
	// keyPrefixIncentive is the root key prefix under which all incentives are stored.
	keyPrefixIncentive = iota
	// keyPrefixCount is the key prefix used to index incentives.
	keyPrefixCount
)

var (
	// KeyPrefixIncentive is the root key prefix under which all incentives are stored.
	KeyPrefixIncentive = []byte{keyPrefixIncentive}
	// KeyPrefixCount is the key prefix used to index incentives.
	KeyPrefixCount = []byte{keyPrefixCount}
)

// GetIncentiveKey gets the store key for an incentive.
func GetIncentiveKey(incentive Incentive) []byte {
	return append(KeyPrefixIncentive, []byte(incentive.Type())...)
}

// GetIncentiveKeyWithIndex gets the store key for an incentive with an index.
func GetIncentiveKeyWithIndex(incentive Incentive, index uint64) []byte {
	return append(GetIncentiveKey(incentive), sdk.Uint64ToBigEndian(index)...)
}

// GetIncentiveCountKey gets the store key for the incentive count.
func GetIncentiveCountKey(incentive Incentive) []byte {
	return append(KeyPrefixCount, []byte(incentive.Type())...)
}
