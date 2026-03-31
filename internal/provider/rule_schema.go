package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// RuleResourceModel is the common Terraform model for all rule resources.
type RuleResourceModel struct {
	ID          types.String    `tfsdk:"id"`
	Name        types.String    `tfsdk:"name"`
	Description types.String    `tfsdk:"description"`
	Order       types.Int64     `tfsdk:"order"`
	Active      types.Bool      `tfsdk:"active"`
	Version     types.Int64     `tfsdk:"version"`
	Policy      types.String    `tfsdk:"policy"`          // JSON string
	Source      []SourceModel   `tfsdk:"source"`
	Objects     types.Map       `tfsdk:"objects"`
}

// SourceModel represents a user/group source entry in Terraform.
type SourceModel struct {
	AssignmentID   types.String `tfsdk:"assignment_id"`
	DisplayName    types.String `tfsdk:"display_name"`
	AssignmentType types.String `tfsdk:"assignment_type"`
}

// CommonRuleAttributes returns the shared schema attributes for all rule resources.
func CommonRuleAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Unique identifier of the rule (UUID). Assigned by the server on creation.",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "Human-readable rule name. Used for display in the Infinity Portal UI and Terraform state.",
			Required:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "Free-text explanation of the rule's purpose.",
			Optional:            true,
			Computed:            true,
			Default:             stringdefault.StaticString(""),
		},
		"order": schema.Int64Attribute{
			MarkdownDescription: "Position in the rulebase that determines evaluation priority. **Lower values are evaluated first**; when two rules match the same request, the lower-order rule takes precedence. The first rule should be `0`.",
			Optional:            true,
			Computed:            true,
		},
		"active": schema.BoolAttribute{
			MarkdownDescription: "Whether the rule is currently enforced. Set to `false` to disable the rule without deleting it. Defaults to `true`.",
			Optional:            true,
			Computed:            true,
			Default:             booldefault.StaticBool(true),
		},
		"version": schema.Int64Attribute{
			MarkdownDescription: "Rule version number, incremented on every modification. Used internally for optimistic concurrency control.",
			Computed:            true,
		},
		"policy": schema.StringAttribute{
			MarkdownDescription: "The enforcement policy configuration as a JSON string (use `jsonencode({...})`). " +
				"Structure varies by rule type — see the [Policy Reference](../guides/policy-reference.md) for the full schema per rule type.\n\n" +
				"**Chats rules** (DLP): `event_type`, `services_and_application`, `data_types`, `action`, `logging`\n\n" +
				"**Access rules**: `services_and_application`, `action`, `logging`\n\n" +
				"**Agents rules** (MCP): `clients`, `servers`, `tooling`, `action`, `logging`\n\n" +
				"**Secure Browsing rules**: `password_reuse`, `zero_phishing`, `safe_search`, `search_reputation`, `domain_exclusions`, `protected_domains`",
			Required: true,
		},
		"objects": schema.MapAttribute{
			MarkdownDescription: "Policy objects attached to this rule. A map of feature type to a list of object IDs.\n\n" +
				"**Access rules**: `FILE_EMULATION`, `DOMAIN_SERVICES`\n\n" +
				"**Chats rules**: `DOMAIN_SERVICES`\n\n" +
				"**Secure Browsing rules**: `PROTECTED_DOMAINS`, `EXCLUDE_DOMAINS`",
			Optional:    true,
			ElementType: types.ListType{ElemType: types.StringType},
		},
	}
}

// SourceBlockSchema returns the schema for the source block.
func SourceBlockSchema() schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		MarkdownDescription: "List of users, groups, or the entire organization this rule applies to (source entries). " +
			"An empty list means the rule applies to everyone in the tenant. For org-wide rules, use `assignment_type = \"ASSIGNMENT_TYPE_ENTIRE_ORG\"`.",
		Optional: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"assignment_id": schema.StringAttribute{
					MarkdownDescription: "UUID of the user or group. Not required when `assignment_type` is `ASSIGNMENT_TYPE_ENTIRE_ORG`.",
					Optional:            true,
				},
				"display_name": schema.StringAttribute{
					MarkdownDescription: "Display name of the user or group. Computed by the server if not provided.",
					Optional:            true,
					Computed:            true,
				},
				"assignment_type": schema.StringAttribute{
					MarkdownDescription: "Assignment scope. Valid values:\n" +
						"  - `ASSIGNMENT_TYPE_USER` — individual user\n" +
						"  - `ASSIGNMENT_TYPE_GROUP` — user group\n" +
						"  - `ASSIGNMENT_TYPE_ENTIRE_ORG` — entire organization",
					Required: true,
				},
			},
		},
	}
}
