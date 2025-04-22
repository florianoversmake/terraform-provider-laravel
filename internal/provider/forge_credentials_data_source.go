package provider

import (
	"context"
	"fmt"
	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ForgeCredentialsDataSource{}

func NewForgeCredentialsDataSource() datasource.DataSource {
	return &ForgeCredentialsDataSource{}
}

type ForgeCredentialsDataSource struct {
	client *forge_client.Client
}

type ForgeCredentialsDataSourceModel struct {
	Filters     []Filter               `tfsdk:"filter"`
	Credentials []ForgeCredentialModel `tfsdk:"credentials"`
}

type ForgeCredentialModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
	Name types.String `tfsdk:"name"`
}

func (d *ForgeCredentialsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_credentials"
}

func (d *ForgeCredentialsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for listing Forge credentials. Use the `filter` block to specify the criteria for filtering credentials.",
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The field name to filter by (e.g., 'name' or 'type')",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "The list of values to match for the specified field",
						},
					},
				},
				Description: "Filter block for selecting specific credentials.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"credentials": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of credentials available in Forge",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":   schema.Int64Attribute{Computed: true},
						"type": schema.StringAttribute{Computed: true},
						"name": schema.StringAttribute{Computed: true},
					},
				},
			},
		},
	}
}

func (d *ForgeCredentialsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerConfig, ok := req.ProviderData.(*providerConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configure Type",
			fmt.Sprintf("Expected *providerConfig, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	if providerConfig.Forge == nil {
		resp.Diagnostics.AddError(
			"Forge Client Not Configured",
			"This resource requires the Forge API token to be configured in the provider. "+
				"Please set the 'forge_api_token' attribute in the provider configuration.",
		)
		return
	}

	d.client = providerConfig.Forge
}

func (d *ForgeCredentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ForgeCredentialsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	credentials, err := d.client.ListCredentials(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading actions", err.Error())
		return
	}

	filteredCredentials := filterForgeCredentials(credentials, state.Filters)

	var credentialModels []ForgeCredentialModel
	for _, c := range filteredCredentials {
		credentialModels = append(credentialModels, ForgeCredentialModel{
			ID:   types.Int64Value(c.ID),
			Type: types.StringValue(c.Type),
			Name: types.StringValue(c.Name),
		})
	}
	state.Credentials = credentialModels

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func filterForgeCredentials(credentials []forge_client.Credential, filters []Filter) []forge_client.Credential {
	if len(filters) == 0 {
		return credentials
	}

	var filtered []forge_client.Credential

	for _, c := range credentials {
		match := true
		for _, f := range filters {
			switch f.Name.ValueString() {
			case "name":
				if !matchesFilter(c.Name, f.Values) {
					match = false
				}
			case "type":
				if !matchesFilter(c.Type, f.Values) {
					match = false
				}
			default:
				// Ignore unknown filters
				match = false
			}
		}

		if match {
			filtered = append(filtered, c)
		}
	}

	return filtered
}
