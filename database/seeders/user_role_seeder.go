package seeders

import (
	"log"
	"template-go-api/app/rbac"
	"template-go-api/app/user"
	"template-go-api/enums"

	"gorm.io/gorm"
)

type UserRoleSeeder struct{}

func (s *UserRoleSeeder) Run(db *gorm.DB) error {
	var adminUser user.User
	result := db.Where("email = ?", "iqbalatma@gmail.com").First(&adminUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Println("Admin user not found, skipping role assignment")
			return nil
		}
		return result.Error
	}

	var adminRole rbac.Role
	result = db.Where("name = ?", string(enums.ADMIN)).First(&adminRole)
	if result.Error != nil {
		return result.Error
	}

	var superadminRole rbac.Role
	result = db.Where("name = ?", string(enums.SUPERADMIN)).First(&superadminRole)
	if result.Error != nil {
		return result.Error
	}

	if err := db.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
		return err
	}
	log.Printf("Attached admin role to admin user")

	if err := db.Model(&adminUser).Association("Roles").Append(&superadminRole); err != nil {
		return err
	}
	log.Printf("Attached superadmin role to admin user")

	return nil
}
