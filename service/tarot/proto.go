package tarot

import "web3Tarot-backend/models"

type GetAllTarotCards struct {
	CardName   string `json:"card_name"`
	CardImage  string `json:"card_image"`
	CardNumber int    `json:"card_number"`
}

type UploadTarotCard struct {
	CardName       string `json:"name"`
	CardNameShort  string `json:"name_short"`
	CardValue      string `json:"value"`
	CardSuit       string `json:"suit"`
	CardType       string `json:"type"`
	CardImage      string `json:"card_image"`
	CardNumber     int    `json:"value_int"`
	CardMeaningUp  string `json:"meaning_up"`
	CardMeaningRev string `json:"meaning_rev"`
	CardDescribe   string `json:"desc"`
}

type UploadTarot struct {
	Cards []UploadTarotCard `json:"cards"`
}

func (t *GetAllTarotCards) FromModel(tarot *models.Tarot) {
	t.CardName = tarot.CardName
	t.CardImage = tarot.CardImage
	t.CardNumber = tarot.ID
}
