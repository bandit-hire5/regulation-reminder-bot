package handlers

import (
	"context"

	"github.com/bot-api/telegram/telebot"
	"github.com/telegram-bot/regulation-reminder-bot/bot/handlers/messages"

	"github.com/telegram-bot/regulation-reminder-bot/conf"
	"github.com/telegram-bot/regulation-reminder-bot/db/resources"
)

var actionMap map[int64]string
var jobMap map[int64]*resources.Job

func Handle(config conf.Config) telebot.HandlerFunc {
	actionMap = make(map[int64]string)
	jobMap = make(map[int64]*resources.Job)

	return func(ctx context.Context) error {
		message := messages.New(ctx, config, actionMap, jobMap)

		return message.Run()
	}
}
