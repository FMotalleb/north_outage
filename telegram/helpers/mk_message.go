package helpers

import (
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func MakeMessage(update *models.Update) *bot.SendMessageParams {
	mp := new(bot.SendMessageParams)
	input := update.Message
	mp.ChatID = input.Chat.ID
	mp.ParseMode = models.ParseModeHTML
	return mp
}
