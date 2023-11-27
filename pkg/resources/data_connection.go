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
	_ resource.Resource              = &DataConnectionResource{}
	_ resource.ResourceWithConfigure = &DataConnectionResource{}
	// _ resource.ResourceWithImportState = &spaceResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewDataConnectionResource() resource.Resource {
	return &DataConnectionResource{}
}

// orderResource is the resource implementation.
type DataConnectionResource struct {
	client *qlikcloud.Client
}

// orderResourceModel maps the resource schema data.
type DataConnectionResourceModel struct {
	ID                   types.String             `tfsdk:"id"`
	Name                 types.String             `tfsdk:"name"`
	SpaceID              types.String             `tfsdk:"space_id"`
	GatewayID            types.String             `tfsdk:"gateway_id"`
	ConnectionParameters DataConnectionParameters `tfsdk:"connection_parameters"`
	Type                 types.String             `tfsdk:"type"`
	Driver               types.String             `tfsdk:"driver"`
	EngineID             types.String             `tfsdk:"engine_id"`
	ConnectStatement     types.String             `tfsdk:"connect_statement"`
	CredentialsID        types.String             `tfsdk:"credentials_id"`
	CredentialsName      types.String             `tfsdk:"credentials_name"`
}

type DataConnectionParameters struct {
	Server         types.String `tfsdk:"server"`
	Username       types.String `tfsdk:"username"`
	Warehouse      types.String `tfsdk:"warehouse"`
	Database       types.String `tfsdk:"database"`
	MetadataSchema types.String `tfsdk:"metadata_schema"`
	SapClient      types.String `tfsdk:"sap_client"`
	Password       types.String `tfsdk:"password"`
}

// Metadata returns the resource type name.
func (r *DataConnectionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_connection"
}

// Schema defines the schema for the resource.
func (r *DataConnectionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"space_id": schema.StringAttribute{
				Required: true,
			},
			"gateway_id": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Required: true,
			},
			"driver": schema.StringAttribute{
				Computed: true,
			},
			"engine_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connect_statement": schema.StringAttribute{
				Computed: true,
			},
			"credentials_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"credentials_name": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_parameters": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"server": schema.StringAttribute{
						Optional: true,
					},
					"username": schema.StringAttribute{
						Optional: true,
					},
					"warehouse": schema.StringAttribute{
						Optional: true,
					},
					"database": schema.StringAttribute{
						Optional: true,
					},
					"metadata_schema": schema.StringAttribute{
						Optional: true,
					},
					"password": schema.StringAttribute{
						Optional:  true,
						Sensitive: true,
					},
					"sap_client": schema.StringAttribute{
						Optional: true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *DataConnectionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *DataConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DataConnectionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var driver string

	if plan.Type.ValueString() == "reptgt_qdisnowflake" {
		driver = "QlikConnectorsCommonService.exe"

	} else if plan.Type.ValueString() == "SAP_APPLICATION" {
		driver = "QlikConnectorsCommonService.exe"
	}

	src := plan.Type.ValueString()
	c, err := r.GetConnectionString(src, plan)

	newDataConnection := models.ConnectionCreate{
		Name:             plan.Name.ValueString(),
		SpaceID:          plan.SpaceID.ValueString(),
		LogOn:            1,
		ConnectStatement: c.ConnectionString,
		DataSourceID:     src,
		Type:             driver,
		Username:         c.UserID,
		Password:         c.CredentialsConnectionString,
	}

	// Create new space
	DataConnection, err := r.client.CreateDataConnection(newDataConnection)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating data connection",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(DataConnection.ID)
	plan.EngineID = types.StringValue(DataConnection.EngineID)
	plan.ConnectStatement = types.StringValue(DataConnection.ConnectStatement)
	plan.Driver = types.StringValue(DataConnection.Type)
	plan.CredentialsID = types.StringValue(DataConnection.CredentialsID)
	plan.CredentialsName = types.StringValue(DataConnection.CredentialsName)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *DataConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state DataConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	connection, err := r.client.GetDataConnection(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Data Connection",
			"Could not read Connection ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	state.Name = types.StringValue(connection.Name)
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DataConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan DataConnectionResourceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var driver string

	if plan.Type.ValueString() == "reptgt_qdisnowflake" {
		driver = "QlikConnectorsCommonService.exe"

	} else if plan.Type.ValueString() == "SAP_APPLICATION" {
		driver = "QlikConnectorsCommonService.exe"
	}

	src := plan.Type.ValueString()
	c, err := r.GetConnectionString(src, plan)
	updateDataConnection := models.ConnectionUpdate{
		ID:               plan.ID.ValueString(),
		SpaceID:          plan.SpaceID.ValueString(),
		Name:             plan.Name.ValueString(),
		EngineID:         plan.EngineID.ValueString(),
		ConnectStatement: c.ConnectionString,
		DataSourceID:     src,
		Type:             driver,
		Username:         c.UserID,
		Password:         c.CredentialsConnectionString,
	}
	// Create new space
	err = r.client.UpdateDataConnection(plan.ID.ValueString(), updateDataConnection)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating data connection",
			"Could not update data connection, unexpected error: "+err.Error(),
		)
		return
	}

	plan.ConnectStatement = types.StringValue(updateDataConnection.ConnectStatement)
	plan.Driver = types.StringValue(updateDataConnection.Type)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *DataConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state DataConnectionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteDataConnection(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Data Connection",
			"Could not delete Data Connection, unexpected error: "+err.Error(),
		)
		return
	}
}

// func (r *orderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	// Retrieve import ID and save to id attribute
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }

func (r *DataConnectionResource) GetConnectionString(src string, props DataConnectionResourceModel) (*models.GetConnectionStringResponse, error) {

	var conn models.GetConnectionString
	var crd models.GetConnectionString
	if src == "reptgt_qdisnowflake" {
		conn = models.GetConnectionString{
			PropertiesList: []models.ConnectionProperties{
				models.ConnectionProperties{
					Name:  "sourceType",
					Value: "reptgt_qdisnowflake",
				},
				models.ConnectionProperties{
					Name:  "agentId",
					Value: props.GatewayID.ValueString(),
				},
				models.ConnectionProperties{
					Name:  "endpointTypePrefix",
					Value: "reptgt_",
				},
				models.ConnectionProperties{
					Name:  "useDbCommandForTest",
					Value: "true",
				},
				models.ConnectionProperties{
					Name:  "replicateEndpointType",
					Value: "snowflake",
				},
				models.ConnectionProperties{
					Name:  "server",
					Value: props.ConnectionParameters.Server.ValueString(),
				},
				models.ConnectionProperties{
					Name:  "port",
					Value: "443",
				},
				models.ConnectionProperties{
					Name:  "username",
					Value: props.ConnectionParameters.Username.ValueString(),
				},
				models.ConnectionProperties{
					Name:  "warehouse",
					Value: props.ConnectionParameters.Warehouse.ValueString(),
				},
				models.ConnectionProperties{
					Name:  "database",
					Value: props.ConnectionParameters.Database.ValueString(),
				},
				models.ConnectionProperties{
					Name:  "metadataschema",
					Value: props.ConnectionParameters.MetadataSchema.ValueString(),
				},
				models.ConnectionProperties{
					Name:  "stagingtype",
					Value: "SNOWFLAKE_STAGE",
				},
				models.ConnectionProperties{
					Name:  "proxySettingsOrigin",
					Value: "ENDPOINT",
				},
				models.ConnectionProperties{
					Name:  "useProxyServer",
					Value: "false",
				},
			},
		}
		crd = models.GetConnectionString{
			PropertiesList: []models.ConnectionProperties{
				models.ConnectionProperties{
					Name:  "password",
					Value: props.ConnectionParameters.Password.ValueString(),
				},
			},
		}

	} else if src == "SAP_APPLICATION" {
		return nil, nil
	}

	c, _ := r.client.GetConnectionString(src, conn, crd)

	return c, nil
}
