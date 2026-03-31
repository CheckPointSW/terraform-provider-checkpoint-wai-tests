package browse

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
)

// GetPredefinedDLPDataTypes retrieves all predefined DLP data types.
func (c *Client) GetPredefinedDLPDataTypes(ctx context.Context) ([]sdk.DLPDataType, error) {
	resp, err := c.Get(ctx, policyBasePath+"/dlp-datatypes/predefined")
	if err != nil {
		return nil, fmt.Errorf("getting predefined DLP datatypes: %w", err)
	}
	var types []sdk.DLPDataType
	if err := json.Unmarshal(resp, &types); err != nil {
		return nil, fmt.Errorf("parsing predefined DLP datatypes: %w", err)
	}
	return types, nil
}

// GetCustomDLPDataTypes retrieves all custom DLP data types.
func (c *Client) GetCustomDLPDataTypes(ctx context.Context) ([]sdk.DLPDataType, error) {
	resp, err := c.Get(ctx, policyBasePath+"/dlp-datatypes/custom")
	if err != nil {
		return nil, fmt.Errorf("getting custom DLP datatypes: %w", err)
	}
	var types []sdk.DLPDataType
	if err := json.Unmarshal(resp, &types); err != nil {
		return nil, fmt.Errorf("parsing custom DLP datatypes: %w", err)
	}
	return types, nil
}

// GetAllDLPDataTypes retrieves all DLP data types (predefined and custom).
func (c *Client) GetAllDLPDataTypes(ctx context.Context) ([]sdk.DLPDataType, error) {
	resp, err := c.Get(ctx, policyBasePath+"/dlp-datatypes/all")
	if err != nil {
		return nil, fmt.Errorf("getting all DLP datatypes: %w", err)
	}
	var types []sdk.DLPDataType
	if err := json.Unmarshal(resp, &types); err != nil {
		return nil, fmt.Errorf("parsing all DLP datatypes: %w", err)
	}
	return types, nil
}

// GetTenantDLPDataTypes retrieves tenant-specific DLP data type configuration.
func (c *Client) GetTenantDLPDataTypes(ctx context.Context) (*sdk.TenantDLPConfig, error) {
	resp, err := c.Get(ctx, policyBasePath+"/dlp-datatypes")
	if err != nil {
		return nil, fmt.Errorf("getting tenant DLP datatypes: %w", err)
	}
	var config sdk.TenantDLPConfig
	if err := json.Unmarshal(resp, &config); err != nil {
		return nil, fmt.Errorf("parsing tenant DLP datatypes: %w", err)
	}
	return &config, nil
}

// SetTenantDLPDataTypes sets the tenant-specific DLP data type configuration.
func (c *Client) SetTenantDLPDataTypes(ctx context.Context, config sdk.TenantDLPConfig) error {
	_, err := c.Post(ctx, policyBasePath+"/dlp-datatypes", config)
	if err != nil {
		return fmt.Errorf("setting tenant DLP datatypes: %w", err)
	}
	return nil
}
