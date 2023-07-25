package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	apiErr "web3Tarot-backend/errors"
	"web3Tarot-backend/service/user"
	"web3Tarot-backend/util"
)

// Login
// @Tags           User
// @Summary        binds the user's wallet address and create a user
// @Produce        json
// @Param          address    path        string true "address"
// @Param          nonce      body        string true "nonce"
// @Param          signature  body        string true "signature"
// @Success        200        {object}    user.LoginData
// @Failure        400        {object}    apiErr.ErrorInfo
// @Failure        500        {object}    apiErr.ErrorInfo
// @Router         /users/{address}/actions/sign-in       [PUT]
func Login(ctx *gin.Context) {
	params := user.LoginParam{}
	if err := ctx.ShouldBindJSON(&params); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	params.Address = ctx.Param("address")
	if err := params.Validate(); err != nil {
		util.EncodeError(ctx, err)
		return
	}
	data, err := user.Login(ctx, &params)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	//if data.IsNew {
	//	go c.svc.SendWelcome(ctx, params.Address)
	//}
	ctx.JSON(http.StatusOK, data.Token)
}

// GetUser
// @Tags           User
// @Summary        gets all the user's info from db
// @Produce        json
// @Param          address    path        string true "address"
// @Success        200        {object}    user.GetUserData
// @Failure        400        {object}    apiErr.ErrorInfo
// @Failure        500        {object}    apiErr.ErrorInfo
// @Router         /users/{address}       [GET]
func GetUser(ctx *gin.Context) {
	address := ctx.Param("address")
	if len(address) == 0 {
		err := apiErr.ErrInvalidParameter("invalid address")
		util.EncodeError(ctx, err)
		return
	}
	data, err := user.GetUser(ctx, address)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// SetKeyInfo
// @Tags         User
// @Summary      add user's public key and key store info into db
// @Produce      json
// @Param        address    path        string true "address"
// @Param        publicKey  body        string true "publicKey"
// @Param        keyStore   body        string true "keyStore"
// @Success      200        {object}    nil
// @Failure      400        {object}    apiErr.ErrorInfo
// @Failure      500        {object}    apiErr.ErrorInfo
// @Router       /users/{address}/actions/set-e2ee-key      [PUT]
func SetKeyInfo(ctx *gin.Context) {
	param := user.SetKeyInfoParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	param.Address = ctx.Param("address")
	// validate params
	if err := param.Validate(); err != nil {
		util.EncodeError(ctx, err)
		return
	}
	err := user.SetKeyInfo(ctx, &param)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, nil)
}

// GetUserPublicInfo
// @Tags         User
// @Summary      gets the user's public info from db
// @Produce      json
// @Param        addresses  body        []string true "addresses"
// @Success      200        {object}    user.GetPublicUserData
// @Failure      400        {object}    apiErr.ErrorInfo
// @Failure      500        {object}    apiErr.ErrorInfo
// @Router       /public-users/actions/get      [PUT]
func GetUserPublicInfo(ctx *gin.Context) {
	param := user.GetPublicUserParam{}
	if err := ctx.ShouldBindJSON(&param); err != nil {
		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
		return
	}
	data, err := user.GetUserPublicInfo(ctx, param.Addresses)
	if err != nil {
		util.EncodeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, data)
}
