package datasources

import (
	"context"
	"fmt"

	qlikcloud "github.com/daniepett/qlik-cloud-client-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &SpaceDataSource{}
	_ datasource.DataSourceWithConfigure = &SpaceDataSource{}
)

// NewSpacesDataSource is a helper function to simplify the provider implementation.
func NewSpaceDataSource() datasource.DataSource {
	return &SpaceDataSource{}
}

// SpaceDataSource is the data source implementation.
type SpaceDataSource struct {
	client *qlikcloud.Client
}

// SpaceModel maps coffees schema data.
type SpaceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the data source type name.
func (d *SpaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_space"
}

// Schema defines the schema for the data source.
func (d *SpaceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *SpaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SpaceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	Space, err := d.client.GetSpace(state.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Qlik Cloud Space",
			err.Error(),
		)
		return
	}

	state.ID = types.StringValue(Space.ID)
	state.Name = types.StringValue(Space.Name)
	state.Type = types.StringValue(Space.Type)
	state.Description = types.StringValue(Space.Description)

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *SpaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
