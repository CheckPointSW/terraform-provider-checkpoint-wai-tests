package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CheckPointSW/terraform-provider-checkpoint-workforce-ai/internal/sdk"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &BrowseFileProtectionObjectResource{}
	_ resource.ResourceWithImportState = &BrowseFileProtectionObjectResource{}
)

type BrowseFileProtectionObjectResource struct {
	clients *ProviderClients
}

type FileProtectionObjectResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ObjectData  types.String `tfsdk:"object_data"` // JSON string
}

func NewBrowseFileProtectionObjectResource() resource.Resource {
	return &BrowseFileProtectionObjectResource{}
}

func (r *BrowseFileProtectionObjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_browse_file_protection_object"
}

func (r *BrowseFileProtectionObjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a File Protection policy object in Check Point Browse Security.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the object (UUID).",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Human-readable object name.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Explanation of the object's purpose.",
				Optional:            true,
				Computed:            true,
			},
			"object_data": schema.StringAttribute{
				MarkdownDescription: "The object configuration as a JSON string.",
				Required:            true,
			},
		},
	}
}

func (r *BrowseFileProtectionObjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BrowseFileProtectionObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FileProtectionObjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := sdk.FileProtectionObject{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		ObjectData:  json.RawMessage(plan.ObjectData.ValueString()),
	}

	created, err := r.clients.Browse.CreateFileProtectionObject(ctx, obj)
	if err != nil {
		resp.Diagnostics.AddError("Error creating file protection object", err.Error())
		return
	}

	tflog.Trace(ctx, "created browse file protection object", map[string]interface{}{"object_id": created.ObjectID})

	// Use plan values for state after create, since the API may normalize object_data
	// (e.g. drop fields it doesn't recognize). Only take the ID from the API response.
	resp.Diagnostics.Append(resp.State.Set(ctx, &FileProtectionObjectResourceModel{
		ID:          types.StringValue(created.ObjectID),
		Name:        plan.Name,
		Description: plan.Description,
		ObjectData:  plan.ObjectData,
	})...)
}

func (r *BrowseFileProtectionObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FileProtectionObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := r.clients.Browse.GetFileProtectionObject(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading file protection object", err.Error())
		return
	}

	model, err := flattenFileProtectionObjectToModel(obj)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping object to state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *BrowseFileProtectionObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state FileProtectionObjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := sdk.FileProtectionObject{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		ObjectData:  json.RawMessage(plan.ObjectData.ValueString()),
	}

	_, err := r.clients.Browse.UpdateFileProtectionObject(ctx, state.ID.ValueString(), obj)
	if err != nil {
		resp.Diagnostics.AddError("Error updating file protection object", err.Error())
		return
	}

	// Use plan values for state after update, since the API may normalize object_data.
	resp.Diagnostics.Append(resp.State.Set(ctx, &FileProtectionObjectResourceModel{
		ID:          state.ID,
		Name:        plan.Name,
		Description: plan.Description,
		ObjectData:  plan.ObjectData,
	})...)
}

func (r *BrowseFileProtectionObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FileProtectionObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.clients.Browse.DeleteObject(ctx, sdk.DeleteObjectRequest{
		ObjectID: state.ID.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error deleting file protection object", err.Error())
		return
	}
}

func (r *BrowseFileProtectionObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	obj, err := r.clients.Browse.GetFileProtectionObject(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing file protection object", err.Error())
		return
	}

	model, err := flattenFileProtectionObjectToModel(obj)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping imported object", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func flattenFileProtectionObjectToModel(obj *sdk.FileProtectionObject) (*FileProtectionObjectResourceModel, error) {
	cleaned, err := stripJSONNulls(obj.ObjectData)
	if err != nil {
		return nil, fmt.Errorf("cleaning object_data: %w", err)
	}
	return &FileProtectionObjectResourceModel{
		ID:          types.StringValue(obj.ObjectID),
		Name:        types.StringValue(obj.Name),
		Description: types.StringValue(obj.Description),
		ObjectData:  types.StringValue(string(cleaned)),
	}, nil
}
