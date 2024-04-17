package types

import (
	"crypto/sha256"
	"encoding/binary"
)

func CurrencyPairToID(currencyPair string) uint64 {
	hash := sha256.New()
	hash.Write([]byte(currencyPair))
	md := hash.Sum(nil)
	return binary.LittleEndian.Uint64(md)
}
