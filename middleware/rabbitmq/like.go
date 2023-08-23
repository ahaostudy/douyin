package rabbitmq

import (
	"github.com/streadway/amqp"
	"main/dao"
	"strconv"
	"strings"
)

// LikeAdd 点赞业务
func LikeAdd(msg *amqp.Delivery) error {
	// TODO 校验参数是否正确、dao层调用失败要重试
	params := strings.Split(string(msg.Body), " ")
	_uid, _ := strconv.ParseUint(params[0], 10, 32)
	_vid, _ := strconv.ParseUint(params[1], 10, 32)
	uid, vid := uint(_uid), uint(_vid)

	if err := dao.InsertLike(uid, vid); err != nil {
		return err
	}

	return nil
}

// LikeDel 取消点赞
func LikeDel(msg *amqp.Delivery) error {
	// TODO 校验参数是否正确、dao层调用失败要重试
	params := strings.Split(string(msg.Body), " ")
	_uid, _ := strconv.ParseUint(params[0], 10, 32)
	_vid, _ := strconv.ParseUint(params[1], 10, 32)
	uid, vid := uint(_uid), uint(_vid)

	if err := dao.DeleteLike(uid, vid); err != nil {
		return err
	}

	return nil
}
