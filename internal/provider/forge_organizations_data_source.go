package provider

import (
	"context"
	"fmt"
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
	Filters       []Filter                  `tfsdk:"filter"`
	Organizations []ForgeOrganizationModel  `tfsdk:"organizations"`
}

type ForgeOrganizationModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Slug        types.String `tfsdk:"slug"`
	Type        types.String `tfsdk:"type"`
	Avatar      types.String `tfsdk:"avatar"`
	TotalSeats  types.Int64  `tfsdk:"total_seats"`
	UsedSeats   types.Int64  `tfsdk:"used_seats"`
	CanManage   types.Bool   `tfsdk:"can_manage"`
	BillingLink types.String `tfsdk:"billing_link"`
	CreatedAt   types.String `tfsdk:"created_at"`
}

func (d *ForgeOrganizationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_organizations"
}

func (d *ForgeOrganizationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for listing Forge organizations. Use the `filter` block to specify the criteria for filtering organizations.",
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The field name to filter by (e.g., 'name', 'slug', or 'type')",
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
				Description: "List of organizations available in Forge",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":           schema.Int64Attribute{Computed: true, Description: "Organization ID"},
						"name":         schema.StringAttribute{Computed: true, Description: "Organization name"},
						"slug":         schema.StringAttribute{Computed: true, Description: "Organization slug"},
						"type":         schema.StringAttribute{Computed: true, Description: "Organization type"},
						"avatar":       schema.StringAttribute{Computed: true, Description: "Organization avatar URL"},
						"total_seats":  schema.Int64Attribute{Computed: true, Description: "Total seats in organization"},
						"used_seats":   schema.Int64Attribute{Computed: true, Description: "Used seats in organization"},
						"can_manage":   schema.BoolAttribute{Computed: true, Description: "Whether user can manage this organization"},
						"billing_link": schema.StringAttribute{Computed: true, Description: "Billing link for the organization"},
						"created_at":   schema.StringAttribute{Computed: true, Description: "Organization creation timestamp"},
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

	filteredOrganizations := filterForgeOrganizations(organizations, state.Filters)

	var orgModels []ForgeOrganizationModel
	for _, org := range filteredOrganizations {
		orgModels = append(orgModels, ForgeOrganizationModel{
			ID:          types.Int64Value(org.ID),
			Name:        types.StringValue(org.Name),
			Slug:        types.StringValue(org.Slug),
			Type:        types.StringValue(org.Type),
			Avatar:      types.StringValue(org.Avatar),
			TotalSeats:  types.Int64Value(int64(org.TotalSeats)),
			UsedSeats:   types.Int64Value(int64(org.UsedSeats)),
			CanManage:   types.BoolValue(org.CanManage),
			BillingLink: types.StringValue(org.BillingLink),
			CreatedAt:   types.StringValue(org.CreatedAt),
		})
	}
	state.Organizations = orgModels

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func filterForgeOrganizations(organizations []forge_client.Organization, filters []Filter) []forge_client.Organization {
	if len(filters) == 0 {
		return organizations
	}

	var filtered []forge_client.Organization

	for _, org := range organizations {
		match := true
		for _, f := range filters {
			switch f.Name.ValueString() {
			case "name":
				if !matchesFilter(org.Name, f.Values) {
					match = false
				}
			case "slug":
				if !matchesFilter(org.Slug, f.Values) {
					match = false
				}
			case "type":
				if !matchesFilter(org.Type, f.Values) {
					match = false
				}
			default:
				// Ignore unknown filters
				match = false
			}
		}

		if match {
			filtered = append(filtered, org)
		}
	}

	return filtered
}
