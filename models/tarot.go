package models

import "fmt"

const tarotTable = "tarot"

type Tarot struct {
	Model
	CardName       string ` gorm:"column:card_name;type:varchar(255);not null;unique"`
	CardNameShort  string `gorm:"column:card_name_short;type:varchar(255);not null;unique"`
	CardValue      string `gorm:"column:card_value;type:varchar(255);not null"`
	CardSuit       string `gorm:"column:card_suit;type:varchar(255);not null"`
	CardType       string `gorm:"column:card_type;type:varchar(255);not null;default:false"`
	CardMeaning    string `gorm:"column:card_meaning;type:text;not null"`
	CardImage      string `gorm:"column:card_image;type:varchar(255);not null"`
	CardNumber     int    `gorm:"column:card_number;type:int(4);not null"`
	CardMeaningUp  string `gorm:"column:card_meaning_up;type:text;not null"`
	CardMeaningRev string `gorm:"column:card_meaning_rev;type:text;not null"`
	CardDescribe   string `gorm:"column:card_describe;type:text;not null"`
}

func GetTarotByCardNumber(cardNumber int) *Tarot {
	var tarot Tarot
	Db.Table(tarotTable).Where("id=?", cardNumber).First(&tarot)
	fmt.Println(tarot)
	return &tarot
}

func GetTarots() (tarots []*Tarot) {
	Db.Table(tarotTable).Find(&tarots)
	return tarots
}

func CreateTarot(data interface{}) error {
	ddb := Db.Table(tarotTable).Create(data)
	if ddb.Error != nil {
		return ddb.Error
	}
	return nil
}

func EditTarot(cardNumber int, data interface{}) bool {
	ddb := Db.Table(tarotTable).Model(&Tarot{}).Where("id = ?", cardNumber).Updates(data)
	if ddb.Error != nil {
		return false
	}
	return true
}
