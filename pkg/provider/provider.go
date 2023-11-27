package provider

import (
	"context"
	"os"

	qlikcloud "github.com/daniepett/qlik-cloud-client-go"
	"github.com/daniepett/terraform-provider-qlik/pkg/datasources"
	"github.com/daniepett/terraform-provider-qlik/pkg/resources"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &qlikProvider{}
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &qlikProvider{
			version: version,
		}
	}
}

type qlikProvider struct {
	version string
}

type qlikProviderModel struct {
	TenantID     types.String `tfsdk:"tenant_id"`
	Region       types.String `tfsdk:"region"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (p *qlikProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "qlik"
	resp.Version = p.version
}

func (p *qlikProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Required: true,
			},
			"region": schema.StringAttribute{
				Required: true,
			},
			"client_id": schema.StringAttribute{
				Optional: true,
			},
			"client_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *qlikProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config qlikProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	client_id := os.Getenv("QLIK_CLIENT_ID")
	client_secret := os.Getenv("QLIK_CLIENT_SECRET")
	tenant_id := os.Getenv("QLIK_TENANT_ID")
	region := os.Getenv("QLIK_REGION")

	if !config.TenantID.IsNull() {
		tenant_id = config.TenantID.ValueString()
	}

	if !config.Region.IsNull() {
		region = config.Region.ValueString()
	}

	if !config.ClientID.IsNull() {
		client_id = config.ClientID.ValueString()
	}

	if !config.ClientSecret.IsNull() {
		client_secret = config.ClientSecret.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if tenant_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("tenant_id"),
			"Missing Qlik Cloud Tenant ID",
			"The provider cannot create the Qlik Cloud API client as there is a missing or empty value for the Qlik Cloud Tenant ID. "+
				"Set the tenant_id value in the configuration or use the QLIK_CLOUD_TENANT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if region == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Missing Qlik Cloud Region",
			"The provider cannot create the Qlik Cloud API client as there is a missing or empty value for the Qlik Cloud region. "+
				"Set the region value in the configuration or use the QLIK_CLOUD_REGION environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if client_id == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_id"),
			"Missing Qlik Cloud Client ID",
			"The provider cannot create the Qlik Cloud API client as there is a missing or empty value for the Qlik Cloud Client ID. "+
				"Set the client_id value in the configuration or use the QLIK_CLOUD_CLIENT_ID environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if client_secret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("client_secret"),
			"Missing Qlik Cloud Client Secret",
			"The provider cannot create the Qlik Cloud API client as there is a missing or empty value for the Qlik Cloud Client Secret. "+
				"Set the client_secret value in the configuration or use the QLIK_CLOUD_CLIENT_SECRET environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Qlik Cloud client using the configuration values
	client, err := qlikcloud.NewClient(&tenant_id, &region, &client_id, &client_secret)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Qlik Cloud Client",
			"An unexpected error occurred when creating the Qlik Cloud client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Qlik Cloud Client Error: "+err.Error(),
		)
		return
	}

	// Make the Qlik client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *qlikProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewSpacesDataSource,
		datasources.NewSpaceDataSource,
		datasources.NewDataGatewayDataSource,
		datasources.NewDataConnectionsDataSource,
		datasources.NewSourceEntitiesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *qlikProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewSpaceResource,
		resources.NewDataConnectionResource,
		resources.NewDataProjectResource,
		resources.NewDataAppResource,
		resources.NewDataAppSourceSelectionResource,
	}
}
