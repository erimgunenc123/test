package user

import (
	"errors"
	"genericAPI/api/database_connection"
	"genericAPI/internal/models/permission_group"
	"genericAPI/internal/models/refresh_token"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	ID                uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt         datatypes.Date `gorm:"column:created_at"`
	Mail              string         `gorm:"column:mail"`
	Name              string
	Surname           string
	Username          string `gorm:"uniqueIndex;not null;size:32;column:username"`
	PublicID          string `gorm:"uniqueIndex;not null;size:64;column:public_id"`
	Bio               string `gorm:"default:NULL;size:255;column:bio"`
	DisplayName       string `gorm:"default:NULL;column:display_name"`
	Password          string
	IsVerified        bool
	PermissionGroupID uint
	RefreshToken      refresh_token.RefreshToken
	PermissionGroup   permission_group.PermissionGroup `gorm:"foreignKey:PermissionGroupID"`
}

func (u *User) Save() error {
	return database_connection.DB.Create(&u).Error
}

func (u *User) Exists() bool {
	err := database_connection.DB.Take(&u).Error
	return !errors.Is(err, gorm.ErrRecordNotFound)
}

func (u *User) Find() *User {
	err := database_connection.DB.Take(&u).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Unexpected error at 'user.Find': %s", err.Error())
		}
		return nil
	}
	return u
}
