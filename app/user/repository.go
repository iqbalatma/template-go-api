package user

import (
	"errors"
	"mime/multipart"
	"template-go-api/app/media"
	"template-go-api/app/rbac"
	errors2 "template-go-api/errors"
	"template-go-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Repository struct {
	db    *gorm.DB
	media *media.Repository
}

func NewRepository(db *gorm.DB, mediaRepository *media.Repository) *Repository {
	return &Repository{
		db:    db,
		media: mediaRepository,
	}
}

func (r *Repository) GetAllPaginated(c *gin.Context) ([]User, *utils.PaginationMeta, error) {
	var users []User
	meta, err := utils.Paginate(c, r.db.WithContext(c).Preload("Roles"), &users)
	if err != nil {
		return nil, meta, err
	}

	if err := r.loadAvatars(c, users); err != nil {
		return nil, meta, err
	}
	return users, meta, nil
}

func (r *Repository) GetById(c *gin.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(c).Preload("Roles").Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors2.DataNotFoundException()
		}
		return nil, err
	}

	if err := r.LoadAvatar(c, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfile mengubah data diri user yang sedang login. Kolom yang disentuh
// dibatasi secara eksplisit supaya role, password, dan asosiasi lain tidak ikut
// tertulis lewat endpoint ini.
func (r *Repository) UpdateProfile(c *gin.Context, user *User, request UpdateProfileRequest) error {
	err := r.db.WithContext(c).
		Model(&User{}).
		Where("id = ?", user.ID).
		Updates(map[string]any{
			"first_name":   request.FirstName,
			"last_name":    request.LastName,
			"email":        request.Email,
			"phone_number": request.PhoneNumber,
		}).Error
	if err != nil {
		return err
	}

	user.FirstName = request.FirstName
	user.LastName = request.LastName
	user.Email = request.Email
	user.PhoneNumber = request.PhoneNumber
	return nil
}

// UpdateAvatar mengganti avatar user. Collection avatar bersifat single file,
// jadi avatar lama otomatis dihapus beserta filenya.
func (r *Repository) UpdateAvatar(c *gin.Context, user *User, fileHeader *multipart.FileHeader) (*media.Media, error) {
	avatar, err := r.media.Attach(c, user, MediaCollectionAvatar, fileHeader)
	if err != nil {
		return nil, err
	}

	user.Avatar = avatar
	return avatar, nil
}

// LoadAvatar memuat avatar satu user.
func (r *Repository) LoadAvatar(c *gin.Context, user *User) error {
	avatar, err := r.media.GetFirst(c, user, MediaCollectionAvatar)
	if err != nil {
		return err
	}
	user.Avatar = avatar
	return nil
}

// loadAvatars memuat avatar banyak user dalam satu query agar tidak N+1.
func (r *Repository) loadAvatars(c *gin.Context, users []User) error {
	if len(users) == 0 {
		return nil
	}

	ids := make([]string, len(users))
	for i := range users {
		ids[i] = users[i].GetMediaModelID()
	}

	grouped, err := r.media.GetForModels(c, MediaModelType, ids, MediaCollectionAvatar)
	if err != nil {
		return err
	}

	for i := range users {
		if items := grouped[users[i].GetMediaModelID()]; len(items) > 0 {
			users[i].Avatar = &items[0]
		}
	}
	return nil
}

func (r *Repository) AddNew(c *gin.Context, request StoreRequest) (*User, error) {
	hashedPassword, err := utils.MakeHash(request.Password)
	if err != nil {
		return nil, err
	}

	user := User{
		Email:       request.Email,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: request.PhoneNumber,
		Password:    &hashedPassword,
	}

	err = r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if len(request.RoleIDs) > 0 {
			if err := r.syncRoles(tx, &user, request.RoleIDs); err != nil {
				return err
			}
		}

		return tx.Preload("Roles").First(&user, "id = ?", user.ID).Error
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateById(c *gin.Context, id string, request UpdateRequest) (*User, error) {
	var user User

	err := r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors2.DataNotFoundException()
			}
			return err
		}

		user.Email = request.Email
		user.FirstName = request.FirstName
		user.LastName = request.LastName
		user.PhoneNumber = request.PhoneNumber

		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		if request.RoleIDs != nil {
			if err := r.syncRoles(tx, &user, request.RoleIDs); err != nil {
				return err
			}
		}

		return tx.Preload("Roles").First(&user, "id = ?", id).Error
	})
	if err != nil {
		return nil, err
	}

	if err := r.LoadAvatar(c, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) syncRoles(tx *gorm.DB, user *User, roleIDs []string) error {
	var roles []rbac.Role
	if len(roleIDs) > 0 {
		if err := tx.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return err
		}
	}
	return tx.Model(user).Association("Roles").Replace(&roles)
}

func (r *Repository) DeleteById(c *gin.Context, id string) error {
	var user User
	if err := r.db.WithContext(c).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors2.DataNotFoundException()
		}
		return err
	}

	if err := r.db.WithContext(c).Delete(&User{}, "id = ?", id).Error; err != nil {
		return err
	}

	// Media tidak punya foreign key ke users karena relasinya polymorphic,
	// jadi pembersihannya dilakukan di sini.
	return r.media.DeleteForModel(c, &user)
}
