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

var Bits int = 5120

var PrivateKey *rsa.PrivateKey = generateKeyPair(Bits)

var PublicKey *rsa.PublicKey = &PrivateKey.PublicKey

func generateKeyPair(bits int) *rsa.PrivateKey {

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

	data, err := DecryptOAEP(PrivateKey, cipherText)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	return string(data)

}

func DecryptOAEP(private *rsa.PrivateKey, msg []byte) ([]byte, error) {

	msgLen := len(msg)
	step := private.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptPKCS1v15(rand.Reader, PrivateKey, msg[start:finish])
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
