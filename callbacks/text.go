package callbacks

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"html"
	"log"
	"math"
	"strconv"
	"strings"
	"web3Tarot-backend/models"
	v1 "web3Tarot-backend/routers/api/v1"
	"web3Tarot-backend/service/tarot"
	"web3Tarot-backend/setting"
	"web3Tarot-backend/util"
)

// process user's name
func ProcessQuestion(api tgbotapi.BotAPI, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("your question is:\n<b> %s</b>\nplease select a card your want:", html.EscapeString(update.Message.Text)))
	msg.ReplyMarkup = setting.CardMarkup
	msg.ParseMode = "html"

	_, err := api.Send(msg)

	if err != nil {
		log.Panic(err)
	}

	models.UpdateState(update, setting.CARD)
	models.UpdateData(update, map[string]interface{}{"question": update.Message.Text})
}

func ProcessAddressBind(api tgbotapi.BotAPI, update tgbotapi.Update) {
	isSuccess := v1.UpdateUserAddress(update.Message.From.ID, update.Message.Text)
	if isSuccess {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("bind address success！ \nYou can submit txhash which is transfer to EQDG5ALYstmhiMhpSN9j6tEHmLq2Owx7yyfTH1W15wS0NqZH\n to get Token"))
		msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: setting.ReturnButtion,
		}
		msg.ParseMode = "html"

		_, err := api.Send(msg)
		if err != nil {
			log.Panic(err)
		}
		models.UpdateState(update, setting.TxSubmit)
	} else {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("Binding failed, please try again"))
		msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: setting.ReturnButtion,
		}
		msg.ParseMode = "html"

		_, err := api.Send(msg)
		if err != nil {
			log.Panic(err)
		}
	}
}

func ProcessTxSubmit(api tgbotapi.BotAPI, update tgbotapi.Update) {
	address := v1.GetUserAddress(update, 0)

	amount, err := v1.Recharge(address, update.Message.Text)
	if err != nil {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("This transaction has error :\n %v", err))
		msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: setting.ReturnButtion,
		}
		msg.ParseMode = "html"

		_, err := api.Send(msg)
		if err != nil {
			log.Panic(err)
		}
	} else {
		msg := tgbotapi.NewMessage(
			update.Message.Chat.ID,
			fmt.Sprintf("Recharge successful! Your balance is :\n%v", float64(amount)/math.Pow10(6)))
		msg.ReplyMarkup = &tgbotapi.InlineKeyboardMarkup{
			InlineKeyboard: setting.ReturnButtion,
		}
		msg.ParseMode = "html"
		_, err := api.Send(msg)
		if err != nil {
			log.Panic(err)
		}
	}
}

// process click on card button
func ProcessCard(api tgbotapi.BotAPI, update tgbotapi.Update) {
	text := strings.ToLower(update.Message.Text)
	msg1 := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		HtmlFmt("The fortuneteller has received your message", "code"))

	msg1.ParseMode = "html"
	msg1.ReplyMarkup = setting.RemoveKeyboard
	_, err := api.Send(msg1)

	if err != nil {
		log.Panic(err)
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		HtmlFmt("Tarot master is divination...", "code"))

	msg.ParseMode = "html"
	message, err := api.Send(msg)

	if err != nil {
		log.Panic(err)
	}
	result, _ := strconv.Atoi(strings.TrimPrefix(text, "card"))

	cardNumber := randInt(1, 78, result)
	card, err := tarot.GetTarotCard(cardNumber)
	savedData := models.GetData(update)
	savedDataString := fmt.Sprintf(
		"Your Question is : \n   %s \nThe Card you select is : \n  %s\n",
		HtmlFmt(fmt.Sprintf("%v", savedData["question"]), "b"),
		HtmlFmt(card.CardName, "code"))
	reply, err := util.CreateChatGPTResponse(card.CardName, savedData["question"].(string))
	if err != nil {
		reply = "network error ，please try again"
	}
	if v1.IsEnoughBalance(update.Message.From.ID, savedData["question"].(string), reply) {
		v1.CreateDivination(update.Message.From.ID, cardNumber, savedData["question"].(string), reply)
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    message.Chat.ID,
				MessageID: message.MessageID,
			},
			Text: fmt.Sprintf(
				"%sThe fortune teller's answer is: \n%s",
				savedDataString,
				HtmlFmt(fmt.Sprintf("%s", reply), "code")),
			ParseMode: "html",
		}
		_, err = api.Send(edit)
		if err != nil {
			log.Panic(err)
		}
	} else {
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    message.Chat.ID,
				MessageID: message.MessageID,
			},
			Text:      HtmlFmt(fmt.Sprintf("your balance is not enough"), "b"),
			ParseMode: "html",
		}
		_, err = api.Send(edit)
		if err != nil {
			log.Panic(err)
		}
	}

	models.UpdateState(update, models.InitialState)
}
