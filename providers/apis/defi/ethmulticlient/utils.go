package ethmulticlient

import "github.com/ethereum/go-ethereum/rpc"

func EthBlockNumberBatchElem() rpc.BatchElem {
	var result string
	return rpc.BatchElem{
		Method: "eth_blockNumber",
		Result: &result,
	}
}
