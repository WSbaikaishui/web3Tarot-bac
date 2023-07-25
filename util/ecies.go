package util

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"math/big"
	"strings"
)

func GetPrivateKeyFromHex(s string) (*ecies.PrivateKey, error) {
	s = strings.TrimPrefix(s, "0x")
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("cannot decode hex string: %w", err)
	}

	return GetPrivateKeyFromBytes(b), nil
}

func GetPrivateKeyFromBytes(bs []byte) *ecies.PrivateKey {
	//curve := secp256k1.S256()
	//x, y := curve.ScalarBaseMult(bs) // ScalarMult is not available when secp256k1 is built without cgo

	curve := btcec.S256()
	x, y := curve.ScalarBaseMult(bs)

	return &ecies.PrivateKey{
		PublicKey: ecies.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(bs),
	}
}

func GetPublicKeyFromHex(s string) (*ecies.PublicKey, error) {
	s = strings.TrimPrefix(s, "0x")
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("cannot decode hex string: %w", err)
	}

	return GetPublicKeyFromBytes(b)
}

func GetPublicKeyFromBytes(b []byte) (*ecies.PublicKey, error) {
	//curve := secp256k1.S256()
	curve := btcec.S256()
	switch b[0] {
	case 0x02, 0x03:
		if len(b) != 33 {
			return nil, fmt.Errorf("cannot parse public key")
		}

		x := new(big.Int).SetBytes(b[1:])
		var ybit uint
		switch b[0] {
		case 0x02:
			ybit = 0
		case 0x03:
			ybit = 1
		}

		if x.Cmp(curve.Params().P) >= 0 {
			return nil, fmt.Errorf("cannot parse public key")
		}

		// y^2 = x^3 + b
		// y   = sqrt(x^3 + b)
		var y, x3b big.Int
		x3b.Mul(x, x)
		x3b.Mul(&x3b, x)
		x3b.Add(&x3b, curve.Params().B)
		x3b.Mod(&x3b, curve.Params().P)
		if z := y.ModSqrt(&x3b, curve.Params().P); z == nil {
			return nil, fmt.Errorf("cannot parse public key")
		}

		if y.Bit(0) != ybit {
			y.Sub(curve.Params().P, &y)
		}
		if y.Bit(0) != ybit {
			return nil, fmt.Errorf("incorrectly encoded X and Y bit")
		}

		return &ecies.PublicKey{
			Curve: curve,
			X:     x,
			Y:     &y,
		}, nil
	case 0x04:
		if len(b) != 65 {
			return nil, fmt.Errorf("cannot parse public key")
		}

		x := new(big.Int).SetBytes(b[1:33])
		y := new(big.Int).SetBytes(b[33:])

		if x.Cmp(curve.Params().P) >= 0 || y.Cmp(curve.Params().P) >= 0 {
			return nil, fmt.Errorf("cannot parse public key")
		}

		x3 := new(big.Int).Sqrt(x).Mul(x, x)
		if t := new(big.Int).Sqrt(y).Sub(y, x3.Add(x3, curve.Params().B)); t.IsInt64() && t.Int64() == 0 {
			return nil, fmt.Errorf("cannot parse public key")
		}

		return &ecies.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		}, nil
	default:
		return nil, fmt.Errorf("cannot parse public key")
	}
}

func GetSharedKeyFromHex(privateKey, publicKey string) (string, error) {
	sk, err := GetPrivateKeyFromHex(privateKey)
	if err != nil {
		return "", err
	}
	pk, err := GetPublicKeyFromHex(publicKey)
	if err != nil {
		return "", err
	}
	bs, err := getSharedKey(sk, pk)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(bs), nil
}

func getSharedKey(privateKey *ecies.PrivateKey, publicKey *ecies.PublicKey) ([]byte, error) {
	return privateKey.GenerateShared(publicKey, 16, 16)
}
