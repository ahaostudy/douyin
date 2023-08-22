package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

// LikeAdd 点赞业务
func LikeAdd(msg *amqp.Delivery) error {
	// TODO
	fmt.Println("like add :", string(msg.Body))

	return nil
}
