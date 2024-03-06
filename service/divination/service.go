package divination

import (
	"web3Tarot-backend/models"
	"web3Tarot-backend/util"
)

func SetDivinationMessage(param *CreateDivination) error {
	var divination models.Divination
	divination.UserID = param.UserID
	divination.Card1 = param.Card1
	divination.Card1Status = param.Card1Status
	divination.Question = param.Question
	divination.Content = param.Content
	divination.Uuid = util.UuidV4()
	err := models.CreateDivination(&divination)

	return err
}
