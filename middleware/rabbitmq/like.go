package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"main/dao"
	"strconv"
	"strings"
)

// GenerateLikeMQParam 生成传入 Like MQ 的参数
func GenerateLikeMQParam(uid, vid uint) []byte {
	return []byte(fmt.Sprintf("%d %d", uid, vid))
}

// GenerateUnLikeMQParam 生成传入 UnLike MQ 的参数
func GenerateUnLikeMQParam(uid, vid uint) []byte {
	return []byte(fmt.Sprintf("%d %d", uid, vid))
}

// Like 点赞业务
func Like(msg *amqp.Delivery) error {
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

// UnLike 取消点赞
func UnLike(msg *amqp.Delivery) error {
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
