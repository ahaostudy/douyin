package model

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
