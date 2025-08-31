package telegram

import (
	"context"
	"time"

	"github.com/fmotalleb/go-tools/log"
	"go.uber.org/zap"

	"github.com/fmotalleb/north_outage/config"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Run(ctx context.Context, cfg *config.Config) error {
	l := log.Of(ctx).Named("Telegram")
	tel := cfg.Telegram
	if tel.BotKey == "" {
		l.Warn("telegram bot token is not set")
		return nil
	}

	var opts []bot.Option
	client := httpClient(&tel.Proxy)

	hc := bot.WithHTTPClient(time.Second*30, client)
	opts = append(opts, hc)

	b, err := bot.New(tel.BotKey, opts...)
	if err != nil {
		l.Error("failed to connect to telegram bot", zap.Error(err))
		return err
	}
	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, func(ctx context.Context, b *bot.Bot, update *models.Update) {
		msg := new(bot.SendMessageParams)
		msg.ChatID = update.Message.Chat.ID
		msg.Text = update.Message.Text
		b.SendMessage(ctx, msg)
	})
	b.Start(ctx)
	return nil
}
