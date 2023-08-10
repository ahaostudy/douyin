package model

type Like struct {
	ID      uint `gorm:"primarykey" json:"id"`
	UserID  uint `gorm:"index" json:"user_id"`
	VideoID uint `gorm:"index" json:"video_id"`
}
