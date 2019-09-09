package messages

import (
	"fmt"
	"time"

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
				"%d. %s (after %d km)\n",
				job.ID, job.Name, job.LastOdometer+job.PerOdometer-m.user.Odometer,
			)
		} else {
			date := job.LastDate.AddDate(0, 0, int(job.PerDays))
			duration := date.Sub(time.Now())

			displayMsg += fmt.Sprintf(
				"%d. %s (after %d days)\n",
				job.ID, job.Name, int(duration.Hours()/24),
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
