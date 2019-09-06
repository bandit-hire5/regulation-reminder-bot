package helpers

import (
	"github.com/bot-api/telegram"
	"github.com/telegram-bot/regulation-reminder-bot/db/resources"
)

func GenerateMainButtons(user *resources.User) telegram.InlineKeyboardMarkup {
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

func GenerateJobTypeButtons() *telegram.InlineKeyboardMarkup {
	text := []string{"Per Odometer", "Per Days"}
	data := []string{"per-odometer", "per-date"}

	return &telegram.InlineKeyboardMarkup{
		InlineKeyboard: telegram.NewHInlineKeyboard(
			"/",
			text,
			data,
		),
	}
}
