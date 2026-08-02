package enums

type Permission string

const (
	// User permissions
	PermUserIndex  Permission = "user.index"
	PermUserShow   Permission = "user.show"
	PermUserCreate Permission = "user.create"
	PermUserUpdate Permission = "user.update"
	PermUserDelete Permission = "user.delete"

	// Role permissions
	PermRoleIndex  Permission = "role.index"
	PermRoleShow   Permission = "role.show"
	PermRoleCreate Permission = "role.create"
	PermRoleUpdate Permission = "role.update"
	PermRoleDelete Permission = "role.delete"

	// Permission permissions
	PermPermissionIndex  Permission = "permission.index"
	PermPermissionShow   Permission = "permission.show"
	PermPermissionCreate Permission = "permission.create"
	PermPermissionUpdate Permission = "permission.update"
	PermPermissionDelete Permission = "permission.delete"
)

type permissionMeta struct {
	Name        Permission
	Group       string
	Description string
}

var AllPermissions = []permissionMeta{
	{PermUserIndex, "User", "View user list"},
	{PermUserShow, "User", "View user detail"},
	{PermUserCreate, "User", "Create new user"},
	{PermUserUpdate, "User", "Update user"},
	{PermUserDelete, "User", "Delete user"},

	{PermRoleIndex, "Role", "View role list"},
	{PermRoleShow, "Role", "View role detail"},
	{PermRoleCreate, "Role", "Create new role"},
	{PermRoleUpdate, "Role", "Update role"},
	{PermRoleDelete, "Role", "Delete role"},

	{PermPermissionIndex, "Permission", "View permission list"},
	{PermPermissionShow, "Permission", "View permission detail"},
	{PermPermissionCreate, "Permission", "Create new permission"},
	{PermPermissionUpdate, "Permission", "Update permission"},
	{PermPermissionDelete, "Permission", "Delete permission"},
}
