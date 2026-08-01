package auth

import (
	"fmt"
	"template-go-api/app/user"
)

type Resource struct {
	Id          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    *string `json:"last_name"`
	FullName    string  `json:"full_name"`
	Email       string  `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Tokens      struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	} `json:"tokens"`
}

func NewResource(u *user.User, accessToken, refreshToken string) *Resource {
	lastName := ""
	if u.LastName != nil {
		lastName = *u.LastName
	}

	r := &Resource{
		Id:          u.ID.String(),
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		FullName:    fmt.Sprintf("%s %s", u.FirstName, lastName),
		Email:       u.Email,
		PhoneNumber: u.PhoneNumber,
	}
	r.Tokens.AccessToken = accessToken
	r.Tokens.RefreshToken = refreshToken
	return r
}
