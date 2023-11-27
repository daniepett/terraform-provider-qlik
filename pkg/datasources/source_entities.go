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
	_ datasource.DataSource              = &SourceEntitiesDataSource{}
	_ datasource.DataSourceWithConfigure = &SourceEntitiesDataSource{}
)

// NewSpacesDataSource is a helper function to simplify the provider implementation.
func NewSourceEntitiesDataSource() datasource.DataSource {
	return &SourceEntitiesDataSource{}
}

// spacesDataSource is the data source implementation.
type SourceEntitiesDataSource struct {
	client *qlikcloud.Client
}

// spacesDataSourceModel maps the data source schema data.
type SourceEntitiesDataSourceModel struct {
	Entities           []SourceEntityModel `tfsdk:"entities"`
	ProjectID          types.String        `tfsdk:"project_id"`
	AppID              types.String        `tfsdk:"app_id"`
	SourceConnectionID types.String        `tfsdk:"source_connection_id"`
	Database           types.String        `tfsdk:"database"`
	TablePattern       types.String        `tfsdk:"table_pattern"`
	SchemaPattern      types.String        `tfsdk:"schema_pattern"`
	EntityType         types.String        `tfsdk:"entity_type"`
}

type SourceEntityModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	DataAppID types.String `tfsdk:"data_app_id"`
	Schema    types.String `tfsdk:"schema"`
	Database  types.String `tfsdk:"database"`
	Type      types.String `tfsdk:"type"`
	ProjectID types.String `tfsdk:"project_id"`
}

// Metadata returns the data source type name.
func (d *SourceEntitiesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_entities"
}

// Schema defines the schema for the data source.
func (d *SourceEntitiesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"entities": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"data_app_id": schema.StringAttribute{
							Computed: true,
						},
						"schema": schema.StringAttribute{
							Computed: true,
						},
						"database": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"project_id": schema.StringAttribute{
							Computed: true,
						},
					},
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
			"database": schema.StringAttribute{
				Required: true,
			},
			"table_pattern": schema.StringAttribute{
				Required: true,
			},
			"schema_pattern": schema.StringAttribute{
				Required: true,
			},
			"entity_type": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *SourceEntitiesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SourceEntitiesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	ent := models.GetSourceEntities{
		SourceSelection: models.GetSourceEntitiesSourceSelection{
			DataEntitiesSelection: models.GetSourceEntitiesDataEntitiesSelection{
				SourceConnectionID: state.SourceConnectionID.ValueString(),
				IncludePatterns: []models.GetSourceEntitiesIncludePatterns{
					models.GetSourceEntitiesIncludePatterns{
						ProjectID:     state.ProjectID.ValueString(),
						Database:      state.Database.ValueString(),
						TablePattern:  state.TablePattern.ValueString(),
						SchemaPattern: state.SchemaPattern.ValueString(),
						EntityType:    state.EntityType.ValueString(),
					},
				},
			},
		},
	}

	s, err := d.client.GetSourceEntities(state.ProjectID.ValueString(), state.AppID.ValueString(), ent)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Qlik Cloud Spaces",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, entity := range s.Entities {
		entityState := SourceEntityModel{
			ID:        types.StringValue(entity.ID),
			Name:      types.StringValue(entity.Name),
			DataAppID: types.StringValue(entity.DataAppID),
			Schema:    types.StringValue(entity.Schema),
			Database:  types.StringValue(entity.Database),
			Type:      types.StringValue(entity.Type),
			ProjectID: types.StringValue(entity.ProjectID),
		}

		state.Entities = append(state.Entities, entityState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *SourceEntitiesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
