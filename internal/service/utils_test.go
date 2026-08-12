package service

import (
	"autocft/internal/model"
	"reflect"
	"testing"

	"github.com/cloudflare/cloudflare-go/v6/zero_trust"
)

func TestSortConfigsWildcardLast(t *testing.T) {
	in := []*model.IngressConfig{
		{Hostname: "api.example.com", Path: "/"},
		{Hostname: "*.example.com", Path: "/"},
		{Hostname: "app.example.com", Path: "/"},
		{Hostname: "*.other.com", Path: "/"},
	}
	sortConfigs(in)
	want := []string{"api.example.com", "app.example.com", "*.example.com", "*.other.com"}
	got := make([]string, 0, len(in))
	for _, c := range in {
		got = append(got, c.Hostname)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sortConfigs = %v, want %v", got, want)
	}
}

func TestMatchSNItoHostRoundTrip(t *testing.T) {
	cfConfig := &zero_trust.TunnelCloudflaredConfigurationGetResponseConfigIngress{
		Hostname: "*.example.com",
		Service:  "https://10.0.0.1:443",
		OriginRequest: zero_trust.TunnelCloudflaredConfigurationGetResponseConfigIngressOriginRequest{
			MatchSnItoHost: true,
		},
	}
	modelConfig := cloudflareGetConfigToConfig(cfConfig)
	if !modelConfig.Origin.MatchSNItoHost {
		t.Fatal("cloudflareGetConfigToConfig lost MatchSNItoHost")
	}
	updateConfig := configToCloudflareUpdateConfig(&model.IngressConfig{}, modelConfig)
	if !updateConfig.OriginRequest.Value.MatchSnItoHost.Value {
		t.Fatal("configToCloudflareUpdateConfig lost MatchSNItoHost")
	}
	// zero default + true primary merge keeps true
	merged := mergeOriginConfigs(modelConfig.Origin, &model.IngressOriginConfig{})
	if !merged.MatchSNItoHost {
		t.Fatal("mergeOriginConfigs lost MatchSNItoHost")
	}
}
