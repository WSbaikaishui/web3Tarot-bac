package user

import (
	"context"
	"encoding/hex"
	"fmt"
	apiErr "github.com/WSbaikaishui/web3Tarot-backend/errors"
	"github.com/WSbaikaishui/web3Tarot-backend/log"
	"github.com/WSbaikaishui/web3Tarot-backend/model"
	"github.com/WSbaikaishui/web3Tarot-backend/service/message"
	"github.com/WSbaikaishui/web3Tarot-backend/service/notification"
	"github.com/WSbaikaishui/web3Tarot-backend/util"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"github.com/neo-ngd/nchat-backend/config"
	"gopkg.in/guregu/null.v4"
	"time"
)

const signatureTpl = `Welcome to NeoChat!

Signing is the only way we can truly know that you are the owner of the wallet you are connecting. Signing is a safe, gas-less transaction that does not in any way give NeoChat permission to perform any transactions with your wallet.

Wallet address: %s

Nonce: %s`

const seedMessageTpl = `Enable NeoChat end-to-end encryption.

This signature will be used to encrypt the end-to-end key pair, please do not disclose this signature elsewhere.

Wallet address: %s

Secret code: %s`

const secretCodeLen = 32

type Service struct {
	userDB  *model.UserDB
	nonceDB *model.NonceDB
	chatDB  *model.ChatDB
	msgDB   *model.MessageDB
}

type Option func(service *Service) error

func New(userDB *model.UserDB, nonceDB *model.NonceDB, chatDB *model.ChatDB, msgDB *model.MessageDB, opts ...Option) (*Service, error) {
	srv := &Service{
		userDB:  userDB,
		nonceDB: nonceDB,
		chatDB:  chatDB,
		msgDB:   msgDB,
	}
	for _, opt := range opts {
		if err := opt(srv); err != nil {
			return nil, err
		}
	}
	return srv, nil
}

func (svc *Service) Login(ctx context.Context, param *LoginParam) (*LoginData, error) {
	// get nonce from db
	nonce, ok, err := svc.nonceDB.GetNonce(ctx, param.Nonce)
	if err != nil {
		log.Errorf("get nonce error: %v", err)
		return nil, err
	}
	if !ok { // not found in db
		return nil, apiErr.ErrInvalidParameter("invalid nonce")
	}
	if time.Now().Unix() > int64(nonce.Expiration) {
		return nil, apiErr.ErrInvalidParameter("nonce expired")
	}

	// verify signature
	msg := helper.BytesToHex([]byte(fmt.Sprintf(signatureTpl, param.Address, param.Nonce)))
	msgLength := helper.VarIntFromInt(len(msg) / 2)
	serializedMsg := helper.HexToBytes("010001f0" + hex.EncodeToString(msgLength.Bytes()) + msg + "0000")

	sig := helper.HexToBytes(param.Signature)
	pubKeys, err := util.RecoverPubKeyFromSigOnSecp256r1(serializedMsg, sig)
	if err != nil {
		log.Errorf("verify signature failed, err: %v", err)
		return nil, apiErr.ErrInvalidSignature("invalid signature")
	}
	if !util.VerifyAddress(param.Address, pubKeys) {
		return nil, apiErr.ErrInvalidSignature("recovered address does not match")
	}
	// generate token
	hash, err := crypto.AddressToScriptHash(param.Address, helper.DefaultAddressVersion)
	if err != nil {
		return nil, apiErr.ErrInvalidParameter("invalid Address")
	}
	token, err := util.GenerateToken(hash.String(), time.Now().AddDate(0, 0, 7).Unix())
	if err != nil {
		log.Errorf("generate token failed, err: %v", err)
		return nil, err
	}
	// try to get user from db
	_, ok, err = svc.userDB.GetUser(ctx, param.Address)
	if err != nil {
		log.Errorf("get user failed, err: %v", err)
		return nil, err
	}
	data := &LoginData{
		Token: token,
		IsNew: false,
	}
	if !ok {
		// generate seed msg
		secretCode, err := util.RandomString(secretCodeLen)
		if err != nil {
			log.Errorf("generate random string error: %v", err)
			return nil, err
		}
		user := model.User{
			Address:     param.Address,
			SeedMessage: secretCode,
			PublicKey:   null.String{},
			KeyStore:    null.String{},
			IsOnline:    false,
		}
		// add user to db
		if err := svc.userDB.Create(ctx, &user); err != nil {
			log.Errorf("create user failed, err: %v", err)
			return nil, err
		}
		data.IsNew = true
	}

	// delete nonce
	if err := svc.nonceDB.DeleteNonce(ctx, nonce); err != nil {
		log.Errorf("delete nonce failed, err: %v", err)
	}
	return data, nil
}

func (svc *Service) GetUser(ctx context.Context, address string) (*GetUserData, error) {
	wa := ctx.Value(util.AuthKey).(string)
	if wa != address {
		return nil, apiErr.ErrForbidden("wallet address mismatch")
	}
	data := new(GetUserData)
	user, ok, err := svc.userDB.GetUser(ctx, address)
	if err != nil {
		log.Errorf("find user err: %v", err)
		return nil, err
	}
	if !ok {
		return nil, apiErr.ErrNotFound("user not found")
	}
	data.FromModel(user)
	return data, nil
}

func (svc *Service) GetUserPublicInfo(ctx context.Context, addresses []string) ([]PublicUser, error) {
	users, ok, err := svc.userDB.GetUsers(ctx, addresses)
	if err != nil {
		log.Errorf("find users err: %v", err)
		return nil, err
	}
	if !ok {
		return nil, apiErr.ErrNotFound("users not found")
	}
	publicUsers := make([]PublicUser, len(users))
	for i, _ := range users {
		publicUsers[i].FromModel(users[i])
	}

	return publicUsers, nil
}

func (svc *Service) SetKeyInfo(ctx context.Context, param *SetKeyInfoParam) error {
	wa := ctx.Value(util.AuthKey).(string)
	if wa != param.Address {
		return apiErr.ErrForbidden("Wallet address mismatch")
	}
	// get user
	user, ok, err := svc.userDB.GetUser(ctx, param.Address)
	if err != nil {
		log.Errorf("find user err: %v", err)
		return err
	}
	if !ok {
		return apiErr.ErrNotFound("User not found")
	}

	// only can set the key info once
	if !user.PublicKey.IsZero() || !user.KeyStore.IsZero() {
		return apiErr.ErrForbidden("User key info already set")
	}
	publicKey := null.StringFrom(param.PublicKey)
	keyStore := null.StringFrom(param.KeyStore)
	return svc.userDB.SetKeyInfo(ctx, param.Address, publicKey, keyStore)
}

func (svc *Service) SendWelcome(ctx context.Context, address string) {
	// wait for ws conn?
	time.Sleep(time.Second)

	// create bot welcome msg
	if config.AppCfg.ChatGPT.BotAddress != "" && address != config.AppCfg.ChatGPT.BotAddress {
		sendWelcome(ctx, config.AppCfg.ChatGPT.BotAddress, address, svc.msgDB, svc.chatDB)
		// create manager msg
		if config.AppCfg.ChatGPT.ManagerAddress != "" && address != config.AppCfg.ChatGPT.ManagerAddress {
			sendWelcome(ctx, config.AppCfg.ChatGPT.ManagerAddress, address, svc.msgDB, svc.chatDB)
		}
	}
}

func sendWelcome(ctx context.Context, sender, receiver string, msgDB *model.MessageDB, chatDB *model.ChatDB) {
	msg := new(model.Message)
	{
		msg.Uuid = util.UuidV4()
		msg.Sender = sender
		msg.Receiver = receiver
		msg.Content = message.JsonStringify(message.NewTextContent("Welcome to NeoChat!"))
		msg.IsEncrypted = false
		msg.Time = util.FormatUtcNow()
	}
	err := msgDB.Create(ctx, msg)
	if err != nil {
		log.Errorf("create message err: %v", err)
		return
	}

	// update chat last msg
	senderChat := new(model.Chat)
	rcvChat := new(model.Chat)
	err = message.UpdateChatLastMsg(ctx, chatDB, senderChat, rcvChat, msg)
	if err != nil {
		log.Errorf(fmt.Sprintf("UpdateChatLastMsg err: %v", err))
		return
	}

	rawMsg := notification.MarshalNewMessageEvent(notification.CreateNewMessageEventFromModel(msg))
	notification.SendMsg(receiver, rawMsg)
}
