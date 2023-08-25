package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"main/dao"
	"strconv"
	"strings"
)

// GenerateLikeAddMQParam 生成传入 LikeAdd MQ 的参数
func GenerateLikeAddMQParam(uid, vid uint) []byte {
	return []byte(fmt.Sprintf("%d %d", uid, vid))
}

// GenerateLikeDelMQParam 生成传入 LikeDel MQ 的参数
func GenerateLikeDelMQParam(uid, vid uint) []byte {
	return []byte(fmt.Sprintf("%d %d", uid, vid))
}

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
