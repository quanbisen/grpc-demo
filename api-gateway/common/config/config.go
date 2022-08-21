package config

import (
	"github.com/spf13/viper"
	"os"
)

func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + string(os.PathSeparator) + "config")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
