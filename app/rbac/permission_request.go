package rbac

type PermissionStoreRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=64,unique_column=permissions name"`
	Description *string `json:"description" binding:"omitempty,max=255"`
	Group       string  `json:"group" binding:"required,min=2,max=64"`
}

type PermissionUpdateRequest struct {
	ExceptID    string  `json:"-" binding:"-"`
	Name        string  `json:"name" binding:"required,min=2,max=64,unique_column=permissions name ExceptID"`
	Description *string `json:"description" binding:"omitempty,max=255"`
	Group       string  `json:"group" binding:"required,min=2,max=64"`
}
