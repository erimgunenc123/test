package token_type

import "gorm.io/gorm"

type TokenType struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}
