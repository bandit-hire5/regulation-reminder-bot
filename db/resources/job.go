package resources

import "time"

const JobsTableName = "jobs"

type Job struct {
	ID           int64      `db:"id"`
	UserID       int64      `db:"user_id"`
	Type         string     `db:"type"`
	Name         string     `db:"name"`
	PerDays      int64      `db:"per_days"`
	PerOdometer  int64      `db:"per_odometer"`
	LastDate     *time.Time `db:"last_date"`
	LastOdometer int64      `db:"last_odometer"`
	Periodically bool       `db:"periodically"`
}

func (p *Job) TableName() string {
	return JobsTableName
}
