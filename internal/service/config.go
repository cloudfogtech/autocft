package service

import (
	"autocft/internal/model"
	"autocft/internal/utils"
	"strings"

	"github.com/spf13/viper"
)

const DefaultBasedir = "/app/autocft"
const DefaultCron = "*/10 * * * * *"
const DefaultAdminEmail = "admin@example.com"
const DefaultAdminPassword = "autocft@admin#123"

func LoadConfigFromEnv() (*model.SystemConfig, *model.IngressConfig) {
	viper.SetEnvPrefix("AUTOCFT")
	viper.AutomaticEnv()
	envValueFunc := func(key string) (string, bool) {
		val := viper.GetString(key)
		return val, val != ""
	}
	systemConfig := loadSystemConfigFromEnv(envValueFunc)
	if systemConfig.Basedir == "" {
		systemConfig.Basedir = DefaultBasedir
	} else {
		systemConfig.Basedir = strings.TrimRight(systemConfig.Basedir, "/")
	}
	if systemConfig.Cron == "" {
		systemConfig.Cron = DefaultCron
	}
	if systemConfig.AdminEmail == "" {
		systemConfig.AdminEmail = DefaultAdminEmail
	}
	if systemConfig.AdminPassword == "" {
		systemConfig.AdminPassword = DefaultAdminPassword
	}
	return systemConfig, loadIngressConfigFromEnv(envValueFunc)
}

func loadSystemConfigFromEnv(f func(key string) (string, bool)) *model.SystemConfig {
	systemConfig := &model.SystemConfig{}
	utils.ParseGoTagToStruct("env", f, systemConfig)
	return systemConfig
}

func loadIngressConfigFromEnv(f func(key string) (string, bool)) *model.IngressConfig {
	ingressConfig := &model.IngressConfig{
		Origin: &model.IngressOriginConfig{},
	}
	utils.ParseGoTagToStruct("env", f, ingressConfig)
	utils.ParseGoTagToStruct("env", f, ingressConfig.Origin)
	return ingressConfig
}
