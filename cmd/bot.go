package main

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	repo "github.com/mellomaths/sach-telegram-bot/internal/adapters/postgresql/sqlc"
	"github.com/mellomaths/sach-telegram-bot/internal/sacs"
	"go.uber.org/zap"
)

type bot struct {
	api *tgbotapi.BotAPI
	db  *pgx.Conn
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

func (b *bot) cleanCommandMessage(text string) string {
	text = strings.Replace(text, "/sac", "", 1)
	text = strings.TrimSpace(text)
	return text
}

func (b *bot) start() error {
	sacsService := sacs.NewService(repo.New(b.db))
	zap.L().Info(
		"Starting bot loop",
		zap.String("bot", b.api.Self.UserName),
		zap.Int64("bot_id", b.api.Self.ID),
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
			zap.Int64("bot_id", b.api.Self.ID),
			zap.Int64("user_id", update.Message.From.ID),
			zap.String("text", update.Message.Text),
		)
		u := sacs.User{
			Id:        update.Message.From.ID,
			FirstName: update.Message.From.FirstName,
			LastName:  update.Message.From.LastName,
			UserName:  update.Message.From.UserName,
		}
		if !update.Message.IsCommand() {
			zap.L().Info("Not a command", zap.Int64("user_id", u.Id), zap.String("text", update.Message.Text))
			continue
		}

		// Extract the command from the Message.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		text := b.cleanCommandMessage(update.Message.Text)
		if text == "" {
			zap.L().Info("Command is empty", zap.Int64("user_id", u.Id), zap.String("text", update.Message.Text))
			msg.Text = "Por favor, digite a reclamação."
			b.api.Send(msg)
			continue
		}
		switch update.Message.Command() {
		case "help":
			msg.Text = "I'll keep your complaints on file for review, use /sac to send me one."
		case "status":
			msg.Text = "I'm ok."
		case "sac":
			sacsService.SaveSac(context.Background(), u, text)
			msg.Text = "Reclamação salva com sucesso."
		default:
			zap.L().Info("Command not found", zap.Int64("user_id", u.Id), zap.String("text", update.Message.Text))
			continue
		}

		msg.ReplyToMessageID = update.Message.MessageID
		zap.L().Info(
			"Replying",
			zap.String("bot", b.api.Self.UserName),
			zap.Int64("bot_id", b.api.Self.ID),
			zap.Int64("user_id", update.Message.From.ID),
			zap.String("text", update.Message.Text),
		)
		b.api.Send(msg)
	}
	return nil
}
