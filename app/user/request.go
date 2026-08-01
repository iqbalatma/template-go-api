package user

type StoreRequest struct {
	FirstName   string  `json:"first_name" binding:"required,min=2,max=64"`
	LastName    *string `json:"last_name" binding:"omitempty,max=64"`
	Email       string  `json:"email" binding:"required,min=2,max=64,unique_column=users email"`
	Password    string  `json:"password" binding:"required,min=8,max=255"`
	PhoneNumber *string `json:"phone_number" binding:"omitempty,max=20,unique_column=users phone_number"`
}

type UpdateRequest struct {
	ExceptID    string  `json:"-" binding:"-"`
	FirstName   string  `json:"first_name" binding:"required,min=2,max=64"`
	LastName    *string `json:"last_name" binding:"omitempty,max=64"`
	Email       string  `json:"email" binding:"required,min=2,max=64,unique_column=users email ExceptID"`
	PhoneNumber *string `json:"phone_number" binding:"omitempty,max=20,unique_column=users phone_number ExceptID"`
}
