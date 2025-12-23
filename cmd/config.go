package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigName("configs")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	viper.SetDefault("environment", "production")
}
