package models

type tarot struct {
	Model
	CardName    string ` gorm:"column:card_name;type:varchar(255);not null;unique" `
	CardType    bool   `gorm:"column:card_type;type:bool;not null;default:false"`
	CardMeaning string `gorm:"column:card_meaning;type:text;not null"`
	CardImage   string `gorm:"column:card_image;type:varchar(255);not null"`
	CardNumber  int    `gorm:"column:card_number;type:int(4);not null;unique"`
}

func GetTarotByCardNumber(cardNumber int) (tarot tarot) {
	db.Where("card_number = ?", cardNumber).First(&tarot)
	return tarot
}

func GetTarots() (tarots []tarot) {
	db.Find(&tarots)
	return tarots
}

func CreateTarot(data interface{}) bool {
	ddb := db.Create(data)
	if ddb.Error != nil {
		return false
	}
	return true
}

func EditTarot(cardNumber int, data interface{}) bool {
	ddb := db.Model(&tarot{}).Where("card_number = ?", cardNumber).Updates(data)
	if ddb.Error != nil {
		return false
	}
	return true
}
