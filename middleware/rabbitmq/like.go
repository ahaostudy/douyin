package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"main/config"
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

	// 操作数据库
	for i := 0; i < config.SQLMaxReTryCount; i++ {
		if e := dao.InsertLike(uid, vid); e != nil {
			err = e
		} else {
			break
		}
	}

	return err
}

// UnLike 取消点赞
func UnLike(msg *amqp.Delivery) error {
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

	// 操作数据库
	for i := 0; i < config.SQLMaxReTryCount; i++ {
		if e := dao.DeleteLike(uid, vid); e != nil {
			err = e
		} else {
			break
		}
	}

	return nil
}
