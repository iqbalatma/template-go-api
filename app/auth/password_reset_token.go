package auth

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordResetToken struct {
	ID        uuid.UUID `gorm:"column:id;primary_key;type:uuid"`
	Email     string    `gorm:"column:email"`
	Token     string    `gorm:"column:token;unique"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (p *PasswordResetToken) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}
