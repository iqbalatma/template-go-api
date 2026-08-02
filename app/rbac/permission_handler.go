package rbac

import (
	"template-go-api/enums"
	"template-go-api/utils"
	"template-go-api/validator"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	repository *PermissionRepository
}

func NewPermissionHandler(repository *PermissionRepository) *PermissionHandler {
	return &PermissionHandler{
		repository: repository,
	}
}

func (h *PermissionHandler) Index(c *gin.Context) error {
	permissions, err := h.repository.GetAll(c)
	if err != nil {
		return err
	}
	utils.ResponseJSON(c, enums.SUCCESS, "Permissions retrieved", NewPermissionResourceCollection(permissions))
	return nil
}

func (h *PermissionHandler) Show(c *gin.Context) error {
	permission, err := h.repository.GetById(c, c.Param("id"))
	if err != nil {
		return err
	}
	utils.ResponseJSON(c, enums.SUCCESS, "Get permission by id successfully", NewPermissionResource(permission))
	return nil
}

func (h *PermissionHandler) Store(c *gin.Context) error {
	var request PermissionStoreRequest
	if !validator.BindAndValidate(c, &request) {
		return nil
	}

	permission, err := h.repository.AddNew(c, request)
	if err != nil {
		return err
	}

	utils.ResponseJSON(c, enums.CREATED, "Add new permission successfully", NewPermissionResource(permission))
	return nil
}

func (h *PermissionHandler) Update(c *gin.Context) error {
	var request PermissionUpdateRequest
	request.ExceptID = c.Param("id")
	if !validator.BindAndValidate(c, &request) {
		return nil
	}

	permission, err := h.repository.UpdateById(c, c.Param("id"), request)
	if err != nil {
		return err
	}

	utils.ResponseJSON(c, enums.SUCCESS, "Update permission by id successfully", NewPermissionResource(permission))
	return nil
}

func (h *PermissionHandler) Destroy(c *gin.Context) error {
	if err := h.repository.DeleteById(c, c.Param("id")); err != nil {
		return err
	}
	utils.ResponseJSON(c, enums.SUCCESS, "Delete permission by id successfully", nil)
	return nil
}
