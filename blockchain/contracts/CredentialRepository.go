// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package store

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// RepositoryABI is the input ABI used to generate the binding from.
const RepositoryABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"credentials\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"id\",\"type\":\"bytes32\"},{\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"register\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"by\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"Saved\",\"type\":\"event\"}]"

// RepositoryBin is the compiled bytecode used for deploying new contracts.
const RepositoryBin = `0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506102ab806100606000396000f30060806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063299d0348146100515780632f9267321461009e575b600080fd5b34801561005d57600080fd5b5061008060048036038101908080356000191690602001909291905050506100dd565b60405180826000191660001916815260200191505060405180910390f35b3480156100aa57600080fd5b506100db600480360381019080803560001916906020019092919080356000191690602001909291905050506100f5565b005b60016020528060005260406000206000915090505481565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156101df576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f4f6e6c79206f776e65722063616e2063616c6c20746869732066756e6374696f81526020017f6e2e00000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b80600160008460001916600019168152602001908152602001600020816000191690555081600019167f44a27dbdd2a8591b2851f8c8f004aa50d1081d9b41af8c7522d9d7ea57e9f68f3383604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200182600019166000191681526020019250505060405180910390a250505600a165627a7a72305820f869eca7ef3899967507bc3842444f0355fac3e29ded0dfb1e809990d5135cbe0029`

// DeployRepository deploys a new Ethereum contract, binding an instance of Repository to it.
func DeployRepository(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Repository, error) {
	parsed, err := abi.JSON(strings.NewReader(RepositoryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RepositoryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Repository{RepositoryCaller: RepositoryCaller{contract: contract}, RepositoryTransactor: RepositoryTransactor{contract: contract}, RepositoryFilterer: RepositoryFilterer{contract: contract}}, nil
}

// Repository is an auto generated Go binding around an Ethereum contract.
type Repository struct {
	RepositoryCaller     // Read-only binding to the contract
	RepositoryTransactor // Write-only binding to the contract
	RepositoryFilterer   // Log filterer for contract events
}

// RepositoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type RepositoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RepositoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RepositoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RepositoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RepositoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RepositorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RepositorySession struct {
	Contract     *Repository       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RepositoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RepositoryCallerSession struct {
	Contract *RepositoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// RepositoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RepositoryTransactorSession struct {
	Contract     *RepositoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// RepositoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type RepositoryRaw struct {
	Contract *Repository // Generic contract binding to access the raw methods on
}

// RepositoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RepositoryCallerRaw struct {
	Contract *RepositoryCaller // Generic read-only contract binding to access the raw methods on
}

// RepositoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RepositoryTransactorRaw struct {
	Contract *RepositoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRepository creates a new instance of Repository, bound to a specific deployed contract.
func NewRepository(address common.Address, backend bind.ContractBackend) (*Repository, error) {
	contract, err := bindRepository(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Repository{RepositoryCaller: RepositoryCaller{contract: contract}, RepositoryTransactor: RepositoryTransactor{contract: contract}, RepositoryFilterer: RepositoryFilterer{contract: contract}}, nil
}

// NewRepositoryCaller creates a new read-only instance of Repository, bound to a specific deployed contract.
func NewRepositoryCaller(address common.Address, caller bind.ContractCaller) (*RepositoryCaller, error) {
	contract, err := bindRepository(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RepositoryCaller{contract: contract}, nil
}

// NewRepositoryTransactor creates a new write-only instance of Repository, bound to a specific deployed contract.
func NewRepositoryTransactor(address common.Address, transactor bind.ContractTransactor) (*RepositoryTransactor, error) {
	contract, err := bindRepository(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RepositoryTransactor{contract: contract}, nil
}

// NewRepositoryFilterer creates a new log filterer instance of Repository, bound to a specific deployed contract.
func NewRepositoryFilterer(address common.Address, filterer bind.ContractFilterer) (*RepositoryFilterer, error) {
	contract, err := bindRepository(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RepositoryFilterer{contract: contract}, nil
}

// bindRepository binds a generic wrapper to an already deployed contract.
func bindRepository(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RepositoryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Repository *RepositoryRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Repository.Contract.RepositoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Repository *RepositoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Repository.Contract.RepositoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Repository *RepositoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Repository.Contract.RepositoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Repository *RepositoryCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Repository.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Repository *RepositoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Repository.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Repository *RepositoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Repository.Contract.contract.Transact(opts, method, params...)
}

// Credentials is a free data retrieval call binding the contract method 0x299d0348.
//
// Solidity: function credentials( bytes32) constant returns(bytes32)
func (_Repository *RepositoryCaller) Credentials(opts *bind.CallOpts, arg0 [32]byte) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Repository.contract.Call(opts, out, "credentials", arg0)
	return *ret0, err
}

// Credentials is a free data retrieval call binding the contract method 0x299d0348.
//
// Solidity: function credentials( bytes32) constant returns(bytes32)
func (_Repository *RepositorySession) Credentials(arg0 [32]byte) ([32]byte, error) {
	return _Repository.Contract.Credentials(&_Repository.CallOpts, arg0)
}

// Credentials is a free data retrieval call binding the contract method 0x299d0348.
//
// Solidity: function credentials( bytes32) constant returns(bytes32)
func (_Repository *RepositoryCallerSession) Credentials(arg0 [32]byte) ([32]byte, error) {
	return _Repository.Contract.Credentials(&_Repository.CallOpts, arg0)
}

// Register is a paid mutator transaction binding the contract method 0x2f926732.
//
// Solidity: function register(id bytes32, hash bytes32) returns()
func (_Repository *RepositoryTransactor) Register(opts *bind.TransactOpts, id [32]byte, hash [32]byte) (*types.Transaction, error) {
	return _Repository.contract.Transact(opts, "register", id, hash)
}

// Register is a paid mutator transaction binding the contract method 0x2f926732.
//
// Solidity: function register(id bytes32, hash bytes32) returns()
func (_Repository *RepositorySession) Register(id [32]byte, hash [32]byte) (*types.Transaction, error) {
	return _Repository.Contract.Register(&_Repository.TransactOpts, id, hash)
}

// Register is a paid mutator transaction binding the contract method 0x2f926732.
//
// Solidity: function register(id bytes32, hash bytes32) returns()
func (_Repository *RepositoryTransactorSession) Register(id [32]byte, hash [32]byte) (*types.Transaction, error) {
	return _Repository.Contract.Register(&_Repository.TransactOpts, id, hash)
}

// RepositorySavedIterator is returned from FilterSaved and is used to iterate over the raw logs and unpacked data for Saved events raised by the Repository contract.
type RepositorySavedIterator struct {
	Event *RepositorySaved // Event containing the contract specifics and raw log

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
func (it *RepositorySavedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RepositorySaved)
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
		it.Event = new(RepositorySaved)
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
func (it *RepositorySavedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RepositorySavedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RepositorySaved represents a Saved event raised by the Repository contract.
type RepositorySaved struct {
	Id   [32]byte
	By   common.Address
	Hash [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSaved is a free log retrieval operation binding the contract event 0x44a27dbdd2a8591b2851f8c8f004aa50d1081d9b41af8c7522d9d7ea57e9f68f.
//
// Solidity: e Saved(id indexed bytes32, by address, hash bytes32)
func (_Repository *RepositoryFilterer) FilterSaved(opts *bind.FilterOpts, id [][32]byte) (*RepositorySavedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Repository.contract.FilterLogs(opts, "Saved", idRule)
	if err != nil {
		return nil, err
	}
	return &RepositorySavedIterator{contract: _Repository.contract, event: "Saved", logs: logs, sub: sub}, nil
}

// WatchSaved is a free log subscription operation binding the contract event 0x44a27dbdd2a8591b2851f8c8f004aa50d1081d9b41af8c7522d9d7ea57e9f68f.
//
// Solidity: e Saved(id indexed bytes32, by address, hash bytes32)
func (_Repository *RepositoryFilterer) WatchSaved(opts *bind.WatchOpts, sink chan<- *RepositorySaved, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _Repository.contract.WatchLogs(opts, "Saved", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RepositorySaved)
				if err := _Repository.contract.UnpackLog(event, "Saved", log); err != nil {
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
