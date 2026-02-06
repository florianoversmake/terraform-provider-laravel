package provider

import (
	"context"
	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ForgeRegionsDataSource{}

func NewForgeRegionsDataSource() datasource.DataSource {
	return &ForgeRegionsDataSource{}
}

type ForgeRegionsDataSource struct {
	client *forge_client.Client
}

type ForgeRegionsDataSourceModel struct {
	ProviderID types.Int64        `tfsdk:"provider_id"`
	Filters    []Filter           `tfsdk:"filter"`
	Regions    []ForgeRegionModel `tfsdk:"regions"`
}

type ForgeRegionModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Code          types.String `tfsdk:"code"`
	AlternateCode types.String `tfsdk:"alternate_code"`
}

func (d *ForgeRegionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_regions"
}

func (d *ForgeRegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for listing Forge regions for a specific cloud provider.",
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The field name to filter by (e.g., 'name', 'code')",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "The list of values to match for the specified field",
						},
					},
				},
				Description: "Filter block for selecting specific regions.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"provider_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the cloud provider to get regions for",
			},
			"regions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of regions available for the specified provider",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":             schema.Int64Attribute{Computed: true, Description: "The region ID"},
						"name":           schema.StringAttribute{Computed: true, Description: "The region name"},
						"code":           schema.StringAttribute{Computed: true, Description: "The region code"},
						"alternate_code": schema.StringAttribute{Computed: true, Description: "The alternate region code"},
					},
				},
			},
		},
	}
}

func (d *ForgeRegionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ForgeRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ForgeRegionsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	regions, err := d.client.ListProviderRegions(ctx, state.ProviderID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading cloud provider regions", err.Error())
		return
	}

	filteredRegions := filterForgeRegions(regions, state.Filters)

	var regionModels []ForgeRegionModel
	for _, r := range filteredRegions {
		regionModels = append(regionModels, ForgeRegionModel{
			ID:            types.Int64Value(r.ID),
			Name:          types.StringValue(r.Name),
			Code:          types.StringValue(r.Code),
			AlternateCode: types.StringValue(r.AlternateCode),
		})
	}
	state.Regions = regionModels

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func filterForgeRegions(regions []forge_client.ProviderRegionInfo, filters []Filter) []forge_client.ProviderRegionInfo {
	if len(filters) == 0 {
		return regions
	}

	var filtered []forge_client.ProviderRegionInfo
	for _, r := range regions {
		match := true
		for _, f := range filters {
			switch f.Name.ValueString() {
			case "name":
				if !matchesFilter(r.Name, f.Values) {
					match = false
				}
			case "code":
				if !matchesFilter(r.Code, f.Values) {
					match = false
				}
			case "alternate_code":
				if !matchesFilter(r.AlternateCode, f.Values) {
					match = false
				}
			}
			if !match {
				break
			}
		}
		if match {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
