package permission_group

import (
	"genericAPI/internal/models/permission"
	"gorm.io/gorm"
)

type PermissionGroup struct {
	gorm.Model
	ID          uint `gorm:"primaryKey"`
	Name        string
	Permissions []permission.Permission `gorm:"many2many:permission_group_permissions"`
}
