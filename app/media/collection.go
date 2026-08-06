package media

import (
	"fmt"
	"strings"
	"sync"
	errors2 "template-go-api/errors"
	"template-go-api/utils"
)

// Collection mendefinisikan aturan sebuah koleksi file pada satu model,
// analog dengan registerMediaCollections() di spatie/laravel-medialibrary.
type Collection struct {
	Name string
	// Disk kosong berarti memakai disk default.
	Disk string
	// SingleFile true berarti file lama otomatis dihapus saat file baru masuk.
	SingleFile bool
	// MaxSize dalam byte. Nilai 0 berarti tanpa batas.
	MaxSize int64
	// AcceptedMimeTypes kosong berarti semua mime diterima. Mendukung
	// wildcard sufiks seperti "image/*".
	AcceptedMimeTypes []string
	// Conversions digenerate sinkron saat file diunggah.
	Conversions []Conversion
}

func (c Collection) accepts(mimeType string) bool {
	if len(c.AcceptedMimeTypes) == 0 {
		return true
	}
	for _, accepted := range c.AcceptedMimeTypes {
		if strings.HasSuffix(accepted, "/*") {
			if strings.HasPrefix(mimeType, strings.TrimSuffix(accepted, "*")) {
				return true
			}
			continue
		}
		if strings.EqualFold(accepted, mimeType) {
			return true
		}
	}
	return false
}

// validate memeriksa file terhadap aturan collection dan mengembalikan
// HTTPError agar bisa langsung dikembalikan ke client.
func (c Collection) validate(size int64, mimeType string) *utils.HTTPError {
	if c.MaxSize > 0 && size > c.MaxSize {
		return errors2.InvalidAction(fmt.Sprintf(
			"File terlalu besar, maksimal %s untuk collection %s",
			utils.HumanFileSize(c.MaxSize),
			c.Name,
		))
	}
	if !c.accepts(mimeType) {
		return errors2.InvalidAction(fmt.Sprintf(
			"Tipe file %s tidak diizinkan untuk collection %s, yang diizinkan: %s",
			mimeType,
			c.Name,
			strings.Join(c.AcceptedMimeTypes, ", "),
		))
	}
	return nil
}

// Registry menyimpan definisi collection per model type. Setiap package model
// mendaftarkan collection miliknya sendiri sehingga package media tidak perlu
// tahu apa pun tentang model tersebut.
type Registry struct {
	mu          sync.RWMutex
	collections map[string]Collection
}

func NewRegistry() *Registry {
	return &Registry{
		collections: make(map[string]Collection),
	}
}

func registryKey(modelType, collectionName string) string {
	return modelType + ":" + collectionName
}

func (r *Registry) Register(modelType string, collections ...Collection) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, collection := range collections {
		r.collections[registryKey(modelType, collection.Name)] = collection
	}
}

// Get mengembalikan definisi collection. Collection yang tidak terdaftar
// ditolak supaya client tidak bisa mengunggah ke collection sembarangan.
func (r *Registry) Get(modelType, collectionName string) (Collection, *utils.HTTPError) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	collection, ok := r.collections[registryKey(modelType, collectionName)]
	if !ok {
		return Collection{}, errors2.InvalidAction(fmt.Sprintf(
			"Collection %s tidak terdaftar untuk %s",
			collectionName,
			modelType,
		))
	}
	return collection, nil
}
