package v1

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"web3Tarot-backend/models"
	"web3Tarot-backend/pkg/logging"
	"web3Tarot-backend/service/divination"
	"web3Tarot-backend/service/transaction"
	"web3Tarot-backend/service/user"
	"web3Tarot-backend/util"
)

func CreateDivination(ID int, cardNumber int, question string, answer string) {
	var param divination.CreateDivination
	param.UserID = ID
	param.Card1 = cardNumber
	param.Card1Status = true
	param.Question = question
	param.Content = answer
	err := divination.SetDivinationMessage(&param)
	if err != nil {
		logging.Error("create fail")
	}
}

func IsEnoughBalance(ID int, question string, answer string) bool {
	user, err := user.GetUser(ID)
	if err != nil {
		return false
	}
	count, _ := util.TokenCaculate(question)
	ans, _ := util.TokenCaculate(answer)
	if user.Token > count+ans {
		user.Token -= count + ans

		err := models.UpdateUserBalance(ID, user.Token)
		if err != nil {
			return false
		}
		return true
	}
	return false
}

func GetDivinations(update tgbotapi.Update, method int) []string {

	if method == 0 {
		return transaction.GetTransactions(update.Message.From.ID)
	}
	return transaction.GetTransactions(update.CallbackQuery.From.ID)
}
