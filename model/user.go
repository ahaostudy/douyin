package model

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type User struct {
	ID              uint   `gorm:"primarykey" json:"id"`
	Name            string `json:"name"`
	Username        string `gorm:"unique" json:"username"`
	Password        string `json:"password"`
	Avatar          string `json:"avatar"`
	Signature       string `json:"signature"`
	BackgroundImage string `json:"background_image"`

	FollowCount    int  `gorm:"-:migration;<-:false" json:"follow_count"`
	FollowerCount  int  `gorm:"-:migration;<-:false" json:"follower_count"`
	IsFollow       bool `gorm:"-:migration;<-:false" json:"is_follow"`
	TotalFavorited int  `gorm:"-:migration;<-:false" json:"total_favorited"`
	WorkCount      int  `gorm:"-:migration;<-:false" json:"work_count"`
	FavoriteCount  int  `gorm:"-:migration;<-:false" json:"favorite_count"`
}

func (u *User) AfterFind(tx *gorm.DB) error {
	staticDir := viper.GetString("server.static")
	u.Avatar = staticDir + u.Avatar
	u.BackgroundImage = staticDir + u.BackgroundImage
	return nil
}
