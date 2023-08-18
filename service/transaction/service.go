package transaction

import (
	apiErr "web3Tarot-backend/errors"
	"web3Tarot-backend/log"
	"web3Tarot-backend/models"
	"web3Tarot-backend/setting"
)

func IsTransactionExist(txHash string) bool {
	return models.IsTransactionExist(txHash)
}

func ReCharge(data ChargeData) error {
	user, ok, err := models.GetUser(data.Address)
	if err != nil {
		log.Errorf("find user err: %v", err)
		return err
	}
	if !ok {
		return apiErr.ErrNotFound("User not found")
	}
	transaction := models.Transaction{
		Address: data.Address,
		TxHash:  data.TxHash,
		Amount:  data.Amount,
	}
	user.Count = user.Count + transaction.Amount/setting.AppSetting.GasPerTarot
	err = models.CreateTransaction(&transaction, user.Count)
	if err != nil {
		return err
	}
	return nil
}
