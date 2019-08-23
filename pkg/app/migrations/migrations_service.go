package migrations

import (
	"fmt"

	"kube-job-runner/pkg/app/config"
	"kube-job-runner/pkg/app/reporter"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/go_bindata"
)

type Service struct {
	Reporter *reporter.Reporter
	Config   config.Config
}

func (service *Service) RunMigrations() error {
	assetSource := bindata.Resource(AssetNames(), func(name string) ([]byte, error) {
		return Asset(name)
	})
	driver, _ := bindata.WithInstance(assetSource)
	postgresURL := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=disable",
		service.Config.PgUser,
		service.Config.PgPassword,
		service.Config.PgHost,
		service.Config.PgPort,
		service.Config.PgDatabase,
	)
	migrater, err := migrate.NewWithSourceInstance("go-bindata", driver, postgresURL)
	if err != nil {
		service.Reporter.Error("migrate.create.error", err, map[string]interface{}{})
		return err
	}
	err = migrater.Up()
	if err != nil {
		service.Reporter.Error("migrate.run.error", err, map[string]interface{}{})
		return err
	}

	return nil
}
