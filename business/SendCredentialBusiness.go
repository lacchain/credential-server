package business

import (
	"encoding/json"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"

	bl "github.com/lacchain/credential-server/blockchain"
	"github.com/lacchain/credential-server/models"

	"github.com/ethereum/go-ethereum/common"
)

//SendCredential saved into blockchain
func SendCredential(credentials []*models.Credential, nodeURL string, publicAddress string, verificationContract string) (*models.VerifyResponse, error) {
	client := new(bl.Client)
	err := client.Connect(nodeURL)

	if err != nil {
		return nil, err
	}

	defer client.Close()

	response := new(models.VerifyResponse)
	response.Valid = true
	errorResponse := new(models.Error)
	errorResponse.Code = "200"
	errorResponse.Message = "OK"
	response.Error = errorResponse

	contractAddress := common.HexToAddress(verificationContract)
	address := common.HexToAddress(publicAddress)

//iterate credential and verify
for _, credential := range credentials {

	//verify credential and send email
	log.Printf("Verifying credential ID: %s", *credential.CredentialData.ID)
	rawCredential, err := json.Marshal(credential.CredentialData)
	if err != nil {
		log.Println("Credential isn't a json format")
		return nil, err
	}
	result, err := client.VerifyCredential(contractAddress, rawCredential, address)
	if err != nil {
		return nil, err
	}
	if !result {
		response.Valid = false
		errorResponse.Code = "400"
		errorResponse.Message = "Credential is invalid"
	}

	credentialHash := sha256.Sum256(rawCredential)

	//credentialMetadata.BlockNumber = blockNumber.String()
	//credentialMetadata.Timestamp = timestamp
	//credentialMetadata.Transaction = tx.Hex()

	//credential.Metadata = credentialMetadata

	//obtener el content de credentialsubject

	qrFile, err := generateQR("http://credentialserver.iadb.org/", credentialHash, getNameSubject(credential.CredentialData.CredentialSubject))
	if err != nil {
		fmt.Printf("Failed generate QR: %s", err)
	}

	//wait by new email in parameter

	blockNumber := new(big.Int)
    blockNumber, _ = blockNumber.SetString(credential.Metadata.BlockNumber, 10)
	err = sendCredentialByEmail(getNameSubject(credential.CredentialData.CredentialSubject), credential.Metadata.Email, []byte(string(rawCredential)), credential.Metadata.Transaction,blockNumber,credential.Metadata.Timestamp, getExpirationDate(credential.CredentialData.CredentialSubject), qrFile)
	if err != nil {
		fmt.Printf("Failed to send email: %s", err)
	}
}

	return response, nil
}
