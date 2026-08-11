package service

import (
	"autocft/internal/connector"
	"autocft/internal/model"
	"fmt"
	"strings"
)

func (as *AutoCFTService) getCloudflareConfig() ([]*model.IngressConfig, error) {
	cfg, err := as.cloudflareClient.GetConfiguration()
	if err != nil {
		return nil, fmt.Errorf("get cloudflare ingress config failed: %w", err)
	}
	cloudflareConfig := make([]*model.IngressConfig, 0)
	for _, ingress := range cfg.Config.Ingress {
		if ingress.Service == connector.FallbackService && ingress.Hostname == "" {
			continue
		}
		cloudflareConfig = append(cloudflareConfig, cloudflareGetConfigToConfig(&ingress))
	}
	return cloudflareConfig, nil
}

func (as *AutoCFTService) getHistoryConfig() ([]*model.IngressConfig, bool, error) {
	return readHistory(as.systemConfig.Basedir + "/" + HistoryFile)
}

func (as *AutoCFTService) getContainerConfig() ([]*model.IngressConfig, int, error) {
	containers, err := as.dockerClient.GetContainers()
	if err != nil {
		return nil, 0, fmt.Errorf("get container list failed: %w", err)
	}
	containerConfig := make([]*model.IngressConfig, 0)
	for _, ct := range containers {
		if len(ct.Labels) == 0 {
			continue
		}
		ing := containerLabelsToConfig(ct.Labels)
		if !ing.Enabled {
			as.logger.Debug("Container skipped", "name", strings.Join(ct.Names, ","), "id", ct.ID[:12])
			continue
		}
		if errList := verifyIngressConfig(ing); errList != nil {
			as.logger.Warn("Verifying ingress config failed, Container skipped",
				"name", strings.Join(ct.Names, ","),
				"id", ct.ID[:12])
			for _, err := range errList {
				as.logger.Warn(fmt.Sprintf("    %s", err))
			}
			continue
		}
		containerConfig = append(containerConfig, ing)
	}
	return containerConfig, len(containers), nil
}

func (as *AutoCFTService) calculateUpdateConfig(cloudflareConfig, historyConfig, containerConfig []*model.IngressConfig) ([]*model.IngressConfig, model.SyncCount) {
	// calculate webManagedIngressConfig = cloudflareConfig - historyConfig (by hostname)
	cfConfigmap := toMapByHost(cloudflareConfig)
	historyConfigmap := toMapByHost(historyConfig)
	var webManagedIngressConfig []*model.IngressConfig
	// For first execution, treat all config from cloudflare as web managed config
	if len(historyConfigmap) == 0 {
		webManagedIngressConfig = cloneIngressList(cloudflareConfig)
	} else {
		for host, cfConfig := range cfConfigmap {
			// if historyConfig is not contains a config, it means it's managed by cloudflare web
			if _, ok := historyConfigmap[host]; !ok {
				webManagedIngressConfig = append(webManagedIngressConfig, cfConfig)
			}
		}
	}
	updateMap := toMapByHost(webManagedIngressConfig)
	for _, config := range containerConfig {
		updateMap[config.Hostname] = config
	}
	var updateConfig []*model.IngressConfig
	for _, v := range updateMap {
		updateConfig = append(updateConfig, v)
	}
	sortConfigs(updateConfig)
	return updateConfig, calculateCount(webManagedIngressConfig, historyConfig, updateConfig)
}
