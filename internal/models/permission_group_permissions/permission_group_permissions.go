package permission_group_permissions

import "gorm.io/gorm"

type PermissionGroupPermission struct {
	gorm.Model
	PermissionGroupID uint `gorm:"primaryKey"`
	PermissionID      uint `gorm:"primaryKey"`
}
