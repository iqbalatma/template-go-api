package rbac

import (
	"template-go-api/enums"
	"template-go-api/utils"

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
