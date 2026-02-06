package provider

import (
	"context"
	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ForgeOrganizationsDataSource{}

func NewForgeOrganizationsDataSource() datasource.DataSource {
	return &ForgeOrganizationsDataSource{}
}

type ForgeOrganizationsDataSource struct {
	client *forge_client.Client
}

type ForgeOrganizationsDataSourceModel struct {
	Filters       []Filter                 `tfsdk:"filter"`
	Organizations []ForgeOrganizationModel `tfsdk:"organizations"`
}

type ForgeOrganizationModel struct {
	ID        types.Int64  `tfsdk:"id"`
	Slug      types.String `tfsdk:"slug"`
	Name      types.String `tfsdk:"name"`
	Owner     types.Bool   `tfsdk:"owner"`
	CreatedAt types.String `tfsdk:"created_at"`
}

func (d *ForgeOrganizationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_organizations"
}

func (d *ForgeOrganizationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for listing Forge organizations the user has access to.",
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
				Description: "Filter block for selecting specific organizations.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"organizations": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of organizations the user has access to",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":         schema.Int64Attribute{Computed: true, Description: "The organization ID"},
						"slug":       schema.StringAttribute{Computed: true, Description: "The organization slug"},
						"name":       schema.StringAttribute{Computed: true, Description: "The organization name"},
						"owner":      schema.BoolAttribute{Computed: true, Description: "Whether the user is the owner"},
						"created_at": schema.StringAttribute{Computed: true, Description: "The creation timestamp"},
					},
				},
			},
		},
	}
}

func (d *ForgeOrganizationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ForgeOrganizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ForgeOrganizationsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	organizations, err := d.client.ListOrganizations(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading organizations", err.Error())
		return
	}

	filteredOrgs := filterForgeOrganizations(organizations, state.Filters)

	var orgModels []ForgeOrganizationModel
	for _, o := range filteredOrgs {
		model := ForgeOrganizationModel{
			ID:    types.Int64Value(o.ID),
			Slug:  types.StringValue(o.Slug),
			Name:  types.StringValue(o.Name),
			Owner: types.BoolValue(o.Owner),
		}
		if o.CreatedAt != nil {
			model.CreatedAt = types.StringValue(*o.CreatedAt)
		} else {
			model.CreatedAt = types.StringNull()
		}
		orgModels = append(orgModels, model)
	}
	state.Organizations = orgModels

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func filterForgeOrganizations(orgs []forge_client.Organization, filters []Filter) []forge_client.Organization {
	if len(filters) == 0 {
		return orgs
	}

	var filtered []forge_client.Organization
	for _, o := range orgs {
		match := true
		for _, f := range filters {
			switch f.Name.ValueString() {
			case "name":
				if !matchesFilter(o.Name, f.Values) {
					match = false
				}
			case "slug":
				if !matchesFilter(o.Slug, f.Values) {
					match = false
				}
			}
			if !match {
				break
			}
		}
		if match {
			filtered = append(filtered, o)
		}
	}
	return filtered
}
