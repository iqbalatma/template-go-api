package utils

import (
	"fmt"
	"math"
	"time"
)

const (
	dateFormat     = "2006-01-02"
	dateTimeFormat = "02/01/2006 15:04:05"
)

func FormatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(dateFormat)
	return &s
}

func FormatDateTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(dateTimeFormat)
	return &s
}

func FormatDateTimeVal(t time.Time) string {
	return t.Format(dateTimeFormat)
}

func HumanFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}

	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGT"[exp])
}

func FormatRupiah(amount float64) string {
	rounded := math.Round(amount)
	negative := rounded < 0
	if negative {
		rounded = -rounded
	}

	intPart := int64(rounded)
	result := ""
	for i, digit := range fmt.Sprintf("%d", intPart) {
		if i > 0 && (len(fmt.Sprintf("%d", intPart))-i)%3 == 0 {
			result += "."
		}
		result += string(digit)
	}

	if negative {
		return "Rp -" + result
	}
	return "Rp " + result
}
