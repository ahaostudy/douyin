package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"main/config"
	"main/dao"
	"strconv"
	"strings"
)

// GenerateFollowMQParam 生成传入 Follow MQ 的参数
func GenerateFollowMQParam(uid, tid uint) []byte {
	return []byte(fmt.Sprintf("%d %d", uid, tid))
}

// GenerateUnFollowMQParam 生成传入 UnFollow MQ 的参数
func GenerateUnFollowMQParam(uid, tid uint) []byte {
	return []byte(fmt.Sprintf("%d %d", uid, tid))
}

// Follow 关注
func Follow(msg *amqp.Delivery) error {
	// 解析参数
	params := strings.Split(string(msg.Body), " ")
	_uid, err := strconv.ParseUint(params[0], 10, 32)
	if err != nil {
		return err
	}
	_tid, err := strconv.ParseUint(params[1], 10, 32)
	if err != nil {
		return err
	}
	uid, tid := uint(_uid), uint(_tid)

	// 操作数据库
	for i := 0; i < config.SQLMaxReTryCount; i++ {
		if e := dao.InsertFollow(uid, tid); e != nil {
			err = e
		} else {
			break
		}
	}

	return err
}

// UnFollow 取消关注
func UnFollow(msg *amqp.Delivery) error {
	// 解析参数
	params := strings.Split(string(msg.Body), " ")
	_uid, err := strconv.ParseUint(params[0], 10, 32)
	if err != nil {
		return err
	}
	_tid, err := strconv.ParseUint(params[1], 10, 32)
	if err != nil {
		return err
	}
	uid, tid := uint(_uid), uint(_tid)

	// 操作数据库
	for i := 0; i < config.SQLMaxReTryCount; i++ {
		if e := dao.DeleteFollow(uid, tid); e != nil {
			err = e
		} else {
			break
		}
	}

	return err
}
