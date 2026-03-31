package browse

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
)

// GetAgentsRulebase retrieves the full Agents rulebase.
func (c *Client) GetAgentsRulebase(ctx context.Context) (*sdk.Rulebase, error) {
	resp, err := c.Get(ctx, policyBasePath+"/agents/rulebase")
	if err != nil {
		return nil, fmt.Errorf("getting agents rulebase: %w", err)
	}
	var rb sdk.Rulebase; err = json.Unmarshal(resp, &rb)
	if err != nil {
		return nil, fmt.Errorf("parsing agents rulebase: %w", err)
	}
	return &rb, nil
}

// GetAgentsRule retrieves a single rule from the Agents rulebase by ID.
func (c *Client) GetAgentsRule(ctx context.Context, ruleID string) (*sdk.Rule, error) {
	rb, err := c.GetAgentsRulebase(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range rb.Rules {
		if r.RuleID == ruleID {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("agents rule %s not found", ruleID)
}

// CreateAgentsRule creates a new rule in the Agents rulebase.
func (c *Client) CreateAgentsRule(ctx context.Context, req sdk.AddRuleRequest) (*sdk.Rule, error) {
	resp, err := c.Post(ctx, policyBasePath+"/agents/rule"+productQuery, req)
	if err != nil {
		return nil, fmt.Errorf("creating agents rule: %w", err)
	}
	var rule sdk.Rule
	if err := json.Unmarshal(resp, &rule); err != nil {
		return nil, fmt.Errorf("parsing created agents rule: %w", err)
	}
	return &rule, nil
}

// SetAgentsPolicy replaces the entire policy of an Agents rule.
func (c *Client) SetAgentsPolicy(ctx context.Context, req sdk.SetPolicyRequest) error {
	_, err := c.Put(ctx, policyBasePath+"/agents/rule/set-policy"+productQuery, req)
	if err != nil {
		return fmt.Errorf("setting agents policy: %w", err)
	}
	return nil
}

// PatchAgentsPolicy partially updates an Agents rule's policy.
func (c *Client) PatchAgentsPolicy(ctx context.Context, req sdk.PatchPolicyRequest) error {
	_, err := c.Patch(ctx, policyBasePath+"/agents/rule/patch-policy"+productQuery, req)
	if err != nil {
		return fmt.Errorf("patching agents policy: %w", err)
	}
	return nil
}
