package seeders

import (
	"log"
	"template-go-api/app/user"
	"template-go-api/utils"

	"gorm.io/gorm"
)

type UserSeeder struct{}

var seedUsers = []struct {
	FirstName   string
	LastName    string
	Email       string
	Password    string
	PhoneNumber string
}{
	{"Admin", "User", "iqbalatma@gmail.com", "admin", "+628100000001"},
	{"John", "Doe", "john.doe@example.com", "password", "+628100000002"},
	{"Jane", "Doe", "jane.doe@example.com", "password", "+628100000003"},
}

func (s *UserSeeder) Run(db *gorm.DB) error {
	for _, u := range seedUsers {
		hashed, err := utils.MakeHash(u.Password)
		if err != nil {
			return err
		}

		lastName := u.LastName
		phoneNumber := u.PhoneNumber
		record := user.User{
			FirstName:   u.FirstName,
			LastName:    &lastName,
			Email:       u.Email,
			Password:    &hashed,
			PhoneNumber: &phoneNumber,
		}

		result := db.Where(user.User{Email: u.Email}).FirstOrCreate(&record)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected > 0 {
			log.Printf("Seeded user: %s", u.Email)
		} else {
			log.Printf("Skipped (already exists): %s", u.Email)
		}
	}
	return nil
}
