package workforceai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
)

// GetChatsRulebase retrieves the full Chats rulebase.
func (c *Client) GetChatsRulebase(ctx context.Context) (*sdk.Rulebase, error) {
	resp, err := c.Get(ctx, policyBasePath+"/chats/rulebase")
	if err != nil {
		return nil, fmt.Errorf("getting chats rulebase: %w", err)
	}
	var rb sdk.Rulebase; err = json.Unmarshal(resp, &rb)
	if err != nil {
		return nil, fmt.Errorf("parsing chats rulebase: %w", err)
	}
	return &rb, nil
}

// GetChatsRule retrieves a single rule from the Chats rulebase by ID.
func (c *Client) GetChatsRule(ctx context.Context, ruleID string) (*sdk.Rule, error) {
	rb, err := c.GetChatsRulebase(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range rb.Rules {
		if r.RuleID == ruleID {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("chats rule %s not found", ruleID)
}

// CreateChatsRule creates a new rule in the Chats rulebase.
func (c *Client) CreateChatsRule(ctx context.Context, req sdk.AddRuleRequest) (*sdk.Rule, error) {
	resp, err := c.Post(ctx, policyBasePath+"/chats/rule"+productQuery, req)
	if err != nil {
		return nil, fmt.Errorf("creating chats rule: %w", err)
	}
	var rule sdk.Rule
	if err := json.Unmarshal(resp, &rule); err != nil {
		return nil, fmt.Errorf("parsing created chats rule: %w", err)
	}
	return &rule, nil
}

// SetChatsPolicy replaces the entire policy of a Chats rule.
func (c *Client) SetChatsPolicy(ctx context.Context, req sdk.SetPolicyRequest) error {
	_, err := c.Put(ctx, policyBasePath+"/chats/rule/set-policy"+productQuery, req)
	if err != nil {
		return fmt.Errorf("setting chats policy: %w", err)
	}
	return nil
}

// PatchChatsPolicy partially updates a Chats rule's policy.
func (c *Client) PatchChatsPolicy(ctx context.Context, req sdk.PatchPolicyRequest) error {
	_, err := c.Patch(ctx, policyBasePath+"/chats/rule/patch-policy"+productQuery, req)
	if err != nil {
		return fmt.Errorf("patching chats policy: %w", err)
	}
	return nil
}
