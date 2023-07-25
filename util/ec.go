package util

import (
	"crypto/elliptic"
	"fmt"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/joeqian10/neo3-gogogo/keys"
	"math/big"
)

var curve = crypto.P256
var bitLen = (curve.Params().BitSize + 7) / 8

// calculate h = (Q + 1 + 2 * sqrt(Q)) / N
var h = new(big.Int).Div(
	new(big.Int).Add(
		new(big.Int).Add(curve.Params().P, big.NewInt(1)),
		new(big.Int).Mul(big.NewInt(2), new(big.Int).Sqrt(curve.Params().P))),
	curve.Params().N)

// Signature is a type representing an ecdsa signature.
type Signature struct {
	R *big.Int
	S *big.Int
}

// RecoverPubKeyFromSigOnSecp256r1 recover on Secp256r1 from a signature
func RecoverPubKeyFromSigOnSecp256r1(message, signature []byte) ([]*crypto.ECPoint, error) {
	hash := crypto.Sha256(message)

	R := new(big.Int).SetBytes(signature[0:bitLen])
	S := new(big.Int).SetBytes(signature[bitLen:])
	sig := &Signature{
		R: R,
		S: S,
	}

	result := []*crypto.ECPoint{}
	hh := h.Uint64()
	for i := 0; uint64(i) <= hh; i++ {
		// 1.1 x = (n * i) + r
		Rx := new(big.Int).Mul(curve.Params().N, new(big.Int).SetInt64(int64(i)))
		Rx.Add(Rx, sig.R)
		if Rx.Cmp(curve.Params().P) != -1 {
			continue
		}

		for j := 0; j < 2; j++ {
			// convert 02<Rx> to point R. (step 1.2 and 1.3). If we are on an odd
			// iteration then 1.6 will be done with -R, so we calculate the other
			// term when uncompressing the point.
			Ry, err := decodeCompressedY(curve, Rx, uint(j))
			if err != nil {
				continue
			}

			// 1.4 Check n*R is point at infinity
			nRx, nRy := curve.ScalarMult(Rx, Ry, curve.Params().N.Bytes())
			if nRx.Sign() != 0 || nRy.Sign() != 0 {
				continue
			}

			// 1.5 calculate e from message using the same algorithm as ecdsa
			// signature calculation.
			e := hashToInt(hash, curve)

			// Step 1.6.1:
			// We calculate the two terms sR and eG separately multiplied by the
			// inverse of r (from the signature). We then add them to calculate
			// Q = r^-1(sR-eG)
			invr := new(big.Int).ModInverse(sig.R, curve.Params().N)

			// first term.
			invrS := new(big.Int).Mul(invr, sig.S)
			invrS.Mod(invrS, curve.Params().N)
			sRx, sRy := curve.ScalarMult(Rx, Ry, invrS.Bytes())

			// second term.
			e.Neg(e)
			e.Mod(e, curve.Params().N)
			e.Mul(e, invr)
			e.Mod(e, curve.Params().N)
			minuseGx, minuseGy := curve.ScalarBaseMult(e.Bytes())

			// this would be faster if we did a mult and add in one
			// step to prevent the jacobian conversion back and forth.
			Qx, Qy := curve.Add(sRx, sRy, minuseGx, minuseGy)
			Q, err := crypto.CreateECPoint(Qx, Qy, &curve)
			if err != nil {
				return nil, err
			}
			result = append(result, Q)
		}
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("invalid signature")
	}
	return result, nil
}

// hashToInt converts a hash value to an integer. There is some disagreement
// about how this is done. [NSA] suggests that this is done in the obvious
// manner, but [SECG] truncates the hash to the bit-length of the curve order
// first. We follow [SECG] because that's what OpenSSL does. Additionally,
// OpenSSL right shifts excess bits from the number if the hash is too large
// and we mirror that too.
// This is borrowed from crypto/ecdsa.
func hashToInt(hash []byte, c elliptic.Curve) *big.Int {
	orderBits := c.Params().N.BitLen()
	orderBytes := (orderBits + 7) / 8
	if len(hash) > orderBytes {
		hash = hash[:orderBytes]
	}

	ret := new(big.Int).SetBytes(hash)
	excess := len(hash)*8 - orderBits
	if excess > 0 {
		ret.Rsh(ret, uint(excess))
	}
	return ret
}

// decodeCompressedY performs decompression of Y coordinate for given X and Y's least significant bit
func decodeCompressedY(curve elliptic.Curve, x *big.Int, ylsb uint) (*big.Int, error) {
	//c := elliptic.P256()
	cp := curve.Params()
	three := big.NewInt(3)

	xCubed := new(big.Int).Exp(x, three, cp.P)
	threeX := new(big.Int).Mul(x, three)
	threeX.Mod(threeX, cp.P)
	ySquared := new(big.Int).Sub(xCubed, threeX)
	ySquared.Add(ySquared, cp.B)
	ySquared.Mod(ySquared, cp.P)
	y := new(big.Int).ModSqrt(ySquared, cp.P)
	if y == nil {
		return nil, fmt.Errorf("error computing Y for compressed point")
	}
	if y.Bit(0) != ylsb {
		y.Neg(y)
		y.Mod(y, cp.P)
	}
	return y, nil
}

func VerifyAddress(address string, pubKeys []*crypto.ECPoint) bool {
	for _, pubKey := range pubKeys {
		if address == keys.PublicKeyToAddress(pubKey, helper.DefaultAddressVersion) {
			return true
		}
	}
	return false
}
