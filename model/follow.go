package model

type Follow struct {
	ID         uint `gorm:"primarykey" json:"id"`
	UserID     uint `gorm:"index" json:"user_id"`
	FollowerID uint `gorm:"index" json:"follower_id"`
}
