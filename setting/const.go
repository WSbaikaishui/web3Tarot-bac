package setting

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

const (
	ItemChoose              string = "STATE:ITEM_CHOOSE"
	CARD                    string = "STATE:CARD"
	QUESION                 string = "STATE:QUESTION"
	WaitReturn              string = "STATE:WAITRETURN"
	AddressBind             string = "STATE:ADDRESS_BIND"
	TxSubmit                string = "STATE:TX_SUBMIT"
	StatsCallbackDataPrefix string = "ukey"
)

var ItemsKeyboard = [][]tgbotapi.InlineKeyboardButton{
	{
		tgbotapi.NewInlineKeyboardButtonData("Recharge", "purchase"),
		tgbotapi.NewInlineKeyboardButtonData("Balance", "balance"),
	},
	{
		tgbotapi.NewInlineKeyboardButtonData("Divination", "divine"),
	},
	{
		//tgbotapi.NewInlineKeyboardButtonData("history", "history"),
		tgbotapi.NewInlineKeyboardButtonData("Recharge list", "chargeList"),
	},
}

var ReturnButtion = [][]tgbotapi.InlineKeyboardButton{
	{tgbotapi.NewInlineKeyboardButtonData("< Return", "return")},
}

var CardUp = [][]tgbotapi.InlineKeyboardButton{
	{tgbotapi.NewInlineKeyboardButtonData("card1", "card1"),
		tgbotapi.NewInlineKeyboardButtonData("card2", "card3"),
		tgbotapi.NewInlineKeyboardButtonData("card2", "card3")},
}

var CardMarkup = tgbotapi.ReplyKeyboardMarkup{
	Keyboard: [][]tgbotapi.KeyboardButton{
		{
			{"card1", false, false},
			{"card2", false, false},
			{"card3", false, false},
		},
	},
}

var RemoveKeyboard = tgbotapi.NewRemoveKeyboard(false)
