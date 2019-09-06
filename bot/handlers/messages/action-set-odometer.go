package messages

import (
	"fmt"
	"time"

	"github.com/bot-api/telegram/telebot"

	"github.com/spf13/cast"

	"github.com/bot-api/telegram"
)

func (m *Message) ActionSetOdometer(data string) error {
	db := m.config.DB()

	delete(m.actionMap, m.user.ID)

	odometer, err := cast.ToInt64E(data)
	if err != nil {
		return err
	}

	m.user.Odometer = odometer
	m.user.UpdatedAt = time.Now()

	err = db.UpdateUser(m.user)
	if err != nil {
		return err
	}

	update := telebot.GetUpdate(m.ctx)
	api := telebot.GetAPI(m.ctx)

	msg := telegram.NewMessage(update.Chat().ID, fmt.Sprintf("Your new odometer is %d", odometer))
	_, err = api.SendMessage(m.ctx, msg)

	return err
}
