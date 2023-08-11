package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Comment struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	VideoID     uint      `gorm:"index" json:"video_id"`
	CommentText string    `json:"content"`
	CreatedAt   time.Time `json:"create_date"`

	User User `gorm:"-:migration;<-:false" json:"user"`
}

type DateTime time.Time

func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(t).Format("01-02"))), nil
}

func (c Comment) MarshalJSON() ([]byte, error) {
	type CmtJSON Comment
	return json.Marshal(&struct {
		CmtJSON
		CreatedAt DateTime `json:"create_date"`
	}{
		CmtJSON:   (CmtJSON)(c),
		CreatedAt: DateTime(c.CreatedAt),
	})
}
