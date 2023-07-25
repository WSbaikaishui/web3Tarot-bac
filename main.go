package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"syscall"
	"web3Tarot-backend/models"
	"web3Tarot-backend/pkg/logging"
	"web3Tarot-backend/routers"
	"web3Tarot-backend/setting"
)

func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()

	endless.DefaultReadTimeOut = setting.ServerSetting.ReadTimeout
	endless.DefaultWriteTimeOut = setting.ServerSetting.WriteTimeout
	endless.DefaultMaxHeaderBytes = 1 << 20
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		logging.Info("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		logging.Error("Server err: %v", err)
	}
}
