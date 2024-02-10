package permission_group

import (
	"genericAPI/api/database_connection"
	"genericAPI/internal/models/permission"
	"gorm.io/gorm"
)

type PermissionGroup struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	Name        string
	Permissions []permission.Permission `gorm:"many2many:permission_group_permissions"`
}

func (p *PermissionGroup) Save() error {
	return database_connection.DB.Create(&p).Error

}
