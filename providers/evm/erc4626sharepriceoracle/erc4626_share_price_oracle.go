// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc4626sharepriceoracle

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ERC4626SharePriceOracleMetaData contains all meta data concerning the ERC4626SharePriceOracle contract.
var ERC4626SharePriceOracleMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractERC4626\",\"name\":\"_target\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_heartbeat\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_deviationTrigger\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_gracePeriod\",\"type\":\"uint64\"},{\"internalType\":\"uint16\",\"name\":\"_observationsToUse\",\"type\":\"uint16\"},{\"internalType\":\"address\",\"name\":\"_automationRegistry\",\"type\":\"address\"},{\"internalType\":\"uint216\",\"name\":\"_startingAnswer\",\"type\":\"uint216\"},{\"internalType\":\"uint256\",\"name\":\"_allowedAnswerChangeLower\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_allowedAnswerChangeUpper\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ERC4626SharePriceOracle__ContractKillSwitch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC4626SharePriceOracle__CumulativeTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC4626SharePriceOracle__FuturePerformData\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC4626SharePriceOracle__NoUpkeepConditionMet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC4626SharePriceOracle__OnlyCallableByAutomationRegistry\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC4626SharePriceOracle__SharePriceTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC4626SharePriceOracle__StalePerformData\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reportedAnswer\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"minAnswer\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"maxAnswer\",\"type\":\"uint256\"}],\"name\":\"KillSwitchActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeUpdated\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeAnswerCalculated\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"latestAnswer\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeWeightedAverageAnswer\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isNotSafeToUse\",\"type\":\"bool\"}],\"name\":\"OracleUpdated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ONE_SHARE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"allowedAnswerChangeLower\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"allowedAnswerChangeUpper\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"answer\",\"outputs\":[{\"internalType\":\"uint216\",\"name\":\"\",\"type\":\"uint216\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"automationRegistry\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"checkUpkeep\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"upkeepNeeded\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentIndex\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deviationTrigger\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"ans\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeWeightedAverageAnswer\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"notSafeToUse\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLatestAnswer\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gracePeriod\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"heartbeat\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"killSwitch\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"observations\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint192\",\"name\":\"cumulative\",\"type\":\"uint192\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"observationsLength\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"performData\",\"type\":\"bytes\"}],\"name\":\"performUpkeep\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"target\",\"outputs\":[{\"internalType\":\"contractERC4626\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"targetDecimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ERC4626SharePriceOracleABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC4626SharePriceOracleMetaData.ABI instead.
var ERC4626SharePriceOracleABI = ERC4626SharePriceOracleMetaData.ABI

// ERC4626SharePriceOracle is an auto generated Go binding around an Ethereum contract.
type ERC4626SharePriceOracle struct {
	ERC4626SharePriceOracleCaller     // Read-only binding to the contract
	ERC4626SharePriceOracleTransactor // Write-only binding to the contract
	ERC4626SharePriceOracleFilterer   // Log filterer for contract events
}

// ERC4626SharePriceOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC4626SharePriceOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC4626SharePriceOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC4626SharePriceOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC4626SharePriceOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC4626SharePriceOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC4626SharePriceOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC4626SharePriceOracleSession struct {
	Contract     *ERC4626SharePriceOracle // Generic contract binding to set the session for
	CallOpts     bind.CallOpts            // Call options to use throughout this session
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// ERC4626SharePriceOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC4626SharePriceOracleCallerSession struct {
	Contract *ERC4626SharePriceOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                  // Call options to use throughout this session
}

// ERC4626SharePriceOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC4626SharePriceOracleTransactorSession struct {
	Contract     *ERC4626SharePriceOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                  // Transaction auth options to use throughout this session
}

// ERC4626SharePriceOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC4626SharePriceOracleRaw struct {
	Contract *ERC4626SharePriceOracle // Generic contract binding to access the raw methods on
}

// ERC4626SharePriceOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC4626SharePriceOracleCallerRaw struct {
	Contract *ERC4626SharePriceOracleCaller // Generic read-only contract binding to access the raw methods on
}

// ERC4626SharePriceOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC4626SharePriceOracleTransactorRaw struct {
	Contract *ERC4626SharePriceOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC4626SharePriceOracle creates a new instance of ERC4626SharePriceOracle, bound to a specific deployed contract.
func NewERC4626SharePriceOracle(address common.Address, backend bind.ContractBackend) (*ERC4626SharePriceOracle, error) {
	contract, err := bindERC4626SharePriceOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC4626SharePriceOracle{ERC4626SharePriceOracleCaller: ERC4626SharePriceOracleCaller{contract: contract}, ERC4626SharePriceOracleTransactor: ERC4626SharePriceOracleTransactor{contract: contract}, ERC4626SharePriceOracleFilterer: ERC4626SharePriceOracleFilterer{contract: contract}}, nil
}

// NewERC4626SharePriceOracleCaller creates a new read-only instance of ERC4626SharePriceOracle, bound to a specific deployed contract.
func NewERC4626SharePriceOracleCaller(address common.Address, caller bind.ContractCaller) (*ERC4626SharePriceOracleCaller, error) {
	contract, err := bindERC4626SharePriceOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC4626SharePriceOracleCaller{contract: contract}, nil
}

// NewERC4626SharePriceOracleTransactor creates a new write-only instance of ERC4626SharePriceOracle, bound to a specific deployed contract.
func NewERC4626SharePriceOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC4626SharePriceOracleTransactor, error) {
	contract, err := bindERC4626SharePriceOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC4626SharePriceOracleTransactor{contract: contract}, nil
}

// NewERC4626SharePriceOracleFilterer creates a new log filterer instance of ERC4626SharePriceOracle, bound to a specific deployed contract.
func NewERC4626SharePriceOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC4626SharePriceOracleFilterer, error) {
	contract, err := bindERC4626SharePriceOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC4626SharePriceOracleFilterer{contract: contract}, nil
}

// bindERC4626SharePriceOracle binds a generic wrapper to an already deployed contract.
func bindERC4626SharePriceOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC4626SharePriceOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC4626SharePriceOracle.Contract.ERC4626SharePriceOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC4626SharePriceOracle.Contract.ERC4626SharePriceOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC4626SharePriceOracle.Contract.ERC4626SharePriceOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC4626SharePriceOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC4626SharePriceOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC4626SharePriceOracle.Contract.contract.Transact(opts, method, params...)
}

// ONESHARE is a free data retrieval call binding the contract method 0xb7d122b5.
//
// Solidity: function ONE_SHARE() view returns(uint256)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) ONESHARE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "ONE_SHARE")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// ONESHARE is a free data retrieval call binding the contract method 0xb7d122b5.
//
// Solidity: function ONE_SHARE() view returns(uint256)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) ONESHARE() (*big.Int, error) {
	return _ERC4626SharePriceOracle.Contract.ONESHARE(&_ERC4626SharePriceOracle.CallOpts)
}

// ONESHARE is a free data retrieval call binding the contract method 0xb7d122b5.
//
// Solidity: function ONE_SHARE() view returns(uint256)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) ONESHARE() (*big.Int, error) {
	return _ERC4626SharePriceOracle.Contract.ONESHARE(&_ERC4626SharePriceOracle.CallOpts)
}

// AllowedAnswerChangeLower is a free data retrieval call binding the contract method 0x7167adbc.
//
// Solidity: function allowedAnswerChangeLower() view returns(uint256)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) AllowedAnswerChangeLower(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "allowedAnswerChangeLower")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// AllowedAnswerChangeLower is a free data retrieval call binding the contract method 0x7167adbc.
//
// Solidity: function allowedAnswerChangeLower() view returns(uint256)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) AllowedAnswerChangeLower() (*big.Int, error) {
	return _ERC4626SharePriceOracle.Contract.AllowedAnswerChangeLower(&_ERC4626SharePriceOracle.CallOpts)
}

// AllowedAnswerChangeLower is a free data retrieval call binding the contract method 0x7167adbc.
//
// Solidity: function allowedAnswerChangeLower() view returns(uint256)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) AllowedAnswerChangeLower() (*big.Int, error) {
	return _ERC4626SharePriceOracle.Contract.AllowedAnswerChangeLower(&_ERC4626SharePriceOracle.CallOpts)
}

// AllowedAnswerChangeUpper is a free data retrieval call binding the contract method 0x7d4cdb4f.
//
// Solidity: function allowedAnswerChangeUpper() view returns(uint256)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) AllowedAnswerChangeUpper(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "allowedAnswerChangeUpper")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// AllowedAnswerChangeUpper is a free data retrieval call binding the contract method 0x7d4cdb4f.
//
// Solidity: function allowedAnswerChangeUpper() view returns(uint256)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) AllowedAnswerChangeUpper() (*big.Int, error) {
	return _ERC4626SharePriceOracle.Contract.AllowedAnswerChangeUpper(&_ERC4626SharePriceOracle.CallOpts)
}

// AllowedAnswerChangeUpper is a free data retrieval call binding the contract method 0x7d4cdb4f.
//
// Solidity: function allowedAnswerChangeUpper() view returns(uint256)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) AllowedAnswerChangeUpper() (*big.Int, error) {
	return _ERC4626SharePriceOracle.Contract.AllowedAnswerChangeUpper(&_ERC4626SharePriceOracle.CallOpts)
}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(uint216)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) Answer(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "answer")
	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err
}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(uint216)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) Answer() (*big.Int, error) {
	return _ERC4626SharePriceOracle.Contract.Answer(&_ERC4626SharePriceOracle.CallOpts)
}

// Answer is a free data retrieval call binding the contract method 0x85bb7d69.
//
// Solidity: function answer() view returns(uint216)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) Answer() (*big.Int, error) {
	return _ERC4626SharePriceOracle.Contract.Answer(&_ERC4626SharePriceOracle.CallOpts)
}

// AutomationRegistry is a free data retrieval call binding the contract method 0x5dc228a0.
//
// Solidity: function automationRegistry() view returns(address)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) AutomationRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "automationRegistry")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// AutomationRegistry is a free data retrieval call binding the contract method 0x5dc228a0.
//
// Solidity: function automationRegistry() view returns(address)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) AutomationRegistry() (common.Address, error) {
	return _ERC4626SharePriceOracle.Contract.AutomationRegistry(&_ERC4626SharePriceOracle.CallOpts)
}

// AutomationRegistry is a free data retrieval call binding the contract method 0x5dc228a0.
//
// Solidity: function automationRegistry() view returns(address)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) AutomationRegistry() (common.Address, error) {
	return _ERC4626SharePriceOracle.Contract.AutomationRegistry(&_ERC4626SharePriceOracle.CallOpts)
}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes ) view returns(bool upkeepNeeded, bytes performData)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) CheckUpkeep(opts *bind.CallOpts, arg0 []byte) (struct {
	UpkeepNeeded bool
	PerformData  []byte
}, error,
) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "checkUpkeep", arg0)

	outstruct := new(struct {
		UpkeepNeeded bool
		PerformData  []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.UpkeepNeeded = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.PerformData = *abi.ConvertType(out[1], new([]byte)).(*[]byte)

	return *outstruct, err
}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes ) view returns(bool upkeepNeeded, bytes performData)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) CheckUpkeep(arg0 []byte) (struct {
	UpkeepNeeded bool
	PerformData  []byte
}, error,
) {
	return _ERC4626SharePriceOracle.Contract.CheckUpkeep(&_ERC4626SharePriceOracle.CallOpts, arg0)
}

// CheckUpkeep is a free data retrieval call binding the contract method 0x6e04ff0d.
//
// Solidity: function checkUpkeep(bytes ) view returns(bool upkeepNeeded, bytes performData)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) CheckUpkeep(arg0 []byte) (struct {
	UpkeepNeeded bool
	PerformData  []byte
}, error,
) {
	return _ERC4626SharePriceOracle.Contract.CheckUpkeep(&_ERC4626SharePriceOracle.CallOpts, arg0)
}

// CurrentIndex is a free data retrieval call binding the contract method 0x26987b60.
//
// Solidity: function currentIndex() view returns(uint16)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) CurrentIndex(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "currentIndex")
	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err
}

// CurrentIndex is a free data retrieval call binding the contract method 0x26987b60.
//
// Solidity: function currentIndex() view returns(uint16)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) CurrentIndex() (uint16, error) {
	return _ERC4626SharePriceOracle.Contract.CurrentIndex(&_ERC4626SharePriceOracle.CallOpts)
}

// CurrentIndex is a free data retrieval call binding the contract method 0x26987b60.
//
// Solidity: function currentIndex() view returns(uint16)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) CurrentIndex() (uint16, error) {
	return _ERC4626SharePriceOracle.Contract.CurrentIndex(&_ERC4626SharePriceOracle.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "decimals")
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) Decimals() (uint8, error) {
	return _ERC4626SharePriceOracle.Contract.Decimals(&_ERC4626SharePriceOracle.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) Decimals() (uint8, error) {
	return _ERC4626SharePriceOracle.Contract.Decimals(&_ERC4626SharePriceOracle.CallOpts)
}

// DeviationTrigger is a free data retrieval call binding the contract method 0xfeabaa02.
//
// Solidity: function deviationTrigger() view returns(uint64)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) DeviationTrigger(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "deviationTrigger")
	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err
}

// DeviationTrigger is a free data retrieval call binding the contract method 0xfeabaa02.
//
// Solidity: function deviationTrigger() view returns(uint64)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) DeviationTrigger() (uint64, error) {
	return _ERC4626SharePriceOracle.Contract.DeviationTrigger(&_ERC4626SharePriceOracle.CallOpts)
}

// DeviationTrigger is a free data retrieval call binding the contract method 0xfeabaa02.
//
// Solidity: function deviationTrigger() view returns(uint64)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) DeviationTrigger() (uint64, error) {
	return _ERC4626SharePriceOracle.Contract.DeviationTrigger(&_ERC4626SharePriceOracle.CallOpts)
}

// GetLatest is a free data retrieval call binding the contract method 0xc36af460.
//
// Solidity: function getLatest() view returns(uint256 ans, uint256 timeWeightedAverageAnswer, bool notSafeToUse)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) GetLatest(opts *bind.CallOpts) (struct {
	Ans                       *big.Int
	TimeWeightedAverageAnswer *big.Int
	NotSafeToUse              bool
}, error,
) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "getLatest")

	outstruct := new(struct {
		Ans                       *big.Int
		TimeWeightedAverageAnswer *big.Int
		NotSafeToUse              bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Ans = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TimeWeightedAverageAnswer = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.NotSafeToUse = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err
}

// GetLatest is a free data retrieval call binding the contract method 0xc36af460.
//
// Solidity: function getLatest() view returns(uint256 ans, uint256 timeWeightedAverageAnswer, bool notSafeToUse)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) GetLatest() (struct {
	Ans                       *big.Int
	TimeWeightedAverageAnswer *big.Int
	NotSafeToUse              bool
}, error,
) {
	return _ERC4626SharePriceOracle.Contract.GetLatest(&_ERC4626SharePriceOracle.CallOpts)
}

// GetLatest is a free data retrieval call binding the contract method 0xc36af460.
//
// Solidity: function getLatest() view returns(uint256 ans, uint256 timeWeightedAverageAnswer, bool notSafeToUse)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) GetLatest() (struct {
	Ans                       *big.Int
	TimeWeightedAverageAnswer *big.Int
	NotSafeToUse              bool
}, error,
) {
	return _ERC4626SharePriceOracle.Contract.GetLatest(&_ERC4626SharePriceOracle.CallOpts)
}

// GetLatestAnswer is a free data retrieval call binding the contract method 0x96237c02.
//
// Solidity: function getLatestAnswer() view returns(uint256, bool)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) GetLatestAnswer(opts *bind.CallOpts) (*big.Int, bool, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "getLatestAnswer")
	if err != nil {
		return *new(*big.Int), *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(bool)).(*bool)

	return out0, out1, err
}

// GetLatestAnswer is a free data retrieval call binding the contract method 0x96237c02.
//
// Solidity: function getLatestAnswer() view returns(uint256, bool)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) GetLatestAnswer() (*big.Int, bool, error) {
	return _ERC4626SharePriceOracle.Contract.GetLatestAnswer(&_ERC4626SharePriceOracle.CallOpts)
}

// GetLatestAnswer is a free data retrieval call binding the contract method 0x96237c02.
//
// Solidity: function getLatestAnswer() view returns(uint256, bool)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) GetLatestAnswer() (*big.Int, bool, error) {
	return _ERC4626SharePriceOracle.Contract.GetLatestAnswer(&_ERC4626SharePriceOracle.CallOpts)
}

// GracePeriod is a free data retrieval call binding the contract method 0xa06db7dc.
//
// Solidity: function gracePeriod() view returns(uint64)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) GracePeriod(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "gracePeriod")
	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err
}

// GracePeriod is a free data retrieval call binding the contract method 0xa06db7dc.
//
// Solidity: function gracePeriod() view returns(uint64)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) GracePeriod() (uint64, error) {
	return _ERC4626SharePriceOracle.Contract.GracePeriod(&_ERC4626SharePriceOracle.CallOpts)
}

// GracePeriod is a free data retrieval call binding the contract method 0xa06db7dc.
//
// Solidity: function gracePeriod() view returns(uint64)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) GracePeriod() (uint64, error) {
	return _ERC4626SharePriceOracle.Contract.GracePeriod(&_ERC4626SharePriceOracle.CallOpts)
}

// Heartbeat is a free data retrieval call binding the contract method 0x3defb962.
//
// Solidity: function heartbeat() view returns(uint64)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) Heartbeat(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "heartbeat")
	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err
}

// Heartbeat is a free data retrieval call binding the contract method 0x3defb962.
//
// Solidity: function heartbeat() view returns(uint64)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) Heartbeat() (uint64, error) {
	return _ERC4626SharePriceOracle.Contract.Heartbeat(&_ERC4626SharePriceOracle.CallOpts)
}

// Heartbeat is a free data retrieval call binding the contract method 0x3defb962.
//
// Solidity: function heartbeat() view returns(uint64)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) Heartbeat() (uint64, error) {
	return _ERC4626SharePriceOracle.Contract.Heartbeat(&_ERC4626SharePriceOracle.CallOpts)
}

// KillSwitch is a free data retrieval call binding the contract method 0xada14698.
//
// Solidity: function killSwitch() view returns(bool)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) KillSwitch(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "killSwitch")
	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err
}

// KillSwitch is a free data retrieval call binding the contract method 0xada14698.
//
// Solidity: function killSwitch() view returns(bool)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) KillSwitch() (bool, error) {
	return _ERC4626SharePriceOracle.Contract.KillSwitch(&_ERC4626SharePriceOracle.CallOpts)
}

// KillSwitch is a free data retrieval call binding the contract method 0xada14698.
//
// Solidity: function killSwitch() view returns(bool)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) KillSwitch() (bool, error) {
	return _ERC4626SharePriceOracle.Contract.KillSwitch(&_ERC4626SharePriceOracle.CallOpts)
}

// Observations is a free data retrieval call binding the contract method 0x252c09d7.
//
// Solidity: function observations(uint256 ) view returns(uint64 timestamp, uint192 cumulative)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) Observations(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Timestamp  uint64
	Cumulative *big.Int
}, error,
) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "observations", arg0)

	outstruct := new(struct {
		Timestamp  uint64
		Cumulative *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Timestamp = *abi.ConvertType(out[0], new(uint64)).(*uint64)
	outstruct.Cumulative = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err
}

// Observations is a free data retrieval call binding the contract method 0x252c09d7.
//
// Solidity: function observations(uint256 ) view returns(uint64 timestamp, uint192 cumulative)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) Observations(arg0 *big.Int) (struct {
	Timestamp  uint64
	Cumulative *big.Int
}, error,
) {
	return _ERC4626SharePriceOracle.Contract.Observations(&_ERC4626SharePriceOracle.CallOpts, arg0)
}

// Observations is a free data retrieval call binding the contract method 0x252c09d7.
//
// Solidity: function observations(uint256 ) view returns(uint64 timestamp, uint192 cumulative)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) Observations(arg0 *big.Int) (struct {
	Timestamp  uint64
	Cumulative *big.Int
}, error,
) {
	return _ERC4626SharePriceOracle.Contract.Observations(&_ERC4626SharePriceOracle.CallOpts, arg0)
}

// ObservationsLength is a free data retrieval call binding the contract method 0xcd7f5ce2.
//
// Solidity: function observationsLength() view returns(uint16)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) ObservationsLength(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "observationsLength")
	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err
}

// ObservationsLength is a free data retrieval call binding the contract method 0xcd7f5ce2.
//
// Solidity: function observationsLength() view returns(uint16)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) ObservationsLength() (uint16, error) {
	return _ERC4626SharePriceOracle.Contract.ObservationsLength(&_ERC4626SharePriceOracle.CallOpts)
}

// ObservationsLength is a free data retrieval call binding the contract method 0xcd7f5ce2.
//
// Solidity: function observationsLength() view returns(uint16)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) ObservationsLength() (uint16, error) {
	return _ERC4626SharePriceOracle.Contract.ObservationsLength(&_ERC4626SharePriceOracle.CallOpts)
}

// Target is a free data retrieval call binding the contract method 0xd4b83992.
//
// Solidity: function target() view returns(address)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) Target(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "target")
	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err
}

// Target is a free data retrieval call binding the contract method 0xd4b83992.
//
// Solidity: function target() view returns(address)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) Target() (common.Address, error) {
	return _ERC4626SharePriceOracle.Contract.Target(&_ERC4626SharePriceOracle.CallOpts)
}

// Target is a free data retrieval call binding the contract method 0xd4b83992.
//
// Solidity: function target() view returns(address)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) Target() (common.Address, error) {
	return _ERC4626SharePriceOracle.Contract.Target(&_ERC4626SharePriceOracle.CallOpts)
}

// TargetDecimals is a free data retrieval call binding the contract method 0x36c1387e.
//
// Solidity: function targetDecimals() view returns(uint8)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCaller) TargetDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ERC4626SharePriceOracle.contract.Call(opts, &out, "targetDecimals")
	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err
}

// TargetDecimals is a free data retrieval call binding the contract method 0x36c1387e.
//
// Solidity: function targetDecimals() view returns(uint8)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) TargetDecimals() (uint8, error) {
	return _ERC4626SharePriceOracle.Contract.TargetDecimals(&_ERC4626SharePriceOracle.CallOpts)
}

// TargetDecimals is a free data retrieval call binding the contract method 0x36c1387e.
//
// Solidity: function targetDecimals() view returns(uint8)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleCallerSession) TargetDecimals() (uint8, error) {
	return _ERC4626SharePriceOracle.Contract.TargetDecimals(&_ERC4626SharePriceOracle.CallOpts)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleTransactor) PerformUpkeep(opts *bind.TransactOpts, performData []byte) (*types.Transaction, error) {
	return _ERC4626SharePriceOracle.contract.Transact(opts, "performUpkeep", performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _ERC4626SharePriceOracle.Contract.PerformUpkeep(&_ERC4626SharePriceOracle.TransactOpts, performData)
}

// PerformUpkeep is a paid mutator transaction binding the contract method 0x4585e33b.
//
// Solidity: function performUpkeep(bytes performData) returns()
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleTransactorSession) PerformUpkeep(performData []byte) (*types.Transaction, error) {
	return _ERC4626SharePriceOracle.Contract.PerformUpkeep(&_ERC4626SharePriceOracle.TransactOpts, performData)
}

// ERC4626SharePriceOracleKillSwitchActivatedIterator is returned from FilterKillSwitchActivated and is used to iterate over the raw logs and unpacked data for KillSwitchActivated events raised by the ERC4626SharePriceOracle contract.
type ERC4626SharePriceOracleKillSwitchActivatedIterator struct {
	Event *ERC4626SharePriceOracleKillSwitchActivated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC4626SharePriceOracleKillSwitchActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC4626SharePriceOracleKillSwitchActivated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ERC4626SharePriceOracleKillSwitchActivated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ERC4626SharePriceOracleKillSwitchActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC4626SharePriceOracleKillSwitchActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC4626SharePriceOracleKillSwitchActivated represents a KillSwitchActivated event raised by the ERC4626SharePriceOracle contract.
type ERC4626SharePriceOracleKillSwitchActivated struct {
	ReportedAnswer *big.Int
	MinAnswer      *big.Int
	MaxAnswer      *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterKillSwitchActivated is a free log retrieval operation binding the contract event 0x472eb7b5f33f38b3f139fef7fc88820178d1cf9d0abffee988953791a84d66da.
//
// Solidity: event KillSwitchActivated(uint256 reportedAnswer, uint256 minAnswer, uint256 maxAnswer)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleFilterer) FilterKillSwitchActivated(opts *bind.FilterOpts) (*ERC4626SharePriceOracleKillSwitchActivatedIterator, error) {
	logs, sub, err := _ERC4626SharePriceOracle.contract.FilterLogs(opts, "KillSwitchActivated")
	if err != nil {
		return nil, err
	}
	return &ERC4626SharePriceOracleKillSwitchActivatedIterator{contract: _ERC4626SharePriceOracle.contract, event: "KillSwitchActivated", logs: logs, sub: sub}, nil
}

// WatchKillSwitchActivated is a free log subscription operation binding the contract event 0x472eb7b5f33f38b3f139fef7fc88820178d1cf9d0abffee988953791a84d66da.
//
// Solidity: event KillSwitchActivated(uint256 reportedAnswer, uint256 minAnswer, uint256 maxAnswer)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleFilterer) WatchKillSwitchActivated(opts *bind.WatchOpts, sink chan<- *ERC4626SharePriceOracleKillSwitchActivated) (event.Subscription, error) {
	logs, sub, err := _ERC4626SharePriceOracle.contract.WatchLogs(opts, "KillSwitchActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC4626SharePriceOracleKillSwitchActivated)
				if err := _ERC4626SharePriceOracle.contract.UnpackLog(event, "KillSwitchActivated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseKillSwitchActivated is a log parse operation binding the contract event 0x472eb7b5f33f38b3f139fef7fc88820178d1cf9d0abffee988953791a84d66da.
//
// Solidity: event KillSwitchActivated(uint256 reportedAnswer, uint256 minAnswer, uint256 maxAnswer)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleFilterer) ParseKillSwitchActivated(log types.Log) (*ERC4626SharePriceOracleKillSwitchActivated, error) {
	event := new(ERC4626SharePriceOracleKillSwitchActivated)
	if err := _ERC4626SharePriceOracle.contract.UnpackLog(event, "KillSwitchActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ERC4626SharePriceOracleOracleUpdatedIterator is returned from FilterOracleUpdated and is used to iterate over the raw logs and unpacked data for OracleUpdated events raised by the ERC4626SharePriceOracle contract.
type ERC4626SharePriceOracleOracleUpdatedIterator struct {
	Event *ERC4626SharePriceOracleOracleUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ERC4626SharePriceOracleOracleUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ERC4626SharePriceOracleOracleUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ERC4626SharePriceOracleOracleUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ERC4626SharePriceOracleOracleUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ERC4626SharePriceOracleOracleUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ERC4626SharePriceOracleOracleUpdated represents a OracleUpdated event raised by the ERC4626SharePriceOracle contract.
type ERC4626SharePriceOracleOracleUpdated struct {
	TimeUpdated               *big.Int
	TimeAnswerCalculated      *big.Int
	LatestAnswer              *big.Int
	TimeWeightedAverageAnswer *big.Int
	IsNotSafeToUse            bool
	Raw                       types.Log // Blockchain specific contextual infos
}

// FilterOracleUpdated is a free log retrieval operation binding the contract event 0x8c2ecabee4ec920ce555679e0efd8aa525aee64dc39f4a13a093f9fc72893fdf.
//
// Solidity: event OracleUpdated(uint256 timeUpdated, uint256 timeAnswerCalculated, uint256 latestAnswer, uint256 timeWeightedAverageAnswer, bool isNotSafeToUse)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleFilterer) FilterOracleUpdated(opts *bind.FilterOpts) (*ERC4626SharePriceOracleOracleUpdatedIterator, error) {
	logs, sub, err := _ERC4626SharePriceOracle.contract.FilterLogs(opts, "OracleUpdated")
	if err != nil {
		return nil, err
	}
	return &ERC4626SharePriceOracleOracleUpdatedIterator{contract: _ERC4626SharePriceOracle.contract, event: "OracleUpdated", logs: logs, sub: sub}, nil
}

// WatchOracleUpdated is a free log subscription operation binding the contract event 0x8c2ecabee4ec920ce555679e0efd8aa525aee64dc39f4a13a093f9fc72893fdf.
//
// Solidity: event OracleUpdated(uint256 timeUpdated, uint256 timeAnswerCalculated, uint256 latestAnswer, uint256 timeWeightedAverageAnswer, bool isNotSafeToUse)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleFilterer) WatchOracleUpdated(opts *bind.WatchOpts, sink chan<- *ERC4626SharePriceOracleOracleUpdated) (event.Subscription, error) {
	logs, sub, err := _ERC4626SharePriceOracle.contract.WatchLogs(opts, "OracleUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ERC4626SharePriceOracleOracleUpdated)
				if err := _ERC4626SharePriceOracle.contract.UnpackLog(event, "OracleUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOracleUpdated is a log parse operation binding the contract event 0x8c2ecabee4ec920ce555679e0efd8aa525aee64dc39f4a13a093f9fc72893fdf.
//
// Solidity: event OracleUpdated(uint256 timeUpdated, uint256 timeAnswerCalculated, uint256 latestAnswer, uint256 timeWeightedAverageAnswer, bool isNotSafeToUse)
func (_ERC4626SharePriceOracle *ERC4626SharePriceOracleFilterer) ParseOracleUpdated(log types.Log) (*ERC4626SharePriceOracleOracleUpdated, error) {
	event := new(ERC4626SharePriceOracleOracleUpdated)
	if err := _ERC4626SharePriceOracle.contract.UnpackLog(event, "OracleUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
