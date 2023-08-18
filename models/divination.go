package models

import "time"

const divinationTable = "divinations"

type Divination struct {
	Model
	Uuid        string    `gorm:"column:uuid;type:varchar(255);NOT NULL"`
	UserAddress string    `gorm:"column:user_address;type:varchar(255);NOT NULL"`
	Card1       int       `gorm:"column:card1;type:varchar(255);NOT NULL"`
	Card1Status bool      `gorm:"column:card1_status;type:bool;NOT NULL"`
	Card2       int       `gorm:"column:card2;type:varchar(255);NOT NULL"`
	Card2Status bool      `gorm:"column:card2_status;type:bool;NOT NULL"`
	Card3       int       `gorm:"column:card3;type:varchar(255);NOT NULL"`
	Card3Status bool      `gorm:"column:card3_status;type:bool;NOT NULL"`
	Content     string    `gorm:"column:content;type:TEXT;NOT NULL"`
	IsEncrypted bool      `gorm:"column:is_encrypted;type:bool;NOT NULL;DEFAULT:false"`
	Time        string    `gorm:"column:time;type:varchar(255);NOT NULL"`
	TxID        string    `gorm:"column:tx_id;type:varchar(255);NOT NULL"`
	CreatedAt   time.Time `gorm:"column:created_at;<-:false;type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updated_at;<-:false;type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func GetDivinationByUserAddress(userAddress string) (divination []*Divination) {
	Db.Table(divinationTable).Where("user_address = ?", userAddress).First(&divination)
	return divination
}

func GetDivinations() (divinations []*Divination) {
	Db.Table(divinationTable).Find(&divinations)
	return divinations
}

func GetDivinationByUuid(uuid string) (divination *Divination) {
	Db.Table(divinationTable).Where("uuid = ?", uuid).First(&divination)
	return divination
}

func CreateDivination(data interface{}) error {
	ddb := Db.Table(divinationTable).Create(data)
	return ddb.Error
}

func EditDivination(userAddress string, data interface{}) bool {
	ddb := Db.Table(divinationTable).Model(&Divination{}).Where("user_address = ?", userAddress).Updates(data)
	if ddb.Error != nil {
		return false
	}
	return true
}
