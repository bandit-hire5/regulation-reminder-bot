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

func (d *DB) ComingJobListByUser(userID int64) ([]resources.Job, error) {
	var jobs []resources.Job

	err := d.db.Select("jobs.*").
		From(resources.JobsTableName).
		InnerJoin("users", dbx.NewExp("jobs.user_id = users.id")).
		Where(dbx.And(
			dbx.HashExp{
				"jobs.user_id": userID,
			},
			dbx.Or(
				dbx.And(
					dbx.HashExp{
						"jobs.type": "odometer",
					},
					dbx.NewExp("users.odometer >= jobs.last_odometer + jobs.per_odometer - {:reserve_odometer}", dbx.Params{"reserve_odometer": 200}),
					dbx.NewExp("users.odometer < jobs.last_odometer + jobs.per_odometer"),
				),
				dbx.And(
					dbx.HashExp{
						"jobs.type": "date",
					},
					dbx.NewExp("current_date >= jobs.last_date + interval '1 day' * (jobs.per_days::integer - {:reserve_days})", dbx.Params{"reserve_days": 1}),
					dbx.NewExp("current_date < jobs.last_date + interval '1 day' * jobs.per_days::integer"),
				),
			),
		)).
		All(&jobs)

	return jobs, err
}

func (d *DB) ExpiredJobListByUser(userID int64) ([]resources.Job, error) {
	var jobs []resources.Job

	err := d.db.Select("jobs.*").
		From(resources.JobsTableName).
		InnerJoin("users", dbx.NewExp("jobs.user_id = users.id")).
		Where(dbx.Or(
			dbx.And(
				dbx.HashExp{
					"jobs.user_id": userID,
				},
				dbx.HashExp{
					"jobs.type": "odometer",
				},
				dbx.NewExp("users.odometer >= jobs.last_odometer + jobs.per_odometer"),
			),
			dbx.And(
				dbx.HashExp{
					"jobs.user_id": userID,
				},
				dbx.HashExp{
					"jobs.type": "date",
				},
				dbx.NewExp("current_date >= jobs.last_date + interval '1 day' * jobs.per_days::integer"),
			),
		)).
		All(&jobs)

	return jobs, err
}
