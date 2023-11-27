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
	_ resource.Resource              = &DataAppSourceSelectionResource{}
	_ resource.ResourceWithConfigure = &DataAppSourceSelectionResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewDataAppSourceSelectionResource() resource.Resource {
	return &DataAppSourceSelectionResource{}
}

// orderResource is the resource implementation.
type DataAppSourceSelectionResource struct {
	client *qlikcloud.Client
}

// orderResourceModel maps the resource schema data.
type DataAppSourceSelectionResourceModel struct {
	ID                 types.String           `tfsdk:"id"`
	ProjectID          types.String           `tfsdk:"project_id"`
	AppID              types.String           `tfsdk:"app_id"`
	SourceConnectionID types.String           `tfsdk:"source_connection_id"`
	SourceSelection    []SourceSelectionModel `tfsdk:"source_selection"`
}

type SourceSelectionModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	DataAppID types.String `tfsdk:"data_app_id"`
	Schema    types.String `tfsdk:"schema"`
	Database  types.String `tfsdk:"database"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
}

// Metadata returns the resource type name.
func (r *DataAppSourceSelectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_app_source_selection"
}

// Schema defines the schema for the resource.
func (r *DataAppSourceSelectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				Required: true,
			},
			"app_id": schema.StringAttribute{
				Required: true,
			},
			"source_connection_id": schema.StringAttribute{
				Required: true,
			},
			"source_selection": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Required: true,
						},
						"name": schema.StringAttribute{
							Required: true,
						},
						"data_app_id": schema.StringAttribute{
							Required: true,
						},
						"schema": schema.StringAttribute{
							Required: true,
						},
						"database": schema.StringAttribute{
							Required: true,
						},
						"type": schema.StringAttribute{
							Required: true,
						},
						"project_id": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *DataAppSourceSelectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *DataAppSourceSelectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DataAppSourceSelectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	a := models.SourceSelectionPut{
		Data: models.SourceSelectionPutData{
			DataEntitiesSelection: models.SourceSelectionDataEntitiesSelection{
				SourceConnectionID: plan.SourceConnectionID.ValueString(),
			},
		},
	}

	for _, source := range plan.SourceSelection {
		entity := models.SourceEntity{
			ID:        source.ID.ValueString(),
			Name:      source.Name.ValueString(),
			DataAppID: source.DataAppID.ValueString(),
			Schema:    source.Schema.ValueString(),
			Database:  source.Database.ValueString(),
			Type:      source.Type.ValueString(),
			ProjectID: source.ProjectID.ValueString(),
		}

		a.Data.DataEntitiesSelection.DataEntities = append(a.Data.DataEntitiesSelection.DataEntities, entity)
	}

	// Create new space
	s, err := r.client.PutSourceSelection(plan.ProjectID.ValueString(), plan.AppID.ValueString(), a)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data connection",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(s.Key)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read resource information.
func (r *DataAppSourceSelectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DataAppSourceSelectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	s, err := r.client.GetSourceSelection(state.ProjectID.ValueString(), state.AppID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Data App Source Selection",
			"Could not read Data App Source Selection"+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to model
	for _, source := range s.SourceSelection.DataEntitiesSelection.DataEntities {
		sourceState := SourceSelectionModel{
			ID:        types.StringValue(source.ID),
			Name:      types.StringValue(source.Name),
			DataAppID: types.StringValue(source.DataAppID),
			Schema:    types.StringValue(source.Schema),
			Database:  types.StringValue(source.Database),
			Type:      types.StringValue(source.Type),
			ProjectID: types.StringValue(source.ProjectID),
		}
		state.SourceSelection = append(state.SourceSelection, sourceState)
	}

	state.ID = types.StringValue(s.Key)
	state.SourceConnectionID = types.StringValue(s.SourceSelection.DataEntitiesSelection.SourceConnectionID)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DataAppSourceSelectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan DataAppSourceSelectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	a := models.SourceSelectionPut{
		Data: models.SourceSelectionPutData{
			DataEntitiesSelection: models.SourceSelectionDataEntitiesSelection{
				SourceConnectionID: plan.SourceConnectionID.ValueString(),
			},
		},
	}

	for _, source := range plan.SourceSelection {
		entity := models.SourceEntity{
			ID:        source.ID.ValueString(),
			Name:      source.Name.ValueString(),
			DataAppID: source.DataAppID.ValueString(),
			Schema:    source.Schema.ValueString(),
			Database:  source.Database.ValueString(),
			Type:      source.Type.ValueString(),
			ProjectID: source.ProjectID.ValueString(),
		}

		a.Data.DataEntitiesSelection.DataEntities = append(a.Data.DataEntitiesSelection.DataEntities, entity)
	}

	// Create new space
	_, err := r.client.PutSourceSelection(plan.ProjectID.ValueString(), plan.AppID.ValueString(), a)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data connection",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DataAppSourceSelectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	// var state DataAppSourceSelectionResourceModel
	// diags := req.State.Get(ctx, &state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// // Delete existing order
	// err := r.client.DeleteDataAppSourceSelection(state.ProjectID.ValueString(), state.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Deleting Data Project",
	// 		"Could not delete Data Project, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }
}

// func (r *orderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
