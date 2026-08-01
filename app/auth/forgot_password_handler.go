package auth

type ForgotPasswordHandler struct{}

func NewForgotPasswordHandler() *ForgotPasswordHandler {
	return &ForgotPasswordHandler{}
}

//
//func (h *ForgotPasswordHandler) ForgotPassword(c *gin.Context) error {
//	var request ForgotPasswordRequest
//	if !validator.BindAndValidate(c, &request) {
//		return nil
//	}
//
//	var u user.User
//	result := config.DB.Where("email = ?", request.Email).First(&u)
//	if result.Error != nil {
//		utils.ResponseJSON(c, enums.SUCCESS, "If that email exists, a reset link has been sent", nil)
//		return nil
//	}
//
//	config.DB.Where("email = ?", request.Email).Delete(&PasswordResetToken{})
//
//	b := make([]byte, 32)
//	if _, err := rand.Read(b); err != nil {
//		return err
//	}
//	plainToken := hex.EncodeToString(b)
//
//	resetToken := PasswordResetToken{
//		Email:     request.Email,
//		Token:     plainToken,
//		ExpiresAt: time.Now().Add(time.Hour),
//	}
//	if err := config.DB.Create(&resetToken).Error; err != nil {
//		return err
//	}
//
//	if err := sendResetEmail(request.Email, plainToken, resetToken.ExpiresAt); err != nil {
//		log.Printf("[ForgotPassword] failed to send email to %s: %v", request.Email, err)
//	}
//
//	utils.ResponseJSON(c, enums.SUCCESS, "If that email exists, a reset link has been sent", nil)
//	return nil
//}
//
//func sendResetEmail(to, token string, expiresAt time.Time) error {
//	subject := "Password Reset Request"
//	body := fmt.Sprintf(`
//<html>
//<body>
//  <p>Hi,</p>
//  <p>You requested a password reset. Use the token below to reset your password.</p>
//  <p><strong>Token:</strong> %s</p>
//  <p>This token will expire at <strong>%s</strong>.</p>
//  <p>If you did not request this, you can safely ignore this email.</p>
//</body>
//</html>`, token, expiresAt.Format("2006-01-02 15:04:05 UTC"))
//
//	return utils.SendMail(utils.Mail{
//		To:      []string{to},
//		Subject: subject,
//		Body:    body,
//	})
//}
