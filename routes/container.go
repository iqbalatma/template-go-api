package routes

import (
	"template-go-api/app/auth"
	"template-go-api/app/media"
	"template-go-api/app/rbac"
	"template-go-api/app/user"

	"gorm.io/gorm"
)

type Container struct {
	AuthHandler           *auth.Handler
	ForgotPasswordHandler *auth.ForgotPasswordHandler
	ResetPasswordHandler  *auth.ResetPasswordHandler
	UserHandler           *user.Handler
	ProfileHandler        *user.ProfileHandler
	PermissionHandler     *rbac.PermissionHandler
	RoleHandler           *rbac.RoleHandler
}

func NewContainer(db *gorm.DB) *Container {
	mediaRepository := media.NewRepository(db)
	userRepository := user.NewRepository(db, mediaRepository)
	permissionRepository := rbac.NewPermissionRepository(db)
	roleRepository := rbac.NewRoleRepository(db)

	// Setiap model mendaftarkan collection miliknya sendiri.
	user.RegisterMediaCollections()

	return &Container{
		AuthHandler:           auth.NewHandler(),
		ForgotPasswordHandler: auth.NewForgotPasswordHandler(),
		ResetPasswordHandler:  auth.NewResetPasswordHandler(),
		UserHandler:           user.NewHandler(userRepository),
		ProfileHandler:        user.NewProfileHandler(userRepository),
		PermissionHandler:     rbac.NewPermissionHandler(permissionRepository),
		RoleHandler:           rbac.NewRoleHandler(roleRepository),
	}
}
