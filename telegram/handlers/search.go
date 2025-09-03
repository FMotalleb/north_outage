package handlers

import (
	"context"
	"strings"

	"github.com/fmotalleb/go-tools/log"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.uber.org/zap"

	"github.com/fmotalleb/north_outage/database"
	im "github.com/fmotalleb/north_outage/models"
	"github.com/fmotalleb/north_outage/telegram/helpers"
	"github.com/fmotalleb/north_outage/telegram/template"
)

const (
	maxSearchResult = 10
	searchCMD       = "/search"
)

func init() {
	register(
		func(_ context.Context, b *bot.Bot) {
			b.RegisterHandler(bot.HandlerTypeMessageText, searchCMD, bot.MatchTypePrefix, search)
			b.RegisterHandlerMatchFunc(shouldSearch, search)
		},
	)
}

func search(ctx context.Context, b *bot.Bot, update *models.Update) {
	l := log.Of(ctx).Named("search")
	// mp := new(bot.SendMessageParams)
	input := update.Message
	search := strings.TrimPrefix(input.Text, searchCMD)
	events, err := fetchEvents(search)

	mp := helpers.MakeMessage(update)
	if err != nil {
		l.Error("failed to fetch data from db", zap.Error(err))
		mp.Text = "خطا در دریافت داده"
	} else {
		data := map[string]any{
			"results": events,
		}
		var out string
		out, err = template.EvaluateTemplate(template.Search, data, update)
		if err != nil {
			l.Error("failed to evaluate template", zap.Error(err), zap.Any("chat", update.Message.Chat))
			mp.Text = "خطایی در نمایش خروجی پیش اومده"
		} else {
			mp.Text = out
		}
	}

	msg, err := b.SendMessage(ctx, mp)
	if err != nil {
		l.Error("failed to send hello message", zap.Error(err))
		return
	}
	l.Debug("sent message", zap.Int("id", msg.ID))
}

func shouldSearch(update *models.Update) bool {
	query := update.Message.Text
	var exists bool
	err := database.Get().
		Table("events").
		Select("1").
		Where("address LIKE ?", "%"+query+"%").
		Limit(1).
		Scan(&exists).Error
	return err == nil && exists
}

func fetchEvents(search string) ([]im.Event, error) {
	out := make([]im.Event, 0, maxSearchResult)
	err := database.Get().
		Table("events").
		Where("address LIKE ?", "%"+search+"%").
		Limit(maxSearchResult).
		Find(&out).Error
	return out, err
}
