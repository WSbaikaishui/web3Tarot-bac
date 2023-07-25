package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

// // Migrate do migration with versioning style instead of auto style.
//
//	func Migrate(db *gorrm.DB) error {
//		m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{})
//		m.InitSchema(initSchema)
//		return m.Migrate()
//	}
func initSchema(db *gorm.DB) error {
	if err := db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin").AutoMigrate(
		&Divination{},
		&Tarot{},
		&Nonce{},
		&User{},
	); err != nil {
		fmt.Println("migrate tables error", err)
		return nil
	}
	return nil
}
