package connector

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/option"
	"github.com/cloudflare/cloudflare-go/v6/packages/pagination"
	"github.com/cloudflare/cloudflare-go/v6/shared"
	"github.com/cloudflare/cloudflare-go/v6/zero_trust"
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

func (c *CloudflareClient) handleError(err error) error {
	if err == nil {
		return nil
	}
	var apiErr *cloudflare.Error
	if errors.As(err, &apiErr) {
		c.logger.Debug("connector client error", "details", string(apiErr.DumpRequest(true)))
	}
	errMsg := fmt.Sprintf("connector client Error: %s", apiErr.Error())
	c.logger.Error(errMsg)
	return errors.New(errMsg)
}
