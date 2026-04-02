package service

import "autocft/internal/model"

func (as *AutoCFTService) calculateDiffConfigs(cloudflareConfigs, historyConfigs, containerConfigs []*model.IngressConfig) ([]*model.IngressConfig, []*model.IngressConfig, model.SyncCount) {
	// calculate webManagedIngressConfigs = cloudflareConfigs - historyConfigs (by hostname)
	cloudflareConfigMap := toMapByHost(cloudflareConfigs)
	historyConfigMap := toMapByHost(historyConfigs)
	var webManagedIngressConfigs []*model.IngressConfig
	// For first execution, treat all containerConfig from cloudflare as web managed containerConfig
	if len(historyConfigMap) == 0 {
		webManagedIngressConfigs = cloneIngressList(cloudflareConfigs)
	} else {
		for host, cloudflareConfig := range cloudflareConfigMap {
			// if historyConfigs is not contains a containerConfig, it means it's managed by cloudflare web
			if _, ok := historyConfigMap[host]; !ok {
				webManagedIngressConfigs = append(webManagedIngressConfigs, cloudflareConfig)
			}
		}
	}
	webManagedIngressConfigMap := toMapByHost(webManagedIngressConfigs)
	for _, containerConfig := range containerConfigs {
		webManagedIngressConfigMap[containerConfig.Hostname] = containerConfig
	}
	var updateConfigs []*model.IngressConfig
	for _, v := range webManagedIngressConfigMap {
		updateConfigs = append(updateConfigs, v)
	}
	// sort (Hostname, Path)
	sortConfigs(updateConfigs)
	deletedConfig := calculateDeletedConfig(webManagedIngressConfigs, historyConfigs, updateConfigs)
	return updateConfigs, deletedConfig, calculateCount(webManagedIngressConfigs, historyConfigs, updateConfigs)
}

func calculateCount(webManagedConfigs, historyConfigs, updateConfigs []*model.IngressConfig) model.SyncCount {
	historyConfigMap := toMapByHost(historyConfigs)
	updateConfigMap := toMapByHost(updateConfigs)
	for _, webManagedConfig := range webManagedConfigs {
		delete(updateConfigMap, webManagedConfig.Hostname)
	}
	addedMap := make(map[string]bool)
	unchangedMap := make(map[string]bool)
	updatedMap := make(map[string]bool)
	deletedMap := make(map[string]bool)
	for host, historyConfig := range historyConfigMap {
		if _, ok := updateConfigMap[host]; !ok {
			deletedMap[host] = true
		} else {
			updateConfig := updateConfigMap[host]
			if ingressEqual(historyConfig, updateConfig) {
				unchangedMap[host] = true
			} else {
				updatedMap[host] = true
			}
		}
	}
	for newHost := range updateConfigMap {
		if _, ok := historyConfigMap[newHost]; !ok {
			addedMap[newHost] = true
		}
	}
	return model.SyncCount{
		WebManaged: len(webManagedConfigs),
		Unchanged:  len(unchangedMap),
		Updated:    len(updatedMap),
		Added:      len(addedMap),
		Deleted:    len(deletedMap),
	}
}

func calculateDeletedConfig(webManagedConfigs, historyConfigs, updateConfigs []*model.IngressConfig) []*model.IngressConfig {
	historyConfigMap := toMapByHost(historyConfigs)
	updateConfigMap := toMapByHost(updateConfigs)
	for _, webManagedConfig := range webManagedConfigs {
		delete(updateConfigMap, webManagedConfig.Hostname)
	}

	deletedConfigs := make([]*model.IngressConfig, 0)
	for host, historyConfig := range historyConfigMap {
		if _, ok := updateConfigMap[host]; !ok {
			deletedConfigs = append(deletedConfigs, historyConfig)
		}
	}
	sortConfigs(deletedConfigs)
	return deletedConfigs
}
