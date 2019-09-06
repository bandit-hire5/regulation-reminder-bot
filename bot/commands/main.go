package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"github.com/spf13/cast"

	"github.com/telegram-bot/regulation-reminder-bot/conf"
	"github.com/telegram-bot/regulation-reminder-bot/db/resources"
)

var actionMap map[int64]string
var jobMap map[int64]*resources.Job

func Handle(config conf.Config) func(ctx context.Context) error {
	actionMap = make(map[int64]string)
	jobMap = make(map[int64]*resources.Job)

	return func(ctx context.Context) error {
		db := config.DB()

		update := telebot.GetUpdate(ctx) // take update from context
		api := telebot.GetAPI(ctx)       // take api from context

		from := update.From()

		user, err := db.GetUserByTelegramID(from.ID)
		if err != nil {
			return err
		}

		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data

			switch data {
			case "/set-odometer":
				_, err := api.EditMessageText(ctx, telegram.NewEditMessageText(
					update.Chat().ID,
					update.CallbackQuery.Message.MessageID,
					"Please enter you current odometer",
				))

				actionMap[user.ID] = data

				return err
			case "/add-job":
				_, err := api.EditMessageText(ctx, telegram.NewEditMessageText(
					update.Chat().ID,
					update.CallbackQuery.Message.MessageID,
					"Please enter job name",
				))

				actionMap[user.ID] = data

				return err
			case "/remove-job":
				_, err := api.EditMessageText(ctx, telegram.NewEditMessageText(
					update.Chat().ID,
					update.CallbackQuery.Message.MessageID,
					"Please enter job ID",
				))

				actionMap[user.ID] = data

				return err
			case "/show-odometer":
				_, err := api.EditMessageText(ctx, telegram.NewEditMessageText(
					update.Chat().ID,
					update.CallbackQuery.Message.MessageID,
					fmt.Sprintf("Your current odometer is %d", user.Odometer),
				))

				return err
			case "/job-list":
				list, err := db.JobListByUser(user.ID)
				if err != nil {
					return err
				}

				var displayMsg string

				for _, job := range list {
					displayMsg += fmt.Sprintf(
						"%d. %s (%d)\n",
						job.ID, job.Name, job.NextOdometer,
					)
				}

				_, err = api.EditMessageText(ctx, telegram.NewEditMessageText(
					update.Chat().ID,
					update.CallbackQuery.Message.MessageID,
					displayMsg,
				))

				return err
			default:
				delete(actionMap, user.ID)

				_, err := api.EditMessageText(ctx, telegram.NewEditMessageText(
					update.Chat().ID,
					update.CallbackQuery.Message.MessageID,
					"Command doesn't work",
				))

				if err != nil {
					return err
				}
			}
		} else if update.Message != nil {
			data := update.Message.Text

			if action, ok := actionMap[user.ID]; ok {
				switch action {
				case "/set-odometer":
					delete(actionMap, user.ID)

					odometer, err := cast.ToInt64E(data)
					if err != nil {
						return err
					}

					user.Odometer = odometer

					err = db.UpdateUser(user)
					if err != nil {
						return err
					}

					msg := telegram.NewMessage(update.Chat().ID, fmt.Sprintf("Your new odometer is %d", odometer))
					_, err = api.SendMessage(ctx, msg)

					return err
				case "/remove-job":
					delete(actionMap, user.ID)

					id, err := cast.ToInt64E(data)
					if err != nil {
						return err
					}

					err = db.RemoveJob(id, user.ID)
					if err != nil {
						return err
					}

					msg := telegram.NewMessage(update.Chat().ID, fmt.Sprintf("Job %d was deleted", id))
					_, err = api.SendMessage(ctx, msg)

					return err
				case "/add-job":
					job, ok := jobMap[user.ID]
					if !ok {
						job = &resources.Job{
							UserID: user.ID,
						}
					}

					step, nextMessage := getJobStep(job)

					switch step {
					case 1:
						job.Name = data
						jobMap[user.ID] = job
					case 2:
						regulation, err := cast.ToInt64E(data)
						if err != nil {
							return err
						}

						job.Regulation = regulation
						jobMap[user.ID] = job
					case 3:
						lastUpdatedAt, err := time.Parse("2006-01-02", data)
						if err != nil {
							return err
						}

						job.LastUpdatedAt = &lastUpdatedAt
						jobMap[user.ID] = job
					case 4:
						odometer, err := cast.ToInt64E(data)
						if err != nil {
							return err
						}

						job.LastOdometer = odometer
						jobMap[user.ID] = job
					case 5:
						odometer, err := cast.ToInt64E(data)
						if err != nil {
							return err
						}

						job.NextOdometer = odometer
						job.LeftOdometer = job.NextOdometer - user.Odometer

						err = db.AddJob(job)
						if err != nil {
							return err
						}

						delete(jobMap, user.ID)
						delete(actionMap, user.ID)
					}

					msg := telegram.NewMessage(update.Chat().ID, nextMessage)
					_, err = api.SendMessage(ctx, msg)

					return err
				}
			}
		}

		msg := telegram.NewMessage(update.Chat().ID, "Choose action:")
		msg.ReplyMarkup = generateButtons(user)

		_, err = api.SendMessage(ctx, msg)

		return err
	}
}

func getJobStep(job *resources.Job) (int, string) {
	if job.Name == "" {
		return 1, "Please enter job regulation"
	}

	if job.Regulation == 0 {
		return 2, "Please enter job last regulation date (yyyy-mm-dd)"
	}

	if job.LastUpdatedAt == nil {
		return 3, "Please enter job last odometer"
	}

	if job.LastOdometer == 0 {
		return 4, "Please enter job next odometer"
	}

	return 5, "Job successfully created!"
}

func generateButtons(user *resources.User) telegram.InlineKeyboardMarkup {
	text := [][]string{
		[]string{"Show odometer", "Job list", "Remove Job"},
		[]string{"Set odometer", "Add new job"},
	}
	data := [][]string{
		[]string{"show-odometer", "job-list", "remove-job"},
		[]string{"set-odometer", "add-job"},
	}

	if user.EnableMonitoring {
		text[1] = append(text[1], "Stop monitoring")
		data[1] = append(data[1], "stop-monitoring")
	} else {
		text[1] = append(text[1], "Start monitoring")
		data[1] = append(data[1], "start-monitoring")
	}

	prefix := "/"
	inlineKeyboards := make([][]telegram.InlineKeyboardButton, len(text))

	for x, items := range text {
		buttons := make([]telegram.InlineKeyboardButton, len(items))

		for y, item := range items {
			buttons[y] = telegram.InlineKeyboardButton{
				Text:         item,
				CallbackData: prefix + data[x][y],
			}
		}

		inlineKeyboards[x] = buttons
	}

	return telegram.InlineKeyboardMarkup{
		InlineKeyboard: inlineKeyboards,
	}
}
