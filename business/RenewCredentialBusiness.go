/*
	Credential Server
	version 0.9
	author: Adrian Pareja Abarca
	email: adriancc5.5@gmail.com
*/

package business

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	bl "github.com/lacchain/credential-server/blockchain"
	"github.com/lacchain/credential-server/models"
)

//RenewCredential revoking old credential and saving the new credential hash into blockchain
func RenewCredential(credentialID string, subject *models.CredentialSubject, nodeURL string, issuer string, privateKey string, verificationContract string, repositoryContract string) (*models.Credential, error) {
	//revoke old credential
	err := RevokeCredential(credentialID, nodeURL, privateKey, verificationContract, repositoryContract)
	if err != nil {
		return nil, err
	}

	client := new(bl.Client)
	err = client.Connect(nodeURL)
	if err != nil {
		return nil, err
	}
	options, err := client.ConfigTransaction(privateKey)
	if err != nil {
		return nil, err
	}
	address := common.HexToAddress(verificationContract)

	//generate new credential
	credential := new(models.Credential)
	credentialData := new(models.CredentialData)
	credentialData.ID = &credentialID
	types := make([]string, 0, 4)
	types = append(types, "VerifiableCredential")
	types = append(types, subject.Type)
	credentialData.Type = types
	credentialData.CredentialSubject = subject.Content
	credentialData.Evidence = subject.Evidence
	credentialData.Issuer = issuer
	credentialData.Proof = getProof("SmartContract", verificationContract)

	rawCredential, err := json.Marshal(credentialData)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(rawCredential))

	fmt.Println("fecha:", subject.ExpirationDate.String())

	date := time.Time(subject.ExpirationDate)

	fmt.Println("tiempo seconds:", date.Unix())

	credentialHash := sha256.Sum256(rawCredential)

	client.SignCredential(address, options, credentialHash, big.NewInt(date.Unix()))

	options, err = client.ConfigTransaction(privateKey)
	if err != nil {
		return nil, err
	}
	idHash := sha256.Sum256([]byte(*credentialData.ID))

	var hashito = make([]byte, 32, 64)

	for i, j := range idHash {
		hashito[i] = j
	}

	fmt.Println("hash credential id:", hex.EncodeToString(hashito))

	client.SetCredential(common.HexToAddress(repositoryContract), options, idHash, credentialHash)

	client.Close()

	return credential, nil
}
