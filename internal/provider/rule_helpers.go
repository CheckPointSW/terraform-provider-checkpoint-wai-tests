package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// expandSource converts Terraform source models to SDK assignments.
// For ASSIGNMENT_TYPE_ENTIRE_ORG, the API still requires assignment_id and display_name
// so we provide sensible defaults.
func expandSource(models []SourceModel) []sdk.Assignment {
	if len(models) == 0 {
		return nil
	}
	assignments := make([]sdk.Assignment, len(models))
	for i, m := range models {
		a := sdk.Assignment{
			AssignmentID:   m.AssignmentID.ValueString(),
			DisplayName:    m.DisplayName.ValueString(),
			AssignmentType: m.AssignmentType.ValueString(),
		}
		// API requires these fields even for entire org assignments
		if a.AssignmentType == "ASSIGNMENT_TYPE_ENTIRE_ORG" {
			if a.AssignmentID == "" {
				a.AssignmentID = "entire-org"
			}
			if a.DisplayName == "" {
				a.DisplayName = "Entire Organization"
			}
		}
		assignments[i] = a
	}
	return assignments
}

// flattenSource converts SDK assignments to Terraform source models,
// preserving null semantics from the plan to avoid inconsistent state.
func flattenSource(assignments []sdk.Assignment, plan []SourceModel) []SourceModel {
	if len(assignments) == 0 {
		return plan
	}
	models := make([]SourceModel, len(assignments))
	for i, a := range assignments {
		m := SourceModel{
			AssignmentType: types.StringValue(a.AssignmentType),
		}
		// Preserve null vs empty from plan for optional fields
		if i < len(plan) {
			if plan[i].AssignmentID.IsNull() {
				m.AssignmentID = types.StringNull()
			} else {
				m.AssignmentID = types.StringValue(a.AssignmentID)
			}
			if plan[i].DisplayName.IsNull() {
				m.DisplayName = types.StringNull()
			} else {
				m.DisplayName = types.StringValue(a.DisplayName)
			}
		} else {
			m.AssignmentID = types.StringValue(a.AssignmentID)
			m.DisplayName = types.StringValue(a.DisplayName)
		}
		models[i] = m
	}
	return models
}

// stripJSONNulls removes keys with null values from a JSON byte slice.
// This prevents Terraform from detecting spurious diffs when the API returns
// explicit nulls for fields the user did not set.
func stripJSONNulls(data []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return data, nil // not an object, return as-is
	}
	for k, v := range m {
		if v == nil {
			delete(m, k)
		}
	}
	return json.Marshal(m)
}

// sourceChanged returns true if the plan source differs from the state source.
func sourceChanged(plan, state []SourceModel) bool {
	if len(plan) != len(state) {
		return true
	}
	for i := range plan {
		if plan[i].AssignmentType.ValueString() != state[i].AssignmentType.ValueString() {
			return true
		}
		if plan[i].AssignmentID.ValueString() != state[i].AssignmentID.ValueString() {
			return true
		}
	}
	return false
}

// expandPolicy converts a JSON string to a map for the API.
func expandPolicy(policyJSON string) (map[string]interface{}, error) {
	var policy map[string]interface{}
	if err := json.Unmarshal([]byte(policyJSON), &policy); err != nil {
		return nil, fmt.Errorf("parsing policy JSON: %w", err)
	}
	return policy, nil
}

// flattenPolicy converts a policy map to a JSON string for Terraform state.
func flattenPolicy(policy map[string]interface{}) (string, error) {
	bytes, err := json.Marshal(policy)
	if err != nil {
		return "", fmt.Errorf("serializing policy: %w", err)
	}
	return string(bytes), nil
}

// expandObjects converts Terraform map to SDK objects map.
func expandObjects(ctx context.Context, m types.Map) (map[string][]string, error) {
	if m.IsNull() || m.IsUnknown() {
		return nil, nil
	}
	result := make(map[string][]string)
	elements := m.Elements()
	for k, v := range elements {
		listVal, ok := v.(types.List)
		if !ok {
			continue
		}
		var ids []string
		diags := listVal.ElementsAs(ctx, &ids, false)
		if diags.HasError() {
			return nil, fmt.Errorf("expanding objects list for key %s", k)
		}
		result[k] = ids
	}
	return result, nil
}

// flattenRuleToModel converts an SDK Rule to the Terraform RuleResourceModel,
// preserving plan/state values for fields the API doesn't round-trip cleanly.
func flattenRuleToModel(rule *sdk.Rule, plan *RuleResourceModel) (*RuleResourceModel, error) {
	model := &RuleResourceModel{
		ID:      types.StringValue(rule.RuleID),
		Name:    types.StringValue(rule.Name),
		Order:   types.Int64Value(int64(rule.Order)),
		Active:  types.BoolValue(rule.Active),
		Version: types.Int64Value(int64(rule.Version)),
	}

	// Description: use API value
	if rule.Description != "" {
		model.Description = types.StringValue(rule.Description)
	} else {
		model.Description = types.StringValue("")
	}

	// Policy: preserve plan value to avoid API round-trip differences (e.g. wrapping in "default" key)
	if plan != nil {
		model.Policy = plan.Policy
	} else {
		policyJSON, err := flattenPolicy(rule.Policy)
		if err != nil {
			return nil, err
		}
		model.Policy = types.StringValue(policyJSON)
	}

	// Source: flatten from API but preserve null semantics from plan
	var planSource []SourceModel
	if plan != nil {
		planSource = plan.Source
	}
	model.Source = flattenSource(rule.Source, planSource)

	// Objects: preserve plan value
	if plan != nil {
		model.Objects = plan.Objects
	} else {
		model.Objects = types.MapNull(types.ListType{ElemType: types.StringType})
	}

	return model, nil
}

// buildAddRuleRequest constructs an AddRuleRequest from Terraform model.
func buildAddRuleRequest(model *RuleResourceModel) (sdk.AddRuleRequest, error) {
	policy, err := expandPolicy(model.Policy.ValueString())
	if err != nil {
		return sdk.AddRuleRequest{}, err
	}

	order := 0
	if !model.Order.IsNull() && !model.Order.IsUnknown() {
		order = int(model.Order.ValueInt64())
	}

	return sdk.AddRuleRequest{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		Order:       order,
		Policy:      policy,
		Source:      expandSource(model.Source),
	}, nil
}
