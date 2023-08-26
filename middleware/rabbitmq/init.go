package rabbitmq

var (
	RMQLike     *RabbitMQ
	RMQUnLike   *RabbitMQ
	RMQFollow   *RabbitMQ
	RMQUnFollow *RabbitMQ
)

// InitRabbitMQ 初始化RabbitMQ
func InitRabbitMQ() {
	// 创建MQ并启动消费者
	// 无论调用多少次 NewWorkRabbitMQ，只会创建一次连接
	// 不同队列共用一个连接，可以保持不同队列消费消息的顺序

	RMQLike = NewWorkRabbitMQ("like")
	go RMQLike.Consume(Like) // 传入一个业务函数，用于消费消息

	RMQUnLike = NewWorkRabbitMQ("unlike")
	go RMQUnLike.Consume(UnLike)

	RMQFollow = NewWorkRabbitMQ("follow")
	go RMQFollow.Consume(Follow)

	RMQUnFollow = NewWorkRabbitMQ("unfollow")
	go RMQUnFollow.Consume(UnFollow)
}

// DestroyRabbitMQ 销毁RabbitMQ
func DestroyRabbitMQ() {
	RMQLike.Destroy()
	RMQUnLike.Destroy()
	RMQFollow.Destroy()
	RMQUnFollow.Destroy()
}
