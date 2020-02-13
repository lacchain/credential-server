package business

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	bl "github.com/lacchain/credential-server/blockchain"
	"github.com/lacchain/credential-server/util"
)

//RevokeCredential into blockchain
func RevokeCredential(credentialID string, nodeURL string, privateKey string, revocationContract string, repositoryContract string) error {
	client := new(bl.Client)
	err := client.Connect(nodeURL)
	defer client.Close()
	if err != nil {
		return err
	}
	options, err := client.ConfigTransaction(privateKey)
	if err != nil {
		return err
	}
	idHash := sha256.Sum256([]byte(credentialID))

	credential, err := client.GetCredential(common.HexToAddress(repositoryContract), idHash)
	if err != nil {
		return err
	}

	var emptyCredential [32]byte

	if util.EqualBytes(credential, emptyCredential) {
		fmt.Println("Credential not found")
		return errors.New("Credential not found")
	}

	client.RevokeCredential(common.HexToAddress(revocationContract), options, credential)

	return nil
}
