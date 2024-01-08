package api

import (
	"genericAPI/api/api_config"
	"genericAPI/api/environment"
	"genericAPI/internal/models/auth_token"
	"genericAPI/internal/models/permission"
	"genericAPI/internal/models/permission_group"
	"genericAPI/internal/models/permission_group_permissions"
	"genericAPI/internal/models/token_type"
	"genericAPI/internal/models/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var DB *gorm.DB

func getNow() time.Time {
	return time.Now().UTC()
}

func InitDB() {
	dbConf := gorm.Config{
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
	if environment.IsTestEnvironment() {
		err := DB.AutoMigrate(
			&user.User{},
			&token_type.TokenType{},
			&auth_token.AuthToken{},
			&permission.Permission{},
			&permission_group.PermissionGroup{},
			&permission_group_permissions.PermissionGroupPermission{},
		)
		if err != nil {
			panic("Failed migrating. Reason: " + err.Error())
		}
	}
}
