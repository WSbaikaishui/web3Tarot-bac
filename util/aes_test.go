package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAes(t *testing.T) {
	p := "Hello world!"
	key := "0x615469e0655c6fc2e9b35a42a6024079e68f71a8fcdfd22d6a83dbc80476eec3"
	c, err := AesEncrypt(p, key)
	assert.Nil(t, err)
	plain, err := AesDecrypt(c, key)
	assert.Nil(t, err)
	assert.Equal(t, p, plain)
}

func TestAesEncrypt(t *testing.T) {
	p := "How to make Moutai?"
	key := "0x615469e0655c6fc2e9b35a42a6024079e68f71a8fcdfd22d6a83dbc80476eec3"
	c, err := AesEncrypt(p, key)
	assert.Nil(t, err)
	fmt.Println(c)
}

//func TestEncrypt(t *testing.T) {
//	prv1, err := ecies.GenerateKey(rand.Reader, secp256k1.S256(), nil)
//	assert.Nil(t, err)
//
//	message := []byte("Hello, world.")
//	ct, err := Encrypt(rand.Reader, nil, &prv1.PublicKey, message, nil, nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	fmt.Println(helper.BytesToHex(ct))
//
//	pt, err := prv1.Decrypt(ct, nil, nil)
//	fmt.Println(string(pt))
//}
