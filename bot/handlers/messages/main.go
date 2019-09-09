package messages

import (
	"context"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"

	"github.com/telegram-bot/regulation-reminder-bot/bot/helpers"
	"github.com/telegram-bot/regulation-reminder-bot/conf"
	"github.com/telegram-bot/regulation-reminder-bot/db/resources"
)

type Message struct {
	ctx       context.Context
	config    conf.Config
	actionMap map[int64]string
	jobMap    map[int64]*resources.Job
	chanMap   map[int64]chan bool
	user      *resources.User
}

func New(
	ctx context.Context,
	config conf.Config,
	actionMap map[int64]string,
	jobMap map[int64]*resources.Job,
	chanMap map[int64]chan bool) *Message {

	return &Message{
		ctx:       ctx,
		config:    config,
		actionMap: actionMap,
		jobMap:    jobMap,
		chanMap:   chanMap,
	}
}

func (m *Message) Run() error {
	db := m.config.DB()

	update := telebot.GetUpdate(m.ctx)
	api := telebot.GetAPI(m.ctx)

	from := update.From()

	user, err := db.GetUserByTelegramID(from.ID)
	if err != nil {
		return err
	}

	m.user = user

	if update.CallbackQuery != nil {
		data := update.CallbackQuery.Data

		switch data {
		case "/set-odometer":
			return m.QuerySetOdometer()
		case "/add-job":
			return m.QueryAddJob()
		case "/remove-job":
			return m.QueryRemoveJob()
		case "/show-odometer":
			return m.QueryShowOdometer()
		case "/per-odometer":
			return m.QueryPerOdometer()
		case "/per-date":
			return m.QueryPerDate()
		case "/job-list":
			return m.QueryJoBList()
		case "/start-monitoring":
			return m.QueryStartMonitoring()
		case "/stop-monitoring":
			return m.QueryStopMonitoring()
		default:
			delete(m.actionMap, user.ID)

			_, err := api.EditMessageText(m.ctx, telegram.NewEditMessageText(
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

		if action, ok := m.actionMap[user.ID]; ok {
			switch action {
			case "/set-odometer":
				return m.ActionSetOdometer(data)
			case "/remove-job":
				return m.ActionRemoveJob(data)
			case "/add-job":
				return m.ActionAddJob(data)
			}
		}
	}

	msg := telegram.NewMessage(update.Chat().ID, "Choose action:")
	msg.ReplyMarkup = helpers.GenerateMainButtons(user)

	_, err = api.SendMessage(m.ctx, msg)

	return err
}
