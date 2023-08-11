package model

import (
	"encoding/json"
	"time"
)

type Message struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	FromUserID uint      `gorm:"index" json:"from_user_id"`
	ToUserID   uint      `gorm:"index" json:"to_user_id"`
	Content    string    `gorm:"type:text" json:"content"`
	CreatedAt  time.Time `json:"create_time"`
}

func (m Message) MarshalJSON() ([]byte, error) {
	type MsgJSON Message
	return json.Marshal(&struct {
		MsgJSON
		CreatedAt int64 `json:"create_time"`
	}{
		MsgJSON:   (MsgJSON)(m),
		CreatedAt: m.CreatedAt.UnixMilli(),
	})
}
