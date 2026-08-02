package seeders

import (
	"log"
	"template-go-api/app/rbac"
	"template-go-api/enums"

	"gorm.io/gorm"
)

type PermissionSeeder struct{}

func (s *PermissionSeeder) Run(db *gorm.DB) error {
	for _, p := range enums.AllPermissions {
		name := string(p.Name)
		description := p.Description
		record := rbac.Permission{
			Name:        name,
			Group:       p.Group,
			Description: &description,
		}

		result := db.Where(rbac.Permission{Name: name}).FirstOrCreate(&record)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected > 0 {
			log.Printf("Seeded permission: %s", name)
		} else {
			log.Printf("Skipped (already exists): %s", name)
		}
	}
	return nil
}
