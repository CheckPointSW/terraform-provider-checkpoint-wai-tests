package browse

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
)

// GetSecureBrowsingRulebase retrieves the full Secure Browsing rulebase.
func (c *Client) GetSecureBrowsingRulebase(ctx context.Context) (*sdk.Rulebase, error) {
	resp, err := c.Get(ctx, policyBasePath+"/secure-browsing/rulebase")
	if err != nil {
		return nil, fmt.Errorf("getting secure browsing rulebase: %w", err)
	}
	var rb sdk.Rulebase; err = json.Unmarshal(resp, &rb)
	if err != nil {
		return nil, fmt.Errorf("parsing secure browsing rulebase: %w", err)
	}
	return &rb, nil
}

// GetSecureBrowsingRule retrieves a single rule from the Secure Browsing rulebase by ID.
func (c *Client) GetSecureBrowsingRule(ctx context.Context, ruleID string) (*sdk.Rule, error) {
	rb, err := c.GetSecureBrowsingRulebase(ctx)
	if err != nil {
		return nil, err
	}
	for _, r := range rb.Rules {
		if r.RuleID == ruleID {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("secure browsing rule %s not found", ruleID)
}

// CreateSecureBrowsingRule creates a new rule in the Secure Browsing rulebase.
func (c *Client) CreateSecureBrowsingRule(ctx context.Context, req sdk.AddRuleRequest) (*sdk.Rule, error) {
	resp, err := c.Post(ctx, policyBasePath+"/secure-browsing/rule"+productQuery, req)
	if err != nil {
		return nil, fmt.Errorf("creating secure browsing rule: %w", err)
	}
	var rule sdk.Rule
	if err := json.Unmarshal(resp, &rule); err != nil {
		return nil, fmt.Errorf("parsing created secure browsing rule: %w", err)
	}
	return &rule, nil
}

// SetSecureBrowsingPolicy replaces the entire policy of a Secure Browsing rule.
func (c *Client) SetSecureBrowsingPolicy(ctx context.Context, req sdk.SetPolicyRequest) error {
	_, err := c.Put(ctx, policyBasePath+"/secure-browsing/rule/set-policy"+productQuery, req)
	if err != nil {
		return fmt.Errorf("setting secure browsing policy: %w", err)
	}
	return nil
}

// PatchSecureBrowsingPolicy partially updates a Secure Browsing rule's policy.
func (c *Client) PatchSecureBrowsingPolicy(ctx context.Context, req sdk.PatchPolicyRequest) error {
	_, err := c.Patch(ctx, policyBasePath+"/secure-browsing/rule/patch-policy"+productQuery, req)
	if err != nil {
		return fmt.Errorf("patching secure browsing policy: %w", err)
	}
	return nil
}

// SetSecureBrowsingObjects sets (replaces) all objects for a Secure Browsing rule for a given feature.
func (c *Client) SetSecureBrowsingObjects(ctx context.Context, req sdk.SetObjectsRequest) error {
	_, err := c.Put(ctx, policyBasePath+"/secure-browsing/rule/set-objects"+productQuery, req)
	if err != nil {
		return fmt.Errorf("setting secure browsing objects: %w", err)
	}
	return nil
}
