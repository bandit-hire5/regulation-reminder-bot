package conf

import (
	"github.com/caarlos0/env"
)

type Bot struct {
	Token string `env:"TELEGRAM_BOT_TOKEN,required"`
	Debug bool   `env:"TELEGRAM_BOT_DEBUG" envDefault:"true"`
}

func (c *ConfigImpl) Bot() *Bot {
	if c.bot != nil {
		return c.bot
	}

	bot := &Bot{}
	if err := env.Parse(bot); err != nil {
		panic(err)
	}

	c.bot = bot

	return c.bot
}
