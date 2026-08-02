package rbac

import (
	"template-go-api/app"
)

type Role struct {
	app.BaseModel `gorm:"embedded"`
	Name          string       `json:"name" gorm:"column:name;unique;not null"`
	IsMutable     bool         `json:"is_mutable" gorm:"column:is_mutable;not null"`
	Description   *string      `json:"description" gorm:"column:description"`
	Permissions   []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}
