package media

import (
	goerrors "errors"
	"io"
	"mime/multipart"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	errors2 "template-go-api/errors"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// sniffLength adalah jumlah byte awal yang dibaca untuk mendeteksi mime type
// dari isi file, bukan dari header Content-Type yang dikirim client.
const sniffLength = 3072

var unsafeFileNameChars = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

type attachOptions struct {
	name             string
	customProperties JSONMap
}

type AttachOption func(*attachOptions)

// WithName mengganti nama tampilan media. Default-nya nama file asli.
func WithName(name string) AttachOption {
	return func(o *attachOptions) {
		o.name = name
	}
}

// WithCustomProperties menyimpan metadata bebas pada kolom custom_properties.
func WithCustomProperties(properties JSONMap) AttachOption {
	return func(o *attachOptions) {
		o.customProperties = properties
	}
}

// Attach menyimpan file yang diunggah ke sebuah collection milik owner.
// Alurnya: validasi aturan collection, tulis file asli, generate conversion,
// simpan record, lalu bersihkan file lama bila collection bersifat single file.
func (r *Repository) Attach(
	c *gin.Context,
	owner HasMedia,
	collectionName string,
	fileHeader *multipart.FileHeader,
	options ...AttachOption,
) (*Media, error) {
	modelType := owner.GetMediaModelType()

	collection, httpErr := collections.Get(modelType, collectionName)
	if httpErr != nil {
		return nil, httpErr
	}

	opts := attachOptions{}
	for _, option := range options {
		option(&opts)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	mimeType, err := detectMimeType(file)
	if err != nil {
		return nil, err
	}

	if httpErr := collection.validate(fileHeader.Size, mimeType); httpErr != nil {
		return nil, httpErr
	}

	disk, err := disks.Disk(collection.Disk)
	if err != nil {
		return nil, err
	}

	fileName := sanitizeFileName(fileHeader.Filename)
	name := opts.name
	if name == "" {
		name = strings.TrimSuffix(fileHeader.Filename, filepath.Ext(fileHeader.Filename))
	}

	newMedia := Media{
		ModelType:        modelType,
		ModelID:          owner.GetMediaModelID(),
		CollectionName:   collection.Name,
		Disk:             disk.Name(),
		FileName:         fileName,
		Name:             name,
		MimeType:         mimeType,
		Size:             fileHeader.Size,
		CustomProperties: opts.customProperties,
	}
	// ID dibuat di awal karena dipakai sebagai nama directory file.
	newMedia.GenerateUUID()
	newMedia.Directory = path.Join(modelType, newMedia.ID.String())

	if err := disk.Put(newMedia.Path(), file); err != nil {
		return nil, err
	}

	// Mulai titik ini file sudah ada di disk, jadi setiap kegagalan harus
	// membersihkan directory yang baru dibuat.
	cleanup := func(cause error) error {
		_ = disk.DeleteDirectory(newMedia.Directory)
		return cause
	}

	conversions, err := generateConversions(disk, &newMedia, collection.Conversions)
	if err != nil {
		return nil, cleanup(err)
	}
	newMedia.Conversions = conversions

	var replaced []Media
	err = r.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if collection.SingleFile {
			if err := tx.Where(
				"model_type = ? AND model_id = ? AND collection_name = ?",
				newMedia.ModelType, newMedia.ModelID, newMedia.CollectionName,
			).Find(&replaced).Error; err != nil {
				return err
			}
		} else {
			order, err := r.nextOrderColumn(tx, newMedia.ModelType, newMedia.ModelID, newMedia.CollectionName)
			if err != nil {
				return err
			}
			newMedia.OrderColumn = order
		}

		if err := tx.Create(&newMedia).Error; err != nil {
			return err
		}

		if len(replaced) > 0 {
			ids := make([]string, len(replaced))
			for i, old := range replaced {
				ids[i] = old.ID.String()
			}
			if err := tx.Where("id IN ?", ids).Delete(&Media{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, cleanup(err)
	}

	// File lama baru dihapus setelah transaksi sukses. Kegagalan di sini hanya
	// menyisakan file yatim, tidak membatalkan upload yang sudah tercatat.
	for i := range replaced {
		r.removeFiles(&replaced[i])
	}

	return &newMedia, nil
}

// GetCollection mengembalikan seluruh media pada sebuah collection, terurut.
func (r *Repository) GetCollection(c *gin.Context, owner HasMedia, collectionName string) ([]Media, error) {
	var items []Media
	err := r.db.WithContext(c).
		Where(
			"model_type = ? AND model_id = ? AND collection_name = ?",
			owner.GetMediaModelType(), owner.GetMediaModelID(), collectionName,
		).
		Order("order_column ASC, created_at ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetFirst mengembalikan media pertama pada sebuah collection, atau nil bila kosong.
func (r *Repository) GetFirst(c *gin.Context, owner HasMedia, collectionName string) (*Media, error) {
	var item Media
	err := r.db.WithContext(c).
		Where(
			"model_type = ? AND model_id = ? AND collection_name = ?",
			owner.GetMediaModelType(), owner.GetMediaModelID(), collectionName,
		).
		Order("order_column ASC, created_at ASC").
		First(&item).Error
	if err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetForModels mengambil media beberapa model sekaligus dan memetakannya per
// model id, dipakai untuk menghindari query N+1 pada endpoint listing.
func (r *Repository) GetForModels(c *gin.Context, modelType string, modelIDs []string, collectionName string) (map[string][]Media, error) {
	result := make(map[string][]Media)
	if len(modelIDs) == 0 {
		return result, nil
	}

	var items []Media
	err := r.db.WithContext(c).
		Where(
			"model_type = ? AND model_id IN ? AND collection_name = ?",
			modelType, modelIDs, collectionName,
		).
		Order("order_column ASC, created_at ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		result[item.ModelID] = append(result[item.ModelID], item)
	}
	return result, nil
}

// DeleteByID menghapus satu media beserta filenya. ownerScope boleh nil bila
// pemanggil sudah memastikan kepemilikan.
func (r *Repository) DeleteByID(c *gin.Context, id string, ownerScope HasMedia) error {
	var item Media

	query := r.db.WithContext(c).Where("id = ?", id)
	if ownerScope != nil {
		query = query.Where(
			"model_type = ? AND model_id = ?",
			ownerScope.GetMediaModelType(), ownerScope.GetMediaModelID(),
		)
	}

	if err := query.First(&item).Error; err != nil {
		if goerrors.Is(err, gorm.ErrRecordNotFound) {
			return errors2.DataNotFoundException("Media not found")
		}
		return err
	}

	if err := r.db.WithContext(c).Delete(&Media{}, "id = ?", item.ID).Error; err != nil {
		return err
	}

	r.removeFiles(&item)
	return nil
}

// ClearCollection menghapus seluruh media pada satu collection milik owner.
func (r *Repository) ClearCollection(c *gin.Context, owner HasMedia, collectionName string) error {
	items, err := r.GetCollection(c, owner, collectionName)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.ID.String()
	}

	if err := r.db.WithContext(c).Where("id IN ?", ids).Delete(&Media{}).Error; err != nil {
		return err
	}

	for i := range items {
		r.removeFiles(&items[i])
	}
	return nil
}

// DeleteForModel membersihkan seluruh media milik sebuah model, dipakai saat
// model induknya dihapus.
func (r *Repository) DeleteForModel(c *gin.Context, owner HasMedia) error {
	var items []Media
	err := r.db.WithContext(c).
		Where("model_type = ? AND model_id = ?", owner.GetMediaModelType(), owner.GetMediaModelID()).
		Find(&items).Error
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	err = r.db.WithContext(c).
		Where("model_type = ? AND model_id = ?", owner.GetMediaModelType(), owner.GetMediaModelID()).
		Delete(&Media{}).Error
	if err != nil {
		return err
	}

	for i := range items {
		r.removeFiles(&items[i])
	}
	return nil
}

func (r *Repository) nextOrderColumn(tx *gorm.DB, modelType, modelID, collectionName string) (int, error) {
	var current int
	err := tx.Model(&Media{}).
		Where(
			"model_type = ? AND model_id = ? AND collection_name = ?",
			modelType, modelID, collectionName,
		).
		Select("COALESCE(MAX(order_column), 0)").
		Scan(&current).Error
	if err != nil {
		return 0, err
	}
	return current + 1, nil
}

// removeFiles menghapus directory milik satu media secara best effort.
func (r *Repository) removeFiles(m *Media) {
	disk, err := disks.Disk(m.Disk)
	if err != nil {
		return
	}
	_ = disk.DeleteDirectory(m.Directory)
}

// detectMimeType membaca sebagian awal file lalu mengembalikan posisi baca ke
// awal supaya file tetap utuh saat ditulis ke disk.
func detectMimeType(file multipart.File) (string, error) {
	head := make([]byte, sniffLength)
	n, err := file.Read(head)
	if err != nil && !goerrors.Is(err, io.EOF) {
		return "", err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}
	return mimetype.Detect(head[:n]).String(), nil
}

// sanitizeFileName membuang karakter yang tidak aman untuk nama file di disk.
func sanitizeFileName(original string) string {
	base := filepath.Base(original)
	ext := strings.ToLower(filepath.Ext(base))
	name := strings.TrimSuffix(base, filepath.Ext(base))

	name = unsafeFileNameChars.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-.")
	if name == "" {
		name = "file"
	}
	if len(name) > 100 {
		name = name[:100]
	}

	ext = unsafeFileNameChars.ReplaceAllString(ext, "")
	return name + ext
}
