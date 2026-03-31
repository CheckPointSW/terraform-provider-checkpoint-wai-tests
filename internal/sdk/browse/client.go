package browse

import (
	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk/auth"
)

const (
	policyBasePath = "/app/browse-protect-policy/external/v1"
	assetsBasePath = "/app/browse-protect-assets/external/v1"
	productQuery   = "?product=browse"
)

// Client provides access to the Check Point Browse Security API.
type Client struct {
	*sdk.BaseClient
}

// NewClient creates a new Browse Security API client.
func NewClient(baseURL, token string, authClient *auth.Client, creds auth.Credentials) *Client {
	return &Client{
		BaseClient: sdk.NewBaseClient(baseURL, token, authClient, creds),
	}
}
