package handlers

import (
	"context"

	"github.com/fmotalleb/go-tools/log"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/fmotalleb/north_outage/database"
	im "github.com/fmotalleb/north_outage/models"
	"github.com/fmotalleb/north_outage/telegram/helpers"
	"github.com/fmotalleb/north_outage/telegram/template"
)

func init() {
	register(
		func(b *bot.Bot) {
			b.RegisterHandler(bot.HandlerTypeMessageText, "/search", bot.MatchTypePrefix, search)
			b.RegisterHandlerMatchFunc(shouldSearch, search)
		},
	)
}

func search(ctx context.Context, b *bot.Bot, update *models.Update) {
	l := log.Of(ctx).Named("search")
	// mp := new(bot.SendMessageParams)
	input := update.Message
	events := fetchEvents(input.Text)
	data := map[string]any{
		"results": events,
	}
	mp := helpers.MakeMessage(update)

	out, err := template.EvaluateTemplate(template.Search, data, update)
	if err != nil {
		l.Error("failed to evaluate template", zap.Error(err), zap.Any("chat", update.Message.Chat))
		mp.Text = "خطایی در نمایش خروجی پیش اومده"
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

func shouldSearch(update *models.Update) bool {
	search := update.Message.Text
	out := make([]im.Event, 0, 1)
	database.Get().
		Table("events").
		Where("address LIKE ?", "%"+search+"%").
		Limit(1).
		Find(&out)
	if len(out) == 0 {
		return false
	}
	return true
}

func fetchEvents(search string) []im.Event {
	out := make([]im.Event, 0, 10)
	database.Get().
		Table("events").
		Where("address LIKE ?", "%"+search+"%").
		Limit(10).
		Find(&out)
	return out
}
