package models

import "time"

const divinationTable = "divination"

type Divination struct {
	Model
	Uuid        string    `gorm:"column:uuid;type:varchar(255);NOT NULL"`
	UserID      int       `gorm:"column:user_id;type:bigint;NOT NULL"`
	Card1       int       `gorm:"column:card1;type:varchar(255);NOT NULL"`
	Card1Status bool      `gorm:"column:card1_status;type:bool;"`
	Card2       int       `gorm:"column:card2;type:varchar(255);"`
	Card2Status bool      `gorm:"column:card2_status;type:bool;"`
	Card3       int       `gorm:"column:card3;type:varchar(255);"`
	Card3Status bool      `gorm:"column:card3_status;type:bool;"`
	Question    string    `gorm:"column:question;type text;NOT NULL"`
	Content     string    `gorm:"column:content;type:TEXT;NOT NULL"`
	IsEncrypted bool      `gorm:"column:is_encrypted;type:bool;NOT NULL;DEFAULT:false"`
	Time        string    `gorm:"column:time;type:varchar(255);NOT NULL"`
	TxID        string    `gorm:"column:tx_id;type:varchar(255);NOT NULL"`
	CreatedAt   time.Time `gorm:"column:created_at;<-:false;type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updated_at;<-:false;type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func GetDivinationByUserID(user_id int) (divination []*Divination) {
	Db.Table(divinationTable).Where("user_id = ?", user_id).First(&divination).Limit(10)
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

//func EditDivination(userAddress string, data interface{}) bool {
//	ddb := Db.Table(divinationTable).Model(&Divination{}).Where("user_address = ?", userAddress).Updates(data)
//	if ddb.Error != nil {
//		return false
//	}
//	return true
//}
