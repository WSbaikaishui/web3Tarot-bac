package callbacks

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"web3Tarot-backend/models"
	"web3Tarot-backend/setting"
)

// handle Emoji clicks
// TODO  这里需要处理逻辑，就是选择的逻辑，选啥，然后对应的操作是啥，比如purchase，和history 就是不一样的
func ClickOnItem(api tgbotapi.BotAPI, update tgbotapi.Update) {

	switch update.CallbackQuery.Data {
	case "balance":
		CmdBalance(api, update, 1)

	case "history":
		CmdHistory(api, update, 1)

	case "chargeList":
		CmdChargeList(api, update, 1)

	case "purchase":
		CmdPurchase(api, update, 1)

	case "divine":
		text := fmt.Sprintf("This is  %s", update.CallbackQuery.Data)
		action := tgbotapi.NewCallback(update.CallbackQuery.ID, text)
		_, _ = api.AnswerCallbackQuery(action)

		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				ChatID:    update.CallbackQuery.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.MessageID,
			},
			Text:      HtmlFmt("Now please write a question your want to ask :", "b"),
			ParseMode: "html",
		}
		_, _ = api.Send(edit)
		models.UpdateState(update, setting.QUESION)
	}

}

// process after-/stats click
// 这个是处理 /stats的。TODO 我觉得可以改成balance的
func ProcessAbout(api tgbotapi.BotAPI, update tgbotapi.Update) {
	info := models.Get(strings.Replace(update.CallbackQuery.Data, "ukey:", "", 1), false)
	savedDataString := fmt.Sprintf(
		"Name: %s\nGender: %s\nFavorite item: %s\nAge ~ %s",
		fmt.Sprintf("%v", info["name"]),
		fmt.Sprintf("%v", info["gender"]),
		fmt.Sprintf("%v", info["fav_item"]),
		fmt.Sprintf("%v", info["age"]),
	)

	action := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, savedDataString)
	_, _ = api.AnswerCallbackQuery(action)
}
