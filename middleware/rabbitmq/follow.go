package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

// GenerateFollowMQParam 生成传入 Follow MQ 的参数
func GenerateFollowMQParam(uid, tid uint) []byte {
	return []byte(fmt.Sprintf("%d %d", uid, tid))
}

// GenerateUnFollowMQParam 生成传入 UnFollow MQ 的参数
func GenerateUnFollowMQParam(uid, tid uint) []byte {
	return []byte(fmt.Sprintf("%d %d", uid, tid))
}

func Follow(msg *amqp.Delivery) error {
	return nil
}

func UnFollow(msg *amqp.Delivery) error {
	return nil
}
