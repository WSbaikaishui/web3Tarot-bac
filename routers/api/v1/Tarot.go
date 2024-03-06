package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	apiErr "web3Tarot-backend/errors"
	"web3Tarot-backend/pkg/logging"
	"web3Tarot-backend/service/tarot"
	"web3Tarot-backend/util"
)

// GetTarotCards
// @Tags           Tarot
// @Summary        gets all tarot cards for display
// @Produce        json
// @Success        200        {object}    tarot.GetAllTarot
// @Failure        500        {object}    errors.ErrorInfo
// @Router         /cards/get-all    [GET]
func GetTarotCards(ctx *gin.Context) {
	data, err := tarot.GetAllTarot()
	if err != nil {
		logging.Error(util.EncodeError(err))
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// GetTarotCard
// @Tags           Tarot
// @Summary        gets a tarot card for display
// @Produce        json
// @Success        200        {object}    tarot.GetTarotCard
// @Failure        500        {object}    errors.ErrorInfo
// @Router         /cards/:id    [GET]
func GetTarotCard(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		err := apiErr.ErrInvalidParameter("invalid id")
		logging.Error(util.EncodeError(err))
		return
	}

	data, err := tarot.GetTarotCard(id)
	if err != nil {
		logging.Error(util.EncodeError(err))
		return
	}
	ctx.JSON(http.StatusOK, data)
}

// GetTarotImage
// @Tags           Tarot
// @Summary        gets a tarot image url
// @Produce        json
// @Success        200        {object}    tarot.GetTarotCard
// @Failure        500        {object}    errors.ErrorInfo
// @Router         /cards/:id    [GET]
func GetTarotImage(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		err := apiErr.ErrInvalidParameter("invalid id")
		logging.Error(util.EncodeError(err))
		return
	}

	data, err := tarot.GetTarotCard(id)
	if err != nil {
		logging.Error(util.EncodeError(err))
		return
	}
	ctx.JSON(http.StatusOK, data.CardImage)
}

// UploadTarotCard
// @Tags           Tarot
// @Summary        upload a tarot card
// @Produce        json
// @Success        200        {object}    tarot.UploadTarotCards
// @Failure        500        {object}    errors.ErrorInfo
// @Router         /cards/upload    [POST]
func UploadTarotCard(ctx *gin.Context) {
	var param tarot.UploadTarot
	if err := ctx.ShouldBindJSON(&param); err != nil {
		logging.Error(util.EncodeError(apiErr.ErrInvalidParameter(err.Error())))
		return
	}
	for i, _ := range param.Cards {
		err := tarot.UploadTarotCardMessage(&param.Cards[i])
		if err != nil {
			logging.Error(util.EncodeError(err))
			return
		}
	}

	ctx.JSON(http.StatusOK, "ok")
}
