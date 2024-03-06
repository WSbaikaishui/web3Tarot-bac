package v1

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"web3Tarot-backend/models"
	"web3Tarot-backend/pkg/logging"
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

func Login(update tgbotapi.Update, method int) bool {
	params := user.LoginParam{}
	if method == 0 {
		params.UserID = update.Message.From.ID
		params.Name = update.Message.From.UserName
		params.FirstName = update.Message.From.FirstName
	} else {
		params.UserID = update.CallbackQuery.From.ID
		params.Name = update.CallbackQuery.From.UserName
		params.FirstName = update.CallbackQuery.From.FirstName
	}

	if err := params.Validate(); err != nil {
		logging.Error(util.EncodeError(err))
	}
	data, err := user.Login(&params)
	if err != nil {
		logging.Error(util.EncodeError(err))
	}
	return data.IsNew

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
func GetUserBalance(update tgbotapi.Update, method int) int {
	id := 0
	if method == 0 {
		id = update.Message.From.ID
	} else {
		id = update.CallbackQuery.From.ID
	}
	data, err := user.GetUser(id)
	if err != nil {
		logging.Error(util.EncodeError(err))
		return 0
	}
	return data.Token
}

func GetUserAddress(update tgbotapi.Update, method int) string {
	id := 0
	if method == 0 {
		id = update.Message.From.ID
	} else {
		id = update.CallbackQuery.From.ID
	}
	data, err := user.GetUser(id)
	if err != nil {
		logging.Error(util.EncodeError(err))
		return ""
	}
	return data.Address
}
func IsUserExist(update tgbotapi.Update, method int) bool {
	if method == 0 {
		return models.IsUserExist(update.Message.From.ID)
	}
	return models.IsUserExist(update.CallbackQuery.From.ID)
}

func UpdateUserAddress(ID int, address string) bool {
	err := models.UpdateUserAddress(ID, address)
	if err != nil {
		return false
	}
	return true
}

//// SetKeyInfo
//// @Tags         User
//// @Summary      add user's public key and key store info into db
//// @Produce      json
//// @Param        address    path        string true "address"
//// @Param        publicKey  body        string true "publicKey"
//// @Param        keyStore   body        string true "keyStore"
//// @Success      200        {object}    nil
//// @Failure      400        {object}    apiErr.ErrorInfo
//// @Failure      500        {object}    apiErr.ErrorInfo
//// @Router       /users/{address}/actions/set-e2ee-key      [PUT]
//func SetKeyInfo(ctx *gin.Context) {
//	param := user.SetKeyInfoParam{}
//	if err := ctx.ShouldBindJSON(&param); err != nil {
//		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
//		return
//	}
//	param.Address = ctx.Param("address")
//	// validate params
//	if err := param.Validate(); err != nil {
//		util.EncodeError(ctx, err)
//		return
//	}
//	err := user.SetKeyInfo(ctx, &param)
//	if err != nil {
//		util.EncodeError(ctx, err)
//		return
//	}
//	ctx.JSON(http.StatusOK, nil)
//}
//
//// GetUserPublicInfo
//// @Tags         User
//// @Summary      gets the user's public info from db
//// @Produce      json
//// @Param        addresses  body        []string true "addresses"
//// @Success      200        {object}    user.GetPublicUserData
//// @Failure      400        {object}    apiErr.ErrorInfo
//// @Failure      500        {object}    apiErr.ErrorInfo
//// @Router       /public-users/actions/get      [PUT]
//func GetUserPublicInfo(ctx *gin.Context) {
//	param := user.GetPublicUserParam{}
//	if err := ctx.ShouldBindJSON(&param); err != nil {
//		util.EncodeError(ctx, apiErr.ErrInvalidParameter(err.Error()))
//		return
//	}
//	data, err := user.GetUserPublicInfo(ctx, param.Addresses)
//	if err != nil {
//		util.EncodeError(ctx, err)
//		return
//	}
//	ctx.JSON(http.StatusOK, data)
//}
