package workforceai

import (
	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk/auth"
)

const (
	policyBasePath = "/app/genai-protect-policy/external/v1"
	assetsBasePath = "/app/genai-protect-assets/external/v1"
	productQuery   = "?product=genai"
)

// Client provides access to the Check Point AI Security API.
type Client struct {
	*sdk.BaseClient
}

// NewClient creates a new Workforce AI API client.
func NewClient(baseURL, token string, authClient *auth.Client, creds auth.Credentials) *Client {
	return &Client{
		BaseClient: sdk.NewBaseClient(baseURL, token, authClient, creds),
	}
}
