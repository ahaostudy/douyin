package model

import (
	"time"
)

type Comment struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"require" json:"user_id"`
	VideoID     uint      `gorm:"require;index" json:"video_id"`
	CommentText string    `json:"comment_text"`
	CreatedAt   time.Time `json:"created_at"`
}
