package service

import (
	"main/dao"
)

// FavoriteAction 赞业务
// uid, vid 分别对应用户ID和视频ID
// t 表示操作类型
func FavoriteAction(uid, vid uint, t int) bool {
	switch t {
	case 1:
		// 获取点赞记录
		_, err := dao.GetLike(uid, vid)
		if err == nil {
			return false
		}
		// 插入点赞记录
		return dao.InsertLike(uid, vid) == nil
	case 2:
		// 删除点赞记录
		return dao.DeleteLike(uid, vid) == nil
	default:
		return false
	}
}
