package messages

import (
	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"github.com/telegram-bot/regulation-reminder-bot/db/resources"
)

func (m *Message) QueryPerOdometer() error {
	update := telebot.GetUpdate(m.ctx)
	api := telebot.GetAPI(m.ctx)

	_, err := api.EditMessageText(m.ctx, telegram.NewEditMessageText(
		update.Chat().ID,
		update.CallbackQuery.Message.MessageID,
		"Please enter job name",
	))

	job, ok := m.jobMap[m.user.ID]
	if !ok {
		job = &resources.Job{
			UserID: m.user.ID,
		}
	}

	job.Type = "odometer"
	m.jobMap[m.user.ID] = job

	m.actionMap[m.user.ID] = "/add-job"

	return err
}
