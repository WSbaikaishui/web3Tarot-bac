package models

const tarotTable = "tarots"

type Tarot struct {
	Model
	CardName    string ` gorm:"column:card_name;type:varchar(255);not null;unique" `
	CardType    bool   `gorm:"column:card_type;type:bool;not null;default:false"`
	CardMeaning string `gorm:"column:card_meaning;type:text;not null"`
	CardImage   string `gorm:"column:card_image;type:varchar(255);not null"`
	CardNumber  int    `gorm:"column:card_number;type:int(4);not null;unique"`
}

func GetTarotByCardNumber(cardNumber int) (tarot Tarot) {
	Db.Where("card_number = ?", cardNumber).First(&tarot)
	return tarot
}

func GetTarots() (tarots []Tarot) {
	Db.Find(&tarots)
	return tarots
}

func CreateTarot(data interface{}) bool {
	ddb := Db.Create(data)
	if ddb.Error != nil {
		return false
	}
	return true
}

func EditTarot(cardNumber int, data interface{}) bool {
	ddb := Db.Model(&Tarot{}).Where("card_number = ?", cardNumber).Updates(data)
	if ddb.Error != nil {
		return false
	}
	return true
}
