package business

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	secp256k1 "github.com/ethereum/go-ethereum/crypto/secp256k1"
)

//ValidateDid Signature and tamper proof using JWT format
func ValidateDid(DID string) (bool, error) {
	log.Println("Caller is using JWT did")

	header, payload, signature, err := validateJWT(DID)
	if err != nil {
		log.Println("Error validating JWT")
		return false, err
	}

	log.Println(len(signature))

	parte1 := base64.StdEncoding.EncodeToString(header)
	parte2 := base64.StdEncoding.EncodeToString(payload)
	parte3 := base64.StdEncoding.EncodeToString(signature)

	log.Printf("base64 signature:%v", parte3)

	parte1 = strings.TrimRight(parte1, "=")
	parte2 = strings.TrimRight(parte2, "=")

	hashInput := parte1 + "." + parte2

	hash256 := sha256.Sum256([]byte(hashInput))

	var hashVerify = make([]byte, 32, 64)

	for i, j := range hash256 {
		hashVerify[i] = j
	}

	hexaSig := hex.EncodeToString(signature)

	log.Printf("firmaaHexa:%v", hexaSig)

	var payloadJSON interface{}
	err = json.Unmarshal(payload, &payloadJSON)
	if err != nil {
		return false, err
	}

	data := payloadJSON.(map[string]interface{})

	log.Println("iss:", data["iss"])

	//add 04 before publicKey because Y is into publicKey
	pubKey, err := hex.DecodeString("04" + data["iss"].(string))

	if err != nil {
		return false, err
	}

	validate := secp256k1.VerifySignature(pubKey, hashVerify, signature)

	log.Println("is validate?:", validate)

	return validate, nil
}

// validate JWT format
func validateJWT(DID string) ([]byte, []byte, []byte, error) {
	if strings.Count(DID, ".") != 2 {
		return nil, nil, nil, errors.New("Isn't a JWT format")
	}
	parts := strings.Split(DID, ".")

	header, err := jwt.DecodeSegment(parts[0])
	if err != nil {
		return nil, nil, nil, errors.New("Isn't a JWT format")
	}
	payload, err := jwt.DecodeSegment(parts[1])
	if err != nil {
		return nil, nil, nil, errors.New("Isn't a JWT format")
	}
	signature, err := jwt.DecodeSegment(parts[2])
	if err != nil {
		return nil, nil, nil, errors.New("Isn't a JWT format")
	}

	return header, payload, signature, nil
}
