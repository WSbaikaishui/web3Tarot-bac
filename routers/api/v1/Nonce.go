package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"web3Tarot-backend/pkg/logging"
	"web3Tarot-backend/service/nonce"
	"web3Tarot-backend/util"
)

// GetNonce
// @Tags           Nonce
// @Summary        gets a nonce for a user to sign
// @Produce        json
// @Success        200        {object}    nonce.GetNonceData
// @Failure        500        {object}    errors.ErrorInfo
// @Router         /nonces    [POST]
func GetNonce(ctx *gin.Context) {
	data, err := nonce.GetNonce(ctx)
	if err != nil {
		logging.Error(util.EncodeError(err))
		return
	}
	ctx.JSON(http.StatusOK, data)
}
