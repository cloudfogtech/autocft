package service

import (
	"autocft/internal/model"
	"autocft/internal/utils"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/zero_trust"
)

func containerLabelsToConfig(labels map[string]string) *model.IngressConfig {
	ingressConfig := &model.IngressConfig{}
	ingressConfig.Origin = &model.IngressOriginConfig{}
	valueFunc := func(key string) (string, bool) {
		val, ok := labels[key]
		return val, ok
	}
	utils.ParseGoTagToStruct("tag", valueFunc, ingressConfig)
	utils.ParseGoTagToStruct("tag", valueFunc, ingressConfig.Origin)

	return ingressConfig
}

func cfGetConfigToConfig(cfConfig *zero_trust.TunnelCloudflaredConfigurationGetResponseConfigIngress) *model.IngressConfig {
	if cfConfig == nil {
		return nil
	}

	return &model.IngressConfig{
		Enabled:  true,
		Hostname: cfConfig.Hostname,
		Service:  cfConfig.Service,
		Path:     cfConfig.Path,
		Origin: &model.IngressOriginConfig{
			ConnectTimeout:         cfConfig.OriginRequest.ConnectTimeout,
			DisableChunkedEncoding: cfConfig.OriginRequest.DisableChunkedEncoding,
			HTTP2Origin:            cfConfig.OriginRequest.HTTP2Origin,
			HTTPHostHeader:         cfConfig.OriginRequest.HTTPHostHeader,
			KeepAliveConnections:   cfConfig.OriginRequest.KeepAliveConnections,
			KeepAliveTimeout:       cfConfig.OriginRequest.KeepAliveTimeout,
			NoHappyEyeballs:        cfConfig.OriginRequest.NoHappyEyeballs,
			NoTLSVerify:            cfConfig.OriginRequest.NoTLSVerify,
			OriginServerName:       cfConfig.OriginRequest.OriginServerName,
			ProxyType:              cfConfig.OriginRequest.ProxyType,
			TCPKeepAlive:           cfConfig.OriginRequest.TCPKeepAlive,
			TLSTimeout:             cfConfig.OriginRequest.TLSTimeout,
		},
	}
}

func configToCFUpdateConfig(defaultIngressConfig *model.IngressConfig, updateIngressConfig *model.IngressConfig) *zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress {
	ingressConfig := mergeIngressConfigs(updateIngressConfig, defaultIngressConfig)

	// 构建 Cloudflare API 参数
	cfConfig := &zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress{
		Hostname: cloudflare.F(ingressConfig.Hostname),
		Service:  cloudflare.F(ingressConfig.Service),
		Path:     cloudflare.F(ingressConfig.Path),
		OriginRequest: cloudflare.F(zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngressOriginRequest{
			//Access: cloudflare.F(ingressConfig.Origin.Access),
			//CAPool: cloudflare.F(ingressConfig.Origin.CAPool),
			ConnectTimeout:         cloudflare.F(ingressConfig.Origin.ConnectTimeout),
			DisableChunkedEncoding: cloudflare.F(ingressConfig.Origin.DisableChunkedEncoding),
			HTTP2Origin:            cloudflare.F(ingressConfig.Origin.HTTP2Origin),
			HTTPHostHeader:         cloudflare.F(ingressConfig.Origin.HTTPHostHeader),
			KeepAliveConnections:   cloudflare.F(ingressConfig.Origin.KeepAliveConnections),
			KeepAliveTimeout:       cloudflare.F(ingressConfig.Origin.KeepAliveTimeout),
			NoHappyEyeballs:        cloudflare.F(ingressConfig.Origin.NoHappyEyeballs),
			NoTLSVerify:            cloudflare.F(ingressConfig.Origin.NoTLSVerify),
			OriginServerName:       cloudflare.F(ingressConfig.Origin.OriginServerName),
			ProxyType:              cloudflare.F(ingressConfig.Origin.ProxyType),
			TCPKeepAlive:           cloudflare.F(ingressConfig.Origin.TCPKeepAlive),
			TLSTimeout:             cloudflare.F(ingressConfig.Origin.TLSTimeout),
		}),
	}
	return cfConfig
}

// 合并两个 IngressConfig
func mergeIngressConfigs(primary, fallback *model.IngressConfig) *model.IngressConfig {
	if primary == nil && fallback == nil {
		return &model.IngressConfig{}
	}
	if primary == nil {
		primary = &model.IngressConfig{}
	}
	if fallback == nil {
		fallback = &model.IngressConfig{}
	}

	return &model.IngressConfig{
		Hostname: utils.MergeStringField(primary.Hostname, fallback.Hostname),
		Service:  utils.MergeStringField(primary.Service, fallback.Service),
		Path:     utils.MergeStringField(primary.Path, fallback.Path),
		Origin:   mergeOriginConfigs(primary.Origin, fallback.Origin),
	}
}

// 合并两个 IngressOriginConfig
func mergeOriginConfigs(primary, fallback *model.IngressOriginConfig) *model.IngressOriginConfig {
	if primary == nil && fallback == nil {
		return &model.IngressOriginConfig{}
	}
	if primary == nil {
		primary = &model.IngressOriginConfig{}
	}
	if fallback == nil {
		fallback = &model.IngressOriginConfig{}
	}
	return &model.IngressOriginConfig{
		ConnectTimeout:         utils.MergeInt64Field(primary.ConnectTimeout, fallback.ConnectTimeout),
		DisableChunkedEncoding: utils.MergeBoolField(primary.DisableChunkedEncoding, fallback.DisableChunkedEncoding),
		HTTP2Origin:            utils.MergeBoolField(primary.HTTP2Origin, fallback.HTTP2Origin),
		HTTPHostHeader:         utils.MergeStringField(primary.HTTPHostHeader, fallback.HTTPHostHeader),
		KeepAliveConnections:   utils.MergeInt64Field(primary.KeepAliveConnections, fallback.KeepAliveConnections),
		KeepAliveTimeout:       utils.MergeInt64Field(primary.KeepAliveTimeout, fallback.KeepAliveTimeout),
		NoHappyEyeballs:        utils.MergeBoolField(primary.NoHappyEyeballs, fallback.NoHappyEyeballs),
		NoTLSVerify:            utils.MergeBoolField(primary.NoTLSVerify, fallback.NoTLSVerify),
		OriginServerName:       utils.MergeStringField(primary.OriginServerName, fallback.OriginServerName),
		ProxyType:              utils.MergeStringField(primary.ProxyType, fallback.ProxyType),
		TCPKeepAlive:           utils.MergeInt64Field(primary.TCPKeepAlive, fallback.TCPKeepAlive),
		TLSTimeout:             utils.MergeInt64Field(primary.TLSTimeout, fallback.TLSTimeout),
	}
}
