package routes

import (
	"template-go-api/app/auth"
	"template-go-api/app/user"

	"gorm.io/gorm"
)

type Container struct {
	AuthHandler           *auth.Handler
	ForgotPasswordHandler *auth.ForgotPasswordHandler
	ResetPasswordHandler  *auth.ResetPasswordHandler
	UserHandler           *user.Handler
}

func NewContainer(db *gorm.DB) *Container {
	authHandler := auth.NewHandler()
	forgotPasswordHandler := auth.NewForgotPasswordHandler()
	resetPasswordHandler := auth.NewResetPasswordHandler()

	userRepository := user.NewRepository(db)
	userHandler := user.NewHandler(userRepository)

	return &Container{
		AuthHandler:           authHandler,
		ForgotPasswordHandler: forgotPasswordHandler,
		ResetPasswordHandler:  resetPasswordHandler,
		UserHandler:           userHandler,
	}
}
