package connector

import (
	"autocft/internal/utils"
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/dns"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/packages/pagination"
	"github.com/cloudflare/cloudflare-go/v6/shared"
	"github.com/cloudflare/cloudflare-go/v6/zero_trust"
	"github.com/cloudflare/cloudflare-go/v6/zones"
)

type CloudflareClient struct {
	logger    *slog.Logger
	client    *cloudflare.Client
	accountID string
	tunnelID  string
}

const FallbackService = "http_status:404"

var FallbackIngress = zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress{
	Hostname: cloudflare.F(""),
	Service:  cloudflare.F(FallbackService),
}

func NewCloudflareClient(logger *slog.Logger, apiToken, accountID, tunnelID string) *CloudflareClient {
	client := cloudflare.NewClient(
		option.WithAPIToken(apiToken))
	return &CloudflareClient{
		logger,
		client,
		accountID,
		tunnelID,
	}
}

func (c *CloudflareClient) GetTunnelID() string {
	return c.tunnelID
}

func (c *CloudflareClient) GetTunnelInfo() (*shared.CloudflareTunnel, error) {
	result, err := c.client.ZeroTrust.Tunnels.Cloudflared.Get(context.TODO(), c.tunnelID, zero_trust.TunnelCloudflaredGetParams{
		AccountID: cloudflare.F(c.accountID),
	})
	if c.handleError(err) != nil {
		return nil, err
	}
	return result, nil
}

func (c *CloudflareClient) GetConnection() (*pagination.SinglePage[zero_trust.Client], error) {
	result, err := c.client.ZeroTrust.Tunnels.Cloudflared.Connections.Get(
		context.TODO(),
		c.tunnelID,
		zero_trust.TunnelCloudflaredConnectionGetParams{
			AccountID: cloudflare.F(c.accountID),
		},
	)
	if c.handleError(err) != nil {
		return nil, err
	}
	return result, nil
}

func (c *CloudflareClient) GetConfiguration() (*zero_trust.TunnelCloudflaredConfigurationGetResponse, error) {
	result, err := c.client.ZeroTrust.Tunnels.Cloudflared.Configurations.Get(
		context.Background(),
		c.tunnelID,
		zero_trust.TunnelCloudflaredConfigurationGetParams{
			AccountID: cloudflare.F(c.accountID),
		},
	)
	if c.handleError(err) != nil {
		return nil, err
	}
	return result, nil
}

func (c *CloudflareClient) UpdateConfiguration(ingressConfigs []zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfigIngress) (res *zero_trust.TunnelCloudflaredConfigurationUpdateResponse, err error) {
	result, err := c.client.ZeroTrust.Tunnels.Cloudflared.Configurations.Update(
		context.TODO(),
		c.tunnelID,
		zero_trust.TunnelCloudflaredConfigurationUpdateParams{
			AccountID: cloudflare.F(c.accountID),
			Config: cloudflare.F(zero_trust.TunnelCloudflaredConfigurationUpdateParamsConfig{
				Ingress: cloudflare.F(ingressConfigs),
			}),
		},
	)
	if c.handleError(err) != nil {
		return nil, err
	}
	return result, nil
}

func (c *CloudflareClient) GetZones() ([]zones.Zone, error) {

	zoneList := make([]zones.Zone, 0)

	pager := c.client.Zones.ListAutoPaging(context.TODO(), zones.ZoneListParams{
		Account: cloudflare.F(zones.ZoneListParamsAccount{ID: cloudflare.F(c.accountID)}),
		PerPage: cloudflare.F(utils.PerPage),
	})

	for pager.Next() {
		zone := pager.Current()
		zoneName := utils.NormalizeHostname(zone.Name)
		if zoneName == "" || zone.ID == "" {
			continue
		}
		zoneList = append(zoneList, zone)
	}

	if err := pager.Err(); c.handleError(err) != nil {
		return nil, err
	}
	if len(zoneList) == 0 {
		return nil, fmt.Errorf("no zones available for account %s", c.accountID)
	}
	return zoneList, nil
}

func (c *CloudflareClient) GetRecords(zoneID string) (map[string]*dns.RecordResponse, error) {
	recordMap := make(map[string]*dns.RecordResponse)
	pager := c.client.DNS.Records.ListAutoPaging(context.TODO(), dns.RecordListParams{
		ZoneID:  cloudflare.F(zoneID),
		Type:    cloudflare.F(dns.RecordListParamsTypeCNAME),
		PerPage: cloudflare.F(utils.PerPage),
	})

	for pager.Next() {
		record := pager.Current()
		recordMap[record.Name] = &record
	}

	if err := pager.Err(); c.handleError(err) != nil {
		return nil, err
	}
	return recordMap, nil
}

func (c *CloudflareClient) NewRecord(zoneID string, data dns.CNAMERecordParam) error {
	_, err := c.client.DNS.Records.New(context.TODO(), dns.RecordNewParams{
		ZoneID: cloudflare.F(zoneID),
		Body:   data,
	})
	if c.handleError(err) != nil {
		return err
	}
	return nil
}

func (c *CloudflareClient) UpdateRecord(zoneID, recordID string, data dns.CNAMERecordParam) error {
	_, err := c.client.DNS.Records.Update(context.TODO(), recordID, dns.RecordUpdateParams{
		ZoneID: cloudflare.F(zoneID),
		Body:   data,
	})
	if c.handleError(err) != nil {
		return err
	}
	return nil
}

func (c *CloudflareClient) DeleteRecord(zoneID, recordID string) error {
	_, err := c.client.DNS.Records.Delete(context.TODO(), recordID, dns.RecordDeleteParams{
		ZoneID: cloudflare.F(zoneID),
	})
	if c.handleError(err) != nil {
		return err
	}
	return nil
}

func (c *CloudflareClient) handleError(err error) error {
	if err == nil {
		return nil
	}
	var apiErr *cloudflare.Error
	var errMsg string
	if errors.As(err, &apiErr) {
		c.logger.Debug("connector client error", "details", string(apiErr.DumpRequest(true)))
		errMsg = fmt.Sprintf("connector client Error: %s", apiErr.Error())
	} else {
		errMsg = fmt.Sprintf("connector client Error: %s", err.Error())
	}
	c.logger.Error(errMsg)
	return errors.New(errMsg)
}
