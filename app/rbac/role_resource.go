package rbac

type RoleResource struct {
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	IsMutable   bool                   `json:"is_mutable"`
	Description *string                `json:"description"`
	Permissions []PermissionResource   `json:"permissions"`
}

func NewRoleResource(role *Role) *RoleResource {
	permissions := make([]PermissionResource, len(role.Permissions))
	for i, p := range role.Permissions {
		permissions[i] = *NewPermissionResource(&p)
	}
	return &RoleResource{
		Id:          role.ID.String(),
		Name:        role.Name,
		IsMutable:   role.IsMutable,
		Description: role.Description,
		Permissions: permissions,
	}
}

func NewRoleResourceCollection(roles []Role) []RoleResource {
	result := make([]RoleResource, len(roles))
	for i, role := range roles {
		result[i] = *NewRoleResource(&role)
	}
	return result
}
