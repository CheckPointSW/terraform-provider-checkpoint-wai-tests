package provider

import (
	"context"
	"os"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk/auth"
	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk/browse"
	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk/workforceai"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &CheckPointWorkforceAIProvider{}

// CheckPointWorkforceAIProvider defines the provider implementation.
type CheckPointWorkforceAIProvider struct {
	version string
}

// CheckPointWorkforceAIProviderModel describes the provider data model.
type CheckPointWorkforceAIProviderModel struct {
	Region    types.String `tfsdk:"region"`
	ClientID  types.String `tfsdk:"client_id"`
	AccessKey types.String `tfsdk:"access_key"`
}

// ProviderClients holds the SDK clients shared across resources.
type ProviderClients struct {
	WorkforceAI  *workforceai.Client
	Browse *browse.Client
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CheckPointWorkforceAIProvider{
			version: version,
		}
	}
}

func (p *CheckPointWorkforceAIProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cpwai"
	resp.Version = p.version
}

func (p *CheckPointWorkforceAIProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Check Point Workforce AI provider manages security policies for " +
			"**WorkforceAI Protect** and **Browse Protect** products via the Check Point Infinity Portal API.\n\n" +
			"## Authentication\n\n" +
			"Create an API key in the [Infinity Portal](https://portal.checkpoint.com/dashboard/settings/api-keys) " +
			"(New → New Account API Key → Service: WorkforceAI Protect). This gives you a **Client ID** and **Access Key**.\n\n" +
			"Configure them via provider attributes or environment variables:\n" +
			"- `CHECKPOINT_INFINITY_CLIENT_ID`\n" +
			"- `CHECKPOINT_INFINITY_ACCESS_KEY`\n" +
			"- `CHECKPOINT_INFINITY_REGION` (optional, defaults to `eu`)\n\n" +
			"## Products\n\n" +
			"| Product | Resources | Description |\n" +
			"|---------|-----------|-------------|\n" +
			"| WorkforceAI Protect | `cpwai_workforce_ai_chats_rule`, `cpwai_workforce_ai_access_rule`, `cpwai_workforce_ai_agents_rule` | DLP, access control, and MCP/agent policies for WorkforceAI services |\n" +
			"| Browse Protect | `cpwai_browse_dlp_rule`, `cpwai_browse_access_rule`, `cpwai_browse_secure_browsing_rule` | Same + secure browsing policies for web traffic |\n",
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				MarkdownDescription: "Infinity Portal API region. Determines the API gateway URL:\n" +
					"  - `eu` → `cloudinfra-gw.portal.checkpoint.com`\n" +
					"  - `us` → `cloudinfra-gw-us.portal.checkpoint.com`\n\n" +
					"Can also be set via `CHECKPOINT_INFINITY_REGION` env var. Defaults to `eu`.",
				Optional: true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "Client ID from the Infinity Portal API key. " +
					"Created at **Infinity Portal → Settings → API Keys**. " +
					"Can also be set via `CHECKPOINT_INFINITY_CLIENT_ID` env var.",
				Optional: true,
			},
			"access_key": schema.StringAttribute{
				MarkdownDescription: "Access key (secret) from the Infinity Portal API key. " +
					"Can also be set via `CHECKPOINT_INFINITY_ACCESS_KEY` env var.",
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *CheckPointWorkforceAIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config CheckPointWorkforceAIProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve client_id
	clientID := os.Getenv("CHECKPOINT_INFINITY_CLIENT_ID")
	if !config.ClientID.IsNull() {
		clientID = config.ClientID.ValueString()
	}
	if clientID == "" {
		resp.Diagnostics.AddError(
			"Missing Client ID",
			"The provider requires a client ID. Set it via the `client_id` attribute or the `CHECKPOINT_INFINITY_CLIENT_ID` environment variable.",
		)
		return
	}

	// Resolve access_key
	accessKey := os.Getenv("CHECKPOINT_INFINITY_ACCESS_KEY")
	if !config.AccessKey.IsNull() {
		accessKey = config.AccessKey.ValueString()
	}
	if accessKey == "" {
		resp.Diagnostics.AddError(
			"Missing Access Key",
			"The provider requires an access key. Set it via the `access_key` attribute or the `CHECKPOINT_INFINITY_ACCESS_KEY` environment variable.",
		)
		return
	}

	// Resolve region
	region := "eu"
	if envRegion := os.Getenv("CHECKPOINT_INFINITY_REGION"); envRegion != "" {
		region = envRegion
	}
	if !config.Region.IsNull() {
		region = config.Region.ValueString()
	}

	var baseURL string
	switch region {
	case "eu":
		baseURL = "https://cloudinfra-gw.portal.checkpoint.com"
	case "us":
		baseURL = "https://cloudinfra-gw-us.portal.checkpoint.com"
	default:
		resp.Diagnostics.AddError(
			"Invalid Region",
			"Region must be 'eu' or 'us', got: "+region,
		)
		return
	}

	creds := auth.Credentials{
		ClientID:  clientID,
		AccessKey: accessKey,
	}

	// Create auth client and obtain token
	authClient := auth.NewClient(baseURL)
	token, err := authClient.GetToken(ctx, creds)
	if err != nil {
		resp.Diagnostics.AddError(
			"Authentication Failed",
			"Failed to obtain JWT token from CloudInfra: "+err.Error(),
		)
		return
	}

	clients := &ProviderClients{
		WorkforceAI:  workforceai.NewClient(baseURL, token, authClient, creds),
		Browse: browse.NewClient(baseURL, token, authClient, creds),
	}

	resp.DataSourceData = clients
	resp.ResourceData = clients
}

func (p *CheckPointWorkforceAIProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWorkforceAIChatsRuleResource,
		NewWorkforceAIAccessRuleResource,
		NewWorkforceAIAgentsRuleResource,
		NewBrowseDlpRuleResource,
		NewBrowseAccessRuleResource,
		NewBrowseSecureBrowsingRuleResource,
		NewBrowseFileProtectionObjectResource,
		NewBrowseDomainsObjectResource,
	}
}

func (p *CheckPointWorkforceAIProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewWorkforceAIChatsRulebaseDataSource,
		NewWorkforceAIAccessRulebaseDataSource,
		NewWorkforceAIAgentsRulebaseDataSource,
		NewBrowseDlpRulebaseDataSource,
		NewBrowseAccessRulebaseDataSource,
		NewBrowseSecureBrowsingRulebaseDataSource,
	}
}
