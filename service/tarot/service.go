package tarot

import (
	"fmt"
	"web3Tarot-backend/models"
)

func GetAllTarot() (*[]GetAllTarotCards, error) {
	msgs := models.GetTarots()
	if len(msgs) == 0 {
		return nil, fmt.Errorf("no tarot cards")
	}
	messageArray := make([]GetAllTarotCards, len(msgs))
	for i, _ := range msgs {
		messageArray[i].FromModel(msgs[i])
	}
	return &messageArray, nil
}

func GetTarotCard(id int) (*models.Tarot, error) {
	msg := models.GetTarotByCardNumber(id)
	if msg == nil {
		return nil, fmt.Errorf("no tarot card")
	}
	return msg, nil
}

func UploadTarotCardMessage(param *UploadTarotCard) error {
	tarot := models.Tarot{
		CardName:       param.CardName,
		CardNameShort:  param.CardNameShort,
		CardValue:      param.CardValue,
		CardSuit:       param.CardSuit,
		CardType:       param.CardType,
		CardImage:      param.CardImage,
		CardNumber:     param.CardNumber,
		CardMeaningUp:  param.CardMeaningUp,
		CardMeaningRev: param.CardMeaningRev,
		CardDescribe:   param.CardDescribe,
	}
	fmt.Println(tarot.CardImage)
	err := models.CreateTarot(&tarot)
	if err != nil {
		return err
	}
	return nil
}
