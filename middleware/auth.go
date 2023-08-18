package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"net/http"
	"strings"
	apiErr "web3Tarot-backend/errors"
	"web3Tarot-backend/log"
	"web3Tarot-backend/util"
)

func getBearerToken(req *http.Request) (string, bool) {
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		return "", false
	}
	return authHeaderParts[1], true
}

var skipAuthPrefix = []string{
	"/api/v1/nonces",
	"/api/v1/users/:address/actions/sign-in",
	"/api/v1/public-users",
	"/api/v1/notifications",
}

func canSkipAuth(ctx *gin.Context) bool {
	fullPath := ctx.FullPath()
	for _, p := range skipAuthPrefix {
		if strings.HasPrefix(fullPath, p) {
			return true
		}
	}
	return false
}

func AuthMiddleware(ctx *gin.Context) {
	skipAuth := canSkipAuth(ctx)
	if skipAuth {
		ctx.Next()
		return
	}

	token, ok := getBearerToken(ctx.Request)
	if !ok {
		util.EncodeError(ctx, apiErr.ErrUnauthorized("Empty auth info"))
		return
	}
	if len(token) == 0 {
		util.EncodeError(ctx, apiErr.ErrUnauthorized("Invalid auth info"))
		return
	}
	// check using jwt style
	claims, err := util.VerifyToken(token)
	if err != nil {
		if !strings.Contains(err.Error(), "token is expired") {
			log.Errorf("verifyToken token: %s, err: %v", token, err)
		}
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				util.EncodeError(ctx, apiErr.ErrTokenExpired())
				return
			} else {
				util.EncodeError(ctx, apiErr.ErrUnauthorized("JWT Token invalid"))
				return
			}
		} else {
			util.EncodeError(ctx, apiErr.ErrUnauthorized(err.Error()))
			return
		}
	}
	address160, err := helper.UInt160FromString(claims.UID)
	if err != nil {
		log.Errorf("Transfer address to Uint160: %s, err: %v", claims.UID, err)
	}
	address := crypto.ScriptHashToAddress(address160, helper.DefaultAddressVersion)
	ctx.Set(util.AuthKey, address)
	ctx.Next()
}
