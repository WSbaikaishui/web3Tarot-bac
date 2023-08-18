package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

const nonceTable = "nonce"

type Nonce struct {
	ID uint `gorm:"column:id;type:int(11) unsigned NOT NULL AUTO_INCREMENT;PRIMARY_KEY"`
	//Address    string `gorm:"column:address;type:VARCHAR(255);NOT NULL;DEFAULT:'';unique"`
	Nonce      string `gorm:"column:nonce;type:VARCHAR(255);NOT NULL;DEFAULT:'';unique"`
	Expiration uint   `gorm:"column:expiration;type:int(11) unsigned NOT NULL"`

	CreatedAt time.Time `gorm:"column:created_at;<-:false;type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"column:updated_at;<-:false;type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// CreateNonce creates a Nonce
func CreateNonce(ctx context.Context, nonce *Nonce) error {
	isSuccess := Db.Table(nonceTable).Create(nonce).Error
	if isSuccess != nil {
		return isSuccess
	}
	return nil
}

func GetNonce(ctx context.Context, rawNonce string) (*Nonce, bool, error) {
	nonce := &Nonce{}
	err := Db.Table(nonceTable).Where("nonce = ? ", rawNonce).First(nonce).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return nonce, true, nil
}

func DeleteNonce(ctx context.Context, nonce *Nonce) error {
	if nonce.ID == 0 {
		return fmt.Errorf("id is empty")
	}
	return Db.Table(nonceTable).Delete(nonce).Error
}
