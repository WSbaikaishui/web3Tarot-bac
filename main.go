package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
	"web3Tarot-backend/callbacks"
	"web3Tarot-backend/models"
	"web3Tarot-backend/pkg/logging"
	"web3Tarot-backend/setting"
)

//func main() {
//	setting.Setup()
//	models.Setup()
//	logging.Setup()
//
//	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
//	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
//	endless.DefaultMaxHeaderBytes = 1 << 20
//	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
//
//	server := endless.NewServer(endPoint, routers.InitRouter())
//	server.BeforeBegin = func(add string) {
//		logging.Info("Actual pid is %d", syscall.Getpid())
//	}
//
//	err := server.ListenAndServe()
//	if err != nil {
//		logging.Error("Server err: %v", err)
//	}
//}

func main() {
	// entry and initialization point
	setting.Setup()
	models.Setup()
	models.Init()
	logging.Setup()

	api, err := tgbotapi.NewBotAPI(setting.AppSetting.TgApi)

	if err != nil {
		logging.Error(err)
	}

	api.Debug = true

	logging.Info("Running @%s[%d]", api.Self.UserName, api.Self.ID)

	HandleUpdateBase(*api) // blocking func
}

// all "Messages" handler
func MainMessageProcessor(api tgbotapi.BotAPI, update tgbotapi.Update) {
	state := models.GetCurrentState(update)

	switch {

	// handle commands
	case state == models.InitialState: // empty state not initialized yet

		switch update.Message.Text {
		case "/start":
			go callbacks.CmdStart(api, update, 0)
		case "/chargeList":
			go callbacks.CmdChargeList(api, update, 0)
		case "/history":
			go callbacks.CmdHistory(api, update, 0)
		case "/balance":
			go callbacks.CmdBalance(api, update, 0)
		}

	// handle other text updates
	case state == setting.QUESION:
		go callbacks.ProcessQuestion(api, update)
	case state == setting.CARD:
		go callbacks.ProcessCard(api, update)
	case state == setting.AddressBind:
		go callbacks.ProcessAddressBind(api, update)
	case state == setting.TxSubmit:
		go callbacks.ProcessTxSubmit(api, update)
	}
}

// callback queries handler
func MainCallbackQueryProcessor(api tgbotapi.BotAPI, update tgbotapi.Update) {
	state := models.GetCurrentState(update)

	switch {
	case state == setting.ItemChoose:
		go callbacks.ClickOnItem(api, update)
	case state == setting.WaitReturn || state == setting.TxSubmit || state == setting.AddressBind:
		go callbacks.CmdStart(api, update, 1)

	// check if Data has prefix ukey:
	case state == models.InitialState && strings.HasPrefix(
		update.CallbackQuery.Data,
		setting.StatsCallbackDataPrefix+":",
	):
		go callbacks.ProcessAbout(api, update)
	}
}

// updates entry point
func HandleUpdateBase(api tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := api.GetUpdatesChan(u)

	if err != nil {
		panic(err)
	}

	for update := range updates {
		switch {
		case update.Message != nil:
			go MainMessageProcessor(api, update)
		case update.CallbackQuery != nil:
			go MainCallbackQueryProcessor(api, update)
		}
	}
}
