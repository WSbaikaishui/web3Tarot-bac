package transaction

import (
	"fmt"
	"math"
	apiErr "web3Tarot-backend/errors"
	"web3Tarot-backend/log"
	"web3Tarot-backend/models"
	"web3Tarot-backend/service/user"
	"web3Tarot-backend/setting"
	"web3Tarot-backend/util"
)

func IsTransactionExist(txHash string) bool {
	return models.IsTransactionExist(txHash)
}

func GetTransactions(user_id int) []string {
	data, err := user.GetUser(user_id)
	if err != nil {
		return []string{}
	}
	transaction := models.GetTransactionByAddress(data.Address)
	txList := make([]string, len(transaction))
	for i, _ := range transaction {
		txList[i] = transaction[i].TxHash
	}
	return txList

}

func ReCharge(data ChargeData) (int, error) {
	user, ok, err := models.GetUser(data.Address)
	if err != nil {
		log.Errorf("find user err: %v", err)
		return 0, err
	}
	if !ok {
		return 0, apiErr.ErrNotFound("User not found")
	}
	transaction := models.Transaction{
		Address: data.Address,
		TxHash:  data.TxHash,
		Amount:  data.Amount,
	}
	token, err := util.FetchPrice(setting.AppSetting.OkxApi, setting.AppSetting.JsonPath)
	finalToken := float64(data.Amount) / math.Pow10(5) * (float64(token) / math.Pow10(6))
	user.Token = user.Token + int(finalToken)
	fmt.Println("----------------", finalToken, user.Token)
	err = models.CreateTransaction(&transaction, user.UserId, user.Token)
	if err != nil {
		return 0, err
	}
	return user.Token, nil
}
