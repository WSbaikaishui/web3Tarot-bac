package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

const (
	SignSecretKey = `\x88D\xf09\x6\xa0A\x7\xc5V\xbe\x8b\xef\xd7\xd8\xd3\xe6\x98*4`
	AuthKey       = "AuthKey"
)

type TokenClaims struct {
	UID string `json:"uid"`
	jwt.StandardClaims
}

func GenerateToken(uid string, expiration int64) (string, error) {
	claims := TokenClaims{
		UID: uid,
		StandardClaims: jwt.StandardClaims{
			Id:        uid,
			ExpiresAt: expiration,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	st, err := token.SignedString([]byte(SignSecretKey))
	if err != nil {
		return "", err
	}
	return st, nil
}

func VerifyToken(token string) (*TokenClaims, error) {
	var tokenClaims TokenClaims
	_, err := jwt.ParseWithClaims(token, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		// validate alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}

		return []byte(SignSecretKey), nil
	})
	if err != nil {
		return nil, err
	}
	return &tokenClaims, nil
}
