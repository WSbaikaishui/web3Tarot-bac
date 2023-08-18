package user

import (
	"context"
	"fmt"
	"github.com/joeqian10/neo3-gogogo/crypto"
	"github.com/joeqian10/neo3-gogogo/helper"
	"gopkg.in/guregu/null.v4"
	"time"
	apiErr "web3Tarot-backend/errors"
	"web3Tarot-backend/log"
	"web3Tarot-backend/models"
	"web3Tarot-backend/util"
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

//type Service struct {
//	userDB  *models.UserDB
//	nonceDB *models.NonceDB
//	chatDB  *model.ChatDB
//	msgDB   *model.MessageDB
//}

//type Option func(service *Service) error

//func New(userDB *model.UserDB, nonceDB *model.NonceDB, chatDB *model.ChatDB, msgDB *model.MessageDB, opts ...Option) (*Service, error) {
//	srv := &Service{
//		userDB:  userDB,
//		nonceDB: nonceDB,
//		chatDB:  chatDB,
//		msgDB:   msgDB,
//	}
//	for _, opt := range opts {
//		if err := opt(srv); err != nil {
//			return nil, err
//		}
//	}
//	return srv, nil
//}

func Login(ctx context.Context, param *LoginParam) (*LoginData, error) {
	// get nonce from db
	nonce, ok, err := models.GetNonce(ctx, param.Nonce)
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
	//msgLength := helper.VarIntFromInt(len(msg) / 2)
	//serializedMsg := helper.HexToBytes("010001f0" + hex.EncodeToString(msgLength.Bytes()) + msg + "0000")
	//
	//sig := helper.HexToBytes(param.Signature)
	pubkeyEth, err := util.VerifyMessage(msg, param.Signature)
	if err != nil {
		log.Errorf("verify message failed, err: %v", err)
		return nil, apiErr.ErrInvalidSignature("invalid signature")
	}
	if param.Address != pubkeyEth {
		return nil, apiErr.ErrInvalidSignature("recovered address does not match")
	}
	//pubKeys, err := util.RecoverPubKeyFromSigOnSecp256r1(serializedMsg, sig)
	//if err != nil {
	//	log.Errorf("verify signature failed, err: %v", err)
	//	return nil, apiErr.ErrInvalidSignature("invalid signature")
	//}
	//if !util.VerifyAddress(param.Address, pubKeys) {
	//	return nil, apiErr.ErrInvalidSignature("recovered address does not match")
	//}
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
	_, ok, err = models.GetUser(param.Address)
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
		user := models.User{
			Address:     param.Address,
			SeedMessage: secretCode,
			PublicKey:   null.String{},
			KeyStore:    null.String{},
		}
		// add user to db
		if err := models.CreateUser(&user); err != nil {
			log.Errorf("create user failed, err: %v", err)
			return nil, err
		}
		data.IsNew = true
	}

	// delete nonce
	if err := models.DeleteNonce(ctx, nonce); err != nil {
		log.Errorf("delete nonce failed, err: %v", err)
	}
	return data, nil
}

func GetUser(ctx context.Context, address string) (*GetUserData, error) {
	wa := ctx.Value(util.AuthKey).(string)
	if wa != address {
		return nil, apiErr.ErrForbidden("wallet address mismatch")
	}
	data := new(GetUserData)
	user, ok, err := models.GetUser(address)
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

func GetUserPublicInfo(ctx context.Context, addresses []string) ([]PublicUser, error) {
	users, ok, err := models.GetUsers(addresses)
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

//func SetKeyInfo(ctx context.Context, param *SetKeyInfoParam) error {
//	wa := ctx.Value(util.AuthKey).(string)
//	if wa != param.Address {
//		return apiErr.ErrForbidden("Wallet address mismatch")
//	}
//	// get user
//	user, ok, err := models.GetUser(param.Address)
//	if err != nil {
//		log.Errorf("find user err: %v", err)
//		return err
//	}
//	if !ok {
//		return apiErr.ErrNotFound("User not found")
//	}
//
//	// only can set the key info once
//	if !user.PublicKey.IsZero() || !user.KeyStore.IsZero() {
//		return apiErr.ErrForbidden("User key info already set")
//	}
//	publicKey := null.StringFrom(param.PublicKey)
//	keyStore := null.StringFrom(param.KeyStore)
//	return models.SetKeyInfo(param.Address, publicKey, keyStore)
//}

//func SendWelcome(ctx context.Context, address string) {
//	// wait for ws conn?
//	time.Sleep(time.Second)
//
//	// create bot welcome msg
//	if config.AppCfg.ChatGPT.BotAddress != "" && address != config.AppCfg.ChatGPT.BotAddress {
//		sendWelcome(ctx, config.AppCfg.ChatGPT.BotAddress, address, svc.msgDB, svc.chatDB)
//		// create manager msg
//		if config.AppCfg.ChatGPT.ManagerAddress != "" && address != config.AppCfg.ChatGPT.ManagerAddress {
//			sendWelcome(ctx, config.AppCfg.ChatGPT.ManagerAddress, address, svc.msgDB, svc.chatDB)
//		}
//	}
//}

//func sendWelcome(ctx context.Context, sender, receiver string, msgDB *model.MessageDB, chatDB *model.ChatDB) {
//	msg := new(models.Message)
//	{
//		msg.Uuid = util.UuidV4()
//		msg.Sender = sender
//		msg.Receiver = receiver
//		msg.Content = message.JsonStringify(message.NewTextContent("Welcome to NeoChat!"))
//		msg.IsEncrypted = false
//		msg.Time = util.FormatUtcNow()
//	}
//	err := msgDB.Create(ctx, msg)
//	if err != nil {
//		log.Errorf("create message err: %v", err)
//		return
//	}
//
//	// update chat last msg
//	senderChat := new(model.Chat)
//	rcvChat := new(model.Chat)
//	err = message.UpdateChatLastMsg(ctx, chatDB, senderChat, rcvChat, msg)
//	if err != nil {
//		log.Errorf(fmt.Sprintf("UpdateChatLastMsg err: %v", err))
//		return
//	}
//
//	rawMsg := notification.MarshalNewMessageEvent(notification.CreateNewMessageEventFromModel(msg))
//	notification.SendMsg(receiver, rawMsg)
//}
