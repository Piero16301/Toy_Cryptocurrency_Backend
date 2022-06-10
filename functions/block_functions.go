package functions

import (
	"Toy_Cryptocurrency/models"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func GetProofOfWork(previousProof int) int {
	newProof := 1
	checkProof := false
	var numberOfZeros = 6

	for !checkProof {
		newNumber := newProof*newProof - previousProof*previousProof
		newHash := EncryptSHA256Int(newNumber)
		if newHash[:numberOfZeros] == strings.Repeat("0", numberOfZeros) {
			checkProof = true
		} else {
			newProof += 1
		}
	}

	return newProof
}

func ValidateSignature(stringPrivateKey string, stringSignature string) bool {
	// Convertir la clave privada a formato PEM
	pemPrivateKey := fmt.Sprintf(`-----BEGIN RSA PRIVATE KEY-----
%s
-----END RSA PRIVATE KEY-----`, stringPrivateKey)
	data, _ := pem.Decode([]byte(pemPrivateKey))
	privateKey, _ := x509.ParsePKCS1PrivateKey(data.Bytes)
	publicKey := &privateKey.PublicKey

	// Construir hash del mensaje
	message := []byte("bloque firmado")
	messageHash := sha256.New()
	_, err := messageHash.Write(message)
	if err != nil {
		return false
	}
	messageHashSum := messageHash.Sum(nil)

	// Verificar firma
	signature, _ := base64.StdEncoding.DecodeString(stringSignature)
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, messageHashSum, signature)
	return err == nil
}

func EncryptSHA256Int(number int) string {
	stringNumber := strconv.Itoa(number)
	newHash := sha256.Sum256([]byte(stringNumber))
	return fmt.Sprintf("%x", newHash[:])
}

func EncryptSHA256String(text string) string {
	newHash := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", newHash[:])
}

func EncryptSHA256Block(previousBlock models.Block) string {
	byteBlock, _ := json.Marshal(previousBlock)
	stringBlock := string(byteBlock)
	blockHash := sha256.Sum256([]byte(stringBlock))
	return fmt.Sprintf("%x", blockHash[:])
}

func GetMacAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "00000"
	}
	var macAddresses []string
	for _, interf := range interfaces {
		address := interf.HardwareAddr.String()
		if address != "" {
			macAddresses = append(macAddresses, address)
		}
	}
	return macAddresses[0]
}
