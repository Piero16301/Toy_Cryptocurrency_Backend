package functions

import (
	"Toy_Cryptocurrency/models"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func GetProofOfWork(previousProof int) int {
	newProof := 1
	checkProof := false
	var numberOfZeros = 6

	for checkProof == false {
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
