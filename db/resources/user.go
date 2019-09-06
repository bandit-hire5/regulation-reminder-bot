package resources

import "time"

const UsersTableName = "users"

type User struct {
	ID               int64     `db:"id"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
	TelegramID       int64     `db:"telegram_id"`
	FirstName        string    `db:"first_name"`
	LastName         string    `db:"last_name"`
	Username         string    `db:"username"`
	Odometer         int64     `db:"odometer"`
	EnableMonitoring bool      `db:"enable_monitoring"`
}

func (p *User) TableName() string {
	return UsersTableName
}
