package rbac

type RoleStoreRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=64,unique_column=roles name"`
	Description *string `json:"description" binding:"omitempty,max=255"`
}

type RoleUpdateRequest struct {
	ExceptID    string  `json:"-" binding:"-"`
	Name        string  `json:"name" binding:"required,min=2,max=64,unique_column=roles name ExceptID"`
	Description *string `json:"description" binding:"omitempty,max=255"`
}
