package migrator

import (
	"github.com/sirupsen/logrus"
)

//go:generate go-bindata -nometadata -ignore .+\.go$ -pkg migrator -o bindata.go ./...
//go:generate gofmt -w bindata.go

const (
	migrationsDir = "migrations"
)

var (
	Migrations *MigrationsLoader
)

var log = logrus.New()

type AssetFn func(name string) ([]byte, error)
type AssetDirFn func(name string) ([]string, error)

func init() {
	Migrations = NewMigrationsLoader()
	if err := Migrations.loadDir(migrationsDir); err != nil {
		log.WithField("service", "load-migrations").WithError(err).Fatal("failed to load migrations")
		return
	}
}
