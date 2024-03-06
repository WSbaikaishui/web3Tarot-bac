package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

const userTable = "user"

type User struct {
	Model
	UserId      int       `gorm:"column:user_id; NOT NULL;unique"`
	FirstName   string    `gorm:"column:first_name;type:VARCHAR(255)"`
	Name        string    `gorm:"column:name;type:VARCHAR(255)"`
	Address     string    `gorm:"column:address;type:VARCHAR(255);NOT NULL;DEFAULT:'';unique"`
	SeedMessage string    `gorm:"column:seed_message;type:VARCHAR(255);NOT NULL;DEFAULT:''"` // will store secret code only to save storage
	Token       int       `gorm:"column:token;type:int;NOT NULL;DEFAULT:1000"`
	CreatedAt   time.Time `gorm:"column:created_at;<-:false;type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updated_at;<-:false;type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func CreateUser(user *User) error {
	isSuccess := Db.Table(userTable).Create(user).Error
	if isSuccess != nil {
		return isSuccess
	}
	return nil
}

func GetUser(address string) (*User, bool, error) {
	user := new(User)
	if err := Db.Table(userTable).Where("address = ?", address).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return user, true, nil
}

func GetUsers(addresses []string) ([]*User, bool, error) {
	var users []*User
	if err := Db.Where("address IN ?", addresses).Find(&users).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return users, true, nil
}

func GetUserBalance(user_id int) (*User, bool, error) {
	user := new(User)
	if err := Db.Table(userTable).Where("user_id = ?", user_id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return user, true, nil
}

func UpdateUserBalance(user_id int, count int) error {
	err := Db.Table(userTable).Where("user_id = ?", user_id).Updates(map[string]interface{}{
		"token": count,
	})
	return err.Error
}

func UpdateUserAddress(user_id int, address string) error {
	err := Db.Table(userTable).Where("user_id = ?", user_id).Updates(map[string]interface{}{
		"address": address,
	})
	return err.Error
}

func IsUserExist(user_id int) bool {
	var user User
	Db.Table(userTable).Select("id").Where("user_id = ?", user_id).First(&user)
	fmt.Println("user", user)
	if user.ID > 0 {
		return true
	}
	return false
}

//func SetKeyInfo(address string, publicKey, keyStore null.String) error {
//	isSuccess := Db.Model(&User{}).Where("address = ?", address).Updates(map[string]interface{}{
//		"public_key": publicKey,
//		"key_store":  keyStore,
//	}).Error
//	if isSuccess != nil {
//		return isSuccess
//	}
//	return nil
//}

func UpdateUserCount(address string, count int) error {

	isSuccess := Db.Model(&User{}).Where("address = ?", address).Updates(map[string]interface{}{
		"count": count,
	}).Error
	if isSuccess != nil {
		return isSuccess
	}
	return nil
}

func (u *User) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("CreatedAt", time.Now())
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

func (u *User) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}
