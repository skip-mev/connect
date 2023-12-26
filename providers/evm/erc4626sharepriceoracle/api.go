package erc4626sharepriceoracle

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/skip-mev/slinky/providers/evm"
	oracletypes "github.com/skip-mev/slinky/x/oracle/types"
)

func (h *ERC4626SharePriceAPIHandler) getPriceForPair(pair oracletypes.CurrencyPair) (*big.Int, error) {
	metadata, ok := h.config.TokenNameToMetadata[pair.Quote]
	if !ok {
		return nil, fmt.Errorf("token %s metadata not found", pair.Quote)
	}

	client, err := ethclient.Dial(h.rpcEndpoint)
	if err != nil {
		return nil, err
	}

	contractAddress, found := h.getPairContractAddress(pair)
	if !found {
		return nil, fmt.Errorf("contract address for pair %v not found", pair)
	}

	contract, err := NewERC4626SharePriceOracle(common.HexToAddress(contractAddress), client)
	if err != nil {
		return nil, err
	}

	latest, err := contract.GetLatest(&bind.CallOpts{})
	if err != nil || latest.NotSafeToUse {
		return nil, err
	}

	var _price *big.Int
	if metadata.IsTWAP {
		_price = latest.TimeWeightedAverageAnswer
	} else {
		_price = latest.Ans
	}

	return _price, nil
}

// getPairContractAddress gets the contract address for the pair.
func (h *ERC4626SharePriceAPIHandler) getPairContractAddress(pair oracletypes.CurrencyPair) (string, bool) {
	metadata, found := h.config.TokenNameToMetadata[pair.Quote]
	if found {
		return metadata.Symbol, found
	}

	return "", found
}

// getRPCEndpoint returns the endpoint to fetch prices from.
func getRPCEndpoint(config evm.Config) string {
	return fmt.Sprintf("%s/%s", config.RPCEndpoint, config.APIKey)
}
