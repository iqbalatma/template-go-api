package user

import (
	"errors"
	"template-go-api/app/rbac"
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
	meta, err := utils.Paginate(c, r.db.WithContext(c).Preload("Roles"), &users)
	if err != nil {
		return nil, meta, err
	}
	return users, meta, nil
}

func (r *Repository) GetById(c *gin.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(c).Preload("Roles").Where("id = ?", id).First(&user).Error
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

	err = r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if len(request.RoleIDs) > 0 {
			if err := r.syncRoles(tx, &user, request.RoleIDs); err != nil {
				return err
			}
		}

		return tx.Preload("Roles").First(&user, "id = ?", user.ID).Error
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateById(c *gin.Context, id string, request UpdateRequest) (*User, error) {
	var user User

	err := r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors2.DataNotFoundException()
			}
			return err
		}

		user.Email = request.Email
		user.FirstName = request.FirstName
		user.LastName = request.LastName
		user.PhoneNumber = request.PhoneNumber

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		if request.RoleIDs != nil {
			if err := r.syncRoles(tx, &user, request.RoleIDs); err != nil {
				return err
			}
		}

		return tx.Preload("Roles").First(&user, "id = ?", id).Error
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) syncRoles(tx *gorm.DB, user *User, roleIDs []string) error {
	var roles []rbac.Role
	if len(roleIDs) > 0 {
		if err := tx.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return err
		}
	}
	return tx.Model(user).Association("Roles").Replace(&roles)
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
