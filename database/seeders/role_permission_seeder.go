package seeders

import (
	"log"
	"template-go-api/app/rbac"
	"template-go-api/enums"

	"gorm.io/gorm"
)

type RolePermissionSeeder struct{}

func (s *RolePermissionSeeder) Run(db *gorm.DB) error {
	var adminRole rbac.Role
	result := db.Where("name = ?", string(enums.ADMIN)).First(&adminRole)
	if result.Error != nil {
		return result.Error
	}

	var permissions []rbac.Permission
	result = db.Find(&permissions)
	if result.Error != nil {
		return result.Error
	}

	for _, perm := range permissions {
		if err := db.Model(&adminRole).Association("Permissions").Append(&perm); err != nil {
			return err
		}
	}
	log.Printf("Attached all permissions to admin role")

	return nil
}
