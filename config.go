package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type RoleType string

const RTServer RoleType = "server"
const RTClient RoleType = "client"

type Config struct {
	Role     RoleType
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
