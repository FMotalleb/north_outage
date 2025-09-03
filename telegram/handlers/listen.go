package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fmotalleb/go-tools/log"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/fmotalleb/north_outage/database"
	"github.com/fmotalleb/north_outage/memory"
	im "github.com/fmotalleb/north_outage/models"
	"github.com/fmotalleb/north_outage/telegram/helpers"
)

type listenReq struct {
	city   string
	search string
}

var mem = memory.NewRuntime[listenReq]()

const memoryTTL = time.Minute

func init() {
	register(
		func(_ context.Context, b *bot.Bot) {
			b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "listen:", bot.MatchTypePrefix, listen)
		},
	)
}

func listen(ctx context.Context, b *bot.Bot, update *models.Update) {
	l := log.Of(ctx).Named("listen")
	key := strings.TrimPrefix(update.CallbackQuery.Data, "listen:")
	cp := new(bot.AnswerCallbackQueryParams)
	cp.ShowAlert = true
	cp.CallbackQueryID = update.CallbackQuery.ID
	cp.Text = "لطفا کمی صبر کنید"
	_, _ = b.AnswerCallbackQuery(ctx, cp)
	mp := helpers.MakeMessage(update)
	data, ok := mem.Pop(key)
	if !ok {
		l.Error("failed to retrieve data related to key", zap.String("key", key))
		responseError(ctx, l, mp, b)
		return
	}
	db := database.Get()
	listen := new(im.Listener)
	listen.City = data.city
	listen.SearchTerm = data.search
	listen.TelegramCID = update.CallbackQuery.Message.Message.Chat.ID
	listen.TelegramTID = int64(update.CallbackQuery.Message.Message.MessageThreadID)
	mp.Text = "در صورت دریافت اطلاعات جدید و همچنین بیست دقیقه قبل از قطعی بهت اطلاع میدم"
	if err := db.Save(listen).Error; err != nil {
		l.Error("failed to store request to database", zap.Any("request", listen), zap.Error(err))
		mp.Text = "خطا در ذخیره‌سازی داده، احتمالا داری آیتم تکراری ذخیره میکنی"
	}
	msg, err := b.SendMessage(ctx, mp)
	if err != nil {
		l.Error("failed to send search results", zap.Error(err))
		return
	}
	l.Debug("sent message", zap.Int("id", msg.ID))
}

func responseError(ctx context.Context, l *zap.Logger, mp *bot.SendMessageParams, b *bot.Bot) {
	mp.Text = "درخواست منقضی شد لطفا مجددا جست و جو کنید"
	msg, err := b.SendMessage(ctx, mp)
	if err != nil {
		l.Error("failed to send search results", zap.Error(err))
		return
	}
	l.Debug("sent error message", zap.Int("id", msg.ID))
}

func createRequest(search, city string) string {
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("random uuid generation failed: %w", err))
	}
	req := listenReq{
		city:   city,
		search: search,
	}
	key := uuid.String()
	mem.Put(key, req, memoryTTL)
	return key
}
