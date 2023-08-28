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

func (t *DateTime) UnmarshalJSON(data []byte) error {
	// 假设日期格式是 "01-02"
	parsedTime, err := time.Parse(fmt.Sprintf("\"%s\"", "01-02"), string(data))
	if err != nil {
		return err
	}
	*t = DateTime(parsedTime)
	return nil
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
func (c *Comment) UnmarshalJSON(data []byte) error {
	type CmtJSON Comment
	aux := &struct {
		CmtJSON
		CreatedAt DateTime `json:"create_date"`
	}{
		CmtJSON: (CmtJSON)(*c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.ID = aux.ID
	c.UserID = aux.UserID
	c.VideoID = aux.VideoID
	c.CommentText = aux.CommentText
	c.CreatedAt = time.Time(aux.CreatedAt)
	c.User = aux.User
	return nil
}
