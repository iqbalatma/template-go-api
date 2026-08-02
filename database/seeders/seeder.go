package seeders

import (
	"log"

	"gorm.io/gorm"
)

type Seeder interface {
	Run(db *gorm.DB) error
}

var registry []Seeder

func RegisterAll() {
	registry = []Seeder{
		&RoleSeeder{},
		&PermissionSeeder{},
		&RolePermissionSeeder{},
		&UserSeeder{},
		&UserRoleSeeder{},
	}
}

func RunAll(db *gorm.DB) {
	if len(registry) == 0 {
		RegisterAll()
	}

	for _, s := range registry {
		if err := s.Run(db); err != nil {
			log.Fatalf("Seeder failed: %v", err)
		}
	}
}
