package migrations

import (
	"autocft/internal/service"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		systemConfig, _ := service.LoadConfigFromEnv()
		superusers, err := app.FindCollectionByNameOrId(core.CollectionNameSuperusers)
		if err != nil {
			return err
		}
		record := core.NewRecord(superusers)
		record.Set("email", systemConfig.AdminEmail)
		record.Set("password", systemConfig.AdminPassword)

		return app.Save(record)
	}, func(app core.App) error {
		systemConfig, _ := service.LoadConfigFromEnv()
		record, _ := app.FindAuthRecordByEmail(core.CollectionNameSuperusers, systemConfig.AdminEmail)
		if record == nil {
			return nil
		}
		return app.Delete(record)
	})
}
