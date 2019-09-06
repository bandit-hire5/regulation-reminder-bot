package db

import (
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/telegram-bot/regulation-reminder-bot/db/resources"
)

func (d *DB) AddJob(job *resources.Job) error {
	return d.db.Model(job).Insert()
}

func (d *DB) UpdateJob(job *resources.Job) error {
	return d.db.Model(job).Update()
}

func (d *DB) JobListByUser(userID int64) ([]resources.Job, error) {
	var jobs []resources.Job

	err := d.db.Select("*").
		From(resources.JobsTableName).
		Where(dbx.HashExp{
			"user_id": userID,
		}).
		All(&jobs)

	return jobs, err
}

func (d *DB) RemoveJob(jobID, userID int64) error {
	_, err := d.db.Delete(
		resources.JobsTableName,
		dbx.HashExp{"id": jobID, "user_id": userID},
	).Execute()

	return err
}
