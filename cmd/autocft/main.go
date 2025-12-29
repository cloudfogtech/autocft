package main

import (
	"autocft/internal/model"
	"autocft/internal/service"
	_ "autocft/migrations"
	"log"
	"os"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

func main() {
	systemConfig, ingressConfig := service.LoadConfigFromEnv()
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: systemConfig.Basedir + "/pb_data",
	})
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isGoRun,
	})
	autoCFT := service.NewAutoCFTService(app, systemConfig, ingressConfig)
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Manually trigger a synchronization",
		Run: func(cmd *cobra.Command, args []string) {
			dry, _ := cmd.Flags().GetBool("dry")
			var options *model.SyncOptions
			if dry {
				options = &model.SyncOptions{Dry: true}
			}
			autoCFT.RunSync(options)
		},
	}
	runCmd.Flags().Bool("dry", false, "run in dry-run mode")
	app.RootCmd.AddCommand(runCmd)
	var cronJob *cron.Cron
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		cronJob = cron.New(cron.WithSeconds())
		_, err := cronJob.AddFunc(systemConfig.Cron, func() {
			autoCFT.RunSync()
		})
		if err != nil {
			log.Fatalf("failed to add cron job: %v", err)
			return err
		}
		cronJob.Start()
		return e.Next()
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
	cronJob.Stop()
}
