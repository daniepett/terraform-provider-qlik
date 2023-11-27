package datasources

import (
	"context"
	"fmt"

	qlikcloud "github.com/daniepett/qlik-cloud-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DataGatewayDataSource{}
	_ datasource.DataSourceWithConfigure = &DataGatewayDataSource{}
)

// NewSpacesDataSource is a helper function to simplify the provider implementation.
func NewDataGatewayDataSource() datasource.DataSource {
	return &DataGatewayDataSource{}
}

// DataGatewayDataSource is the data source implementation.
type DataGatewayDataSource struct {
	client *qlikcloud.Client
}

// DataGatewayModel maps coffees schema data.
type DataGatewayModel struct {
	ID types.String `tfsdk:"id"`
	// Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	SpaceID     types.String `tfsdk:"space_id"`
}

// Metadata returns the data source type name.
func (d *DataGatewayDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_gateway"
}

// Schema defines the schema for the data source.
func (d *DataGatewayDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
			},
			// "name": schema.StringAttribute{
			// 	Optional: true,
			// },
			"type": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"space_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *DataGatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DataGatewayModel

	// var DataGatewayName string
	var DataGatewayID string
	// var DataGateway *models.DataGateway
	var err error
	// req.Config.GetAttribute(ctx, path.Root("name"), &DataGatewayName)
	req.Config.GetAttribute(ctx, path.Root("id"), &DataGatewayID)
	// DataGatewayID := state.ID.ValueString()

	// if DataGatewayName != "" {
	// 	DataGateway, err = d.client.GetDataGatewayByName(DataGatewayName)
	// } else if DataGatewayID != "" {
	// 	DataGateway, err = d.client.GetDataGateway(DataGatewayID)
	// } else {
	// 	resp.Diagnostics.AddError(
	// 		"Both name and ID missing",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	dg, err := d.client.GetDataGateway(DataGatewayID)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Qlik Cloud DataGateway",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(dg.ID)
	state.Name = types.StringValue(dg.Name)
	state.Type = types.StringValue(dg.Type)
	state.Description = types.StringValue(dg.Description)
	state.SpaceID = types.StringValue(dg.SpaceID)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *DataGatewayDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = client
}
