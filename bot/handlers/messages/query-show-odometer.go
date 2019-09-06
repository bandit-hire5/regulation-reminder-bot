package messages

import (
	"fmt"

	"github.com/bot-api/telegram/telebot"

	"github.com/bot-api/telegram"
)

func (m *Message) QueryShowOdometer() error {
	update := telebot.GetUpdate(m.ctx)
	api := telebot.GetAPI(m.ctx)

	_, err := api.EditMessageText(m.ctx, telegram.NewEditMessageText(
		update.Chat().ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("Your current odometer is %d", m.user.Odometer),
	))

	return err
}
