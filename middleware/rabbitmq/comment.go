package rabbitmq

import (
	"fmt"
	"main/config"
	"main/dao"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
)

// GenerateUnCommentMQParam 生成传入 UnComment MQ 的参数
func GenerateUnCommentMQParam(commentID uint) []byte {
	return []byte(strconv.FormatUint(uint64(commentID), 10))
}

// GenerateCommentMQParam 生成传入 Comment MQ 的参数
func GenerateCommentMQParam(uid, vid uint, commentText string) []byte {
	return []byte(fmt.Sprintf("%d %d %s", uid, vid, commentText))
}

// Comment 插入评论
func Comment(msg *amqp.Delivery) error {
	// 解析参数
	params := strings.Split(string(msg.Body), " ")
	_uid, err := strconv.ParseUint(params[0], 10, 32)
	if err != nil {
		return err
	}
	_vid, err := strconv.ParseUint(params[1], 10, 32)
	if err != nil {
		return err
	}
	uid, vid := uint(_uid), uint(_vid)
	commentText := params[2]

	// 操作数据库
	for i := 0; i < config.SQLMaxReTryCount; i++ {
		if _, e := dao.InsertComment(uid, vid, commentText); e != nil {
			err = e
		} else {
			break
		}
	}
	return err
}

// UnComment 删除评论
func UnComment(msg *amqp.Delivery) error {
	// 解析参数
	commentID, err := strconv.ParseUint(string(msg.Body), 10, 32)
	if err != nil {
		return err
	}

	// 操作数据库
	for i := 0; i < config.SQLMaxReTryCount; i++ {
		if e := dao.DeleteComment(uint(commentID)); e != nil {
			err = e
		} else {
			break
		}
	}

	return err
}
