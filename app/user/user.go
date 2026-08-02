package user

import (
	"template-go-api/app"
	"template-go-api/app/rbac"
)

type User struct {
	app.BaseModel `gorm:"embedded"`
	FirstName     string      `json:"first_name" gorm:"column:first_name"`
	LastName      *string     `json:"last_name" gorm:"column:last_name"`
	Email         string      `json:"email" gorm:"column:email;unique"`
	Password      *string     `json:"-" gorm:"column:password"`
	PhoneNumber   *string     `json:"phone_number" gorm:"column:phone_number;unique"`
	Roles         []rbac.Role `json:"roles,omitempty" gorm:"many2many:user_roles;"`
}

func (user *User) GetSubjectKey() string {
	return user.ID.String()
}
