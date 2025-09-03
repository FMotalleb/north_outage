package telegram

import (
	"context"
	"time"

	"github.com/fmotalleb/go-tools/log"
	"go.uber.org/zap"

	"github.com/fmotalleb/north_outage/config"
	im "github.com/fmotalleb/north_outage/models"
	"github.com/fmotalleb/north_outage/telegram/handlers"

	"github.com/go-telegram/bot"
)

func Run(ctx context.Context, cfg *config.Config, nc <-chan im.Notification) error {
	l := log.Of(ctx).Named("Telegram")
	ctx = log.WithLogger(ctx, l)
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
	handlers.SetupHandlers(ctx, b)
	go bindToChannel(ctx, b, nc)
	b.Start(ctx)
	return nil
}

func bindToChannel(ctx context.Context, b *bot.Bot, nc <-chan im.Notification) {
	l := log.Of(ctx).Named("binder")
	for {
		select {
		case n := <-nc:
			l.Debug("notification received", zap.Any("event", n))
			sp := new(bot.SendMessageParams)
			sp.ChatID = n.Listener.TelegramCID
			sp.MessageThreadID = int(n.Listener.TelegramTID)
			sp.Text = n.Event.Address
			b.SendMessage(ctx, sp)
		case <-ctx.Done():
			return
		}
	}
}
