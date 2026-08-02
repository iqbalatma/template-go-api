package rbac

import (
	"template-go-api/enums"
	"template-go-api/utils"
	"template-go-api/validator"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	repository *RoleRepository
}

func NewRoleHandler(repository *RoleRepository) *RoleHandler {
	return &RoleHandler{
		repository: repository,
	}
}

func (h *RoleHandler) Index(c *gin.Context) error {
	roles, err := h.repository.GetAll(c, true)
	if err != nil {
		return err
	}
	utils.ResponseJSON(c, enums.SUCCESS, "Roles retrieved", NewRoleResourceCollection(roles))
	return nil
}

func (h *RoleHandler) MasterIndex(c *gin.Context) error {
	roles, err := h.repository.GetAll(c, false)
	if err != nil {
		return err
	}
	utils.ResponseJSON(c, enums.SUCCESS, "Roles retrieved", NewRoleMasterResourceCollection(roles))
	return nil
}

func (h *RoleHandler) Show(c *gin.Context) error {
	role, err := h.repository.GetById(c, c.Param("id"))
	if err != nil {
		return err
	}
	utils.ResponseJSON(c, enums.SUCCESS, "Get role by id successfully", NewRoleResource(role))
	return nil
}

func (h *RoleHandler) Store(c *gin.Context) error {
	var request RoleStoreRequest
	if !validator.BindAndValidate(c, &request) {
		return nil
	}

	role, err := h.repository.AddNew(c, request)
	if err != nil {
		return err
	}

	utils.ResponseJSON(c, enums.CREATED, "Add new role successfully", NewRoleResource(role))
	return nil
}

func (h *RoleHandler) Update(c *gin.Context) error {
	var request RoleUpdateRequest
	request.ExceptID = c.Param("id")
	if !validator.BindAndValidate(c, &request) {
		return nil
	}

	role, err := h.repository.UpdateById(c, c.Param("id"), request)
	if err != nil {
		return err
	}

	utils.ResponseJSON(c, enums.SUCCESS, "Update role by id successfully", NewRoleResource(role))
	return nil
}

func (h *RoleHandler) Destroy(c *gin.Context) error {
	if err := h.repository.DeleteById(c, c.Param("id")); err != nil {
		return err
	}
	utils.ResponseJSON(c, enums.SUCCESS, "Delete role by id successfully", nil)
	return nil
}

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
