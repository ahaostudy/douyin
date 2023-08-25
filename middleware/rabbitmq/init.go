package rabbitmq

var (
	RMQLikeAdd *RabbitMQ
	RMQLikeDel *RabbitMQ
)

// InitRabbitMQ 初始化RabbitMQ
func InitRabbitMQ() {
	// 创建MQ并启动消费者
	// 无论调用多少次 NewWorkRabbitMQ，只会创建一次连接
	// 不同队列共用一个连接，可以保持不同队列消费消息的顺序

	// like_add
	RMQLikeAdd = NewWorkRabbitMQ("like_add")
	go RMQLikeAdd.Consume(LikeAdd) // 传入一个业务函数，用于消费消息
	// like_del
	RMQLikeDel = NewWorkRabbitMQ("like_del")
	go RMQLikeDel.Consume(LikeDel)
}

// DestroyRabbitMQ 销毁RabbitMQ
func DestroyRabbitMQ() {
	RMQLikeAdd.Destroy()
	RMQLikeDel.Destroy()
}
