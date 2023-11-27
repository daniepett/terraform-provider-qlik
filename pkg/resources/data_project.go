package resources

import (
	"context"
	"fmt"

	qlikcloud "github.com/daniepett/qlik-cloud-client-go"
	"github.com/daniepett/qlik-cloud-client-go/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &DataProjectResource{}
	_ resource.ResourceWithConfigure = &DataProjectResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewDataProjectResource() resource.Resource {
	return &DataProjectResource{}
}

// orderResource is the resource implementation.
type DataProjectResource struct {
	client *qlikcloud.Client
}

// orderResourceModel maps the resource schema data.
type DataProjectResourceModel struct {
	ID                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	SpaceID           types.String `tfsdk:"space_id"`
	LakehouseType     types.String `tfsdk:"lakehouse_type"`
	Type              types.String `tfsdk:"type"`
	StorageConnection types.String `tfsdk:"storage_connection"`
	BatchMode         types.Bool   `tfsdk:"batch_mode"`
}

// Metadata returns the resource type name.
func (r *DataProjectResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_project"
}

// Schema defines the schema for the resource.
func (r *DataProjectResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"space_id": schema.StringAttribute{
				Required: true,
			},
			"lakehouse_type": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"storage_connection": schema.StringAttribute{
				Required: true,
			},
			"batch_mode": schema.BoolAttribute{
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *DataProjectResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*qlikcloud.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *qlikcloud.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create a new resource.
func (r *DataProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DataProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newDataProject := models.DataProjectCreate{
		SpaceID: plan.SpaceID.ValueString(),
		Data: models.DataProjectConfiguration{
			Name:              plan.Name.ValueString(),
			LakehouseType:     plan.LakehouseType.ValueString(),
			Type:              plan.Type.ValueString(),
			StorageConnection: plan.StorageConnection.ValueString(),
			Description:       plan.Description.ValueString(),
			BatchMode:         plan.BatchMode.ValueBool(),
		},
	}

	// Create new space
	dataProject, err := r.client.CreateDataProject(newDataProject)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data connection",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(dataProject.DataProject.ID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *DataProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DataProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.GetDataProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Data Project",
			"Could not read Data Project ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(project.DataProject.Name)
	state.Description = types.StringValue(project.DataProject.Description)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DataProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan DataProjectResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := models.DataProjectUpdate{
		SpaceID: plan.SpaceID.ValueString(),
		Data: models.DataProjectConfiguration{
			ID:                plan.ID.ValueString(),
			Name:              plan.Name.ValueString(),
			LakehouseType:     plan.LakehouseType.ValueString(),
			Type:              plan.Type.ValueString(),
			StorageConnection: plan.StorageConnection.ValueString(),
			Description:       plan.Description.ValueString(),
			BatchMode:         plan.BatchMode.ValueBool(),
		},
	}

	// Create new space
	_, err := r.client.UpdateDataProject(plan.ID.ValueString(), project)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating data project",
			"Could not update data project, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DataProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state DataProjectResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteDataProject(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Data Project",
			"Could not delete Data Project, unexpected error: "+err.Error(),
		)
		return
	}
}

// func (r *orderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
