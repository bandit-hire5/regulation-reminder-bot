package db

import (
	"database/sql"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/telegram-bot/regulation-reminder-bot/db/resources"
)

func (d *DB) GetUserByTelegramID(userID int64) (*resources.User, error) {
	var user resources.User

	err := d.db.Select().
		From(resources.UsersTableName).
		Where(dbx.HashExp{
			"telegram_id": userID,
		}).
		One(&user)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &user, err
}

func (d *DB) AddUser(user *resources.User) error {
	return d.db.Model(user).Insert()
}

func (d *DB) UpdateUser(user *resources.User) error {
	return d.db.Model(user).Update()
}
