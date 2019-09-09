package messages

import (
	"fmt"
	"time"

	"github.com/bot-api/telegram/telebot"

	"github.com/bot-api/telegram"
)

func (m *Message) QueryStartMonitoring() error {
	if m.user.EnableMonitoring == true {
		return nil
	}

	db := m.config.DB()

	m.user.EnableMonitoring = true

	err := db.UpdateUser(m.user)
	if err != nil {
		return err
	}

	m.chanMap[m.user.ID] = make(chan bool)

	update := telebot.GetUpdate(m.ctx)
	api := telebot.GetAPI(m.ctx)

	go func() {
		ticker := time.NewTicker(time.Second * 5)

		defer func() {
			ticker.Stop()
		}()

		for {
			select {
			case <-ticker.C:
				duration := time.Now().Sub(m.user.UpdatedAt)

				if duration.Hours()/24 >= 30 {
					m.user.Odometer += m.user.AvgOdometer

					err := db.UpdateUser(m.user)
					if err == nil {
						msg := telegram.NewMessage(update.Chat().ID, "Odometer was changed automatically. Please correct it manually!")
						_, _ = api.SendMessage(m.ctx, msg)
					}
				}

				list, err := db.ComingJobListByUser(m.user.ID)
				fmt.Println("err", err)
				fmt.Println("list", list)
				if err == nil {
					for _, item := range list {
						msg := telegram.NewMessage(update.Chat().ID, fmt.Sprintf("Coming job %s", item.Name))
						_, _ = api.SendMessage(m.ctx, msg)
					}
				}

				list, err = db.ExpiredJobListByUser(m.user.ID)
				if err == nil {
					for _, item := range list {
						msg := telegram.NewMessage(update.Chat().ID, fmt.Sprintf("Expired job %s", item.Name))
						_, _ = api.SendMessage(m.ctx, msg)
					}
				}
			case <-m.chanMap[m.user.ID]:
				return
			}
		}
	}()

	_, err = api.EditMessageText(m.ctx, telegram.NewEditMessageText(
		update.Chat().ID,
		update.CallbackQuery.Message.MessageID,
		"Monitoring was started",
	))

	return nil
}
