package user

import (
	"genericAPI/internal/models/auth_token"
	"genericAPI/internal/models/permission_group"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID                uint64    `gorm:"primaryKey;autoIncrement"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	Name              string
	Surname           string
	Username          string `gorm:"uniqueIndex;not null;size:32;column:username"`
	PublicID          string `gorm:"uniqueIndex;not null;size:64;column:public_id"`
	Bio               string `gorm:"default:NULL;size:255;column:bio"`
	DisplayName       string `gorm:"default:NULL;column:display_name"`
	Password          string
	IsVerified        bool
	PermissionGroupID uint
	AuthTokens        []auth_token.AuthToken           `gorm:"many2many:auth_token"`
	PermissionGroup   permission_group.PermissionGroup `gorm:"foreignKey:PermissionGroupID"`
}
