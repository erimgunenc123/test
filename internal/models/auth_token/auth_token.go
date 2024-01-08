package auth_token

import (
	"genericAPI/internal/models/token_type"
	"gorm.io/gorm"
)

type AuthToken struct {
	gorm.Model
	Token       string `gorm:"not null;size:255"`
	UserID      uint64 `gorm:"uniqueIndex:uidx_userid_token_type"`
	TokenTypeID uint
	TokenType   token_type.TokenType
}
