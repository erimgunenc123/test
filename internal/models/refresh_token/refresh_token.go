package refresh_token

import (
	"errors"
	"genericAPI/api/database_connection"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
)

type RefreshToken struct {
	gorm.Model
	Token  string `gorm:"not null;size:255"`
	UserID uint64 `gorm:"uniqueIndex:uidx_userid_token_type"`
}

func (r *RefreshToken) Save() error {
	return database_connection.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "UserID"}},
		DoUpdates: clause.AssignmentColumns([]string{"token"}),
	}).Create(&r).Error
}

func (r *RefreshToken) Find() *RefreshToken {
	err := database_connection.DB.Take(&r).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Unexpected error at 'refresh_token.Find': %s", err.Error())
		}
		return nil
	}
	return r
}
