package auth

import (
	"template-go-api/app/user"
	"template-go-api/config"
	"template-go-api/enums"
	errors2 "template-go-api/errors"
	"template-go-api/utils"
	"template-go-api/validator"
	"time"

	"github.com/gin-gonic/gin"
)

type ResetPasswordHandler struct{}

func NewResetPasswordHandler() *ResetPasswordHandler {
	return &ResetPasswordHandler{}
}

func (h *ResetPasswordHandler) ResetPassword(c *gin.Context) error {
	var request ResetPasswordRequest
	if !validator.BindAndValidate(c, &request) {
		return nil
	}

	var resetToken PasswordResetToken
	result := config.DB.Where("token = ?", request.Token).First(&resetToken)
	if result.Error != nil {
		httpErr := errors2.InvalidCredential("Invalid or expired reset token")
		c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError(httpErr.Message, httpErr.Code, nil))
		return nil
	}

	if time.Now().After(resetToken.ExpiresAt) {
		config.DB.Delete(&resetToken)
		httpErr := errors2.InvalidCredential("Invalid or expired reset token")
		c.AbortWithStatusJSON(httpErr.StatusCode, utils.NewHttpError(httpErr.Message, httpErr.Code, nil))
		return nil
	}

	hashed, err := utils.MakeHash(request.Password)
	if err != nil {
		return err
	}

	if err := config.DB.Model(&user.User{}).Where("email = ?", resetToken.Email).Update("password", hashed).Error; err != nil {
		return err
	}

	config.DB.Delete(&resetToken)

	utils.ResponseJSON(c, enums.SUCCESS, "Password reset successfully", nil)
	return nil
}
