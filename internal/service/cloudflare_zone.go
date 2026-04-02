package service

import (
	"strings"

	"github.com/cloudflare/cloudflare-go/v6/zones"
)

type CloudflareZone struct {
	cfd   *CloudflareDNS
	zones []zones.Zone
}

func (cfd *CloudflareDNS) NewCFZone() *CloudflareZone {
	return &CloudflareZone{
		cfd:   cfd,
		zones: []zones.Zone{},
	}
}

func (cfz *CloudflareZone) Load() error {
	zoneList, err := cfz.cfd.as.cloudflareClient.GetZones()
	if err != nil {
		return err
	}
	cfz.zones = zoneList
	return nil
}

func (cfz *CloudflareZone) GetZone(domain string) (zones.Zone, bool) {
	for _, zone := range cfz.zones {
		if strings.HasSuffix(domain, zone.Name) {
			return zone, true
		}
	}
	return zones.Zone{}, false
}
