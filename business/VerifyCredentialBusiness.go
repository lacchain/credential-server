package business

import (
	"encoding/json"
	"log"

	bl "github.com/lacchain/credential-server/blockchain"
	"github.com/lacchain/credential-server/models"

	"github.com/ethereum/go-ethereum/common"
)

//VerifyCredential saved into blockchain
func VerifyCredential(credentials []*models.CredentialData, nodeURL string, publicAddress string, verificationContract string) (*models.VerifyResponse, error) {
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
	//failsID := make([]string, 0, 10)
	//var fail error

	//iterate credential and verify
	for _, credential := range credentials {
		//credentialData := credential.CredentialData
		log.Printf("Verifying credential ID: %s", *credential.ID)
		rawCredential, err := json.Marshal(credential)
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
	}

	/*if fail != nil {
		response.Valid = false
		errorResponse.Code = "400"
		errorResponse.Message = "Credential isn't valid"
		return response, fail
	}*/

	return response, nil
}
