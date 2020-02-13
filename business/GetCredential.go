/*
	Credential Server
	version 0.9
	author: Adrian Pareja Abarca
	email: adriancc5.5@gmail.com
*/

package business

import (
	"crypto/sha256"

	bl "github.com/lacchain/credential-server/blockchain"
	"github.com/lacchain/credential-server/models"
	"github.com/ethereum/go-ethereum/common"
)

//getCredential from credentials repository
func getCredential(credentialID string, nodeURL string, repositoryContract string) (*models.Credential, error) {
	client := new(bl.Client)
	client.Connect(nodeURL)
	idHash := sha256.Sum256([]byte(credentialID))

	_, err := client.GetCredential(common.HexToAddress(repositoryContract), idHash)
	if err != nil {
		return nil, err
	}
	client.Close()

	return nil, nil
}
