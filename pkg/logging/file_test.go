package logging

import (
	"fmt"
	"testing"
	"time"
	"web3Tarot-backend/setting"
)

func TestOpenLogFile(t *testing.T) {
	// openLogFile(fileName, filePath)
	setting.Setup()
	fileName := "log-" + time.Now().Format(time.DateTime) + ".log"
	filePath := getLogFilePath()
	file, err := OpenLogFile(fileName, filePath)
	if err != nil {
		t.Errorf("openLogFile err: %v", err)
	}
	fmt.Println(file)
}
