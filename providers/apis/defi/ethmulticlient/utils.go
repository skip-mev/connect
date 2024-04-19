package ethmulticlient

import "github.com/ethereum/go-ethereum/rpc"

// EthBlockNumberBatchElem returns an initialized BatchElem for the eth_blockNumber call.
func EthBlockNumberBatchElem() rpc.BatchElem {
	var result string
	return rpc.BatchElem{
		Method: "eth_blockNumber",
		Result: &result,
	}
}
