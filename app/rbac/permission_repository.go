package rbac

import (
	"errors"
	errors2 "template-go-api/errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

func (r *PermissionRepository) GetAll(c *gin.Context) ([]Permission, error) {
	var permissions []Permission
	err := r.db.WithContext(c).Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) GetById(c *gin.Context, id string) (*Permission, error) {
	var permission Permission
	err := r.db.WithContext(c).Where("id = ?", id).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors2.DataNotFoundException()
		}
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) AddNew(c *gin.Context, request PermissionStoreRequest) (*Permission, error) {
	permission := Permission{
		Name:        request.Name,
		Description: request.Description,
		Group:       request.Group,
	}
	if err := r.db.WithContext(c).Create(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) UpdateById(c *gin.Context, id string, request PermissionUpdateRequest) (*Permission, error) {
	var permission Permission
	if err := r.db.WithContext(c).First(&permission, "id = ?", id).Error; err != nil {
		return nil, errors2.DataNotFoundException()
	}

	permission.Name = request.Name
	permission.Description = request.Description
	permission.Group = request.Group

	if err := r.db.WithContext(c).Save(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) DeleteById(c *gin.Context, id string) error {
	result := r.db.WithContext(c).Delete(&Permission{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors2.DataNotFoundException()
	}
	return nil
}
