package handlers

import (
	"context"

	"github.com/go-telegram/bot"
)

type TGHandlerRegistrant = func(context.Context, *bot.Bot)

var (
	registrants []TGHandlerRegistrant
	finalized   bool
)

// registry    map[string]string

func register(r TGHandlerRegistrant) {
	if finalized {
		panic("registering handlers after finalization of bot is prohibited")
	}
	registrants = append(registrants, r)
}

func SetupHandlers(ctx context.Context, b *bot.Bot) {
	finalized = true
	for _, h := range registrants {
		h(ctx, b)
	}
}
