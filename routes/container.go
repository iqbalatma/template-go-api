package routes

import (
	"template-go-api/app/auth"
	"template-go-api/app/rbac"
	"template-go-api/app/user"

	"gorm.io/gorm"
)

type Container struct {
	AuthHandler           *auth.Handler
	ForgotPasswordHandler *auth.ForgotPasswordHandler
	ResetPasswordHandler  *auth.ResetPasswordHandler
	UserHandler           *user.Handler
	PermissionHandler     *rbac.PermissionHandler
	RoleHandler           *rbac.RoleHandler
}

func NewContainer(db *gorm.DB) *Container {
	userRepository := user.NewRepository(db)
	permissionRepository := rbac.NewPermissionRepository(db)
	roleRepository := rbac.NewRoleRepository(db)

	return &Container{
		AuthHandler:           auth.NewHandler(),
		ForgotPasswordHandler: auth.NewForgotPasswordHandler(),
		ResetPasswordHandler:  auth.NewResetPasswordHandler(),
		UserHandler:           user.NewHandler(userRepository),
		PermissionHandler:     rbac.NewPermissionHandler(permissionRepository),
		RoleHandler:           rbac.NewRoleHandler(roleRepository),
	}
}
