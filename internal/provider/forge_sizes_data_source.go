package provider

import (
	"context"
	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ForgeSizesDataSource{}

func NewForgeSizesDataSource() datasource.DataSource {
	return &ForgeSizesDataSource{}
}

type ForgeSizesDataSource struct {
	client *forge_client.Client
}

type ForgeSizesDataSourceModel struct {
	ProviderID types.Int64      `tfsdk:"provider_id"`
	Filters    []Filter         `tfsdk:"filter"`
	Sizes      []ForgeSizeModel `tfsdk:"sizes"`
}

type ForgeSizeModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Code         types.String `tfsdk:"code"`
	Series       types.String `tfsdk:"series"`
	Category     types.String `tfsdk:"category"`
	Cpus         types.Int64  `tfsdk:"cpus"`
	DiskType     types.String `tfsdk:"disk_type"`
	Architecture types.String `tfsdk:"architecture"`
	Ram          types.Int64  `tfsdk:"ram"`
	Disk         types.Int64  `tfsdk:"disk"`
}

func (d *ForgeSizesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_sizes"
}

func (d *ForgeSizesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for listing Forge server sizes for a specific cloud provider.",
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The field name to filter by (e.g., 'name', 'code', 'category')",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "The list of values to match for the specified field",
						},
					},
				},
				Description: "Filter block for selecting specific sizes.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"provider_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the cloud provider to get sizes for",
			},
			"sizes": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of sizes available for the specified provider",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":           schema.Int64Attribute{Computed: true, Description: "The size ID"},
						"name":         schema.StringAttribute{Computed: true, Description: "The size name"},
						"code":         schema.StringAttribute{Computed: true, Description: "The size code"},
						"series":       schema.StringAttribute{Computed: true, Description: "The size series"},
						"category":     schema.StringAttribute{Computed: true, Description: "The size category"},
						"cpus":         schema.Int64Attribute{Computed: true, Description: "The number of CPUs"},
						"disk_type":    schema.StringAttribute{Computed: true, Description: "The disk type"},
						"architecture": schema.StringAttribute{Computed: true, Description: "The CPU architecture"},
						"ram":          schema.Int64Attribute{Computed: true, Description: "The RAM in MB"},
						"disk":         schema.Int64Attribute{Computed: true, Description: "The disk size in GB"},
					},
				},
			},
		},
	}
}

func (d *ForgeSizesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ForgeSizesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ForgeSizesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sizes, err := d.client.ListProviderSizes(ctx, state.ProviderID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading cloud provider sizes", err.Error())
		return
	}

	filteredSizes := filterForgeSizes(sizes, state.Filters)

	var sizeModels []ForgeSizeModel
	for _, s := range filteredSizes {
		sizeModels = append(sizeModels, ForgeSizeModel{
			ID:           types.Int64Value(s.ID),
			Name:         types.StringValue(s.Name),
			Code:         types.StringValue(s.Code),
			Series:       types.StringValue(s.Series),
			Category:     types.StringValue(s.Category),
			Cpus:         types.Int64Value(int64(s.Cpus)),
			DiskType:     types.StringValue(s.DiskType),
			Architecture: types.StringValue(s.Architecture),
			Ram:          types.Int64Value(int64(s.Ram)),
			Disk:         types.Int64Value(int64(s.Disk)),
		})
	}
	state.Sizes = sizeModels

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func filterForgeSizes(sizes []forge_client.ProviderSizeInfo, filters []Filter) []forge_client.ProviderSizeInfo {
	if len(filters) == 0 {
		return sizes
	}

	var filtered []forge_client.ProviderSizeInfo
	for _, s := range sizes {
		match := true
		for _, f := range filters {
			switch f.Name.ValueString() {
			case "name":
				if !matchesFilter(s.Name, f.Values) {
					match = false
				}
			case "code":
				if !matchesFilter(s.Code, f.Values) {
					match = false
				}
			case "series":
				if !matchesFilter(s.Series, f.Values) {
					match = false
				}
			case "category":
				if !matchesFilter(s.Category, f.Values) {
					match = false
				}
			case "disk_type":
				if !matchesFilter(s.DiskType, f.Values) {
					match = false
				}
			case "architecture":
				if !matchesFilter(s.Architecture, f.Values) {
					match = false
				}
			}
			if !match {
				break
			}
		}
		if match {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
