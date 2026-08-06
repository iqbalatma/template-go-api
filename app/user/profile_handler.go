package user

import (
	"template-go-api/app/media"
	"template-go-api/enums"
	errors2 "template-go-api/errors"
	"template-go-api/utils"
	"template-go-api/validator"

	"github.com/gin-gonic/gin"
)

// avatarFormField adalah nama field multipart yang menampung file avatar.
const avatarFormField = "avatar"

// ProfileHandler melayani route /api/me, yaitu aksi user terhadap dirinya
// sendiri. Ownernya selalu diambil dari token, tidak pernah dari input client.
// Bandingkan dengan Handler yang melayani CRUD user oleh admin.
type ProfileHandler struct {
	repository *Repository
}

func NewProfileHandler(repository *Repository) *ProfileHandler {
	return &ProfileHandler{
		repository: repository,
	}
}

// currentUser mengambil user yang sedang login dari context. Nilai ini diisi
// oleh AuthMiddleware berdasarkan token.
func currentUser(c *gin.Context) (*User, error) {
	raw, exists := c.Get("user")
	if !exists {
		return nil, errors2.UnauthorizedException()
	}

	u, ok := raw.(*User)
	if !ok {
		return nil, errors2.UnauthorizedException()
	}
	return u, nil
}

// Me menampilkan profil user yang sedang login.
// GET /api/me
func (h *ProfileHandler) Me(c *gin.Context) error {
	u, err := currentUser(c)
	if err != nil {
		return err
	}

	// AuthMiddleware hanya memuat kolom user dan rolenya, avatar dimuat di sini
	// supaya client tidak perlu memanggil endpoint terpisah.
	if err := h.repository.LoadAvatar(c, u); err != nil {
		return err
	}

	utils.ResponseJSON(c, enums.SUCCESS, "Authenticated user", NewResource(u))
	return nil
}

// Update mengubah data diri user yang sedang login, termasuk avatarnya.
// Kirim sebagai multipart/form-data bila menyertakan file avatar, atau JSON
// biasa bila hanya mengubah data teks.
// PATCH /api/me
func (h *ProfileHandler) Update(c *gin.Context) error {
	u, err := currentUser(c)
	if err != nil {
		return err
	}

	var request UpdateProfileRequest
	// ExceptID diambil dari token supaya pengecekan unique tidak menabrak
	// data milik user itu sendiri.
	request.ExceptID = u.ID.String()
	if !validator.BindAndValidateForm(c, &request) {
		return nil
	}

	// Avatar bersifat opsional: tanpa file, avatar lama dibiarkan apa adanya.
	// Diproses lebih dulu supaya file yang ditolak tidak menyisakan perubahan
	// data diri yang terlanjur tersimpan.
	if fileHeader, err := c.FormFile(avatarFormField); err == nil {
		if _, err := h.repository.UpdateAvatar(c, u, fileHeader); err != nil {
			return err
		}
	} else if err := h.repository.LoadAvatar(c, u); err != nil {
		return err
	}

	if err := h.repository.UpdateProfile(c, u, request); err != nil {
		return err
	}

	utils.ResponseJSON(c, enums.SUCCESS, "Update profile successfully", NewResource(u))
	return nil
}

// ShowAvatar menampilkan avatar milik user yang sedang login.
// GET /api/me/avatar
func (h *ProfileHandler) ShowAvatar(c *gin.Context) error {
	u, err := currentUser(c)
	if err != nil {
		return err
	}

	if err := h.repository.LoadAvatar(c, u); err != nil {
		return err
	}

	utils.ResponseJSON(c, enums.SUCCESS, "Get avatar successfully", media.NewResource(u.Avatar))
	return nil
}
