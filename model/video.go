package model

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"time"
)

type Video struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	AuthorID  uint      `gorm:"index" json:"author_id"`
	PlayUrl   string    `json:"play_url"`
	CoverUrl  string    `json:"cover_url"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`

	Author        *User `gorm:"-:migration;<-:false" json:"author"`
	FavoriteCount int   `gorm:"-:migration;<-:false" json:"favorite_count"`
	CommentCount  int   `gorm:"-:migration;<-:false" json:"comment_count"`
	IsFavorite    bool  `gorm:"-:migration;<-:false" json:"is_favorite"`
}

func (v *Video) AfterFind(tx *gorm.DB) error {
	staticDir := viper.GetString("server.static")
	v.PlayUrl = staticDir + v.PlayUrl
	v.CoverUrl = staticDir + v.CoverUrl
	return nil
}
