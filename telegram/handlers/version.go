package handlers

import (
	"context"

	"github.com/fmotalleb/go-tools/git"
	"github.com/fmotalleb/go-tools/log"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/fmotalleb/north_outage/telegram/helpers"
)

func init() {
	register(
		func(_ context.Context, b *bot.Bot) {
			b.RegisterHandler(bot.HandlerTypeMessageText, "/version", bot.MatchTypePrefix, version)
		},
	)
}

func version(ctx context.Context, b *bot.Bot, update *models.Update) {
	chat := update.Message.Chat
	l := log.Of(ctx).
		Named("version").
		With(zap.Any("chat", chat))
	mp := helpers.MakeMessage(update)

	mp.Text = git.String()

	msg, err := b.SendMessage(ctx, mp)
	if err != nil {
		l.Error("failed to send version message", zap.Error(err))
		return
	}
	l.Debug("message sent", zap.Int("id", msg.ID))
}
