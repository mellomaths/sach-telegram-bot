package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func init() {
	initConfig()
	initLogger()
}

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, viper.GetString("db.connectionString"))
	if err != nil {
		zap.L().Error("failed to connect to postgres database", zap.Error(err))
		os.Exit(1)
	}
	defer conn.Close(ctx)
	zap.L().Info("connected to database")
	bot := bot{db: conn}
	bot.mount(viper.GetString("bot.apiToken"))
	bot.start()
}
