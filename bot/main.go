package bot

import (
	"log"

	"github.com/telegram-bot/regulation-reminder-bot/bot/commands"
	"github.com/telegram-bot/regulation-reminder-bot/bot/handlers"

	"github.com/bot-api/telegram"
	"github.com/bot-api/telegram/telebot"
	"github.com/sirupsen/logrus"
	"github.com/telegram-bot/regulation-reminder-bot/conf"
	"golang.org/x/net/context"
)

type Bot struct {
	config conf.Config
	log    *logrus.Entry
}

func New(config conf.Config) *Bot {
	return &Bot{
		config: config,
		log:    config.Log(),
	}
}

func (a Bot) Config() conf.Config {
	return a.config
}

func (a *Bot) Start() error {
	config := a.Config()

	token := config.Bot().Token
	debug := config.Bot().Debug

	api := telegram.New(token)
	api.Debug(debug)

	bot := telebot.NewWithAPI(api)
	bot.Use(telebot.Recover()) // recover if handler panic

	netCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bot.HandleFunc(handlers.Handle(config))

	// Use command middleware, that helps to work with commands
	bot.Use(telebot.Commands(map[string]telebot.Commander{
		"start": commands.Start(config),
	}))

	err := bot.Serve(netCtx)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
