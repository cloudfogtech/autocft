package service

import (
	"autocft/internal/model"
	"encoding/json"
	"os"
	"sort"
	"strings"
)

// sortConfigs sorts rules most-specific first so earlier rules never shadow later ones.
// Wildcard hostnames (*.example.com) match every subdomain, so they must come after all exact hostnames.
func sortConfigs(configs []*model.IngressConfig) {
	sort.Slice(configs, func(i, j int) bool {
		wi := strings.HasPrefix(configs[i].Hostname, "*")
		wj := strings.HasPrefix(configs[j].Hostname, "*")
		if wi != wj {
			return !wi
		}
		if configs[i].Hostname == configs[j].Hostname {
			return configs[i].Path < configs[j].Path
		}
		return configs[i].Hostname < configs[j].Hostname
	})
}

func verifyIngressConfig(c *model.IngressConfig) []string {
	errList := make([]string, 0)
	if c.Service == "" {
		errList = append(errList, "'autocft.service' is required")
	}
	if c.Hostname == "" {
		errList = append(errList, "'autocft.hostname' is required")
	}
	if len(errList) > 0 {
		return errList
	}
	return nil
}

func readHistory(path string) ([]*model.IngressConfig, bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, true, err
	}
	if len(b) == 0 {
		return nil, true, nil
	}
	var list []*model.IngressConfig
	if err = json.Unmarshal(b, &list); err != nil {
		return nil, true, err
	}
	return list, false, nil
}

func writePrettyJSON(path string, v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func toMapByHost(list []*model.IngressConfig) map[string]*model.IngressConfig {
	m := make(map[string]*model.IngressConfig, len(list))
	for _, v := range list {
		if v == nil {
			continue
		}
		h := strings.TrimSpace(v.Hostname)
		if h == "" {
			continue
		}
		m[h] = v
	}
	return m
}

func cloneIngressList(list []*model.IngressConfig) []*model.IngressConfig {
	var out []*model.IngressConfig
	for _, v := range list {
		if v == nil {
			continue
		}
		cp := *v
		if v.Origin != nil {
			o := *v.Origin
			cp.Origin = &o
		}
		out = append(out, &cp)
	}
	return out
}

func ingressDeepEqual(a, b []*model.IngressConfig) bool {
	if len(a) != len(b) {
		return false
	}
	am := toMapByHost(a)
	bm := toMapByHost(b)
	if len(am) != len(bm) {
		return false
	}
	for host, av := range am {
		bv, ok := bm[host]
		if !ok {
			return false
		}
		if !ingressEqual(av, bv) {
			return false
		}
	}
	return true
}

func ingressEqual(x, y *model.IngressConfig) bool {
	if x.Hostname != y.Hostname || x.Service != y.Service || x.Path != y.Path {
		return false
	}
	return originEqual(x.Origin, y.Origin)
}

func originEqual(a, b *model.IngressOriginConfig) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.ConnectTimeout == b.ConnectTimeout &&
		a.DisableChunkedEncoding == b.DisableChunkedEncoding &&
		a.HTTP2Origin == b.HTTP2Origin &&
		a.HTTPHostHeader == b.HTTPHostHeader &&
		a.KeepAliveConnections == b.KeepAliveConnections &&
		a.KeepAliveTimeout == b.KeepAliveTimeout &&
		a.MatchSNItoHost == b.MatchSNItoHost &&
		a.NoHappyEyeballs == b.NoHappyEyeballs &&
		a.NoTLSVerify == b.NoTLSVerify &&
		a.OriginServerName == b.OriginServerName &&
		a.ProxyType == b.ProxyType &&
		a.TCPKeepAlive == b.TCPKeepAlive &&
		a.TLSTimeout == b.TLSTimeout
}
