package seeders

import (
	"log"
	"template-go-api/app/rbac"
	"template-go-api/enums"

	"gorm.io/gorm"
)

type RoleSeeder struct{}

var seedRoles = []struct {
	Name        string
	IsMutable   bool
	Description string
}{
	{string(enums.SUPERADMIN), false, "Has full access to all resources"},
	{string(enums.ADMIN), false, "Has administrative access"},
}

func (s *RoleSeeder) Run(db *gorm.DB) error {
	for _, r := range seedRoles {
		description := r.Description
		record := rbac.Role{
			Name:        r.Name,
			IsMutable:   r.IsMutable,
			Description: &description,
		}

		result := db.Where(rbac.Role{Name: r.Name}).FirstOrCreate(&record)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected > 0 {
			log.Printf("Seeded role: %s", r.Name)
		} else {
			log.Printf("Skipped (already exists): %s", r.Name)
		}
	}
	return nil
}
