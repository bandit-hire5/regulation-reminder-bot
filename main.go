package main

import (
	"database/sql"
	"fmt"

	"github.com/telegram-bot/regulation-reminder-bot/bot"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/telegram-bot/regulation-reminder-bot/conf"
	mg "github.com/telegram-bot/regulation-reminder-bot/migrator"
)

type Migrator func(*sql.DB, mg.MigrateDir, int) (int, error)

func MigrateDB(direction string, count int, dbClient *sql.DB, migratorFn Migrator) (int, error) {
	applied, err := migratorFn(dbClient, mg.MigrateDir(direction), count)

	return applied, errors.Wrap(err, "failed to apply migrations")
}

func main() {
	config := conf.New()
	log := config.Log()

	rootCmd := &cobra.Command{
		Use: "api",
	}

	migrateCmd := &cobra.Command{
		Use:   "migrate [up|down|redo] [COUNT]",
		Short: "migrate schema",
		Long:  "performs a schema migration command",
		Run: func(cmd *cobra.Command, args []string) {
			var count int
			var err error

			// Allow invocations with 1 or 2 args. All other args counts are erroneous
			if len(args) < 1 || len(args) > 2 {
				log.WithField("arguments", args).Error("wrong argument count")
				return
			}

			// If a second arg is present, parse it to an int and use it as the count
			// argument to the migration call
			if len(args) == 2 {
				if count, err = cast.ToIntE(args[1]); err != nil {
					log.WithError(err).Error("failed to parse count")
					return
				}
			}

			applied, err := MigrateDB(args[0], count, config.MigratorDB(), mg.Migrations.Migrate)
			log = log.WithField("applied", applied)
			if err != nil {
				log.WithError(err).Error("migration failed")
				return
			}

			log.Info("migrations applied")
		},
	}

	rootCmd.AddCommand(migrateCmd)

	runCmd := &cobra.Command{
		Use: "run",
		Run: func(cmd *cobra.Command, args []string) {
			defer func() {
				if rvr := recover(); rvr != nil {
					log.WithField("panic stack trace", rvr).Error("app panicked")
				}
			}()

			bot := bot.New(config)
			if err := bot.Start(); err != nil {
				panic(errors.Wrap(err, "failed to start telegram bot"))
			}
		},
	}

	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		log.WithField("cobra", "read").Error(fmt.Sprintf("failed to read command %s", err.Error()))
		return
	}
}
