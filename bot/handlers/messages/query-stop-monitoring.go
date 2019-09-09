package messages

import (
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
)

func (m *Message) QueryStopMonitoring() error {
	if m.user.EnableMonitoring == false {
		return nil
	}

	db := m.config.DB()

	m.user.EnableMonitoring = false

	err := db.UpdateUser(m.user)
	if err != nil {
		return err
	}

	if _, ok := m.chanMap[m.user.ID]; ok {
		m.chanMap[m.user.ID] <- true
	}

	update := telebot.GetUpdate(m.ctx)
	api := telebot.GetAPI(m.ctx)

	_, err = api.EditMessageText(m.ctx, telegram.NewEditMessageText(
		update.Chat().ID,
		update.CallbackQuery.Message.MessageID,
		"Monitoring was stopped",
	))

	return nil
}
