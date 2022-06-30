package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"net/http"

	"fmt"
)

var Bits int = 1024

var PrivateKey *rsa.PrivateKey = generateKeyPair(Bits)

var PublicKey *rsa.PublicKey = &PrivateKey.PublicKey

func generateKeyPair(bits int) *rsa.PrivateKey {
	// This method requires a random number of bits.
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	return privateKey
}

func SendPublicKey(w http.ResponseWriter, r *http.Request) {

	derPkix, _ := x509.MarshalPKIXPublicKey(PublicKey)

	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}

	json.NewEncoder(w).Encode(string(pem.EncodeToMemory(block)))

}

func DescryptMessage(passCrypted string) string {

	cipherText, _ := base64.StdEncoding.DecodeString(passCrypted)

	data, err := rsa.DecryptPKCS1v15(rand.Reader, PrivateKey, cipherText)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	return string(data)

}