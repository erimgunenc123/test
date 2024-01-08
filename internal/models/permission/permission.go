package permission

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	ID          int    `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"default:NULL;size:100;column:permission_name"`
	Description string `gorm:"default:NULL;size:255;column:description"`
}
