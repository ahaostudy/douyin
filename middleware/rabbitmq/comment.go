package rabbitmq

import (
	"github.com/streadway/amqp"
	"main/config"
	"main/dao"
	"strconv"
)

// GenerateDelCommentMQParam 生成传入 DelComment MQ 的参数
func GenerateDelCommentMQParam(commentID uint) []byte {
	return []byte(strconv.FormatUint(uint64(commentID), 10))
}

// DelComment 删除评论
func DelComment(msg *amqp.Delivery) error {
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
