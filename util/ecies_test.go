package util

import (
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSharedKey(t *testing.T) {
	// user1
	sk0 := "1828384858687888000000000000000000000000000000009808a8b8c8d8e8f8"
	pk0 := "022064bcf0d97500b3d10553d29e988bb132fceb5213110046e819ea331505689e"

	// bot
	// NVCCXuVjRKt1ac7ri4WuVC8sDrYeZ6NeCx
	sk1 := "1234567800000000000000000000000000000000000000000000000090abcdef"
	pk1 := "03d0e1f2d72cfc94bc23ad236be921bedc8be1affa283f10faffbdea519c8fd1c1"

	dk0, err := GetSharedKeyFromHex(sk0, pk1)
	assert.Nil(t, err)
	assert.Equal(t, "0xcd751129a4dff2f9eb12f966c7d9c6845c4acfd483dd1507be3b790c1c64278e", dk0)
	//fmt.Println(dk0)

	dk1, err := GetSharedKeyFromHex(sk1, pk0)
	assert.Nil(t, err)
	assert.Equal(t, "0xcd751129a4dff2f9eb12f966c7d9c6845c4acfd483dd1507be3b790c1c64278e", dk1)
	//fmt.Println(dk1)

	msg := `{"type":"Text","data":"What is the weather like in Shanghai?"}`
	encrypted, err := AesEncrypt(msg, dk0)
	assert.Nil(t, err)
	fmt.Println(encrypted)

	//e := "U2FsdGVkX19sPpclNZDi/VfmcGS/+CR6XmNa9X3aafu1m0cm6vtqKOFay+bqesJ/1rydpXJYdpNArng2CA+6wT8YrLkps+ABks2Lkeq1OP69STscR0BCmYPfpX/3o/QAYC5D89Vw1n10NxccWSJgz2lfy+Qc06wsPqZ7SAIgNVWGHyTqs0fIfcseNp3QyLUhf5ZW+Enfpm6S5tGGe57mxidDgYJcetAyuUTjtPzY8S2IlLuP9gR/LOsHJRrl0WlyBYJAAlUdCeqsIuEKZS/3jn/4/N6iF3k9XGeaHBPVV7LXj3VTuPQum7adMP6JUrPqgNFNTyKpvKI/0oZ8hQeO1OqhsNK5KdKa2pD4fgL3cBB1bkjG4GY67UfRT6lcBXuUn7YbyRX1CBH/i89GD3u0+jGMqaoaKZxBXXj7AllfeX37Q//YC3e9zDEGBuwWXDR3cp6a/z3gvQe199CXAd6ILRGV9lSDqU/uHLRt9UfL9sqP3BZmQ0c5z6kh+NiZQwZcVYbhPjFsVUUvPq9ra09IJi0KSpAMQ52yf0vmMir/4tJWm/Cy6kVXEX9r/CjVGkxjqrS+xXtwE3U9whrNUSN3fcHl0+4UXo2a98kt6xKSzkBlphAILmgBdzahzQTDNVL9G14SjAnqsptQ/pv9YcaRaMXDitiXUDtLH6oBJ8yKG/gwMHLoczIRPbiDrc0swlh60F+RF5Hj9lJv7g/ks7zkfscAmimaVjPqECC2XR9f9IsPvDZBRcftHKw4W46MAVQvfCVzYOs/y+wsffzwHWQ+X3XAu4no/+XE4fTf0yBoC6T8usyTnuR1/nyEXJOCjRTAHOfVaULhrwkA82AD3HnJHBTzLpfHW7b4iGeyagYvW37CGwSj8RCrxxByw3JYwINPJ5O9QJACS0IL7NshjIxQzbOp0jLoZmtvHT42+614X5Vv/WMDJkMxIJSCcQ0zGeIvd1XW3Fo/B1TFxrpUsip0kOvNPPn275OSpOjkuAdzFA8v+DRqRkAhZuEyhgW6ity8VR/6C7NU0huFpN/vcWBQblIlgNuA9cSRiiFTISMpNr1xOf/s6lKFudu0NBotsHlZuGwhpP/DWwGscZNd9hPj7TBEzsCgVr8Xp0RE0FNCYx6R1dEnFqosiRgZPqDf5GZhR42PKayfposQTIqi4ICktB8an/p+l6PYsxYNJOIFHySzeZ487fBG9vnkHYwYFbkdBsx/w9I+3kBCN+CuGlQgGlxBlBUmWBEtpE3CEb/BxhumJ5RMWfwSS0rfOA8PyJG7n0ylxcGd/UleEw4ae9Ofvw=="
	//d, err := AesDecrypt(e, dk0)
	//assert.Nil(t, err)
	//fmt.Println(d)
}

func Test2(t *testing.T) {
	// user2
	sk0 := "18283848586878888878685848382818f8e8d8c8b8a808989808a8b8c8d8e8f8"
	pk0 := "031f05c49ebf8c3c8ebcb4d89148e0d4d3240eedfd53f08f1f472faafecb0de2c3"

	// bot
	// NVCCXuVjRKt1ac7ri4WuVC8sDrYeZ6NeCx
	sk1 := "1234567800000000000000000000000000000000000000000000000090abcdef"
	pk1 := "03d0e1f2d72cfc94bc23ad236be921bedc8be1affa283f10faffbdea519c8fd1c1"

	dk0, err := GetSharedKeyFromHex(sk0, pk1)
	assert.Nil(t, err)
	assert.Equal(t, "0xc28b2e6275532467869d9b3335d6c54bbd9105f6e3e2f5b6ff14080f1b017622", dk0)
	//fmt.Println(dk0)

	dk1, err := GetSharedKeyFromHex(sk1, pk0)
	assert.Nil(t, err)
	assert.Equal(t, "0xc28b2e6275532467869d9b3335d6c54bbd9105f6e3e2f5b6ff14080f1b017622", dk1)
	//fmt.Println(dk1)
}

func TestGetPrivateKeyFromHex(t *testing.T) {
	s := "18283848586878888878685848382818f8e8d8c8b8a808989808a8b8c8d8e8f8"
	sk, err := GetPrivateKeyFromHex(s)
	assert.Nil(t, err)
	p := elliptic.MarshalCompressed(sk.Curve, sk.X, sk.Y)
	fmt.Println(hex.EncodeToString(p))
}
