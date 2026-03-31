package provider

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &WorkforceAIAccessRuleResource{}
	_ resource.ResourceWithImportState = &WorkforceAIAccessRuleResource{}
)

type WorkforceAIAccessRuleResource struct {
	clients *ProviderClients
}

func NewWorkforceAIAccessRuleResource() resource.Resource {
	return &WorkforceAIAccessRuleResource{}
}

func (r *WorkforceAIAccessRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workforce_ai_access_rule"
}

func (r *WorkforceAIAccessRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	attrs := CommonRuleAttributes()
	attrs["source"] = SourceBlockSchema()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Access rule in the Check Point AI Security (WorkforceAI Protect) policy.",
		Attributes:          attrs,
	}
}

func (r *WorkforceAIAccessRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clients, ok := req.ProviderData.(*ProviderClients)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *ProviderClients, got: %T", req.ProviderData))
		return
	}
	r.clients = clients
}

func (r *WorkforceAIAccessRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addReq, err := buildAddRuleRequest(&plan)
	if err != nil {
		resp.Diagnostics.AddError("Error building request", err.Error())
		return
	}

	rule, err := r.clients.WorkforceAI.CreateAccessRule(ctx, addReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating access rule", err.Error())
		return
	}

	tflog.Trace(ctx, "created genai access rule", map[string]interface{}{"rule_id": rule.RuleID})

	// Handle post-create operations: order
	if !plan.Order.IsNull() && !plan.Order.IsUnknown() && int(plan.Order.ValueInt64()) != rule.Order {
		if err := r.clients.WorkforceAI.ReorderRule(ctx, sdk.ReorderRequest{
			RuleID:   rule.RuleID,
			Version:  rule.Version,
			NewOrder: int(plan.Order.ValueInt64()),
		}); err != nil {
			resp.Diagnostics.AddError("Error setting rule order", err.Error())
			return
		}
	}

	// Handle post-create active=false (API creates rules as active by default)
	if !plan.Active.IsNull() && !plan.Active.IsUnknown() && !plan.Active.ValueBool() {
		cur, _ := r.clients.WorkforceAI.GetAccessRule(ctx, rule.RuleID)
		if cur != nil {
			rule.Version = cur.Version
		}
		if err := r.clients.WorkforceAI.SetRuleActive(ctx, sdk.SetActiveRequest{
			RuleID: rule.RuleID, Version: rule.Version, Active: false,
		}); err != nil {
			resp.Diagnostics.AddError("Error disabling rule after create", err.Error())
			return
		}
	}

	// Read back the final state
	created, err := r.clients.WorkforceAI.GetAccessRule(ctx, rule.RuleID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading created rule", err.Error())
		return
	}

	model, err := flattenRuleToModel(created, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping rule to state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *WorkforceAIAccessRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.clients.WorkforceAI.GetAccessRule(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading access rule", err.Error())
		return
	}

	model, err := flattenRuleToModel(rule, &state)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping rule to state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *WorkforceAIAccessRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state RuleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ruleID := state.ID.ValueString()
	version := int(state.Version.ValueInt64())

	// Update info if changed
	if plan.Name.ValueString() != state.Name.ValueString() ||
		plan.Description.ValueString() != state.Description.ValueString() {
		if err := r.clients.WorkforceAI.SetRuleInfo(ctx, sdk.SetInfoRequest{
			RuleID:      ruleID,
			Version:     version,
			Name:        plan.Name.ValueString(),
			Description: plan.Description.ValueString(),
		}); err != nil {
			resp.Diagnostics.AddError("Error updating rule info", err.Error())
			return
		}
		version++
	}

	// Update policy if changed
	if plan.Policy.ValueString() != state.Policy.ValueString() {
		policy, err := expandPolicy(plan.Policy.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error parsing policy", err.Error())
			return
		}
		if err := r.clients.WorkforceAI.SetAccessPolicy(ctx, sdk.SetPolicyRequest{
			RuleID:  ruleID,
			Version: version,
			Policy:  policy,
		}); err != nil {
			resp.Diagnostics.AddError("Error updating policy", err.Error())
			return
		}
		version++
	}

	// Update source if changed
	if sourceChanged(plan.Source, state.Source) {
		if err := r.clients.WorkforceAI.SetRuleSource(ctx, sdk.SetSourceRequest{
			RuleID: ruleID,
			Source: expandSource(plan.Source),
		}); err != nil {
			resp.Diagnostics.AddError("Error updating source", err.Error())
			return
		}
		version++
	}

	// Update active if changed
	if plan.Active.ValueBool() != state.Active.ValueBool() {
		if err := r.clients.WorkforceAI.SetRuleActive(ctx, sdk.SetActiveRequest{
			RuleID:  ruleID,
			Version: version,
			Active:  plan.Active.ValueBool(),
		}); err != nil {
			resp.Diagnostics.AddError("Error updating active state", err.Error())
			return
		}
		version++
	}

	// Update order if changed
	if plan.Order.ValueInt64() != state.Order.ValueInt64() {
		if err := r.clients.WorkforceAI.ReorderRule(ctx, sdk.ReorderRequest{
			RuleID:  ruleID,
			Version: version,
			NewOrder: int(plan.Order.ValueInt64()),
		}); err != nil {
			resp.Diagnostics.AddError("Error updating order", err.Error())
			return
		}
	}

	// Read back final state
	rule, err := r.clients.WorkforceAI.GetAccessRule(ctx, ruleID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading updated rule", err.Error())
		return
	}

	model, err := flattenRuleToModel(rule, &plan)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping rule to state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *WorkforceAIAccessRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RuleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.clients.WorkforceAI.DeleteRule(ctx, state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting access rule", err.Error())
		return
	}
}

func (r *WorkforceAIAccessRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	rule, err := r.clients.WorkforceAI.GetAccessRule(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing access rule", err.Error())
		return
	}

	model, err := flattenRuleToModel(rule, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping imported rule", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
