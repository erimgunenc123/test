package dbops

import (
	"genericAPI/api"
	"genericAPI/internal/models/user"
	"time"
)

func CreateUser(
	username string,
	password string,
	name string,
	surname string,
	displayName string) (*user.User, error) {
	newUser := user.User{
		Name:        name,
		Surname:     surname,
		Username:    username,
		DisplayName: displayName,
		Password:    password,
		CreatedAt:   time.Now(),
	}

	result := api.DB.Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}

	return &newUser, nil
}

func GetUserByUsername(username string, attachPermissions bool) (user *user.User, err error) {
	err = api.DB.Table("user").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	if attachPermissions {
		result := api.DB.Preload("PermissionGroup").First(&user, "id = ?", user.ID)
		if result.Error != nil {
			return nil, err
		}
	}
	return
}
