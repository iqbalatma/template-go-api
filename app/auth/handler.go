package auth

import (
	"errors"
	"template-go-api/app/user"
	"template-go-api/config"
	"template-go-api/enums"
	errors2 "template-go-api/errors"
	"template-go-api/utils"
	"template-go-api/validator"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iqbalatma/gofortify"
	"gorm.io/gorm"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Login(c *gin.Context) error {
	var request LoginRequest
	if !validator.BindAndValidate(c, &request) {
		return nil
	}

	var u user.User
	result := config.DB.Preload("Roles.Permissions").Where("email = ?", request.Email).First(&u)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			httpErr := errors2.InvalidCredential()
			c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError(httpErr.Message, httpErr.Code, nil))
			return nil
		}
		return result.Error
	}

	if u.Password == nil || !utils.CheckHash(*u.Password, request.Password) {
		httpErr := errors2.InvalidCredential()
		c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError(httpErr.Message, httpErr.Code, nil))
		return nil
	}

	accessTokenJTI := uuid.New().String()
	refreshTokenJTI := uuid.New().String()

	accessToken, atv, err := gofortify.Encode(&u, gofortify.AccessToken, true, c.Request.UserAgent(), accessTokenJTI, refreshTokenJTI)
	if err != nil {
		return err
	}

	refreshToken, _, err := gofortify.Encode(&u, gofortify.RefreshToken, true, c.Request.UserAgent(), refreshTokenJTI, accessTokenJTI)
	if err != nil {
		return err
	}

	c.SetCookie("access_token_verifier", atv, gofortify.Config.AccessTokenTTL*60, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, gofortify.Config.RefreshTokenTTL*60, "/", "", false, true)

	utils.ResponseJSON(c, enums.SUCCESS, "Login successfully", NewResource(&u, accessToken, refreshToken))
	return nil
}

func (h *Handler) Logout(c *gin.Context) error {
	payload, exists := c.Get("payload")
	if !exists {
		httpErr := errors2.UnauthorizedException()
		c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError(httpErr.Message, httpErr.Code, nil))
		return nil
	}

	gofortify.BlacklistTokenAndPairToken(payload.(*gofortify.Payload))
	utils.ResponseJSON(c, enums.SUCCESS, "Logout successfully", nil)
	return nil
}

func (h *Handler) Refresh(c *gin.Context) error {
	u := c.MustGet("user").(*user.User)

	accessTokenJTI := uuid.New().String()
	refreshTokenJTI := uuid.New().String()

	accessToken, atv, err := gofortify.Encode(u, gofortify.AccessToken, true, c.Request.UserAgent(), accessTokenJTI, refreshTokenJTI)
	if err != nil {
		return err
	}

	refreshToken, _, err := gofortify.Encode(u, gofortify.RefreshToken, true, c.Request.UserAgent(), refreshTokenJTI, accessTokenJTI)
	if err != nil {
		return err
	}

	c.SetCookie("access_token_verifier", atv, gofortify.Config.AccessTokenTTL*60, "/", "", true, true)
	c.SetCookie("refresh_token", refreshToken, gofortify.Config.RefreshTokenTTL*60, "/", "", true, true)

	utils.ResponseJSON(c, enums.SUCCESS, "Token refreshed successfully", NewResource(u, accessToken, refreshToken))
	return nil
}
