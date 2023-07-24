package user

import (
	"fmt"
	apiErr "github.com/WSbaikaishui/web3Tarot-backend/errors"
	"github.com/WSbaikaishui/web3Tarot-backend/model"
	"gopkg.in/guregu/null.v4"
)

const NeoSignatureLength = 128

type LoginParam struct {
	Address   string `json:"address"`
	Nonce     string `json:"nonce" form:"nonce"`
	Signature string `json:"signature" form:"signature"`
}

func (param *LoginParam) Validate() error {
	if len(param.Address) == 0 {
		return apiErr.ErrInvalidParameter("invalid address")
	}
	if len(param.Nonce) == 0 {
		return apiErr.ErrInvalidParameter("invalid nonce")
	}
	if len(param.Signature) != NeoSignatureLength {
		return apiErr.ErrInvalidParameter("invalid signature length")
	}
	return nil
}

type LoginData struct {
	Token string
	IsNew bool
}

type GetUserData struct {
	Id          uint        `json:"id"`
	Name        null.String `json:"name"`
	Address     string      `json:"address"`
	SeedMessage string      `json:"seedMessage"`
	PublicKey   null.String `json:"publicKey"`
	KeyStore    null.String `json:"keyStore"`
	IsOnline    bool        `json:"isOnline"`
}

func (u *GetUserData) FromModel(user *model.User) {
	u.Id = user.ID
	u.Name = user.Name
	u.Address = user.Address
	u.SeedMessage = fmt.Sprintf(seedMessageTpl, user.Address, user.SeedMessage)
	u.PublicKey = user.PublicKey
	u.KeyStore = user.KeyStore
	u.IsOnline = user.IsOnline
}

type GetPublicUserParam struct {
	Addresses []string `json:"addresses"`
}

type GetPublicUserData struct {
	PublicUsers []PublicUser
}

type PublicUser struct {
	Id        uint        `json:"id"`
	Name      null.String `json:"name"`
	Address   string      `json:"address"`
	PublicKey null.String `json:"publicKey"`
	IsOnline  bool        `json:"isOnline"`
}

func (pu *PublicUser) FromModel(user *model.User) {
	pu.Id = user.ID
	pu.Name = user.Name
	pu.Address = user.Address
	pu.PublicKey = user.PublicKey
	pu.IsOnline = user.IsOnline
}

type SetKeyInfoParam struct {
	Address   string `json:"address"`
	PublicKey string `json:"publicKey"`
	KeyStore  string `json:"keyStore"`
}

func (param *SetKeyInfoParam) Validate() error {
	if len(param.Address) == 0 {
		return apiErr.ErrInvalidParameter("invalid address")
	}
	if len(param.PublicKey) == 0 {
		return apiErr.ErrInvalidParameter("invalid public key")
	}
	if len(param.KeyStore) == 0 {
		return apiErr.ErrInvalidParameter("invalid key store")
	}
	return nil
}
