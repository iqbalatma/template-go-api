package user

import (
	"errors"
	errors2 "template-go-api/errors"
	"template-go-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAllPaginated(c *gin.Context) ([]User, *utils.PaginationMeta, error) {
	var users []User
	meta, err := utils.Paginate(c, r.db, &users)
	if err != nil {
		return nil, meta, err
	}
	return users, meta, nil
}

func (r *Repository) GetById(c *gin.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(c).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors2.DataNotFoundException()
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) AddNew(c *gin.Context, request StoreRequest) (*User, error) {
	hashedPassword, err := utils.MakeHash(request.Password)
	if err != nil {
		return nil, err
	}

	user := User{
		Email:       request.Email,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: request.PhoneNumber,
		Password:    &hashedPassword,
	}
	if err := r.db.WithContext(c).Create(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateById(c *gin.Context, id string, request UpdateRequest) (*User, error) {
	var user User
	if err := r.db.WithContext(c).First(&user, "id = ?", id).Error; err != nil {
		return nil, errors2.DataNotFoundException()
	}

	user.Email = request.Email
	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.PhoneNumber = request.PhoneNumber

	if err := r.db.WithContext(c).Save(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) DeleteById(c *gin.Context, id string) error {
	result := r.db.WithContext(c).Delete(&User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors2.DataNotFoundException()
	}
	return nil
}
