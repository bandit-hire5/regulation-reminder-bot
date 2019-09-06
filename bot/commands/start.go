package commands

import (
	"context"
	"time"

	"github.com/telegram-bot/regulation-reminder-bot/bot/helpers"
	"github.com/telegram-bot/regulation-reminder-bot/db/resources"

	"github.com/telegram-bot/regulation-reminder-bot/conf"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
)

func Start(config conf.Config) telebot.Commander {
	return telebot.CommandFunc(
		func(ctx context.Context, arg string) error {
			db := config.DB()
			update := telebot.GetUpdate(ctx)

			from := update.From()

			user, err := db.GetUserByTelegramID(from.ID)
			if err != nil {
				return err
			}

			if user == nil {
				user = &resources.User{
					TelegramID: from.ID,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
					FirstName:  from.FirstName,
					LastName:   from.LastName,
					Username:   from.Username,
				}

				err := db.AddUser(user)
				if err != nil {
					return err
				}
			}

			api := telebot.GetAPI(ctx)

			msg := telegram.NewMessage(update.Chat().ID, "Choose action:")
			msg.ReplyMarkup = helpers.GenerateMainButtons(user)

			_, err = api.SendMessage(ctx, msg)

			return err
		})
}
