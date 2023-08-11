package service

import "main/dao"

// RelationAction 关注操作
// uid: 用户ID, tid: 对方ID, t: 操作类型
func RelationAction(uid, tid uint, t int) bool {
	switch t {
	case 1:
		// 获取关注记录
		_, err := dao.GetFollow(uid, tid)
		if err == nil {
			return false
		}
		// 插入关注记录
		return dao.InsertFollow(uid, tid) == nil
	case 2:
		// 删除关注记录
		return dao.DeleteFollow(uid, tid) == nil
	default:
		return false
	}
}
