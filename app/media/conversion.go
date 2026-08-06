package media

import (
	"bytes"
	"fmt"
	"image"
	"path"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"

	// Mendaftarkan decoder WEBP ke image.Decode supaya file webp bisa
	// dibuatkan conversion. Output conversion sendiri tetap jpg/png.
	_ "golang.org/x/image/webp"
)

// conversionsDirectory adalah sub folder tempat hasil conversion disimpan,
// relatif terhadap directory milik satu record media.
const conversionsDirectory = "conversions"

type FitMode string

const (
	// FitContain memperkecil gambar secara proporsional sampai muat di dalam
	// kotak Width x Height. Rasio asli dipertahankan, tidak ada bagian terpotong.
	FitContain FitMode = "contain"
	// FitCover memotong bagian tengah gambar agar tepat mengisi kotak
	// Width x Height. Membutuhkan Width dan Height keduanya terisi.
	FitCover FitMode = "cover"
)

// Conversion adalah turunan gambar yang digenerate saat file diunggah.
type Conversion struct {
	Name   string
	Width  int
	Height int
	Fit    FitMode
	// Quality untuk output JPEG, 1-100. Nilai 0 memakai 85.
	Quality int
	// Format output ("jpg", "png", "gif", ...). Kosong berarti mengikuti
	// format file asli, dengan fallback ke jpg bila formatnya tidak didukung.
	Format string
}

func (c Conversion) quality() int {
	if c.Quality <= 0 || c.Quality > 100 {
		return 85
	}
	return c.Quality
}

// resolveFormat menentukan format encoding beserta ekstensi filenya.
func (c Conversion) resolveFormat(originalFileName string) (imaging.Format, string) {
	if c.Format != "" {
		if format, err := imaging.FormatFromExtension(c.Format); err == nil {
			return format, strings.ToLower(strings.TrimPrefix(c.Format, "."))
		}
	}

	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(originalFileName), "."))
	if format, err := imaging.FormatFromExtension(ext); err == nil {
		return format, ext
	}

	return imaging.JPEG, "jpg"
}

func (c Conversion) apply(src image.Image) image.Image {
	switch c.Fit {
	case FitCover:
		return imaging.Fill(src, c.Width, c.Height, imaging.Center, imaging.Lanczos)
	default:
		// imaging.Fit menjaga rasio dan tidak memperbesar gambar yang sudah kecil.
		return imaging.Fit(src, c.Width, c.Height, imaging.Lanczos)
	}
}

func (c Conversion) validate() error {
	if c.Name == "" {
		return fmt.Errorf("media: conversion name is required")
	}
	if c.Width <= 0 && c.Height <= 0 {
		return fmt.Errorf("media: conversion %q needs width or height", c.Name)
	}
	if c.Fit == FitCover && (c.Width <= 0 || c.Height <= 0) {
		return fmt.Errorf("media: conversion %q with fit cover needs both width and height", c.Name)
	}
	return nil
}

// generateConversions membaca file asli dari disk lalu menulis setiap turunan
// ke {directory}/conversions/{name}.{ext}. Nilai kembaliannya dipetakan ke
// kolom conversions: {"thumb": "thumb.jpg"}.
func generateConversions(disk Disk, m *Media, conversions []Conversion) (JSONMap, error) {
	if len(conversions) == 0 || !m.IsImage() {
		return nil, nil
	}

	source, err := disk.Open(m.Path())
	if err != nil {
		return nil, err
	}
	defer source.Close()

	// AutoOrientation menerapkan rotasi EXIF supaya thumbnail tidak terbalik.
	original, err := imaging.Decode(source, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("media: cannot decode image %s: %w", m.FileName, err)
	}

	generated := JSONMap{}
	for _, conversion := range conversions {
		if err := conversion.validate(); err != nil {
			return nil, err
		}

		format, ext := conversion.resolveFormat(m.FileName)

		var buffer bytes.Buffer
		if err := imaging.Encode(
			&buffer,
			conversion.apply(original),
			format,
			imaging.JPEGQuality(conversion.quality()),
		); err != nil {
			return nil, err
		}

		fileName := conversion.Name + "." + ext
		if err := disk.Put(path.Join(m.Directory, conversionsDirectory, fileName), &buffer); err != nil {
			return nil, err
		}
		generated[conversion.Name] = fileName
	}

	return generated, nil
}
