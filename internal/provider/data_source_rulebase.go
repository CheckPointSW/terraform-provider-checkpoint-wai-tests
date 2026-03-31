package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RulebaseDataSourceModel is the shared model for all rulebase data sources.
type RulebaseDataSourceModel struct {
	RulebaseVersion types.Int64  `tfsdk:"rulebase_version"`
	Rules           types.String `tfsdk:"rules"` // JSON array of rules
}

func rulebaseDataSourceSchema(description string) schema.Schema {
	return schema.Schema{
		MarkdownDescription: description,
		Attributes: map[string]schema.Attribute{
			"rulebase_version": schema.Int64Attribute{
				MarkdownDescription: "Monotonically increasing version of this rulebase.",
				Computed:            true,
			},
			"rules": schema.StringAttribute{
				MarkdownDescription: "JSON array of all rules in the rulebase.",
				Computed:            true,
			},
		},
	}
}

func flattenRulebaseModel(rb *sdk.Rulebase) (*RulebaseDataSourceModel, error) {
	rulesJSON, err := json.Marshal(rb.Rules)
	if err != nil {
		return nil, fmt.Errorf("serializing rules: %w", err)
	}
	return &RulebaseDataSourceModel{
		RulebaseVersion: types.Int64Value(int64(rb.RulebaseVersion)),
		Rules:           types.StringValue(string(rulesJSON)),
	}, nil
}

// --- WorkforceAI Chats Rulebase ---

type WorkforceAIChatsRulebaseDataSource struct{ clients *ProviderClients }

func NewWorkforceAIChatsRulebaseDataSource() datasource.DataSource {
	return &WorkforceAIChatsRulebaseDataSource{}
}
func (d *WorkforceAIChatsRulebaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workforce_ai_chats_rulebase"
}
func (d *WorkforceAIChatsRulebaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = rulebaseDataSourceSchema("Reads the WorkforceAI Protect Chats rulebase.")
}
func (d *WorkforceAIChatsRulebaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*ProviderClients)
}
func (d *WorkforceAIChatsRulebaseDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	rb, err := d.clients.WorkforceAI.GetChatsRulebase(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading chats rulebase", err.Error())
		return
	}
	model, err := flattenRulebaseModel(rb)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping rulebase", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// --- WorkforceAI Access Rulebase ---

type WorkforceAIAccessRulebaseDataSource struct{ clients *ProviderClients }

func NewWorkforceAIAccessRulebaseDataSource() datasource.DataSource {
	return &WorkforceAIAccessRulebaseDataSource{}
}
func (d *WorkforceAIAccessRulebaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workforce_ai_access_rulebase"
}
func (d *WorkforceAIAccessRulebaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = rulebaseDataSourceSchema("Reads the WorkforceAI Protect Access rulebase.")
}
func (d *WorkforceAIAccessRulebaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*ProviderClients)
}
func (d *WorkforceAIAccessRulebaseDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	rb, err := d.clients.WorkforceAI.GetAccessRulebase(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading access rulebase", err.Error())
		return
	}
	model, err := flattenRulebaseModel(rb)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping rulebase", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// --- WorkforceAI Agents Rulebase ---

type WorkforceAIAgentsRulebaseDataSource struct{ clients *ProviderClients }

func NewWorkforceAIAgentsRulebaseDataSource() datasource.DataSource {
	return &WorkforceAIAgentsRulebaseDataSource{}
}
func (d *WorkforceAIAgentsRulebaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workforce_ai_agents_rulebase"
}
func (d *WorkforceAIAgentsRulebaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = rulebaseDataSourceSchema("Reads the WorkforceAI Protect Agents rulebase.")
}
func (d *WorkforceAIAgentsRulebaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*ProviderClients)
}
func (d *WorkforceAIAgentsRulebaseDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	rb, err := d.clients.WorkforceAI.GetAgentsRulebase(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading agents rulebase", err.Error())
		return
	}
	model, err := flattenRulebaseModel(rb)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping rulebase", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// --- Browse DLP Rulebase ---

type BrowseDlpRulebaseDataSource struct{ clients *ProviderClients }

func NewBrowseDlpRulebaseDataSource() datasource.DataSource {
	return &BrowseDlpRulebaseDataSource{}
}
func (d *BrowseDlpRulebaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_browse_dlp_rulebase"
}
func (d *BrowseDlpRulebaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = rulebaseDataSourceSchema("Reads the Browse Security DLP rulebase.")
}
func (d *BrowseDlpRulebaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*ProviderClients)
}
func (d *BrowseDlpRulebaseDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	rb, err := d.clients.Browse.GetChatsRulebase(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading chats rulebase", err.Error())
		return
	}
	model, err := flattenRulebaseModel(rb)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping rulebase", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// --- Browse Access Rulebase ---

type BrowseAccessRulebaseDataSource struct{ clients *ProviderClients }

func NewBrowseAccessRulebaseDataSource() datasource.DataSource {
	return &BrowseAccessRulebaseDataSource{}
}
func (d *BrowseAccessRulebaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_browse_access_rulebase"
}
func (d *BrowseAccessRulebaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = rulebaseDataSourceSchema("Reads the Browse Security Access rulebase.")
}
func (d *BrowseAccessRulebaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*ProviderClients)
}
func (d *BrowseAccessRulebaseDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	rb, err := d.clients.Browse.GetAccessRulebase(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading access rulebase", err.Error())
		return
	}
	model, err := flattenRulebaseModel(rb)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping rulebase", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

// --- Browse Secure Browsing Rulebase ---

type BrowseSecureBrowsingRulebaseDataSource struct{ clients *ProviderClients }

func NewBrowseSecureBrowsingRulebaseDataSource() datasource.DataSource {
	return &BrowseSecureBrowsingRulebaseDataSource{}
}
func (d *BrowseSecureBrowsingRulebaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_browse_secure_browsing_rulebase"
}
func (d *BrowseSecureBrowsingRulebaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = rulebaseDataSourceSchema("Reads the Browse Security Secure Browsing rulebase.")
}
func (d *BrowseSecureBrowsingRulebaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*ProviderClients)
}
func (d *BrowseSecureBrowsingRulebaseDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	rb, err := d.clients.Browse.GetSecureBrowsingRulebase(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading secure browsing rulebase", err.Error())
		return
	}
	model, err := flattenRulebaseModel(rb)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping rulebase", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
