// gorm model for the video table in the database

package model

import (
	"time"

	"github.com/gosimple/slug"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Uuid        string `gorm:"primaryKey"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	Url         string `gorm:"not null"`
	ChannelID   string
	Views       uint64
	Duration    int32
	Thumbnail   string
	PublishedAt time.Time
	Privacy     string
	Category    string
	Language    string
	Tags        []string `gorm:"type:jsonb"`
	isDeleted   bool
	Slug        string `gorm:"uniqueIndex"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// Hook before create to generate uuid and slug
func (u *Video) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.NewV4()
	u.Uuid = uuid.String()
	u.Slug = slug.Make(u.Title)
	u.Url = "https://www.showcase.com/watch?v=" + u.Slug
	return nil
}
