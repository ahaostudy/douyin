package main

import (
	"fmt"
	"github.com/spf13/viper"
	"main/config"
	"main/dao"
	"main/middleware/rabbitmq"
	"main/middleware/redis"
	"main/router"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	if err := dao.InitMySQL(); err != nil {
		panic(err)
	}

	redis.InitRedis()

	rabbitmq.InitRabbitMQ()
	defer rabbitmq.DestroyRabbitMQ()

	r := router.InitRouter()
	if err := r.Run(fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))); err != nil {
		panic(err)
	}
}
