package conf

import (
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/telegram-bot/regulation-reminder-bot/db"
)

type Config interface {
	Log() *logrus.Entry
	DB() *db.DB
	MigratorDB() *sql.DB
	Bot() *Bot
}

type ConfigImpl struct {
	log     *logrus.Entry
	db      *db.DB
	migrate *sql.DB
	bot     *Bot
}

func New() Config {
	return &ConfigImpl{}
}
