package model

type User struct {
	ID              uint   `gorm:"primarykey" json:"id"`
	Name            string `json:"name"`
	Username        string `gorm:"require" json:"username"`
	Password        string `gorm:"require" json:"password"`
	Avatar          string `json:"avatar"`
	Signature       string `json:"signature"`
	BackgroundImage string `json:"background_image"`
}
