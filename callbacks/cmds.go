// handler for slash-commands

// btw, that's a bad idea to store text messages in code

package callbacks

import (
	"fmt"
	"log"
	"math"
	"web3Tarot-backend/models"
	v1 "web3Tarot-backend/routers/api/v1"
	"web3Tarot-backend/setting"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// process /start
func CmdStart(api tgbotapi.BotAPI, update tgbotapi.Update, method int) {
	if !v1.IsUserExist(update, method) {
		v1.Login(update, method)
	}
	if method == 1 {
		action := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		_, _ = api.AnswerCallbackQuery(action)
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.MessageID,
				ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: setting.ItemsKeyboard,
				},
			},
			Text:      HtmlFmt("This is a tarot channel, what do you want to do ?", "b"),
			ParseMode: "html",
		}
		_, _ = api.Send(edit)
	} else {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			HtmlFmt("This is a tarot channel, what do you want to do ?", "b"))
		msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: setting.ItemsKeyboard,
		}
		msg.ParseMode = "html"
		_, err := api.Send(msg)
		if err != nil {
			log.Panic(err)
		}
	}
	models.UpdateState(update, setting.ItemChoose)
}

func CmdChargeList(api tgbotapi.BotAPI, update tgbotapi.Update, method int) {
	if !v1.IsUserExist(update, method) {
		v1.Login(update, method)
	}
	txs := v1.GetTransactions(update, method)
	Text := ""
	if len(txs) == 0 {
		Text = HtmlFmt("Oops! No recharge...", "b")
	} else {
		str := ""
		for _, val := range txs {
			str += val + "\n"
		}
		Text = HtmlFmt(fmt.Sprintf("This shows the last 10 recharge records: \n%s", str), "b")

	}

	if method == 0 {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			Text)
		msg.ParseMode = "html"
		_, _ = api.Send(msg)
	} else {
		action := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		_, _ = api.AnswerCallbackQuery(action)
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.MessageID,
				ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: setting.ReturnButtion,
				},
			},
			Text:      Text,
			ParseMode: "html",
		}
		_, _ = api.Send(edit)
	}

	models.UpdateState(update, setting.WaitReturn)
}

func CmdPurchase(api tgbotapi.BotAPI, update tgbotapi.Update, method int) {
	if !v1.IsUserExist(update, method) {
		v1.Login(update, method)
	}
	address := v1.GetUserAddress(update, method)
	if len(address) == 0 {
		if method == 0 {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				HtmlFmt(fmt.Sprintf("Please bind your ton address"), "b"))
			msg.ParseMode = "html"
			_, _ = api.Send(msg)
		} else {
			action := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			_, _ = api.AnswerCallbackQuery(action)
			edit := tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:    update.CallbackQuery.Message.Chat.ID,
					MessageID: update.CallbackQuery.Message.MessageID,
					ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
						InlineKeyboard: setting.ReturnButtion,
					},
				},
				Text:      HtmlFmt(fmt.Sprintf("Please bind your ton address"), "b"),
				ParseMode: "html",
			}
			_, _ = api.Send(edit)
		}
		models.UpdateState(update, setting.AddressBind)
	} else {
		if method == 0 {
			msg := tgbotapi.NewMessage(
				update.Message.Chat.ID,
				HtmlFmt(fmt.Sprintf("Please submit your pay txhash"), "b"))
			msg.ParseMode = "html"
			_, _ = api.Send(msg)
		} else {
			action := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			_, _ = api.AnswerCallbackQuery(action)
			edit := tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					ChatID:    update.CallbackQuery.Message.Chat.ID,
					MessageID: update.CallbackQuery.Message.MessageID,
					ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
						InlineKeyboard: setting.ReturnButtion,
					},
				},
				Text:      HtmlFmt(fmt.Sprintf("Please submit your pay txhash"), "b"),
				ParseMode: "html",
			}
			_, _ = api.Send(edit)
		}
		models.UpdateState(update, setting.TxSubmit)
	}

}

func CmdBalance(api tgbotapi.BotAPI, update tgbotapi.Update, method int) {
	fmt.Println(update.CallbackQuery.From.ID)
	if !v1.IsUserExist(update, method) {
		v1.Login(update, method)
	}
	token := v1.GetUserBalance(update, method)
	balance := float64(token) / (math.Pow10(6))
	if method == 0 {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			HtmlFmt(fmt.Sprintf("Your balance is : %v $", balance), "b"))
		msg.ParseMode = "html"
		_, _ = api.Send(msg)
	} else {
		action := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		_, _ = api.AnswerCallbackQuery(action)
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.MessageID,
				ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: setting.ReturnButtion,
				},
			},
			Text:      HtmlFmt(fmt.Sprintf("Your balance is : %v $", balance), "b"),
			ParseMode: "html",
		}
		_, _ = api.Send(edit)
	}
	models.UpdateState(update, setting.WaitReturn)
}

func CmdHistory(api tgbotapi.BotAPI, update tgbotapi.Update, method int) {
	if !v1.IsUserExist(update, method) {
		v1.Login(update, method)
	}
	txs := v1.GetTransactions(update, method)
	Text := ""
	if len(txs) == 0 {
		Text = HtmlFmt("Oops! No divinations about ...", "b")
	} else {
		str := ""
		for _, val := range txs {
			str += val + "\n"
		}
		Text = HtmlFmt(fmt.Sprintf("This shows the last 10 divinations: \n %s", str), "b")

	}
	if method == 0 {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			Text)
		msg.ParseMode = "html"
		_, _ = api.Send(msg)
	} else {
		action := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		_, _ = api.AnswerCallbackQuery(action)
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.MessageID,
				ReplyMarkup: &tgbotapi.InlineKeyboardMarkup{
					InlineKeyboard: setting.ReturnButtion,
				},
			},
			Text:      Text,
			ParseMode: "html",
		}
		_, _ = api.Send(edit)
	}

	models.UpdateState(update, setting.WaitReturn)
}
