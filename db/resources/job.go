package resources

import "time"

const JobsTableName = "jobs"

type Job struct {
	ID            int64      `db:"id"`
	UserID        int64      `db:"user_id"`
	Name          string     `db:"name"`
	Regulation    int64      `db:"regulation"`
	LastUpdatedAt *time.Time `db:"last_updated_at"`
	LastOdometer  int64      `db:"last_odometer"`
	NextOdometer  int64      `db:"next_odometer"`
	LeftOdometer  int64      `db:"left_odometer"`
}

func (p *Job) TableName() string {
	return JobsTableName
}
