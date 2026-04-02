package service

import (
	"fmt"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
)

type CloudflareZoneRecord map[string]*dns.RecordResponse
type CloudflareRecord struct {
	cfd     *CloudflareDNS
	records map[string]CloudflareZoneRecord
}

func (cfd *CloudflareDNS) NewCFRecord() *CloudflareRecord {
	return &CloudflareRecord{
		cfd:     cfd,
		records: make(map[string]CloudflareZoneRecord),
	}
}

func (cfr *CloudflareRecord) GetRecords(zoneID string) (CloudflareZoneRecord, error) {
	zoneRecords, ok := cfr.records[zoneID]
	if !ok {
		recordList, err := cfr.cfd.as.cloudflareClient.GetRecords(zoneID)
		if err != nil {
			return nil, err
		}
		zoneRecords = recordList
		cfr.records[zoneID] = zoneRecords
	}
	return zoneRecords, nil
}

func (cfr *CloudflareRecord) Verify(record *dns.RecordResponse) bool {
	return record != nil && record.Type == dns.RecordResponseTypeCNAME
}

func (cfr *CloudflareRecord) Create(zoneID string, hostname string) error {
	return cfr.cfd.as.cloudflareClient.NewRecord(zoneID, dns.CNAMERecordParam{
		Name:    cloudflare.F(hostname),
		Content: cloudflare.F(fmt.Sprintf("%s.cfargotunnel.com", cfr.cfd.as.systemConfig.CFTunnelID)),
		TTL:     cloudflare.F(dns.TTL1),
		Type:    cloudflare.F(dns.CNAMERecordTypeCNAME),
		Proxied: cloudflare.F(true),
		Comment: cloudflare.F(AutoCFTManagedDNSComment),
	})
}

func (cfr *CloudflareRecord) Update(zoneID string, record *dns.RecordResponse) error {
	return cfr.cfd.as.cloudflareClient.UpdateRecord(zoneID, record.ID, dns.CNAMERecordParam{
		Name:    cloudflare.F(record.Name),
		Content: cloudflare.F(fmt.Sprintf("%s.cfargotunnel.com", cfr.cfd.as.systemConfig.CFTunnelID)),
		TTL:     cloudflare.F(dns.TTL1),
		Type:    cloudflare.F(dns.CNAMERecordTypeCNAME),
		Proxied: cloudflare.F(true),
		Comment: cloudflare.F(AutoCFTManagedDNSComment),
	})
}

func (cfr *CloudflareRecord) Delete(zoneID string, recordID string) error {
	return cfr.cfd.as.cloudflareClient.DeleteRecord(zoneID, recordID)
}

func (cfr *CloudflareRecord) Load() {
	cfr.records = make(map[string]CloudflareZoneRecord)
}
