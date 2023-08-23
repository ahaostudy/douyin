package utils

import (
	"main/model"
)

// 用户脱敏信息
type ExclusivePwdUser struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	Avatar          string `json:"avatar"`
	Signature       string `json:"signature"`
	BackgroundImage string `json:"background_image"`

	FollowCount    int  `json:"follow_count"`
	FollowerCount  int  `json:"follower_count"`
	IsFollow       bool `json:"is_follow"`
	TotalFavorited int  `json:"total_favorited"`
	WorkCount      int  `json:"work_count"`
	FavoriteCount  int  `json:"favorite_count"`
}

func GetSaftUser(user *model.User) ExclusivePwdUser {
	exclusiveUser := ExclusivePwdUser{
		ID:              user.ID,
		Name:            user.Username,
		Avatar:          user.Avatar,
		Signature:       user.Signature,
		BackgroundImage: user.BackgroundImage,
		FollowCount:     user.FollowCount,
		FollowerCount:   user.FollowerCount,
		IsFollow:        user.IsFollow,
		TotalFavorited:  user.TotalFavorited,
		WorkCount:       user.WorkCount,
		FavoriteCount:   user.FavoriteCount,
	}
	return exclusiveUser
}
