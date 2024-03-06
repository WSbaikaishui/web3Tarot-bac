package v1

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	apiErr "web3Tarot-backend/errors"
	"web3Tarot-backend/pkg/logging"
	"web3Tarot-backend/service/transaction"
	"web3Tarot-backend/util"
)

func Recharge(address string, txhash string) (int, error) {
	//wa := ctx.Value(util.AuthKey).(string)

	var param transaction.ChargeData
	param.Address = address
	param.TxHash = txhash

	if transaction.IsTransactionExist(param.TxHash) {
		return 0, errors.New(util.EncodeError(apiErr.ErrInvalidParameter("transaction already used")))
	}

	amount, err := CheckTransaction(&param)
	if err != nil {
		logging.Error(util.EncodeError(err))
		return 0, errors.New(util.EncodeError(err))
	}
	param.Amount = amount
	token, err := transaction.ReCharge(param)
	if err != nil {
		logging.Error(util.EncodeError(err))
		return 0, errors.New(util.EncodeError(err))
	}
	return token, nil
}

func CheckTransaction(param *transaction.ChargeData) (int, error) {
	//url := setting.AppSetting.TonTestNetApi + "getTransactions?address=" + setting.AppSetting.TonAddress + "&limit=30&archival=true&api_key=" + setting.AppSetting.TonToken
	amount, err := util.GetTxAmount(param.Address, param.TxHash)
	if err != nil {
		return 0, err
	}
	//TODO  处理登录的逻辑
	return amount, nil

}

func GetTransactions(update tgbotapi.Update, method int) []string {

	if method == 0 {
		return transaction.GetTransactions(update.Message.From.ID)
	}
	return transaction.GetTransactions(update.CallbackQuery.From.ID)
}
