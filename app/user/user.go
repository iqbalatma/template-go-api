package user

import (
	"template-go-api/app"
	"template-go-api/app/media"
	"template-go-api/app/rbac"
)

// MediaModelType adalah nilai kolom media.model_type untuk User.
const MediaModelType = "users"

type User struct {
	app.BaseModel `gorm:"embedded"`
	FirstName     string      `json:"first_name" gorm:"column:first_name"`
	LastName      *string     `json:"last_name" gorm:"column:last_name"`
	Email         string      `json:"email" gorm:"column:email;unique"`
	Password      *string     `json:"-" gorm:"column:password"`
	PhoneNumber   *string     `json:"phone_number" gorm:"column:phone_number;unique"`
	Roles         []rbac.Role `json:"roles,omitempty" gorm:"many2many:user_roles;"`
	// Avatar diisi manual oleh repository karena relasi media bersifat
	// polymorphic sehingga tidak bisa di-Preload oleh GORM.
	Avatar *media.Media `json:"avatar,omitempty" gorm:"-"`
}

func (user *User) GetSubjectKey() string {
	return user.ID.String()
}

func (user *User) GetMediaModelType() string {
	return MediaModelType
}

func (user *User) GetMediaModelID() string {
	return user.ID.String()
}
