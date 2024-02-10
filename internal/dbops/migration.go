package dbops

import (
	"genericAPI/api/database_connection"
	"genericAPI/api/environment"
	"genericAPI/internal/models/permission"
	"genericAPI/internal/models/permission_group"
	"genericAPI/internal/models/permission_group_permissions"
	"genericAPI/internal/models/refresh_token"
	"genericAPI/internal/models/user"
	"genericAPI/internal/utils/authentication_utils"
	"gorm.io/datatypes"
	"log"
	"time"
)

func Migrate() {
	if environment.IsTestEnvironment() && environment.AllowMigrations() {
		log.Print("Auto migration enabled. Migrating...")
		err := database_connection.DB.AutoMigrate(
			&user.User{},
			&refresh_token.RefreshToken{},
			&permission.Permission{},
			&permission_group.PermissionGroup{},
			&permission_group_permissions.PermissionGroupPermission{},
		)
		if err != nil {
			panic("Failed migrating. Reason: " + err.Error())
		}
		createGenericPermissionGroup()
		createGenericAdminUser()
		log.Print("Successfully applied database schema.")
	}
}

func createGenericPermissionGroup() {
	permGroup := permission_group.PermissionGroup{
		ID:          1,
		Name:        "Admin",
		Permissions: []permission.Permission{},
	}
	permGroup.Save()
}

func createGenericAdminUser() {
	refreshToken := refresh_token.RefreshToken{
		Token:  authentication_utils.CreateRefreshToken(1),
		UserID: 1,
	}
	adminUser := user.User{
		ID:                1,
		Mail:              "admin@test.com",
		Name:              "admin",
		Surname:           "test",
		Username:          "admin",
		PublicID:          authentication_utils.CreatePublicID(),
		DisplayName:       "admin",
		Password:          "pbkdf2_sha256$10000$xAAC25Wq116vHuOP$O9zdX6KxOL/F0aObvlT9/LdUUQRWY/CXnzpIgq/5dis=", // 1234
		IsVerified:        true,
		PermissionGroupID: 1,
		CreatedAt:         datatypes.Date(time.Now()),
	}
	adminUser.Save()
	refreshToken.Save()
	log.Printf("Created a generic admin user with username:admin password:1234 refresh_token:%s", refreshToken.Token)
}
