package util

import (
	"fmt"
	"testing"
)

func TestVerifySignature(t *testing.T) {
	initdata := "Hello, World!"
	sign := "0x9c7fda8e5ff8b061c165b1b508dcd4ae4e2b06ae322f0674b904b8ef72363e251133e4dce5460193500ba11b5ba74d37c791893af30b5bddc89eb2c946d068551b"
	publicKey := "0xc7ae166404DfA77D2AF214e04Bb8B930274A02b8"

	publicKey2, err := VerifyMessage(initdata, sign)
	fmt.Println(publicKey == publicKey2, err)
}
