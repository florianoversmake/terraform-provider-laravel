package provider

import (
	"context"
	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ForgeProvidersDataSource{}

func NewForgeProvidersDataSource() datasource.DataSource {
	return &ForgeProvidersDataSource{}
}

type ForgeProvidersDataSource struct {
	client *forge_client.Client
}

type ForgeProvidersDataSourceModel struct {
	Filters   []Filter             `tfsdk:"filter"`
	Providers []ForgeProviderModel `tfsdk:"providers"`
}

type ForgeProviderModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Slug              types.String `tfsdk:"slug"`
	SimpleName        types.String `tfsdk:"simple_name"`
	Currency          types.String `tfsdk:"currency"`
	CurrencySymbol    types.String `tfsdk:"currency_symbol"`
	DefaultSizeCode   types.String `tfsdk:"default_size_code"`
	DefaultRegionCode types.String `tfsdk:"default_region_code"`
}

func (d *ForgeProvidersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_providers"
}

func (d *ForgeProvidersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for listing Forge cloud providers (DigitalOcean, AWS, Hetzner, etc.).",
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The field name to filter by (e.g., 'name', 'slug')",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "The list of values to match for the specified field",
						},
					},
				},
				Description: "Filter block for selecting specific cloud providers.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"providers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of cloud providers available in Forge",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                  schema.Int64Attribute{Computed: true, Description: "The provider ID"},
						"name":                schema.StringAttribute{Computed: true, Description: "The provider name"},
						"slug":                schema.StringAttribute{Computed: true, Description: "The provider slug"},
						"simple_name":         schema.StringAttribute{Computed: true, Description: "The provider simple name"},
						"currency":            schema.StringAttribute{Computed: true, Description: "The currency used by this provider"},
						"currency_symbol":     schema.StringAttribute{Computed: true, Description: "The currency symbol"},
						"default_size_code":   schema.StringAttribute{Computed: true, Description: "The default size code"},
						"default_region_code": schema.StringAttribute{Computed: true, Description: "The default region code"},
					},
				},
			},
		},
	}
}

func (d *ForgeProvidersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerConfig, ok := req.ProviderData.(*providerConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Configure Type",
			"Expected *providerConfig. Please report this issue to the provider developers.",
		)
		return
	}

	if providerConfig.Forge == nil {
		resp.Diagnostics.AddError(
			"Forge Client Not Configured",
			"This data source requires the Forge API token to be configured in the provider.",
		)
		return
	}

	d.client = providerConfig.Forge
}

func (d *ForgeProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ForgeProvidersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	providers, err := d.client.ListProviders(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading cloud providers", err.Error())
		return
	}

	filteredProviders := filterForgeProviders(providers, state.Filters)

	var providerModels []ForgeProviderModel
	for _, p := range filteredProviders {
		providerModels = append(providerModels, ForgeProviderModel{
			ID:                types.Int64Value(p.ID),
			Name:              types.StringValue(p.Name),
			Slug:              types.StringValue(p.Slug),
			SimpleName:        types.StringValue(p.SimpleName),
			Currency:          types.StringValue(p.Currency),
			CurrencySymbol:    types.StringValue(p.CurrencySymbol),
			DefaultSizeCode:   types.StringValue(p.DefaultSizeCode),
			DefaultRegionCode: types.StringValue(p.DefaultRegionCode),
		})
	}
	state.Providers = providerModels

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func filterForgeProviders(providers []forge_client.ProviderInfo, filters []Filter) []forge_client.ProviderInfo {
	if len(filters) == 0 {
		return providers
	}

	var filtered []forge_client.ProviderInfo
	for _, p := range providers {
		match := true
		for _, f := range filters {
			switch f.Name.ValueString() {
			case "name":
				if !matchesFilter(p.Name, f.Values) {
					match = false
				}
			case "slug":
				if !matchesFilter(p.Slug, f.Values) {
					match = false
				}
			case "simple_name":
				if !matchesFilter(p.SimpleName, f.Values) {
					match = false
				}
			}
			if !match {
				break
			}
		}
		if match {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
