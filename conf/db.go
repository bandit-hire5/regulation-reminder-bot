package conf

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/telegram-bot/regulation-reminder-bot/db"

	"github.com/caarlos0/env"
)

type DB struct {
	Name     string `env:"BOT_DB_NAME,required"`
	Host     string `env:"BOT_DB_HOST" envDefault:"0.0.0.0"`
	Port     int    `env:"BOT_DB_PORT" envDefault:"5432"`
	User     string `env:"BOT_DB_USER,required"`
	Password string `env:"BOT_DB_PASSWORD,required"`
	SSL      string `env:"BOT_DB_SSL" envDefault:"disable"`
}

func (d DB) Info() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSL,
	)
}

func (c *ConfigImpl) DB() *db.DB {
	if c.db != nil {
		return c.db
	}

	database := &DB{}
	if err := env.Parse(database); err != nil {
		panic(err)
	}

	repo, err := db.New(database.Info())
	if err != nil {
		panic(errors.Wrap(err, "failed to setup db"))
	}

	c.db = repo

	return c.db
}

func (c *ConfigImpl) MigratorDB() *sql.DB {
	if c.migrate != nil {
		return c.migrate
	}

	database := &DB{}
	if err := env.Parse(database); err != nil {
		panic(err)
	}

	repo, err := sql.Open("postgres", database.Info())
	if err != nil {
		panic(errors.Wrap(err, "failed to setup migrate db"))
	}

	c.migrate = repo

	return c.migrate
}
