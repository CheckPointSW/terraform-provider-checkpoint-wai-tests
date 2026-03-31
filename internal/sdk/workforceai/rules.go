package workforceai

import (
	"context"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
)

// SetRuleInfo updates the display info of a rule.
func (c *Client) SetRuleInfo(ctx context.Context, req sdk.SetInfoRequest) error {
	_, err := c.Put(ctx, policyBasePath+"/rules/set-info"+productQuery, req)
	if err != nil {
		return fmt.Errorf("setting rule info: %w", err)
	}
	return nil
}

// ReorderRule changes a rule's position in its rulebase.
func (c *Client) ReorderRule(ctx context.Context, req sdk.ReorderRequest) error {
	_, err := c.Put(ctx, policyBasePath+"/rules/reorder"+productQuery, req)
	if err != nil {
		return fmt.Errorf("reordering rule: %w", err)
	}
	return nil
}

// SetRuleSource replaces the source (user/group assignments) of a rule.
func (c *Client) SetRuleSource(ctx context.Context, req sdk.SetSourceRequest) error {
	_, err := c.Put(ctx, policyBasePath+"/rules/set-source"+productQuery, req)
	if err != nil {
		return fmt.Errorf("setting rule source: %w", err)
	}
	return nil
}

// SetRuleActive enables or disables a rule.
func (c *Client) SetRuleActive(ctx context.Context, req sdk.SetActiveRequest) error {
	_, err := c.Put(ctx, policyBasePath+"/rules/set-active"+productQuery, req)
	if err != nil {
		return fmt.Errorf("setting rule active: %w", err)
	}
	return nil
}

// DeleteRule permanently deletes a rule from its rulebase.
func (c *Client) DeleteRule(ctx context.Context, ruleID string) error {
	_, err := c.Delete(ctx, policyBasePath+"/rules/"+ruleID+productQuery)
	if err != nil {
		return fmt.Errorf("deleting rule %s: %w", ruleID, err)
	}
	return nil
}
