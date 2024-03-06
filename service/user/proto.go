package user

import (
	apiErr "web3Tarot-backend/errors"
	"web3Tarot-backend/models"
)

const NeoSignatureLength = 128

type LoginParam struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"user_name"`
	FirstName string `json:"first_name"`
}

func (param *LoginParam) Validate() error {
	if param.UserID == 0 {
		return apiErr.ErrInvalidParameter("invalid userID")
	}
	if len(param.Name) == 0 {
		return apiErr.ErrInvalidParameter("invalid Name")
	}
	if len(param.FirstName) != 0 {
		return apiErr.ErrInvalidParameter("invalid firstName length")
	}
	return nil
}

type LoginData struct {
	Token string
	IsNew bool
}

type GetUserData struct {
	Id      uint   `json:"id"`
	UserID  int    `json:"user_id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Token   int    `json:"token              "`
}

func (u *GetUserData) FromModel(user *models.User) {
	u.Id = uint(user.ID)
	u.Name = user.Name
	u.Address = user.Address
	//u.SeedMessage = fmt.Sprintf(seedMessageTpl, user.Address, user.SeedMessage)
	u.Token = user.Token
	u.UserID = user.UserId
}

type GetPublicUserParam struct {
	Addresses []string `json:"addresses"`
}

//type GetPublicUserData struct {
//	PublicUsers []PublicUser
//}

//type PublicUser struct {
//	Id        uint        `json:"id"`
//	Name      null.String `json:"name"`
//	Address   string      `json:"address"`
//	PublicKey null.String `json:"publicKey"`
//	IsOnline  bool        `json:"isOnline"`
//}
//
//func (pu *PublicUser) FromModel(user *models.User) {
//	pu.Id = uint(user.ID)
//	pu.Name = user.Name
//	pu.Address = user.Address
//	pu.PublicKey = user.PublicKey
//}

//type SetKeyInfoParam struct {
//	Address   string `json:"address"`
//	PublicKey string `json:"publicKey"`
//	KeyStore  string `json:"keyStore"`
//}
//
//func (param *SetKeyInfoParam) Validate() error {
//	if len(param.Address) == 0 {
//		return apiErr.ErrInvalidParameter("invalid address")
//	}
//	if len(param.PublicKey) == 0 {
//		return apiErr.ErrInvalidParameter("invalid public key")
//	}
//	if len(param.KeyStore) == 0 {
//		return apiErr.ErrInvalidParameter("invalid key store")
//	}
//	return nil
//}
