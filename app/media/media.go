package media

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"path"
	"strings"
	"template-go-api/app"
)

// HasMedia dipenuhi oleh model apa pun yang ingin punya file terlampir.
// Pola sama seperti model_type/model_id pada spatie/laravel-medialibrary.
type HasMedia interface {
	GetMediaModelType() string
	GetMediaModelID() string
}

// JSONMap dipetakan ke kolom JSON MySQL.
type JSONMap map[string]any

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value any) error {
	if value == nil {
		*m = nil
		return nil
	}

	var source []byte
	switch v := value.(type) {
	case []byte:
		source = v
	case string:
		source = []byte(v)
	default:
		return errors.New("media: cannot scan non-json value into JSONMap")
	}

	if len(source) == 0 {
		*m = nil
		return nil
	}
	return json.Unmarshal(source, m)
}

type Media struct {
	app.BaseModel    `gorm:"embedded"`
	ModelType        string  `json:"model_type" gorm:"column:model_type"`
	ModelID          string  `json:"model_id" gorm:"column:model_id"`
	CollectionName   string  `json:"collection_name" gorm:"column:collection_name"`
	Disk             string  `json:"disk" gorm:"column:disk"`
	Directory        string  `json:"directory" gorm:"column:directory"`
	FileName         string  `json:"file_name" gorm:"column:file_name"`
	Name             string  `json:"name" gorm:"column:name"`
	MimeType         string  `json:"mime_type" gorm:"column:mime_type"`
	Size             int64   `json:"size" gorm:"column:size"`
	Conversions      JSONMap `json:"conversions" gorm:"column:conversions;type:json"`
	CustomProperties JSONMap `json:"custom_properties" gorm:"column:custom_properties;type:json"`
	OrderColumn      int     `json:"order_column" gorm:"column:order_column"`
}

func (Media) TableName() string {
	return "media"
}

// Path mengembalikan lokasi file asli relatif terhadap root disk.
func (m *Media) Path() string {
	return path.Join(m.Directory, m.FileName)
}

// ConversionPath mengembalikan lokasi hasil conversion relatif terhadap root disk.
// Nilai kosong berarti conversion tersebut belum pernah digenerate.
func (m *Media) ConversionPath(conversion string) string {
	if m.Conversions == nil {
		return ""
	}
	fileName, ok := m.Conversions[conversion].(string)
	if !ok || fileName == "" {
		return ""
	}
	return path.Join(m.Directory, conversionsDirectory, fileName)
}

func (m *Media) IsImage() bool {
	return strings.HasPrefix(m.MimeType, "image/")
}
