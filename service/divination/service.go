package divination

import (
	"web3Tarot-backend/models"
	"web3Tarot-backend/util"
)

func SetDivinationMessage(param *CreateDivination) error {
	var divination models.Divination
	divination.UserAddress = param.UserAddress
	divination.Card1 = param.Card1
	divination.Card1Status = param.Card1Status
	divination.Card2 = param.Card2
	divination.Card2Status = param.Card2Status
	divination.Card3 = param.Card3
	divination.Card3Status = param.Card3Status
	divination.Content = param.Content
	divination.Uuid = util.UuidV4()
	err := models.CreateDivination(&divination)

	return err
}
