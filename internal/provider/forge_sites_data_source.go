package provider

import (
	"context"
	"fmt"
	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ForgeSitesDataSource{}

func NewForgeSitesDataSource() datasource.DataSource {
	return &ForgeSitesDataSource{}
}

type ForgeSitesDataSource struct {
	client *forge_client.Client
}

type ForgeSitesDataSourceModel struct {
	ServerID types.Int64      `tfsdk:"server_id"`
	Filters  []Filter         `tfsdk:"filter"`
	Sites    []ForgeSiteModel `tfsdk:"sites"`
}

type ForgeSiteModel struct {
	ID                      types.Int64  `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	WebDirectory            types.String `tfsdk:"web_directory"`
	PHPVersion              types.String `tfsdk:"php_version"`
	AppType                 types.String `tfsdk:"app_type"`
	Status                  types.String `tfsdk:"status"`
	Isolated                types.Bool   `tfsdk:"isolated"`
	User                    types.String `tfsdk:"user"`
	QuickDeploy             types.Bool   `tfsdk:"quick_deploy"`
	DeploymentURL           types.String `tfsdk:"deployment_url"`
	HTTPS                   types.Bool   `tfsdk:"https"`
	ZeroDowntimeDeployments types.Bool   `tfsdk:"zero_downtime_deployments"`
	CreatedAt               types.String `tfsdk:"created_at"`
}

func (d *ForgeSitesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_sites"
}

func (d *ForgeSitesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for listing Forge sites on a server. Use the `filter` block to specify the criteria for filtering sites.",
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The field name to filter by (e.g., 'name', 'status', 'php_version')",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "The list of values to match for the specified field",
						},
					},
				},
				Description: "Filter block for selecting specific sites.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"server_id": schema.Int64Attribute{
				Required:    true,
				Description: "The ID of the server to list sites from",
			},
			"sites": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of sites on the server",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                        schema.Int64Attribute{Computed: true, Description: "The site ID"},
						"name":                      schema.StringAttribute{Computed: true, Description: "The site name (domain)"},
						"web_directory":             schema.StringAttribute{Computed: true, Description: "The web directory path"},
						"php_version":               schema.StringAttribute{Computed: true, Description: "The PHP version"},
						"app_type":                  schema.StringAttribute{Computed: true, Description: "The application type"},
						"status":                    schema.StringAttribute{Computed: true, Description: "The site status"},
						"isolated":                  schema.BoolAttribute{Computed: true, Description: "Whether the site runs in isolation"},
						"user":                      schema.StringAttribute{Computed: true, Description: "The user running the site"},
						"quick_deploy":              schema.BoolAttribute{Computed: true, Description: "Whether quick deploy is enabled"},
						"deployment_url":            schema.StringAttribute{Computed: true, Description: "The deployment URL"},
						"https":                     schema.BoolAttribute{Computed: true, Description: "Whether HTTPS is enabled"},
						"zero_downtime_deployments": schema.BoolAttribute{Computed: true, Description: "Whether zero-downtime deployments are enabled"},
						"created_at":                schema.StringAttribute{Computed: true, Description: "The creation timestamp"},
					},
				},
			},
		},
	}
}

func (d *ForgeSitesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
			"This data source requires the Forge API token to be configured in the provider.",
		)
		return
	}

	d.client = providerConfig.Forge
}

func (d *ForgeSitesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ForgeSitesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverID := int(state.ServerID.ValueInt64())
	sites, err := d.client.ListSites(ctx, serverID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading sites", err.Error())
		return
	}

	filteredSites := filterForgeSites(sites, state.Filters)

	var siteModels []ForgeSiteModel
	for _, s := range filteredSites {
		model := ForgeSiteModel{
			ID:                      types.Int64Value(s.ID),
			Name:                    types.StringValue(s.Name),
			WebDirectory:            types.StringValue(s.WebDirectory),
			PHPVersion:              types.StringValue(s.PHPVersion),
			AppType:                 types.StringValue(s.AppType),
			Status:                  types.StringValue(s.Status),
			Isolated:                types.BoolValue(s.Isolated),
			User:                    types.StringValue(s.User),
			DeploymentURL:           types.StringValue(s.DeploymentURL),
			HTTPS:                   types.BoolValue(s.HTTPS),
			ZeroDowntimeDeployments: types.BoolValue(s.ZeroDowntimeDeployments),
		}
		if s.QuickDeploy != nil {
			model.QuickDeploy = types.BoolValue(*s.QuickDeploy)
		} else {
			model.QuickDeploy = types.BoolNull()
		}
		if s.CreatedAt != nil {
			model.CreatedAt = types.StringValue(*s.CreatedAt)
		} else {
			model.CreatedAt = types.StringNull()
		}
		siteModels = append(siteModels, model)
	}
	state.Sites = siteModels

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func filterForgeSites(sites []forge_client.Site, filters []Filter) []forge_client.Site {
	if len(filters) == 0 {
		return sites
	}

	var filtered []forge_client.Site
	for _, s := range sites {
		match := true
		for _, f := range filters {
			switch f.Name.ValueString() {
			case "name":
				if !matchesFilter(s.Name, f.Values) {
					match = false
				}
			case "status":
				if !matchesFilter(s.Status, f.Values) {
					match = false
				}
			case "php_version":
				if !matchesFilter(s.PHPVersion, f.Values) {
					match = false
				}
			case "app_type":
				if !matchesFilter(s.AppType, f.Values) {
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
