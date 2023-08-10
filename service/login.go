package service

import (
	"fmt"
	"main/dao"
	"main/model"
	"main/utils"
)

// Login 登录业务
func Login(username, password string) (*model.User, bool) {
	// 查询用户是否存在
	user, err := dao.GetUserByUsername(username)
	if err != nil || user == nil {
		return nil, false
	}
	// 判断密码是否正确
	fmt.Println(user.Password, password)
	if user.Password != utils.MD5(password) {
		return nil, false
	}

	return user, true
}
