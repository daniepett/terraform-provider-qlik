package datasources

import (
	"context"
	"fmt"

	qlikcloud "github.com/daniepett/qlik-cloud-client-go"
	"github.com/daniepett/qlik-cloud-client-go/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DataConnectionsDataSource{}
	_ datasource.DataSourceWithConfigure = &DataConnectionsDataSource{}
)

// NewSpacesDataSource is a helper function to simplify the provider implementation.
func NewDataConnectionsDataSource() datasource.DataSource {
	return &DataConnectionsDataSource{}
}

// spacesDataSource is the data source implementation.
type DataConnectionsDataSource struct {
	client *qlikcloud.Client
}

// spacesDataSourceModel maps the data source schema data.
type DataConnectionsDataSourceModel struct {
	DataConnections []DataConnectionsModel `tfsdk:"data_connections"`
}

// spacesModel maps coffees schema data.
type DataConnectionsModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	DataSourceID types.String `tfsdk:"data_source_id"`
}

// Metadata returns the data source type name.
func (d *DataConnectionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_connections"
}

// Schema defines the schema for the data source.
func (d *DataConnectionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"data_connections": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"data_source_id": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *DataConnectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state DataConnectionsDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	filter := models.Filter{
		Limit: 10,
	}
	connections, err := d.client.GetDataConnections(filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Qlik Cloud Spaces",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, connection := range connections.Connections {
		connectionState := DataConnectionsModel{
			ID:           types.StringValue(connection.ID),
			Name:         types.StringValue(connection.Name),
			Type:         types.StringValue(connection.Type),
			DataSourceID: types.StringValue(connection.DataSourceID),
		}

		state.DataConnections = append(state.DataConnections, connectionState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *DataConnectionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
