package v1

import (
	"github.com/gin-gonic/gin"
	apiErr "web3Tarot-backend/errors"
	"web3Tarot-backend/service/divination"
	"web3Tarot-backend/util"
)

func GetDivination(ctx *gin.Context) {
	wa := ctx.Value(util.AuthKey).(string)
	var param divination.CreateDivination
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	if param.UserAddress != wa {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter("invalid address"))
		return
	}
	err := divination.SetDivinationMessage(&param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
}
