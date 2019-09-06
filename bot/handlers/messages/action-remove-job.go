package messages

import (
	"fmt"

	"github.com/bot-api/telegram/telebot"
	"github.com/spf13/cast"

	"github.com/bot-api/telegram"
)

func (m *Message) ActionRemoveJob(data string) error {
	db := m.config.DB()

	delete(m.actionMap, m.user.ID)

	id, err := cast.ToInt64E(data)
	if err != nil {
		return err
	}

	err = db.RemoveJob(id, m.user.ID)
	if err != nil {
		return err
	}

	update := telebot.GetUpdate(m.ctx)
	api := telebot.GetAPI(m.ctx)

	msg := telegram.NewMessage(update.Chat().ID, fmt.Sprintf("Job %d was deleted", id))
	_, err = api.SendMessage(m.ctx, msg)

	return err
}
