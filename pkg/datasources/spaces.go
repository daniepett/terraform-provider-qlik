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
	_ datasource.DataSource              = &SpacesDataSource{}
	_ datasource.DataSourceWithConfigure = &SpacesDataSource{}
)

// NewSpacesDataSource is a helper function to simplify the provider implementation.
func NewSpacesDataSource() datasource.DataSource {
	return &SpacesDataSource{}
}

// SpacesDataSource is the data source implementation.
type SpacesDataSource struct {
	client *qlikcloud.Client
}

// SpacesDataSourceModel maps the data source schema data.
type SpacesDataSourceModel struct {
	Spaces []SpacesModel `tfsdk:"spaces"`
	Name   types.String  `tfsdk:"name"`
}

// SpacesModel maps coffees schema data.
type SpacesModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
}

// Metadata returns the data source type name.
func (d *SpacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_spaces"
}

// Schema defines the schema for the data source.
func (d *SpacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"spaces": schema.ListNestedAttribute{
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
						"description": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *SpacesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SpacesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	filter := models.Filter{
		Name:  state.Name.ValueString(),
		Limit: 10,
	}
	Spaces, err := d.client.GetSpaces(filter)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Qlik Cloud Spaces",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, space := range Spaces.Spaces {
		spaceState := SpacesModel{
			ID:          types.StringValue(space.ID),
			Name:        types.StringValue(space.Name),
			Type:        types.StringValue(space.Type),
			Description: types.StringValue(space.Description),
		}

		state.Spaces = append(state.Spaces, spaceState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *SpacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
