package config

import (
	"github.com/spf13/viper"
	"os"
)

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

	return nil
}
