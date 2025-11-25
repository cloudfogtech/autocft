package service

import (
	"autocft/internal/connector"
	"autocft/internal/model"
	"log/slog"
	"os"
	"sync/atomic"
	"time"

	"github.com/cloudflare/cloudflare-go/v6/zero_trust"
	"github.com/pocketbase/pocketbase"
)

const HistoryFile = "latest.json"

type AutoCFTService struct {
	app              *pocketbase.PocketBase
	logger           *slog.Logger
	cloudflareClient *connector.CloudflareClient
	dockerClient     *connector.DockerClient
	defaultConfig    *model.IngressConfig
	systemConfig     *model.SystemConfig
	// running as am atomic lock：0 ready，1 running
	running int32
}

func NewAutoCFTService(app *pocketbase.PocketBase, systemConfig *model.SystemConfig, defaultConfig *model.IngressConfig) *AutoCFTService {
	logger := app.Logger().WithGroup("autocft")
	return &AutoCFTService{
		app:              app,
		logger:           logger,
		cloudflareClient: connector.NewCloudflareClient(logger, systemConfig.CFAPIToken, systemConfig.CFAccountID, systemConfig.CFTunnelID),
		dockerClient:     connector.NewDockerClient(),
		defaultConfig:    defaultConfig,
		systemConfig:     systemConfig,
	}
}

func (as *AutoCFTService) RunSync(options ...*model.SyncOptions) bool {
	// try to set running flag to 1，if it was not 0 , there's already a job is running, so skip
	if !atomic.CompareAndSwapInt32(&as.running, 0, 1) {
		as.logger.Debug("Previous sync is still running, skipping this run")
		return false
	}
	// release this flag when exit this function
	defer atomic.StoreInt32(&as.running, 0)

	as.runSyncWithOptions(options...)
	return true
}

func (as *AutoCFTService) runSyncWithOptions(options ...*model.SyncOptions) {
	start := time.Now()
	// Create Base dir
	_ = os.MkdirAll(as.systemConfig.Basedir, 0o755)
	as.logger.Debug("Start sync")

	// 1. Get container info to get container config
	containerConfig, totalContainers, err := as.getContainerConfig()
	if err != nil {
		as.logger.Error("Get container list failed", "error", err)
		return
	}
	as.logger.Debug("Get container config", "current", len(containerConfig), "total", totalContainers)

	// 2. Get history config
	historyConfig, notExists, err := as.getHistoryConfig()
	if err != nil {
		if notExists {
			as.logger.Info("History config is null, it might be the first time run AutoCFT")
		} else {
			as.logger.Warn("Read history ingress config failed", "error", err)
		}
	} else {
		as.logger.Debug("History ingress config", "count", len(historyConfig))
		// If equal, skip sync
		if ingressDeepEqual(containerConfig, historyConfig) {
			as.logger.Debug("No diff, Sync skipped", "cost", time.Since(start))
			return
		}
	}

	// 3. Get cloudflare config (Web managed ingress config+ updated ingress config)
	cloudflareConfig, err := as.getCloudflareConfig()
	if err != nil {
		as.logger.Error("Get cloudflare ingress config failed", "error", err)
		return
	}
	as.logger.Debug("Get cloudflare ingress config", "count", len(cloudflareConfig))

	// 4. Calculate update config
	updateConfig, count := as.calculateUpdateConfig(cloudflareConfig, historyConfig, containerConfig)
	cfUpdateConfig := make([]zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress, 0)
	for _, config := range updateConfig {
		cfUpdateConfig = append(cfUpdateConfig, *configToCFUpdateConfig(as.defaultConfig, config))
	}
	// Fallback ingress config is required for ingress config
	cfUpdateConfig = append(cfUpdateConfig, connector.FallbackIngress)
	if len(options) > 0 && options[0].Dry {
		as.logger.Info("Dry run, Sync skipped", "cost", time.Since(start))
	} else {
		_, err = as.cloudflareClient.UpdateConfiguration(cfUpdateConfig)
		if err != nil {
			as.logger.Error("Update cloudflare config failed", "error", err)
			return
		}
	}
	as.logger.Debug("Finished updating cloudflare config.")

	// 5. Write history config (Only include ingress config managed by container labels)
	if err := writePrettyJSON(as.systemConfig.Basedir+"/"+HistoryFile, containerConfig); err != nil {
		as.logger.Error("Write history failed, please check", "error", err)
	} else {
		as.logger.Debug("Write history success", "path", as.systemConfig.Basedir+"/"+HistoryFile)
	}
	as.logger.Info("Sync Success.",
		"cost", time.Since(start),
		"web managed", count.WebManaged,
		"added", count.Added,
		"updated", count.Updated,
		"deleted", count.Deleted,
		"unchanged", count.Unchanged)
}
