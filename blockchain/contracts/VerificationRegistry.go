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

// StoreABI is the input ABI used to generate the binding from.
const StoreABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"validDays\",\"type\":\"uint256\"}],\"name\":\"verify\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"},{\"name\":\"\",\"type\":\"address\"}],\"name\":\"verifications\",\"outputs\":[{\"name\":\"iat\",\"type\":\"uint256\"},{\"name\":\"exp\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"revoke\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"by\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"date\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"expDate\",\"type\":\"uint256\"}],\"name\":\"Verified\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"hash\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"by\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"date\",\"type\":\"uint256\"}],\"name\":\"Revoked\",\"type\":\"event\"}]"

// Store is an auto generated Go binding around an Ethereum contract.
type Store struct {
	StoreCaller     // Read-only binding to the contract
	StoreTransactor // Write-only binding to the contract
	StoreFilterer   // Log filterer for contract events
}

// StoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type StoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StoreSession struct {
	Contract     *Store            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StoreCallerSession struct {
	Contract *StoreCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// StoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StoreTransactorSession struct {
	Contract     *StoreTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type StoreRaw struct {
	Contract *Store // Generic contract binding to access the raw methods on
}

// StoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StoreCallerRaw struct {
	Contract *StoreCaller // Generic read-only contract binding to access the raw methods on
}

// StoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StoreTransactorRaw struct {
	Contract *StoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStore creates a new instance of Store, bound to a specific deployed contract.
func NewStore(address common.Address, backend bind.ContractBackend) (*Store, error) {
	contract, err := bindStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Store{StoreCaller: StoreCaller{contract: contract}, StoreTransactor: StoreTransactor{contract: contract}, StoreFilterer: StoreFilterer{contract: contract}}, nil
}

// NewStoreCaller creates a new read-only instance of Store, bound to a specific deployed contract.
func NewStoreCaller(address common.Address, caller bind.ContractCaller) (*StoreCaller, error) {
	contract, err := bindStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StoreCaller{contract: contract}, nil
}

// NewStoreTransactor creates a new write-only instance of Store, bound to a specific deployed contract.
func NewStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*StoreTransactor, error) {
	contract, err := bindStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StoreTransactor{contract: contract}, nil
}

// NewStoreFilterer creates a new log filterer instance of Store, bound to a specific deployed contract.
func NewStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*StoreFilterer, error) {
	contract, err := bindStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StoreFilterer{contract: contract}, nil
}

// bindStore binds a generic wrapper to an already deployed contract.
func bindStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Store *StoreRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Store.Contract.StoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Store *StoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.Contract.StoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Store *StoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Store.Contract.StoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Store *StoreCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Store.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Store *StoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Store.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Store *StoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Store.Contract.contract.Transact(opts, method, params...)
}

// Verifications is a free data retrieval call binding the contract method 0x95f5114c.
//
// Solidity: function verifications( bytes32,  address) constant returns(iat uint256, exp uint256)
func (_Store *StoreCaller) Verifications(opts *bind.CallOpts, arg0 [32]byte, arg1 common.Address) (struct {
	Iat *big.Int
	Exp *big.Int
}, error) {
	ret := new(struct {
		Iat *big.Int
		Exp *big.Int
	})
	out := ret
	err := _Store.contract.Call(opts, out, "verifications", arg0, arg1)
	return *ret, err
}

// Verifications is a free data retrieval call binding the contract method 0x95f5114c.
//
// Solidity: function verifications( bytes32,  address) constant returns(iat uint256, exp uint256)
func (_Store *StoreSession) Verifications(arg0 [32]byte, arg1 common.Address) (struct {
	Iat *big.Int
	Exp *big.Int
}, error) {
	return _Store.Contract.Verifications(&_Store.CallOpts, arg0, arg1)
}

// Verifications is a free data retrieval call binding the contract method 0x95f5114c.
//
// Solidity: function verifications( bytes32,  address) constant returns(iat uint256, exp uint256)
func (_Store *StoreCallerSession) Verifications(arg0 [32]byte, arg1 common.Address) (struct {
	Iat *big.Int
	Exp *big.Int
}, error) {
	return _Store.Contract.Verifications(&_Store.CallOpts, arg0, arg1)
}

// Revoke is a paid mutator transaction binding the contract method 0xb75c7dc6.
//
// Solidity: function revoke(hash bytes32) returns()
func (_Store *StoreTransactor) Revoke(opts *bind.TransactOpts, hash [32]byte) (*types.Transaction, error) {
	return _Store.contract.Transact(opts, "revoke", hash)
}

// Revoke is a paid mutator transaction binding the contract method 0xb75c7dc6.
//
// Solidity: function revoke(hash bytes32) returns()
func (_Store *StoreSession) Revoke(hash [32]byte) (*types.Transaction, error) {
	return _Store.Contract.Revoke(&_Store.TransactOpts, hash)
}

// Revoke is a paid mutator transaction binding the contract method 0xb75c7dc6.
//
// Solidity: function revoke(hash bytes32) returns()
func (_Store *StoreTransactorSession) Revoke(hash [32]byte) (*types.Transaction, error) {
	return _Store.Contract.Revoke(&_Store.TransactOpts, hash)
}

// Verify is a paid mutator transaction binding the contract method 0x382262fc.
//
// Solidity: function verify(hash bytes32, validDays uint256) returns()
func (_Store *StoreTransactor) Verify(opts *bind.TransactOpts, hash [32]byte, validDays *big.Int) (*types.Transaction, error) {
	return _Store.contract.Transact(opts, "verify", hash, validDays)
}

// Verify is a paid mutator transaction binding the contract method 0x382262fc.
//
// Solidity: function verify(hash bytes32, validDays uint256) returns()
func (_Store *StoreSession) Verify(hash [32]byte, validDays *big.Int) (*types.Transaction, error) {
	return _Store.Contract.Verify(&_Store.TransactOpts, hash, validDays)
}

// Verify is a paid mutator transaction binding the contract method 0x382262fc.
//
// Solidity: function verify(hash bytes32, validDays uint256) returns()
func (_Store *StoreTransactorSession) Verify(hash [32]byte, validDays *big.Int) (*types.Transaction, error) {
	return _Store.Contract.Verify(&_Store.TransactOpts, hash, validDays)
}

// StoreRevokedIterator is returned from FilterRevoked and is used to iterate over the raw logs and unpacked data for Revoked events raised by the Store contract.
type StoreRevokedIterator struct {
	Event *StoreRevoked // Event containing the contract specifics and raw log

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
func (it *StoreRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StoreRevoked)
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
		it.Event = new(StoreRevoked)
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
func (it *StoreRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StoreRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StoreRevoked represents a Revoked event raised by the Store contract.
type StoreRevoked struct {
	Hash [32]byte
	By   common.Address
	Date *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterRevoked is a free log retrieval operation binding the contract event 0xbad7982cc5cd1310cfc1c83b1e7e9bc771332036ea180c362910e4068fc5ecf6.
//
// Solidity: e Revoked(hash indexed bytes32, by address, date uint256)
func (_Store *StoreFilterer) FilterRevoked(opts *bind.FilterOpts, hash [][32]byte) (*StoreRevokedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _Store.contract.FilterLogs(opts, "Revoked", hashRule)
	if err != nil {
		return nil, err
	}
	return &StoreRevokedIterator{contract: _Store.contract, event: "Revoked", logs: logs, sub: sub}, nil
}

// WatchRevoked is a free log subscription operation binding the contract event 0xbad7982cc5cd1310cfc1c83b1e7e9bc771332036ea180c362910e4068fc5ecf6.
//
// Solidity: e Revoked(hash indexed bytes32, by address, date uint256)
func (_Store *StoreFilterer) WatchRevoked(opts *bind.WatchOpts, sink chan<- *StoreRevoked, hash [][32]byte) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _Store.contract.WatchLogs(opts, "Revoked", hashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StoreRevoked)
				if err := _Store.contract.UnpackLog(event, "Revoked", log); err != nil {
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

// StoreVerifiedIterator is returned from FilterVerified and is used to iterate over the raw logs and unpacked data for Verified events raised by the Store contract.
type StoreVerifiedIterator struct {
	Event *StoreVerified // Event containing the contract specifics and raw log

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
func (it *StoreVerifiedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StoreVerified)
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
		it.Event = new(StoreVerified)
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
func (it *StoreVerifiedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StoreVerifiedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StoreVerified represents a Verified event raised by the Store contract.
type StoreVerified struct {
	Hash    [32]byte
	By      common.Address
	Date    *big.Int
	ExpDate *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterVerified is a free log retrieval operation binding the contract event 0xef081e4d12d9d2377ab4c9af5421f0ce40d78d2031cb2cccd1c23eb0a3c6ecce.
//
// Solidity: e Verified(hash indexed bytes32, by address, date uint256, expDate uint256)
func (_Store *StoreFilterer) FilterVerified(opts *bind.FilterOpts, hash [][32]byte) (*StoreVerifiedIterator, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _Store.contract.FilterLogs(opts, "Verified", hashRule)
	if err != nil {
		return nil, err
	}
	return &StoreVerifiedIterator{contract: _Store.contract, event: "Verified", logs: logs, sub: sub}, nil
}

// WatchVerified is a free log subscription operation binding the contract event 0xef081e4d12d9d2377ab4c9af5421f0ce40d78d2031cb2cccd1c23eb0a3c6ecce.
//
// Solidity: e Verified(hash indexed bytes32, by address, date uint256, expDate uint256)
func (_Store *StoreFilterer) WatchVerified(opts *bind.WatchOpts, sink chan<- *StoreVerified, hash [][32]byte) (event.Subscription, error) {

	var hashRule []interface{}
	for _, hashItem := range hash {
		hashRule = append(hashRule, hashItem)
	}

	logs, sub, err := _Store.contract.WatchLogs(opts, "Verified", hashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StoreVerified)
				if err := _Store.contract.UnpackLog(event, "Verified", log); err != nil {
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
