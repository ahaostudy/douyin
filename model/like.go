package model

type Like struct {
	ID      uint `gorm:"primarykey" json:"id"`
	UserID  uint `gorm:"require" json:"user_id"`
	VideoID uint `gorm:"require;index" json:"video_id"`
}
