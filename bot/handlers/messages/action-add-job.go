package messages

import (
	"time"

	"github.com/telegram-bot/regulation-reminder-bot/db/resources"

	"github.com/spf13/cast"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
)

func (m *Message) ActionAddJob(data string) error {
	db := m.config.DB()

	job, ok := m.jobMap[m.user.ID]
	if !ok {
		job = &resources.Job{
			UserID: m.user.ID,
		}
	}

	step, nextMessage := getJobStep(job)

	switch step {
	case 1:
		job.Name = data
		m.jobMap[m.user.ID] = job
	case 2:
		perWhat, err := cast.ToInt64E(data)
		if err != nil {
			return err
		}

		if job.Type == "odometer" {
			job.PerOdometer = perWhat
		} else {
			job.PerDays = perWhat
		}

		m.jobMap[m.user.ID] = job
	case 3:
		if job.Type == "odometer" {
			odometer, err := cast.ToInt64E(data)
			if err != nil {
				return err
			}

			job.LastOdometer = odometer
		} else {
			lastUpdatedAt, err := time.Parse("2006-01-02", data)
			if err != nil {
				return err
			}

			job.LastDate = &lastUpdatedAt
		}

		m.jobMap[m.user.ID] = job
	case 4:
		if data == "no" {
			job.Periodically = false
		} else {
			job.Periodically = true
		}

		err := db.AddJob(job)
		if err != nil {
			return err
		}

		delete(m.jobMap, m.user.ID)
		delete(m.actionMap, m.user.ID)
	}

	update := telebot.GetUpdate(m.ctx)
	api := telebot.GetAPI(m.ctx)

	msg := telegram.NewMessage(update.Chat().ID, nextMessage)
	_, err := api.SendMessage(m.ctx, msg)

	return err
}

func getJobStep(job *resources.Job) (int, string) {
	if job.Type == "odometer" {
		if job.Name == "" {
			return 1, "Please enter job per odometer value"
		}

		if job.PerOdometer == 0 {
			return 2, "Please enter job last odometer"
		}

		if job.LastOdometer == 0 {
			return 3, "Please enter job is periodically (yes|no)"
		}
	} else {
		if job.Name == "" {
			return 1, "Please enter job per days value"
		}

		if job.PerDays == 0 {
			return 2, "Please enter job last date (yyyy-mm-dd)"
		}

		if job.LastDate == nil {
			return 3, "Please enter job is periodically (yes|no)"
		}
	}

	return 4, "Job successfully created!"
}
