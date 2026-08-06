package user

type StoreRequest struct {
	FirstName   string   `json:"first_name" binding:"required,min=2,max=64"`
	LastName    *string  `json:"last_name" binding:"omitempty,max=64"`
	Email       string   `json:"email" binding:"required,min=2,max=64,unique_column=users email"`
	Password    string   `json:"password" binding:"required,min=8,max=255"`
	PhoneNumber *string  `json:"phone_number" binding:"omitempty,max=20,unique_column=users phone_number"`
	RoleIDs     []string `json:"role_ids" binding:"omitempty,dive,uuid"`
}

type UpdateRequest struct {
	ExceptID    string   `json:"-" binding:"-"`
	FirstName   string   `json:"first_name" binding:"required,min=2,max=64"`
	LastName    *string  `json:"last_name" binding:"omitempty,max=64"`
	Email       string   `json:"email" binding:"required,min=2,max=64,unique_column=users email ExceptID"`
	PhoneNumber *string  `json:"phone_number" binding:"omitempty,max=20,unique_column=users phone_number ExceptID"`
	RoleIDs     []string `json:"role_ids" binding:"omitempty,dive,uuid"`
}

// UpdateProfileRequest sengaja tidak memuat role_ids maupun password. Role
// hanya boleh diubah lewat endpoint management agar user tidak bisa menaikkan
// hak aksesnya sendiri, dan password punya alurnya sendiri.
//
// Tag form dibutuhkan karena endpoint ini juga menerima multipart/form-data
// saat avatar ikut diunggah. Filenya sendiri diambil lewat c.FormFile.
type UpdateProfileRequest struct {
	ExceptID    string  `json:"-" form:"-" binding:"-"`
	FirstName   string  `json:"first_name" form:"first_name" binding:"required,min=2,max=64"`
	LastName    *string `json:"last_name" form:"last_name" binding:"omitempty,max=64"`
	Email       string  `json:"email" form:"email" binding:"required,min=2,max=64,unique_column=users email ExceptID"`
	PhoneNumber *string `json:"phone_number" form:"phone_number" binding:"omitempty,max=20,unique_column=users phone_number ExceptID"`
}
