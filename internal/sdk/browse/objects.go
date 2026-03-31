package browse

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
)

// --- File Protection Objects ---

// ListFileProtectionObjects retrieves all file protection objects.
func (c *Client) ListFileProtectionObjects(ctx context.Context) ([]sdk.FileProtectionObject, error) {
	resp, err := c.Get(ctx, policyBasePath+"/objects/file-protection")
	if err != nil {
		return nil, fmt.Errorf("listing file protection objects: %w", err)
	}
	var wrapper struct {
		Objects []sdk.FileProtectionObject `json:"objects"`
	}
	if err := json.Unmarshal(resp, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing file protection objects: %w", err)
	}
	return wrapper.Objects, nil
}

// GetFileProtectionObject retrieves a file protection object by ID.
func (c *Client) GetFileProtectionObject(ctx context.Context, objectID string) (*sdk.FileProtectionObject, error) {
	objects, err := c.ListFileProtectionObjects(ctx)
	if err != nil {
		return nil, err
	}
	for _, o := range objects {
		if o.ObjectID == objectID {
			return &o, nil
		}
	}
	return nil, fmt.Errorf("file protection object %s not found", objectID)
}

// CreateFileProtectionObject creates a new file protection object.
// The API returns a CreateObjectResponse with only the object_id.
func (c *Client) CreateFileProtectionObject(ctx context.Context, obj sdk.FileProtectionObject) (*sdk.CreateObjectResponse, error) {
	resp, err := c.Post(ctx, policyBasePath+"/objects/file-protection"+productQuery, obj)
	if err != nil {
		return nil, fmt.Errorf("creating file protection object: %w", err)
	}
	var created sdk.CreateObjectResponse
	if err := json.Unmarshal(resp, &created); err != nil {
		return nil, fmt.Errorf("parsing created file protection object: %w", err)
	}
	return &created, nil
}

// UpdateFileProtectionObject updates an existing file protection object.
// The API returns an UpdateObjectResponse with minimal metadata.
func (c *Client) UpdateFileProtectionObject(ctx context.Context, objectID string, obj sdk.FileProtectionObject) (*sdk.UpdateObjectResponse, error) {
	req := sdk.UpdateObjectRequest{
		ObjectID:    objectID,
		Name:        obj.Name,
		Description: obj.Description,
		ObjectData:  obj.ObjectData,
	}
	resp, err := c.Put(ctx, policyBasePath+"/objects/file-protection"+productQuery, req)
	if err != nil {
		return nil, fmt.Errorf("updating file protection object: %w", err)
	}
	var updated sdk.UpdateObjectResponse
	if err := json.Unmarshal(resp, &updated); err != nil {
		return nil, fmt.Errorf("parsing updated file protection object: %w", err)
	}
	return &updated, nil
}

// --- Domains Objects ---

// ListDomainsObjects retrieves all domains objects.
func (c *Client) ListDomainsObjects(ctx context.Context) ([]sdk.DomainsObject, error) {
	resp, err := c.Get(ctx, policyBasePath+"/objects/domains")
	if err != nil {
		return nil, fmt.Errorf("listing domains objects: %w", err)
	}
	var wrapper struct {
		Objects []sdk.DomainsObject `json:"objects"`
	}
	if err := json.Unmarshal(resp, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing domains objects: %w", err)
	}
	return wrapper.Objects, nil
}

// GetDomainsObject retrieves a domains object by ID.
func (c *Client) GetDomainsObject(ctx context.Context, objectID string) (*sdk.DomainsObject, error) {
	objects, err := c.ListDomainsObjects(ctx)
	if err != nil {
		return nil, err
	}
	for _, o := range objects {
		if o.ObjectID == objectID {
			return &o, nil
		}
	}
	return nil, fmt.Errorf("domains object %s not found", objectID)
}

// CreateDomainsObject creates a new domains object.
// The API returns a CreateObjectResponse with only the object_id.
func (c *Client) CreateDomainsObject(ctx context.Context, obj sdk.DomainsObject) (*sdk.CreateObjectResponse, error) {
	resp, err := c.Post(ctx, policyBasePath+"/objects/domains"+productQuery, obj)
	if err != nil {
		return nil, fmt.Errorf("creating domains object: %w", err)
	}
	var created sdk.CreateObjectResponse
	if err := json.Unmarshal(resp, &created); err != nil {
		return nil, fmt.Errorf("parsing created domains object: %w", err)
	}
	return &created, nil
}

// UpdateDomainsObject updates an existing domains object.
// The API returns an UpdateObjectResponse with minimal metadata.
func (c *Client) UpdateDomainsObject(ctx context.Context, objectID string, obj sdk.DomainsObject) (*sdk.UpdateObjectResponse, error) {
	req := sdk.UpdateObjectRequest{
		ObjectID:    objectID,
		Name:        obj.Name,
		Description: obj.Description,
		ObjectData:  obj.ObjectData,
	}
	resp, err := c.Put(ctx, policyBasePath+"/objects/domains"+productQuery, req)
	if err != nil {
		return nil, fmt.Errorf("updating domains object: %w", err)
	}
	var updated sdk.UpdateObjectResponse
	if err := json.Unmarshal(resp, &updated); err != nil {
		return nil, fmt.Errorf("parsing updated domains object: %w", err)
	}
	return &updated, nil
}

// --- Delete Object (shared) ---

// DeleteObject deletes a policy object by ID.
func (c *Client) DeleteObject(ctx context.Context, req sdk.DeleteObjectRequest) error {
	_, err := c.DoRequest(ctx, "DELETE", policyBasePath+"/objects/delete"+productQuery, req)
	if err != nil {
		return fmt.Errorf("deleting object: %w", err)
	}
	return nil
}
