package rbac

type PermissionResource struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Group       string  `json:"group"`
}

func NewPermissionResource(p *Permission) *PermissionResource {
	return &PermissionResource{
		Id:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
		Group:       p.Group,
	}
}

func NewPermissionResourceCollection(permissions []Permission) []PermissionResource {
	result := make([]PermissionResource, len(permissions))
	for i, p := range permissions {
		result[i] = *NewPermissionResource(&p)
	}
	return result
}
