package database_connection

import (
	"genericAPI/api/api_config"
	"genericAPI/api/database_logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

func getNow() time.Time {
	return time.Now().UTC()
}

func InitDB(logger *database_logger.DBLogger) {
	dbConf := gorm.Config{
		Logger:                   logger,
		SkipDefaultTransaction:   true,
		NowFunc:                  getNow,
		PrepareStmt:              false,
		DisableNestedTransaction: false,
		AllowGlobalUpdate:        true,
		QueryFields:              true,
		Dialector:                nil,
	}

	var err error
	DB, err = gorm.Open(mysql.Open(api_config.Config.DB.GetConnectionStr()), &dbConf)

	if err != nil {
		panic("Error while connecting database. Reason: " + err.Error())
	}

	log.Print("Successfully initialized database connection")
}
