package blockchain

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/lacchain/credential-server/util"

	store "github.com/lacchain/credential-server/blockchain/contracts"
)

const verificationregistryabi = "[{ 'constant': false, 'inputs': [{ 'name': 'hash', 'type': 'bytes32' }, { 'name': 'validDays', 'type': 'uint256' }], 'name': 'verify', 'outputs': [], 'payable': false, 'stateMutability': 'nonpayable', 'type': 'function' }, { 'constant': true, 'inputs': [{ 'name': '', 'type': 'bytes32' }, { 'name': '', 'type': 'address' }], 'name': 'verifications', 'outputs': [{ 'name': 'iat', 'type': 'uint256' }, { 'name': 'exp', 'type': 'uint256' }], 'payable': false, 'stateMutability': 'view', 'type': 'function' }, { 'constant': false, 'inputs': [{ 'name': 'hash', 'type': 'bytes32' }], 'name': 'revoke', 'outputs': [], 'payable': false, 'stateMutability': 'nonpayable', 'type': 'function' }, { 'anonymous': false, 'inputs': [{ 'indexed': true, 'name': 'hash', 'type': 'bytes32' }, { 'indexed': false, 'name': 'by', 'type': 'address' }, { 'indexed': false, 'name': 'date', 'type': 'uint256' }, { 'indexed': false, 'name': 'expDate', 'type': 'uint256' }], 'name': 'Verified', 'type': 'event' }, { 'anonymous': false, 'inputs': [{ 'indexed': true, 'name': 'hash', 'type': 'bytes32' }, { 'indexed': false, 'name': 'by', 'type': 'address' }, { 'indexed': false, 'name': 'date', 'type': 'uint256' }], 'name': 'Revoked', 'type': 'event' }]"

//Client to manage Connection to Ethereum
type Client struct {
	client *ethclient.Client
}

//Connect to Ethereum
func (ec *Client) Connect(nodeURL string) error {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Connected to Ethereum Node")
	ec.client = client
	return nil
}

//Close ethereum connection
func (ec *Client) Close() {
	ec.client.Close()
}

//ConfigTransaction from ethereum address contract
func (ec *Client) ConfigTransaction(pk string) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
		return nil, errors.New("publicKey isn't type *ecdsa")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := ec.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	gasPrice, err := ec.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Println(err)
		return nil, err
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(4000700) // in units
	auth.GasPrice = gasPrice

	gas := util.CalcGasCost(auth.GasLimit, gasPrice)

	fmt.Println("nonce:", auth.Nonce)
	fmt.Println("gasLimit:", auth.GasLimit)
	fmt.Println("gasPrice:", auth.GasPrice)
	fmt.Println("gas:", gas)

	return auth, nil
}

//SignCredential into blockchain
func (ec *Client) SignCredential(contractAddress common.Address, options *bind.TransactOpts, credentialHash [32]byte, duration *big.Int) (error, *common.Hash) {
	contract, err := store.NewStore(contractAddress, ec.client)
	if err != nil {
		log.Printf("Can't instance contract: %s", err)
		return err, nil
	}

	log.Println("Verification Contract instanced")

	var hashito = make([]byte, 32, 64)

	for i, j := range credentialHash {
		hashito[i] = j
	}

	log.Println("credential key:", hex.EncodeToString(hashito))

	log.Println("duration of contract:", duration)

	tx, err := contract.Verify(options, credentialHash, duration)
	if err != nil {
		log.Println("Error calling contract:", err)
		return err, nil
	}

	log.Printf("Tx sent: %s", tx.Hash().Hex())

	fff := tx.Hash()

	return nil, &fff
}

//RevokeCredential into blockchain
func (ec *Client) RevokeCredential(contractAddress common.Address, options *bind.TransactOpts, credentialHash [32]byte) error {
	contract, err := store.NewStore(contractAddress, ec.client)
	if err != nil {
		log.Printf("Can't instance contract: %s", err)
		return err
	}

	log.Println("Contract instanced")

	var hashito = make([]byte, 32, 64)

	for i, j := range credentialHash {
		hashito[i] = j
	}

	log.Println("credential key:", hex.EncodeToString(hashito))

	tx, err := contract.Revoke(options, credentialHash)
	if err != nil {
		log.Println("Error calling contract:", err)
		return err
	}

	log.Printf("Tx sent: %s", tx.Hash().Hex())
	return nil
}

//VerifyCredential saved into blockchain
func (ec *Client) VerifyCredential(contractAddress common.Address, credential []byte, address common.Address) (bool, error) {
	contract, err := store.NewStore(contractAddress, ec.client)
	if err != nil {
		log.Printf("Can't instance contract: %s", err)
		return false, err
	}

	log.Printf("Contract %s instanced", contractAddress)

	hash := sha256.Sum256(credential)

	var hashito = make([]byte, 32, 64)

	for i, j := range hash {
		hashito[i] = j
	}

	log.Println("credential hash:", hex.EncodeToString(hashito))

	result, err := contract.Verifications(&bind.CallOpts{}, hash, address)
	if err != nil {
		log.Println("Error calling contract:", err)
		return false, err
	}

	log.Println("result:", result)

	cmp := result.Exp.Cmp(big.NewInt(0))

	log.Println("cmp:",cmp)

	if cmp > 0 {
		log.Println("time actual:",time.Now().Unix())
		cmp2 := result.Exp.Cmp(big.NewInt(time.Now().Unix()))
		if cmp2 > 0{
			log.Println("Credential hasn't expired")
			log.Println("Credential is valid")
		}else{
			log.Println("Credential has expired")
			log.Println("Credential is invalid")
			return false, nil
		}
		return true, nil
	}

	log.Println("Credential is invalid")
	return false, nil
}

//SetCredential to repository blockchain with id
func (ec *Client) SetCredential(contractAddress common.Address, options *bind.TransactOpts, credentialID [32]byte, credentialHash [32]byte) error {
	contract, err := store.NewRepository(contractAddress, ec.client)
	if err != nil {
		log.Printf("Can't instance contract repository: %s", err)
		return err
	}

	log.Println("Contract Respository instanced")

	var hashito = make([]byte, 32, 64)
	var hashito2 = make([]byte, 32, 64)

	for i, j := range credentialHash {
		hashito[i] = j
	}

	log.Println("credential key:", hex.EncodeToString(hashito))

	for i, j := range credentialID {
		hashito2[i] = j
	}

	log.Println("credential ID:", hex.EncodeToString(hashito2))

	tx, err := contract.Register(options, credentialID, credentialHash)
	if err != nil {
		log.Println("Error calling contract:", err)
		return err
	}

	log.Printf("Tx sent: %s", tx.Hash().Hex())
	return nil
}

//GetCredential from blockchain by Id
func (ec *Client) GetCredential(contractAddress common.Address, id [32]byte) ([32]byte, error) {
	contract, err := store.NewRepository(contractAddress, ec.client)
	if err != nil {
		log.Printf("Can't instance contract: %s", err)
	}

	log.Println("Repository Contract instanced ...")

	result, err := contract.Credentials(&bind.CallOpts{}, id)
	if err != nil {
		log.Println("Error calling contract:", err)
	}

	log.Println("result:", result)

	return result, nil
}

//DeployRepositoryContract into Ethereum blockchain
func (ec *Client) DeployRepositoryContract(pk string) (string, error) {
	options, err := ec.ConfigTransaction(pk)
	if err != nil {
		return "", err
	}
	contractAddress, tx, _, err := store.DeployRepository(options, ec.client)
	if err != nil {
		log.Println(err)
		return "", err
	}

	log.Println("repository contract address:", contractAddress.Hex())
	log.Println("tx:", tx.Hash().Hex())
	return contractAddress.Hex(), nil
}

//GetTransactionReceipt ...
func (ec *Client) GetTransactionReceipt(transactionHash common.Hash)(*big.Int, string, error){
	receipt, err := ec.client.TransactionReceipt(context.Background(), transactionHash)
        if err != nil {
            log.Println(err)
		}
		
	log.Println("Status:",receipt.Status)
	log.Println("BlockNumber:",receipt.BlockNumber)

	block, err := ec.client.BlockByNumber(context.Background(), receipt.BlockNumber)
    if err != nil {
        log.Println(err)
	}
	
	log.Println("block time:",block.Time())

	ts := time.Unix(int64(block.Time()),0).UTC()

	return receipt.BlockNumber, ts.String(), nil
}
