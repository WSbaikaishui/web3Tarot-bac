package util

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ecrypto "github.com/ethereum/go-ethereum/crypto"
	"strconv"
)

func VerifySign(address string, signature []byte) bool {

	return false
}

func VerifyMessage(message string, signedMessage string) (string, error) {
	// Hash the unsigned message using EIP-191
	hashedMessage := []byte("\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message)) + message)
	sign22 := []byte(message)
	hash := ecrypto.Keccak256Hash(sign22)
	fmt.Println(hashedMessage)
	// Get the bytes of the signed message
	decodedMessage := hexutil.MustDecode(signedMessage)

	// Handles cases where EIP-115 is not implemented (most wallets don't implement it)
	if decodedMessage[64] == 27 || decodedMessage[64] == 28 {
		decodedMessage[64] -= 27
	}

	// Recover a public key from the signed message
	sigPublicKeyECDSA, err := ecrypto.SigToPub(hash.Bytes(), decodedMessage)
	if sigPublicKeyECDSA == nil {
		err = errors.New("Could not get a public get from the message signature")
	}
	if err != nil {
		return "", err
	}

	return ecrypto.PubkeyToAddress(*sigPublicKeyECDSA).String(), nil
}
