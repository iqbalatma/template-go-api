package user

import (
	"template-go-api/enums"
	"template-go-api/utils"
	"template-go-api/validator"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repository *Repository
}

func NewHandler(repository *Repository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) Index(c *gin.Context) error {
	users, meta, err := h.repository.GetAllPaginated(c)
	if err != nil {
		return err
	}
	utils.ResponseJSON(c, enums.SUCCESS, "Users retrieved", utils.Payload{
		Data: NewResourceCollection(users),
		Meta: meta,
	})
	return nil
}

func (h *Handler) Show(c *gin.Context) error {
	id := c.Param("id")
	user, err := h.repository.GetById(c, id)
	if err != nil {
		return err
	}
	utils.ResponseJSON(c, enums.SUCCESS, "Get user by id successfully", NewResource(user))
	return nil
}

func (h *Handler) Store(c *gin.Context) error {
	var request StoreRequest
	if !validator.BindAndValidate(c, &request) {
		return nil
	}

	user, err := h.repository.AddNew(c, request)
	if err != nil {
		return err
	}

	utils.ResponseJSON(c, enums.CREATED, "Add new user successfully", NewResource(user))
	return nil
}

func (h *Handler) Update(c *gin.Context) error {
	var request UpdateRequest
	request.ExceptID = c.Param("id")
	if !validator.BindAndValidate(c, &request) {
		return nil
	}

	user, err := h.repository.UpdateById(c, c.Param("id"), request)
	if err != nil {
		return err
	}

	utils.ResponseJSON(c, enums.SUCCESS, "Update user by id successfully", NewResource(user))
	return nil
}

func (h *Handler) Destroy(c *gin.Context) error {
	err := h.repository.DeleteById(c, c.Param("id"))
	if err != nil {
		return err
	}

	utils.ResponseJSON(c, enums.SUCCESS, "Delete user by id successfully", nil)
	return nil
}
