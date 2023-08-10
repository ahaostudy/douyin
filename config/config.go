package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var StaticDir string

// InitConfig 初始化项目配置
func InitConfig() error {
	workDir, err := os.Getwd()
	if err != nil {
		return err
	}
	viper.SetConfigType("yaml")
	viper.SetConfigFile(workDir + "/config/config.yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	StaticDir = fmt.Sprintf("http://%s:%s/static/", viper.GetString("server.host"), viper.GetString("server.port"))

	return nil
}
