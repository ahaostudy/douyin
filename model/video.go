package model

import (
	"time"
)

type Video struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	AuthorID  uint      `gorm:"index" json:"author_id"`
	PlayUrl   string    `json:"play_url"`
	CoverUrl  string    `json:"cover_url"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}
