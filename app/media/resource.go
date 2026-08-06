package media

import (
	"template-go-api/utils"
)

type Resource struct {
	Id               string            `json:"id"`
	CollectionName   string            `json:"collection_name"`
	Name             string            `json:"name"`
	FileName         string            `json:"file_name"`
	MimeType         string            `json:"mime_type"`
	Size             int64             `json:"size"`
	HumanSize        string            `json:"human_size"`
	URL              string            `json:"url"`
	Conversions      map[string]string `json:"conversions"`
	CustomProperties JSONMap           `json:"custom_properties"`
	OrderColumn      int               `json:"order_column"`
	CreatedAt        string            `json:"created_at"`
}

func NewResource(m *Media) *Resource {
	if m == nil {
		return nil
	}

	conversions := make(map[string]string, len(m.Conversions))
	for name := range m.Conversions {
		if url := ConversionURLOf(m, name); url != "" {
			conversions[name] = url
		}
	}

	return &Resource{
		Id:               m.ID.String(),
		CollectionName:   m.CollectionName,
		Name:             m.Name,
		FileName:         m.FileName,
		MimeType:         m.MimeType,
		Size:             m.Size,
		HumanSize:        utils.HumanFileSize(m.Size),
		URL:              URLOf(m),
		Conversions:      conversions,
		CustomProperties: m.CustomProperties,
		OrderColumn:      m.OrderColumn,
		CreatedAt:        utils.FormatDateTimeVal(m.CreatedAt),
	}
}

func NewResourceCollection(items []Media) []Resource {
	result := make([]Resource, len(items))
	for i := range items {
		result[i] = *NewResource(&items[i])
	}
	return result
}
