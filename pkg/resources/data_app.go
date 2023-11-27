package resources

import (
	"context"
	"fmt"

	qlikcloud "github.com/daniepett/qlik-cloud-client-go"
	"github.com/daniepett/qlik-cloud-client-go/models"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &DataAppResource{}
	_ resource.ResourceWithConfigure = &DataAppResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewDataAppResource() resource.Resource {
	return &DataAppResource{}
}

// orderResource is the resource implementation.
type DataAppResource struct {
	client *qlikcloud.Client
}

// orderResourceModel maps the resource schema data.
type DataAppResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	ProjectID   types.String `tfsdk:"project_id"`
}

// Metadata returns the resource type name.
func (r *DataAppResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_app"
}

// Schema defines the schema for the resource.
func (r *DataAppResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"type": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *DataAppResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *DataAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DataAppResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	a := models.DataAppCreate{
		Data: models.DataAppCreateData{
			Name:        plan.Name.ValueString(),
			Type:        plan.Type.ValueString(),
			Description: plan.Description.ValueString(),
		},
	}

	// Create new space
	app, err := r.client.CreateDataApp(plan.ProjectID.ValueString(), a)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data connection",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(app.DataApp.ID)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *DataAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DataAppResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.GetDataApp(state.ProjectID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Qlik Cloud Data App",
			"Could not read App ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(a.DataApp.Name)
	state.Description = types.StringValue(a.DataApp.Description)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DataAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan DataAppResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := models.DataAppUpdate{
		Data: models.DataAppUpdateData{
			Name:        plan.Name.ValueString(),
			Type:        plan.Type.ValueString(),
			Description: plan.Description.ValueString(),
		},
	}

	// Create new space
	a, err := r.client.UpdateDataApp(plan.ProjectID.ValueString(), project)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating data app",
			"Could not update data app, unexpected error: "+err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(a.DataApp.Name)
	plan.Description = types.StringValue(a.DataApp.Description)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DataAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state DataAppResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteDataApp(state.ProjectID.ValueString(), state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Data App",
			"Could not delete Data App, unexpected error: "+err.Error(),
		)
		return
	}
}

// func (r *orderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
