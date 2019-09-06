package messages

import (
	"github.com/bot-api/telegram/telebot"
	"github.com/telegram-bot/regulation-reminder-bot/bot/helpers"

	"github.com/bot-api/telegram"
)

func (m *Message) QueryAddJob() error {
	update := telebot.GetUpdate(m.ctx)
	api := telebot.GetAPI(m.ctx)

	_, err := api.EditMessageText(m.ctx, telegram.NewEditMessageText(
		update.Chat().ID,
		update.CallbackQuery.Message.MessageID,
		"Please select job type:",
	))

	_, err = api.EditMessageReplyMarkup(m.ctx, telegram.NewEditMessageReplyMarkup(
		update.Chat().ID,
		update.CallbackQuery.Message.MessageID,
		helpers.GenerateJobTypeButtons(),
	))

	return err
}
