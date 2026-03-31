package workforceai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
)

// GetAccessRulebase retrieves the full Access rulebase.
func (c *Client) GetAccessRulebase(ctx context.Context) (*sdk.Rulebase, error) {
	resp, err := c.Get(ctx, policyBasePath+"/ai-access/rulebase")
	if err != nil {
		return nil, fmt.Errorf("getting access rulebase: %w", err)
	}
	var rb sdk.Rulebase; err = json.Unmarshal(resp, &rb)
	if err != nil {
		return nil, fmt.Errorf("parsing access rulebase: %w", err)
	}
	return &rb, nil
}

// GetAccessRule retrieves a single rule from the Access rulebase by ID.
func (c *Client) GetAccessRule(ctx context.Context, ruleID string) (*sdk.Rule, error) {
	rb, err := c.GetAccessRulebase(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range rb.Rules {
		if r.RuleID == ruleID {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("access rule %s not found", ruleID)
}

// CreateAccessRule creates a new rule in the Access rulebase.
func (c *Client) CreateAccessRule(ctx context.Context, req sdk.AddRuleRequest) (*sdk.Rule, error) {
	resp, err := c.Post(ctx, policyBasePath+"/ai-access/rule"+productQuery, req)
	if err != nil {
		return nil, fmt.Errorf("creating access rule: %w", err)
	}
	var rule sdk.Rule
	if err := json.Unmarshal(resp, &rule); err != nil {
		return nil, fmt.Errorf("parsing created access rule: %w", err)
	}
	return &rule, nil
}

// SetAccessPolicy replaces the entire policy of an Access rule.
func (c *Client) SetAccessPolicy(ctx context.Context, req sdk.SetPolicyRequest) error {
	_, err := c.Put(ctx, policyBasePath+"/ai-access/rule/set-policy"+productQuery, req)
	if err != nil {
		return fmt.Errorf("setting access policy: %w", err)
	}
	return nil
}

// PatchAccessPolicy partially updates an Access rule's policy.
func (c *Client) PatchAccessPolicy(ctx context.Context, req sdk.PatchPolicyRequest) error {
	_, err := c.Patch(ctx, policyBasePath+"/ai-access/rule/patch-policy"+productQuery, req)
	if err != nil {
		return fmt.Errorf("patching access policy: %w", err)
	}
	return nil
}
