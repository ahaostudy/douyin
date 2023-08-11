package service

import (
	"fmt"
	"main/dao"
	"main/model"
)

func SendComment(uId uint, vId uint, commentText string) (*model.Comment, bool) {
	cId := dao.AddComment(uId, vId, commentText)
	fmt.Println(cId)
	comment, err := dao.SendComment(uId, cId)
	if err != nil {
		return nil, false
	}
	return comment, true
}

func DelComment(commentId uint) bool {
	dao.DelCommon(commentId)
	return true
}

func GetListComment(vId uint, uId uint) ([]*model.Comment, bool) {
	commentList, err := dao.GetListComment(vId, uId)
	if err != nil {
		return nil, false
	}
	return commentList, true
}
