package messages

import (
	"github.com/bot-api/telegram/telebot"

	"github.com/bot-api/telegram"
)

func (m *Message) QueryRemoveJob() error {
	update := telebot.GetUpdate(m.ctx)
	api := telebot.GetAPI(m.ctx)

	_, err := api.EditMessageText(m.ctx, telegram.NewEditMessageText(
		update.Chat().ID,
		update.CallbackQuery.Message.MessageID,
		"Please enter job ID",
	))

	m.actionMap[m.user.ID] = "/remove-job"

	return err
}
