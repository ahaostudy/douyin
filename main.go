package main

import (
	"main/config"
	"main/dao"
	"main/router"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err)
	}

	if err := dao.InitMySQL(); err != nil {
		panic(err)
	}

	r := router.InitRouter()
	if err := r.Run(); err != nil {
		panic(err)
	}
}
