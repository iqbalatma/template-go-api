package rbac

import (
	"template-go-api/app"
)

type Permission struct {
	app.BaseModel `gorm:"embedded"`
	Name          string  `json:"name" gorm:"column:name;unique;not null"`
	Description   *string `json:"description" gorm:"column:description"`
	Group         string  `json:"group" gorm:"column:group;not null"`
	Roles         []Role  `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}
