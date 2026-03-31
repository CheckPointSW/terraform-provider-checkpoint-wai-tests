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
	_ resource.Resource                = &BrowseDomainsObjectResource{}
	_ resource.ResourceWithImportState = &BrowseDomainsObjectResource{}
)

type BrowseDomainsObjectResource struct {
	clients *ProviderClients
}

type DomainsObjectResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	ObjectData  types.String `tfsdk:"object_data"` // JSON string
}

func NewBrowseDomainsObjectResource() resource.Resource {
	return &BrowseDomainsObjectResource{}
}

func (r *BrowseDomainsObjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_browse_domains_object"
}

func (r *BrowseDomainsObjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a Domains policy object in Check Point Browse Security.",
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

func (r *BrowseDomainsObjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *BrowseDomainsObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DomainsObjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := sdk.DomainsObject{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		ObjectData:  json.RawMessage(plan.ObjectData.ValueString()),
	}

	created, err := r.clients.Browse.CreateDomainsObject(ctx, obj)
	if err != nil {
		resp.Diagnostics.AddError("Error creating domains object", err.Error())
		return
	}

	tflog.Trace(ctx, "created browse domains object", map[string]interface{}{"object_id": created.ObjectID})

	// Use plan values for state after create, since the API may normalize object_data
	// (e.g. drop fields it doesn't recognize). Only take the ID from the API response.
	resp.Diagnostics.Append(resp.State.Set(ctx, &DomainsObjectResourceModel{
		ID:          types.StringValue(created.ObjectID),
		Name:        plan.Name,
		Description: plan.Description,
		ObjectData:  plan.ObjectData,
	})...)
}

func (r *BrowseDomainsObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DomainsObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj, err := r.clients.Browse.GetDomainsObject(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading domains object", err.Error())
		return
	}

	model, err := flattenDomainsObjectToModel(obj)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping object to state", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *BrowseDomainsObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state DomainsObjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := sdk.DomainsObject{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		ObjectData:  json.RawMessage(plan.ObjectData.ValueString()),
	}

	_, err := r.clients.Browse.UpdateDomainsObject(ctx, state.ID.ValueString(), obj)
	if err != nil {
		resp.Diagnostics.AddError("Error updating domains object", err.Error())
		return
	}

	// Use plan values for state after update, since the API may normalize object_data.
	resp.Diagnostics.Append(resp.State.Set(ctx, &DomainsObjectResourceModel{
		ID:          state.ID,
		Name:        plan.Name,
		Description: plan.Description,
		ObjectData:  plan.ObjectData,
	})...)
}

func (r *BrowseDomainsObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DomainsObjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.clients.Browse.DeleteObject(ctx, sdk.DeleteObjectRequest{
		ObjectID: state.ID.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Error deleting domains object", err.Error())
		return
	}
}

func (r *BrowseDomainsObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	obj, err := r.clients.Browse.GetDomainsObject(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Error importing domains object", err.Error())
		return
	}

	model, err := flattenDomainsObjectToModel(obj)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping imported object", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func flattenDomainsObjectToModel(obj *sdk.DomainsObject) (*DomainsObjectResourceModel, error) {
	cleaned, err := stripJSONNulls(obj.ObjectData)
	if err != nil {
		return nil, fmt.Errorf("cleaning object_data: %w", err)
	}
	return &DomainsObjectResourceModel{
		ID:          types.StringValue(obj.ObjectID),
		Name:        types.StringValue(obj.Name),
		Description: types.StringValue(obj.Description),
		ObjectData:  types.StringValue(string(cleaned)),
	}, nil
}
