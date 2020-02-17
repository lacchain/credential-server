/*
	Credential Server
	version 0.9
	author: Adrian Pareja Abarca
	email: adriancc5.5@gmail.com
*/

package business

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/smtp"
	"net/mail"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/common"
	strfmt "github.com/go-openapi/strfmt"
	uid "github.com/segmentio/ksuid"

	bl "github.com/lacchain/credential-server/blockchain"
	"github.com/lacchain/credential-server/models"

	qrcode "github.com/skip2/go-qrcode"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "adrianp@iadb.org"

	// Specify a configuration set. If you do not want to use a configuration
	// set, comment out the following constant and the
	// ConfigurationSetName: aws.String(ConfigurationSet) argument below
	//ConfigurationSet = "ConfigSet"

	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	AwsRegion = "us-east-1"

	// The subject line for the email.
	Subject = "Your Credential Verifiable"

	//The email body for recipients with non-HTML email clients.
	TextBody = "This is your credential verifiable which contains information about the document that was hashed"

	// The character encoding for the email.
	CharSet = "UTF-8"
)

//CreateCredential saving the hash into blockchain
func CreateCredential(subjects []*models.CredentialSubject, nodeURL string, issuer string, privateKey string, verificationContract string, repositoryContract string) ([]*models.Credential, error) {
	client := new(bl.Client)
	err := client.Connect(nodeURL)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	options, err := client.ConfigTransaction(privateKey)
	if err != nil {
		return nil, err
	}
	address := common.HexToAddress(verificationContract)

	//Iterate subject and generate Credentials
	credentials := make([]*models.Credential, 0, 50)
	for _, subject := range subjects {
		credential := new(models.Credential)
		credentialData := new(models.CredentialData)
		credentialMetadata := new(models.CredentialMetadata)
		id := generateID()
		credentialData.ID = &id
		types := make([]string, 0, 4)
		types = append(types, "VerifiableCredential")
		types = append(types, subject.Type)
		credentialData.Type = types
		credentialData.CredentialSubject = subject.Content
		credentialData.IssuanceDate = subject.IssuanceDate
		credentialData.Evidence = subject.Evidence
		credentialData.Issuer = issuer
		credentialData.Proof = getProof("SmartContract", verificationContract)
		credential.CredentialData = credentialData
		credentials = append(credentials, credential)

		rawCredential, err := json.Marshal(credentialData)
		if err != nil {
			return nil, errors.New("Credential isn't Json format")
		}

		fmt.Println("####Credential####",string(rawCredential))

		fmt.Println("IssuanceDate:", subject.IssuanceDate.String())
		fmt.Println("ExpirationDate:", subject.ExpirationDate.String())

		date := time.Time(subject.ExpirationDate)

		credentialHash := sha256.Sum256(rawCredential)

		fmt.Println("###Date###",date)
		fmt.Println("###Date Location###:", date.Local())
		fmt.Println("###Expiration Millis###",big.NewInt(date.Unix()))

		err, tx := client.SignCredential(address, options, credentialHash, big.NewInt(date.Unix()))
		if err != nil {
			fmt.Println("Transaction wasn't sent")
		}
		//waiting for mining
		time.Sleep(4 * time.Second)

		blockNumber, timestamp, err:=client.GetTransactionReceipt(*tx)
		if err != nil{
			fmt.Println("Failed to get Receipt:",err)
		} 

		fmt.Println("BlockkkNumber:",blockNumber)
		fmt.Println("BlockkkNumber:",timestamp)

		credentialMetadata.BlockNumber = blockNumber.String()
		credentialMetadata.Timestamp = timestamp
		credentialMetadata.Transaction = tx.Hex()

		credential.Metadata = credentialMetadata

		/*qrFile, err := generateQR("http://credentialserver.iadb.org/", credentialHash, getNameSubject(subject.Content))
		if err != nil {
			fmt.Printf("Failed generate QR: %s", err)
		}

		err = sendCredentialByEmail(getNameSubject(subject.Content), subject.Email, []byte(string(rawCredential)), tx.Hex(),blockNumber,timestamp, subject.ExpirationDate.String(), qrFile)
		if err != nil {
			fmt.Printf("Failed to send email: %s", err)
		}*/

		//Deprecated code to save credential into blockchain

		/*options, err = client.ConfigTransaction(privateKey)
		if err != nil {
			return nil, err
		}
		idHash := sha256.Sum256([]byte(*credential.ID))

		err = client.SetCredential(common.HexToAddress(repositoryContract), options, idHash, credentialHash)
		if err!=nil{
			fmt.Println("Transaction wasn't sent")
		}*/
	}

	return credentials, nil

}

func getProof(typeProof string, verificationMethod string) *models.Proof {
	var proof = new(models.Proof)
	proof.Type = typeProof
	proof.VerificationMethod = verificationMethod
	proof.Created = strfmt.DateTime(time.Now())

	return proof
}

func sendCredentialByEmail(name string, destination string, credential []byte, tx string, blockNumber *big.Int, timestamp string, expirationDate string, qrFile []byte) error {
	from := mail.Address{"", "notarization@lacchain.net"}
    to   := mail.Address{"", destination}
    subj := "LACCchain - Your document's credential"

    // Setup headers
    headers := make(map[string]string)
    headers["From"] = from.String()
    headers["To"] = to.String()
    headers["Subject"] = subj

    // Setup message
    message := ""
    for k,v := range headers {
        message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	
	log.Println("Mark content to accept multiple contents")
	message += "MIME-Version: 1.0\r\n"
	message += fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", "**=myohmy689407924327")

    //place HTML message
	log.Println("Put HTML message")
	message += fmt.Sprintf("\r\n--%s\r\n", "**=myohmy689407924327")
	message += "Content-Type: text/html; charset=\"utf-8\"\r\n"
	message += "Content-Transfer-Encoding: 7bit\r\n"
	message += fmt.Sprintf("\r\n%s", "<html><body><p>Dear "+name+"</p><p>Congratulations! The hash of your file has been registered successfully in the LACChain Blockchain Network. The hash was registered at "+timestamp+" in the transaction "+tx+", that is in the block "+blockNumber.String()+". Attached is your verifiable credential, that will be valid until the expiration date "+expirationDate+" set by yourself.</p>" +
		"<p>If you have any questions, please do not hesitate to reach our to us at info@lacchain.net.</p>"+
		"<p>Best,</p><p>The LACChain Alliance</p></body></html>\r\n")

	//put QR credential
	log.Println("Put HTML message")
	message += fmt.Sprintf("\r\n--%s\r\n", "**=myohmy689407924327")
	message += "Content-Type: image/png;\r\n"
	message += "Content-Transfer-Encoding: base64\r\n"
	message += "\r\n" + base64.StdEncoding.EncodeToString(qrFile)		

	log.Println("Put file attachment")
	message += fmt.Sprintf("\r\n--%s\r\n", "**=myohmy689407924327")
	message += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
	message += "Content-Transfer-Encoding: base64\r\n"
	message += "Content-Disposition: attachment;filename=\"" + "credential.json" + "\"\r\n"
	message += "\r\n" + base64.StdEncoding.EncodeToString(credential)

    // Connect to the SMTP Server
    servername := "smtp.serviciodecorreo.es:587"

    host, _, _ := net.SplitHostPort(servername)

	auth := smtp.PlainAuth("", "notarization@lacchain.net", "Hashing00", host)

	recipients := []string{destination}

	err1 := smtp.SendMail(servername, auth, "notarization@lacchain.net", recipients, []byte(message))
	if err1 != nil {
		fmt.Println("Error:",err1)
	}

	log.Print("Your mail was sent")
	return nil
}

func generateQR(url string, hashCredential [32]byte, filename string) ([]byte, error) {
	var hash = make([]byte, 32, 64)

	for i, j := range hashCredential {
		hash[i] = j
	}

	log.Println("QR hexa:", hex.EncodeToString(hash))
	log.Println("QR bytes :", hash)
	qrFile, err := qrcode.Encode(url+hex.EncodeToString(hash), qrcode.Medium, 256)
	return qrFile, err
}

func generateID() string {
	id := uid.New()
	return id.String()
}

func getNameSubject(contentSubject interface{}) string {
	content := contentSubject.(map[string]interface{})
	name := content["author"]

	return fmt.Sprintf("%v", name)
}

func getLastNameSubject(contentSubject interface{}) string {
	content := contentSubject.(map[string]interface{})
	lastname := content["lastname"]

	return fmt.Sprintf("%v", lastname)
}

func getReceiverMail(contentSubject interface{}) string {
	content := contentSubject.(map[string]interface{})
	email := content["email"]

	return fmt.Sprintf("%v", email)
}

func getFullURL(contentSubject interface{}) string {
	content := contentSubject.(map[string]interface{})
	fullURL := content["full_size_url"]

	return fmt.Sprintf("%v", fullURL)
}

func getExpirationDate(contentSubject interface{}) string {
	content := contentSubject.(map[string]interface{})
	name := content["expirationDate"]

	return fmt.Sprintf("%v", name)
}