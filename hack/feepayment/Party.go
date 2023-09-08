// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main 


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
	_ = abi.ConvertType
)

// TeaPartyMetaData contains all meta data concerning the TeaParty contract.
var TeaPartyMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"transactionId\",\"type\":\"uint256\"}],\"name\":\"TransactionCreated\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"addHolder\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"createTransaction\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getHolders\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getTeaPartyTransactions\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"participant\",\"type\":\"address\"}],\"name\":\"getTransaction\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"openTransactions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"participant\",\"type\":\"address\"}],\"name\":\"removeTransaction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// TeaPartyABI is the input ABI used to generate the binding from.
// Deprecated: Use TeaPartyMetaData.ABI instead.
var TeaPartyABI = TeaPartyMetaData.ABI

// TeaParty is an auto generated Go binding around an Ethereum contract.
type TeaParty struct {
	TeaPartyCaller     // Read-only binding to the contract
	TeaPartyTransactor // Write-only binding to the contract
	TeaPartyFilterer   // Log filterer for contract events
}

// TeaPartyCaller is an auto generated read-only Go binding around an Ethereum contract.
type TeaPartyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TeaPartyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TeaPartyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TeaPartyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TeaPartyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TeaPartySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TeaPartySession struct {
	Contract     *TeaParty         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TeaPartyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TeaPartyCallerSession struct {
	Contract *TeaPartyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TeaPartyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TeaPartyTransactorSession struct {
	Contract     *TeaPartyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TeaPartyRaw is an auto generated low-level Go binding around an Ethereum contract.
type TeaPartyRaw struct {
	Contract *TeaParty // Generic contract binding to access the raw methods on
}

// TeaPartyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TeaPartyCallerRaw struct {
	Contract *TeaPartyCaller // Generic read-only contract binding to access the raw methods on
}

// TeaPartyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TeaPartyTransactorRaw struct {
	Contract *TeaPartyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTeaParty creates a new instance of TeaParty, bound to a specific deployed contract.
func NewTeaParty(address common.Address, backend bind.ContractBackend) (*TeaParty, error) {
	contract, err := bindTeaParty(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TeaParty{TeaPartyCaller: TeaPartyCaller{contract: contract}, TeaPartyTransactor: TeaPartyTransactor{contract: contract}, TeaPartyFilterer: TeaPartyFilterer{contract: contract}}, nil
}

// NewTeaPartyCaller creates a new read-only instance of TeaParty, bound to a specific deployed contract.
func NewTeaPartyCaller(address common.Address, caller bind.ContractCaller) (*TeaPartyCaller, error) {
	contract, err := bindTeaParty(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TeaPartyCaller{contract: contract}, nil
}

// NewTeaPartyTransactor creates a new write-only instance of TeaParty, bound to a specific deployed contract.
func NewTeaPartyTransactor(address common.Address, transactor bind.ContractTransactor) (*TeaPartyTransactor, error) {
	contract, err := bindTeaParty(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TeaPartyTransactor{contract: contract}, nil
}

// NewTeaPartyFilterer creates a new log filterer instance of TeaParty, bound to a specific deployed contract.
func NewTeaPartyFilterer(address common.Address, filterer bind.ContractFilterer) (*TeaPartyFilterer, error) {
	contract, err := bindTeaParty(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TeaPartyFilterer{contract: contract}, nil
}

// bindTeaParty binds a generic wrapper to an already deployed contract.
func bindTeaParty(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TeaPartyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TeaParty *TeaPartyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TeaParty.Contract.TeaPartyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TeaParty *TeaPartyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TeaParty.Contract.TeaPartyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TeaParty *TeaPartyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TeaParty.Contract.TeaPartyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TeaParty *TeaPartyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TeaParty.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TeaParty *TeaPartyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TeaParty.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TeaParty *TeaPartyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TeaParty.Contract.contract.Transact(opts, method, params...)
}

// GetHolders is a free data retrieval call binding the contract method 0x5fe8e7cc.
//
// Solidity: function getHolders() view returns(address[])
func (_TeaParty *TeaPartyCaller) GetHolders(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _TeaParty.contract.Call(opts, &out, "getHolders")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetHolders is a free data retrieval call binding the contract method 0x5fe8e7cc.
//
// Solidity: function getHolders() view returns(address[])
func (_TeaParty *TeaPartySession) GetHolders() ([]common.Address, error) {
	return _TeaParty.Contract.GetHolders(&_TeaParty.CallOpts)
}

// GetHolders is a free data retrieval call binding the contract method 0x5fe8e7cc.
//
// Solidity: function getHolders() view returns(address[])
func (_TeaParty *TeaPartyCallerSession) GetHolders() ([]common.Address, error) {
	return _TeaParty.Contract.GetHolders(&_TeaParty.CallOpts)
}

// GetTeaPartyTransactions is a free data retrieval call binding the contract method 0xa398fcea.
//
// Solidity: function getTeaPartyTransactions() view returns(address[], uint256[])
func (_TeaParty *TeaPartyCaller) GetTeaPartyTransactions(opts *bind.CallOpts) ([]common.Address, []*big.Int, error) {
	var out []interface{}
	err := _TeaParty.contract.Call(opts, &out, "getTeaPartyTransactions")

	if err != nil {
		return *new([]common.Address), *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)

	return out0, out1, err

}

// GetTeaPartyTransactions is a free data retrieval call binding the contract method 0xa398fcea.
//
// Solidity: function getTeaPartyTransactions() view returns(address[], uint256[])
func (_TeaParty *TeaPartySession) GetTeaPartyTransactions() ([]common.Address, []*big.Int, error) {
	return _TeaParty.Contract.GetTeaPartyTransactions(&_TeaParty.CallOpts)
}

// GetTeaPartyTransactions is a free data retrieval call binding the contract method 0xa398fcea.
//
// Solidity: function getTeaPartyTransactions() view returns(address[], uint256[])
func (_TeaParty *TeaPartyCallerSession) GetTeaPartyTransactions() ([]common.Address, []*big.Int, error) {
	return _TeaParty.Contract.GetTeaPartyTransactions(&_TeaParty.CallOpts)
}

// GetTransaction is a free data retrieval call binding the contract method 0x7bb86379.
//
// Solidity: function getTransaction(address participant) view returns(uint256)
func (_TeaParty *TeaPartyCaller) GetTransaction(opts *bind.CallOpts, participant common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TeaParty.contract.Call(opts, &out, "getTransaction", participant)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTransaction is a free data retrieval call binding the contract method 0x7bb86379.
//
// Solidity: function getTransaction(address participant) view returns(uint256)
func (_TeaParty *TeaPartySession) GetTransaction(participant common.Address) (*big.Int, error) {
	return _TeaParty.Contract.GetTransaction(&_TeaParty.CallOpts, participant)
}

// GetTransaction is a free data retrieval call binding the contract method 0x7bb86379.
//
// Solidity: function getTransaction(address participant) view returns(uint256)
func (_TeaParty *TeaPartyCallerSession) GetTransaction(participant common.Address) (*big.Int, error) {
	return _TeaParty.Contract.GetTransaction(&_TeaParty.CallOpts, participant)
}

// OpenTransactions is a free data retrieval call binding the contract method 0xf9f38bdf.
//
// Solidity: function openTransactions() view returns(uint256)
func (_TeaParty *TeaPartyCaller) OpenTransactions(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TeaParty.contract.Call(opts, &out, "openTransactions")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OpenTransactions is a free data retrieval call binding the contract method 0xf9f38bdf.
//
// Solidity: function openTransactions() view returns(uint256)
func (_TeaParty *TeaPartySession) OpenTransactions() (*big.Int, error) {
	return _TeaParty.Contract.OpenTransactions(&_TeaParty.CallOpts)
}

// OpenTransactions is a free data retrieval call binding the contract method 0xf9f38bdf.
//
// Solidity: function openTransactions() view returns(uint256)
func (_TeaParty *TeaPartyCallerSession) OpenTransactions() (*big.Int, error) {
	return _TeaParty.Contract.OpenTransactions(&_TeaParty.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TeaParty *TeaPartyCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TeaParty.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TeaParty *TeaPartySession) Owner() (common.Address, error) {
	return _TeaParty.Contract.Owner(&_TeaParty.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TeaParty *TeaPartyCallerSession) Owner() (common.Address, error) {
	return _TeaParty.Contract.Owner(&_TeaParty.CallOpts)
}

// AddHolder is a paid mutator transaction binding the contract method 0x355624dc.
//
// Solidity: function addHolder() returns()
func (_TeaParty *TeaPartyTransactor) AddHolder(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TeaParty.contract.Transact(opts, "addHolder")
}

// AddHolder is a paid mutator transaction binding the contract method 0x355624dc.
//
// Solidity: function addHolder() returns()
func (_TeaParty *TeaPartySession) AddHolder() (*types.Transaction, error) {
	return _TeaParty.Contract.AddHolder(&_TeaParty.TransactOpts)
}

// AddHolder is a paid mutator transaction binding the contract method 0x355624dc.
//
// Solidity: function addHolder() returns()
func (_TeaParty *TeaPartyTransactorSession) AddHolder() (*types.Transaction, error) {
	return _TeaParty.Contract.AddHolder(&_TeaParty.TransactOpts)
}

// CreateTransaction is a paid mutator transaction binding the contract method 0x97b54cc2.
//
// Solidity: function createTransaction() payable returns(uint256)
func (_TeaParty *TeaPartyTransactor) CreateTransaction(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TeaParty.contract.Transact(opts, "createTransaction")
}

// CreateTransaction is a paid mutator transaction binding the contract method 0x97b54cc2.
//
// Solidity: function createTransaction() payable returns(uint256)
func (_TeaParty *TeaPartySession) CreateTransaction() (*types.Transaction, error) {
	return _TeaParty.Contract.CreateTransaction(&_TeaParty.TransactOpts)
}

// CreateTransaction is a paid mutator transaction binding the contract method 0x97b54cc2.
//
// Solidity: function createTransaction() payable returns(uint256)
func (_TeaParty *TeaPartyTransactorSession) CreateTransaction() (*types.Transaction, error) {
	return _TeaParty.Contract.CreateTransaction(&_TeaParty.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_TeaParty *TeaPartyTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TeaParty.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_TeaParty *TeaPartySession) Deposit() (*types.Transaction, error) {
	return _TeaParty.Contract.Deposit(&_TeaParty.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() payable returns()
func (_TeaParty *TeaPartyTransactorSession) Deposit() (*types.Transaction, error) {
	return _TeaParty.Contract.Deposit(&_TeaParty.TransactOpts)
}

// RemoveTransaction is a paid mutator transaction binding the contract method 0x0ddb3d23.
//
// Solidity: function removeTransaction(address participant) returns()
func (_TeaParty *TeaPartyTransactor) RemoveTransaction(opts *bind.TransactOpts, participant common.Address) (*types.Transaction, error) {
	return _TeaParty.contract.Transact(opts, "removeTransaction", participant)
}

// RemoveTransaction is a paid mutator transaction binding the contract method 0x0ddb3d23.
//
// Solidity: function removeTransaction(address participant) returns()
func (_TeaParty *TeaPartySession) RemoveTransaction(participant common.Address) (*types.Transaction, error) {
	return _TeaParty.Contract.RemoveTransaction(&_TeaParty.TransactOpts, participant)
}

// RemoveTransaction is a paid mutator transaction binding the contract method 0x0ddb3d23.
//
// Solidity: function removeTransaction(address participant) returns()
func (_TeaParty *TeaPartyTransactorSession) RemoveTransaction(participant common.Address) (*types.Transaction, error) {
	return _TeaParty.Contract.RemoveTransaction(&_TeaParty.TransactOpts, participant)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TeaParty *TeaPartyTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TeaParty.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TeaParty *TeaPartySession) RenounceOwnership() (*types.Transaction, error) {
	return _TeaParty.Contract.RenounceOwnership(&_TeaParty.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TeaParty *TeaPartyTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TeaParty.Contract.RenounceOwnership(&_TeaParty.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TeaParty *TeaPartyTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TeaParty.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TeaParty *TeaPartySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TeaParty.Contract.TransferOwnership(&_TeaParty.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TeaParty *TeaPartyTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TeaParty.Contract.TransferOwnership(&_TeaParty.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_TeaParty *TeaPartyTransactor) Withdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _TeaParty.contract.Transact(opts, "withdraw", amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_TeaParty *TeaPartySession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _TeaParty.Contract.Withdraw(&_TeaParty.TransactOpts, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 amount) returns()
func (_TeaParty *TeaPartyTransactorSession) Withdraw(amount *big.Int) (*types.Transaction, error) {
	return _TeaParty.Contract.Withdraw(&_TeaParty.TransactOpts, amount)
}

// TeaPartyOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TeaParty contract.
type TeaPartyOwnershipTransferredIterator struct {
	Event *TeaPartyOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TeaPartyOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TeaPartyOwnershipTransferred)
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
		it.Event = new(TeaPartyOwnershipTransferred)
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
func (it *TeaPartyOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TeaPartyOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TeaPartyOwnershipTransferred represents a OwnershipTransferred event raised by the TeaParty contract.
type TeaPartyOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TeaParty *TeaPartyFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TeaPartyOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TeaParty.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TeaPartyOwnershipTransferredIterator{contract: _TeaParty.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TeaParty *TeaPartyFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TeaPartyOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TeaParty.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TeaPartyOwnershipTransferred)
				if err := _TeaParty.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TeaParty *TeaPartyFilterer) ParseOwnershipTransferred(log types.Log) (*TeaPartyOwnershipTransferred, error) {
	event := new(TeaPartyOwnershipTransferred)
	if err := _TeaParty.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TeaPartyTransactionCreatedIterator is returned from FilterTransactionCreated and is used to iterate over the raw logs and unpacked data for TransactionCreated events raised by the TeaParty contract.
type TeaPartyTransactionCreatedIterator struct {
	Event *TeaPartyTransactionCreated // Event containing the contract specifics and raw log

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
func (it *TeaPartyTransactionCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TeaPartyTransactionCreated)
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
		it.Event = new(TeaPartyTransactionCreated)
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
func (it *TeaPartyTransactionCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TeaPartyTransactionCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TeaPartyTransactionCreated represents a TransactionCreated event raised by the TeaParty contract.
type TeaPartyTransactionCreated struct {
	From          common.Address
	TransactionId *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTransactionCreated is a free log retrieval operation binding the contract event 0xbee0a36f8f2ce3c4d9ded6f7c349c1f9d380c787784d84b2bead8455b7d7f92f.
//
// Solidity: event TransactionCreated(address indexed from, uint256 transactionId)
func (_TeaParty *TeaPartyFilterer) FilterTransactionCreated(opts *bind.FilterOpts, from []common.Address) (*TeaPartyTransactionCreatedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _TeaParty.contract.FilterLogs(opts, "TransactionCreated", fromRule)
	if err != nil {
		return nil, err
	}
	return &TeaPartyTransactionCreatedIterator{contract: _TeaParty.contract, event: "TransactionCreated", logs: logs, sub: sub}, nil
}

// WatchTransactionCreated is a free log subscription operation binding the contract event 0xbee0a36f8f2ce3c4d9ded6f7c349c1f9d380c787784d84b2bead8455b7d7f92f.
//
// Solidity: event TransactionCreated(address indexed from, uint256 transactionId)
func (_TeaParty *TeaPartyFilterer) WatchTransactionCreated(opts *bind.WatchOpts, sink chan<- *TeaPartyTransactionCreated, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _TeaParty.contract.WatchLogs(opts, "TransactionCreated", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TeaPartyTransactionCreated)
				if err := _TeaParty.contract.UnpackLog(event, "TransactionCreated", log); err != nil {
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

// ParseTransactionCreated is a log parse operation binding the contract event 0xbee0a36f8f2ce3c4d9ded6f7c349c1f9d380c787784d84b2bead8455b7d7f92f.
//
// Solidity: event TransactionCreated(address indexed from, uint256 transactionId)
func (_TeaParty *TeaPartyFilterer) ParseTransactionCreated(log types.Log) (*TeaPartyTransactionCreated, error) {
	event := new(TeaPartyTransactionCreated)
	if err := _TeaParty.contract.UnpackLog(event, "TransactionCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
