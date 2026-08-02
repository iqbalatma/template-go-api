package rbac

import (
	"errors"
	errors2 "template-go-api/errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) GetAll(c *gin.Context, withPermissions bool) ([]Role, error) {
	var roles []Role
	query := r.db.WithContext(c)
	if withPermissions {
		query = query.Preload("Permissions")
	}
	err := query.Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetById(c *gin.Context, id string) (*Role, error) {
	var role Role
	err := r.db.WithContext(c).Preload("Permissions").Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors2.DataNotFoundException()
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) AddNew(c *gin.Context, request RoleStoreRequest) (*Role, error) {
	role := Role{
		Name:        request.Name,
		IsMutable:   request.IsMutable,
		Description: request.Description,
	}
	if err := r.db.WithContext(c).Create(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) UpdateById(c *gin.Context, id string, request RoleUpdateRequest) (*Role, error) {
	var role Role
	if err := r.db.WithContext(c).First(&role, "id = ?", id).Error; err != nil {
		return nil, errors2.DataNotFoundException()
	}

	role.Name = request.Name
	role.IsMutable = request.IsMutable
	role.Description = request.Description

	if err := r.db.WithContext(c).Save(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) DeleteById(c *gin.Context, id string) error {
	result := r.db.WithContext(c).Delete(&Role{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors2.DataNotFoundException()
	}
	return nil
}
