package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"hash"
)

func AesEncrypt(plainText, password string) (string, error) {
	salt := make([]byte, 8)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	key, iv, err := DefaultEvpKDF([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	cipherBytes := PKCS5Padding([]byte(plainText), aes.BlockSize)
	mode.CryptBlocks(cipherBytes, cipherBytes)

	data := make([]byte, 16+len(cipherBytes))
	copy(data[:8], []byte("Salted__"))
	copy(data[8:16], salt)
	copy(data[16:], cipherBytes)

	cipherText := base64.StdEncoding.EncodeToString(data)
	return cipherText, nil
}

func AesDecrypt(cipherText, password string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}
	if string(data[:8]) != "Salted__" {
		return "", fmt.Errorf("invalid crypto js aes encryption")
	}

	salt := data[8:16]
	cipherBytes := data[16:]
	key, iv, err := DefaultEvpKDF([]byte(password), salt)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherBytes, cipherBytes)

	result := PKCS5UnPadding(cipherBytes)
	return string(result), nil
}

// https://stackoverflow.com/questions/27677236/encryption-in-javascript-and-decryption-with-php/27678978#27678978
// https://github.com/brix/crypto-js/blob/8e6d15bf2e26d6ff0af5277df2604ca12b60a718/src/evpkdf.js#L55
func EvpKDF(password []byte, salt []byte, keySize int, iterations int, hashAlgorithm string) ([]byte, error) {
	var block []byte
	var hasher hash.Hash
	derivedKeyBytes := make([]byte, 0)
	switch hashAlgorithm {
	case "md5":
		hasher = md5.New()
	default:
		return []byte{}, fmt.Errorf("not implement hasher algorithm")
	}
	for len(derivedKeyBytes) < keySize*4 {
		if len(block) > 0 {
			hasher.Write(block)
		}
		hasher.Write(password)
		hasher.Write(salt)
		block = hasher.Sum([]byte{})
		hasher.Reset()

		for i := 1; i < iterations; i++ {
			hasher.Write(block)
			block = hasher.Sum([]byte{})
			hasher.Reset()
		}
		derivedKeyBytes = append(derivedKeyBytes, block...)
	}
	return derivedKeyBytes[:keySize*4], nil
}

func DefaultEvpKDF(password []byte, salt []byte) (key []byte, iv []byte, err error) {
	// https://github.com/brix/crypto-js/blob/8e6d15bf2e26d6ff0af5277df2604ca12b60a718/src/cipher-core.js#L775
	keySize := 256 / 32
	ivSize := 128 / 32
	derivedKeyBytes, err := EvpKDF(password, salt, keySize+ivSize, 1, "md5")
	if err != nil {
		return []byte{}, []byte{}, err
	}
	return derivedKeyBytes[:keySize*4], derivedKeyBytes[keySize*4:], nil
}

// https://stackoverflow.com/questions/41579325/golang-how-do-i-decrypt-with-des-cbc-and-pkcs7
func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

//// AesDecrypt : cipherText is base64 encoded OpenSSL format
//func AesDecrypt(cipherText, cipher string) (string, error) {
//	cipherBytes, err := crypto.Base64Decode(cipherText)
//	if err != nil {
//		return "", err
//	}
//	if len(cipherBytes) < 16 {
//		return "", fmt.Errorf("invalid cipherText")
//	}
//	key := []byte(cipher)
//	salt := cipherBytes[8:16]
//	raw := cipherBytes[16:]
//	decrypted, err := aesDecrypt(raw, salt, key)
//	if err != nil {
//		return "", err
//	}
//	return string(decrypted), nil
//}
//
//func aesEncrypt(plain, key, iv []byte) ([]byte, error) {
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return nil, err
//	}
//	blockSize := block.BlockSize()
//	origData := PKCS7Padding(plain, blockSize)
//	blockMode := cipher.NewCBCEncrypter(block, iv)
//	cryted := make([]byte, len(origData))
//	blockMode.CryptBlocks(cryted, origData)
//	return cryted, nil
//}
//
//func aesDecrypt(encrypted, salt, key []byte) ([]byte, error) {
//	fmt.Println("key: ", helper.BytesToHex(key))
//	fmt.Println("salt: ", helper.BytesToHex(salt))
//	//dk, iv, err := DefaultEvpKDF(key, salt)
//	dk, iv := EVP_BytesToKeyIV(salt, key)
//	fmt.Println("dk: ", helper.BytesToHex(dk))
//	fmt.Println("iv: ", helper.BytesToHex(iv))
//
//	block, err := aes.NewCipher(dk) // 32 bytes, so using AES-256
//	if err != nil {
//		return nil, err
//	}
//	blockMode := cipher.NewCBCDecrypter(block, iv) // iv is 16 bytes
//	origData := make([]byte, len(encrypted))
//	blockMode.CryptBlocks(origData, encrypted)
//	origData, err = PKCS7UnPadding(origData, block.BlockSize())
//	if err != nil {
//		return nil, fmt.Errorf("PKCS7UnPadding error: %v", err)
//	}
//	return origData, nil
//}
//
//// EVP_BytesToKeyIV gets key and iv using MD5, key is 32 bytes, iv is 16 bytes
//func EVP_BytesToKeyIV(salt []byte, key []byte) ([]byte, []byte) {
//	//var derived []byte
//	//var tmp [16]byte
//	//for len(derived) < 48 {
//	//	tmp = md5.Sum(append(append(tmp[:], key...), salt...))
//	//	derived = append(derived, tmp[:]...)
//	//}
//	//return derived[:32], derived[32:]
//
//	derived := make([]byte, 0)
//	hasher := md5.New()
//	tmp := make([]byte, 0)
//	for len(derived) < 48 {
//		if len(tmp) > 0 {
//			hasher.Write(tmp)
//		}
//		hasher.Write(key)
//		hasher.Write(salt)
//		tmp = hasher.Sum([]byte{})
//		hasher.Reset()
//		derived = append(derived, tmp...)
//	}
//	return derived[:32], derived[32:]
//}
//
//func PKCS7Padding(data []byte, blockSize int) []byte {
//	padLen := blockSize - len(data)%blockSize
//	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
//	return append(data, padding...)
//}
//
//func PKCS7UnPadding(data []byte, blockSize int) ([]byte, error) {
//	length := len(data)
//	if length == 0 {
//		return nil, fmt.Errorf("data is empty")
//	}
//	if length%blockSize != 0 {
//		return nil, fmt.Errorf("data is not block-aligned")
//	}
//	padLen := int(data[length-1])
//	ref := bytes.Repeat([]byte{byte(padLen)}, padLen)
//	if padLen > blockSize || padLen == 0 || !bytes.HasSuffix(data, ref) {
//		return nil, fmt.Errorf("invalid padding")
//	}
//	return data[:length-padLen], nil
//}
