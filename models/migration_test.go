package models

import (
	"fmt"
	"web3Tarot-backend/setting"

	//"web3Tarot-backend/config"
	//"web3Tarot-backend/models"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var (
	envKey    = "TEST_MYSQL_DSN"
	envValue1 = "root:rwRkh5AAbJmDoKMLAS4g@tcp(127.0.0.1:3306)/nchat?charset=utf8mb4&parseTime=True&loc=Local"
	envValue2 = "root:rwRkh5AAbJmDoKMLAS4g@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
)

func TestMigration(t *testing.T) {
	assert := require.New(t)
	err := os.Setenv(envKey, envValue1)
	if err != nil {
		fmt.Println(err.Error())
	}
	dsn := os.Getenv(envKey)
	if dsn == "" {
		t.Skip("set TEST_MYSQL_DSN to run this test")
	}
	setting.Setup()
	Setup()
	assert.Nil(err)

	assert.Nil(Db.Exec("DROP DATABASE IF EXISTS tarot;").Error)
	assert.Nil(Db.Exec("CREATE DATABASE tarot DEFAULT CHARACTER SET utf8mb4;").Error)
	assert.Nil(Db.Exec("USE tarot;").Error)
	assert.Nil(initSchema(Db))

	//if db.Migrator().HasTable("chats") {
	//	assert.Nil(db.Migrator().DropTable("chats"))
	//	assert.Nil(db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin").Migrator().CreateTable(&model.Chat{}))
	//}
	//if db.Migrator().HasTable("messages") {
	//	assert.Nil(db.Migrator().DropTable("messages"))
	//	assert.Nil(db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin").Migrator().CreateTable(&model.Message{}))
	//}
}
