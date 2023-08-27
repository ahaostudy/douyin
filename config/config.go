package config

import (
	"github.com/spf13/viper"
	"os"
	"path"
	"time"
)

var (
	StaticPath string

	RedisKeyLock       = "lock"
	RedisKeyTTL        = 24 * time.Hour
	RedisKeyOfLike     = "like"
	RedisKeyOfAuthor   = "author"
	RedisKeyOfFollow   = "follow"
	RedisKeyOfFollower = "follower"
	RedisKeyOfUser     = "user"
	RedisKeyOfMessage  = "message"
	RedisKeyOfComment  = "comment"
	RedisValueOfNULL   = "NULL"

	SQLMaxReTryCount = 5
)

// InitConfig 初始化项目配置
func InitConfig() error {
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// 初始化配置文件
	viper.SetConfigType("yaml")
	viper.SetConfigFile(workDir + "/config/config.yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	// 初始化全局变量
	StaticPath = path.Join(workDir, viper.GetString("static.path"))

	return nil
}
