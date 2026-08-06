package user

import (
	"template-go-api/app/media"
	"template-go-api/app/rbac"
	"template-go-api/utils"
)

type Resource struct {
	Id          string                    `json:"id"`
	FirstName   string                    `json:"first_name"`
	LastName    *string                   `json:"last_name"`
	Email       string                    `json:"email"`
	PhoneNumber *string                   `json:"phone_number"`
	CreatedAt   string                    `json:"created_at"`
	Roles       []rbac.RoleMasterResource `json:"roles"`
	Avatar      *media.Resource           `json:"avatar"`
}

func NewResource(user *User) *Resource {
	return &Resource{
		Id:          user.ID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		CreatedAt:   utils.FormatDateTimeVal(user.CreatedAt),
		Roles:       rbac.NewRoleMasterResourceCollection(user.Roles),
		Avatar:      media.NewResource(user.Avatar),
	}
}

func NewResourceCollection(users []User) []Resource {
	result := make([]Resource, len(users))
	for i, user := range users {
		result[i] = *NewResource(&user)
	}
	return result
}
