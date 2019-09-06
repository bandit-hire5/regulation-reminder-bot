package messages

import (
	"fmt"

	"github.com/bot-api/telegram/telebot"

	"github.com/bot-api/telegram"
)

func (m *Message) QueryJoBList() error {
	db := m.config.DB()

	list, err := db.JobListByUser(m.user.ID)
	if err != nil {
		return err
	}

	var displayMsg string

	for _, job := range list {
		if job.Type == "odometer" {
			displayMsg += fmt.Sprintf(
				"%d. %s (%d)\n",
				job.ID, job.Name, job.LastOdometer+job.PerOdometer-m.user.Odometer,
			)
		}
	}

	update := telebot.GetUpdate(m.ctx)
	api := telebot.GetAPI(m.ctx)

	_, err = api.EditMessageText(m.ctx, telegram.NewEditMessageText(
		update.Chat().ID,
		update.CallbackQuery.Message.MessageID,
		displayMsg,
	))

	return err
}
