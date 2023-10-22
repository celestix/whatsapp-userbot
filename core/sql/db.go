package sql

import (
	sqlOrig "database/sql"

	waLogger "github.com/celestix/whatsapp-userbot/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var SESSION *gorm.DB

const DEFAULE_USER_ID = 777000

func LoadDB(LOGGER *waLogger.Logger) *sqlOrig.DB {
	LOGGER = LOGGER.Create("database")
	db, err := gorm.Open(sqlite.Open("waub.db"), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		LOGGER.Println("failed to start db:", err.Error())
	}
	SESSION = db

	dB, _ := db.DB()
	dB.SetMaxOpenConns(100)

	LOGGER.Println("Database connected")

	// Create tables if they don't exist
	_ = SESSION.AutoMigrate(&Afk{}, &Note{}, &Filter{}, &ChatSettings{})
	LOGGER.Println("Auto-migrated database schema")
	return dB
}
