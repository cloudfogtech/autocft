package service

import "errors"

type CloudflareDNS struct {
	as  *AutoCFTService
	cfz *CloudflareZone
	cfr *CloudflareRecord
}

func NewCloudflareDNS(as *AutoCFTService) *CloudflareDNS {
	cfd := &CloudflareDNS{
		as: as,
	}
	cfd.cfz = cfd.NewCFZone()
	cfd.cfr = cfd.NewCFRecord()
	return cfd
}

func (cfd *CloudflareDNS) Load() error {
	cfd.cfr.Load()
	if err := cfd.cfz.Load(); err != nil {
		return errors.New("CF Zone Load Error: " + err.Error())
	}
	return nil
}

func (cfd *CloudflareDNS) Sync(hostname string) error {
	zone, ok := cfd.cfz.GetZone(hostname)
	if !ok {
		return errors.New("Can't found Cloudflare Zone for domain: " + hostname)
	}
	records, err := cfd.cfr.GetRecords(zone.ID)
	if err != nil {
		return err
	}
	record, ok := records[hostname]
	if ok {
		if cfd.cfr.Verify(record) {
			if err = cfd.cfr.Update(zone.ID, record); err != nil {
				return err
			}
		} else {
			if err = cfd.cfr.Delete(zone.ID, record.ID); err != nil {
				return err
			}
			if err = cfd.cfr.Create(zone.ID, hostname); err != nil {
				return err
			}
		}
	} else {
		if err = cfd.cfr.Create(zone.ID, hostname); err != nil {
			return err
		}
	}
	return nil
}

func (cfd *CloudflareDNS) Delete(hostname string) error {
	zone, ok := cfd.cfz.GetZone(hostname)
	if !ok {
		return errors.New("Can't found Cloudflare Zone for domain: " + hostname)
	}
	records, err := cfd.cfr.GetRecords(zone.ID)
	if err != nil {
		return err
	}
	record, ok := records[hostname]
	if ok {
		return cfd.cfr.Delete(zone.ID, record.ID)
	}
	return nil
}
