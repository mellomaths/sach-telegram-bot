package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mellomaths/sach-telegram-bot/internal/sac"
	"go.uber.org/zap"
)

type bot struct {
	api *tgbotapi.BotAPI
}

func (b *bot) mount(apiToken string) error {
	zap.L().Info("Mounting bot")
	defer zap.L().Sync()
	bot, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		zap.L().Error("Error mounting bot")
		return err
	}
	zap.L().Info("Authorized on account", zap.String("bot", bot.Self.UserName))
	b.api = bot
	return nil
}

func (b *bot) start() error {
	zap.L().Info(
		"Starting bot loop",
		zap.String("bot", b.api.Self.UserName),
		zap.Int64("botId", b.api.Self.ID),
	)
	uc := tgbotapi.NewUpdate(0)
	uc.Timeout = 60
	updates := b.api.GetUpdatesChan(uc)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		zap.L().Info(
			"Received a message",
			zap.String("bot", b.api.Self.UserName),
			zap.Int64("botId", b.api.Self.ID),
			zap.Int64("from", update.Message.From.ID),
			zap.String("text", update.Message.Text),
		)
		u := sac.User{
			Id:        update.Message.From.ID,
			FirstName: update.Message.From.FirstName,
			LastName:  update.Message.From.LastName,
			UserName:  update.Message.From.UserName,
		}
		if !update.Message.IsCommand() {
			sac.SaveSAC(u, update.Message.Text)
		}

		// Extract the command from the Message.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		switch update.Message.Command() {
		case "help":
			msg.Text = "I'll keep your complaints on file for review, use /sac to send me one."
		case "status":
			msg.Text = "I'm ok."
		case "sac":
			sac.SaveSAC(u, update.Message.Text)
		default:
			continue
		}

		msg.ReplyToMessageID = update.Message.MessageID
		zap.L().Info(
			"Replying",
			zap.String("bot", b.api.Self.UserName),
			zap.Int64("botId", b.api.Self.ID),
			zap.Int64("to", update.Message.From.ID),
			zap.String("text", update.Message.Text),
		)
		b.api.Send(msg)
	}
	return nil
}
