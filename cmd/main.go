package main

import (
	"github.com/spf13/viper"
)

func init() {
	initConfig()
	initLogger()
}

func main() {
	bot := bot{}
	bot.mount(viper.GetString("bot.apiToken"))
	bot.start()
}
