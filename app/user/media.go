package user

import (
	"template-go-api/app/media"
)

const MediaCollectionAvatar = "avatar"

// avatarCollection: satu file per user, hanya gambar, maksimal 2 MB.
var avatarCollection = media.Collection{
	Name:       MediaCollectionAvatar,
	SingleFile: true,
	MaxSize:    2 << 20,
	AcceptedMimeTypes: []string{
		"image/jpeg",
		"image/png",
		"image/webp",
	},
	Conversions: []media.Conversion{
		{Name: "thumb", Width: 150, Height: 150, Fit: media.FitCover, Format: "jpg"},
		{Name: "preview", Width: 600, Height: 600, Fit: media.FitContain, Format: "jpg"},
	},
}

// RegisterMediaCollections mendaftarkan collection milik User.
// Dipanggil sekali saat container dibangun.
func RegisterMediaCollections() {
	media.Register(MediaModelType, avatarCollection)
}
