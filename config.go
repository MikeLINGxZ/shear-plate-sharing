package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Role     string
	Host     string
	Port     string
	Password string
}

var config Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(fmt.Errorf("fatal to unmarshal config: %w", err))
	}
}
