package provider

import (
	"context"
	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ForgePHPVersionsDataSource{}

func NewForgePHPVersionsDataSource() datasource.DataSource {
	return &ForgePHPVersionsDataSource{}
}

type ForgePHPVersionsDataSource struct {
	client *forge_client.Client
}

type ForgePHPVersionsDataSourceModel struct {
	ServerID    types.Int64            `tfsdk:"server_id"`
	PHPVersions []ForgePHPVersionModel `tfsdk:"php_versions"`
}

type ForgePHPVersionModel struct {
	ID                 types.Int64  `tfsdk:"id"`
	Version            types.String `tfsdk:"version"`
	Status             types.String `tfsdk:"status"`
	DisplayableVersion types.String `tfsdk:"displayable_version"`
	BinaryName         types.String `tfsdk:"binary_name"`
	UsedOnCLI          types.Bool   `tfsdk:"used_on_cli"`
	UsedAsDefault      types.Bool   `tfsdk:"used_as_default"`
}

func (d *ForgePHPVersionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_php_versions"
}

func (d *ForgePHPVersionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for listing PHP versions installed on a Forge server.",
		Attributes: map[string]schema.Attribute{
			"server_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the server to list PHP versions from",
			},
			"php_versions": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of PHP versions on the server",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                  schema.Int64Attribute{Computed: true, Description: "The PHP version ID"},
						"version":             schema.StringAttribute{Computed: true, Description: "The PHP version identifier (e.g., php83)"},
						"status":              schema.StringAttribute{Computed: true, Description: "The installation status"},
						"displayable_version": schema.StringAttribute{Computed: true, Description: "The human-readable version name (e.g., PHP 8.3)"},
						"binary_name":         schema.StringAttribute{Computed: true, Description: "The PHP binary name"},
						"used_on_cli":         schema.BoolAttribute{Computed: true, Description: "Whether this version is used on CLI"},
						"used_as_default":     schema.BoolAttribute{Computed: true, Description: "Whether this version is the default"},
					},
				},
			},
		},
	}
}

func (d *ForgePHPVersionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ForgePHPVersionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ForgePHPVersionsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := int(state.ServerID.ValueInt64())
	versions, err := d.client.ListPHPVersions(ctx, serverID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading PHP versions", err.Error())
		return
	}

	var versionModels []ForgePHPVersionModel
	for _, v := range versions {
		model := ForgePHPVersionModel{
			ID:                 types.Int64Value(int64(v.ID)),
			Version:            types.StringValue(v.Version),
			Status:             types.StringValue(v.Status),
			DisplayableVersion: types.StringValue(v.DisplayableVersion),
			BinaryName:         types.StringValue(v.BinaryName),
			UsedOnCLI:          types.BoolValue(v.UsedOnCLI),
			UsedAsDefault:      types.BoolValue(v.UsedAsDefault),
		}
		versionModels = append(versionModels, model)
	}
	state.PHPVersions = versionModels

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}
