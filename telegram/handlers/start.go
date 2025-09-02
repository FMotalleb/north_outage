package handlers

import (
	"context"

	"github.com/fmotalleb/go-tools/log"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/fmotalleb/north_outage/telegram/template"
)

func init() {
	register(
		func(b *bot.Bot) {
			b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypePrefix, hello)
			b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypePrefix, hello)
		},
	)
}

func hello(ctx context.Context, b *bot.Bot, update *models.Update) {
	l := log.Of(ctx).Named("hello")
	mp := new(bot.SendMessageParams)
	mp.ChatID = update.Message.Chat.ID
	out, err := template.EvaluateTemplate(template.Help, nil, update)
	if err != nil {
		mp.Text = "سلام"
	} else {
		mp.Text = out
	}

	msg, err := b.SendMessage(ctx, mp)
	if err != nil {
		l.Error("failed to send hello message", zap.Error(err))
		return
	}
	l.Debug("sent message", zap.Int("id", msg.ID))
}
