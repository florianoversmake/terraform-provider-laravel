package provider

import (
	"context"
	"fmt"
	"terraform-provider-laravel/internal/forge_client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ForgeServersDataSource{}

func NewForgeServersDataSource() datasource.DataSource {
	return &ForgeServersDataSource{}
}

type ForgeServersDataSource struct {
	client *forge_client.Client
}

type ForgeServersDataSourceModel struct {
	Filters []Filter           `tfsdk:"filter"`
	Servers []ForgeServerModel `tfsdk:"servers"`
}

type ForgeServerModel struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Provider         types.String `tfsdk:"provider"`
	Region           types.String `tfsdk:"region"`
	Size             types.String `tfsdk:"size"`
	IPAddress        types.String `tfsdk:"ip_address"`
	PrivateIPAddress types.String `tfsdk:"private_ip_address"`
	PHPVersion       types.String `tfsdk:"php_version"`
	DatabaseType     types.String `tfsdk:"database_type"`
	SSHPort          types.Int64  `tfsdk:"ssh_port"`
	IsReady          types.Bool   `tfsdk:"is_ready"`
	Revoked          types.Bool   `tfsdk:"revoked"`
	CreatedAt        types.String `tfsdk:"created_at"`
}

func (d *ForgeServersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_forge_servers"
}

func (d *ForgeServersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for listing Forge servers. Use the `filter` block to specify the criteria for filtering servers.",
		Blocks: map[string]schema.Block{
			"filter": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:    true,
							Description: "The field name to filter by (e.g., 'name', 'provider', 'region')",
						},
						"values": schema.ListAttribute{
							ElementType: types.StringType,
							Required:    true,
							Description: "The list of values to match for the specified field",
						},
					},
				},
				Description: "Filter block for selecting specific servers.",
			},
		},
		Attributes: map[string]schema.Attribute{
			"servers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of servers in Forge organization",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id":                 schema.Int64Attribute{Computed: true, Description: "The server ID"},
						"name":               schema.StringAttribute{Computed: true, Description: "The server name"},
						"provider":           schema.StringAttribute{Computed: true, Description: "The cloud provider (e.g., aws, ocean2, vultr2)"},
						"region":             schema.StringAttribute{Computed: true, Description: "The region code"},
						"size":               schema.StringAttribute{Computed: true, Description: "The server size"},
						"ip_address":         schema.StringAttribute{Computed: true, Description: "The public IP address"},
						"private_ip_address": schema.StringAttribute{Computed: true, Description: "The private IP address"},
						"php_version":        schema.StringAttribute{Computed: true, Description: "The PHP version"},
						"database_type":      schema.StringAttribute{Computed: true, Description: "The database type"},
						"ssh_port":           schema.Int64Attribute{Computed: true, Description: "The SSH port"},
						"is_ready":           schema.BoolAttribute{Computed: true, Description: "Whether the server is ready"},
						"revoked":            schema.BoolAttribute{Computed: true, Description: "Whether the server is revoked"},
						"created_at":         schema.StringAttribute{Computed: true, Description: "The creation timestamp"},
					},
				},
			},
		},
	}
}

func (d *ForgeServersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ForgeServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ForgeServersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	servers, err := d.client.ListServers(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error reading servers", err.Error())
		return
	}

	filteredServers := filterForgeServers(servers, state.Filters)

	var serverModels []ForgeServerModel
	for _, s := range filteredServers {
		model := ForgeServerModel{
			ID:           types.Int64Value(s.ID),
			Name:         types.StringValue(s.Name),
			Provider:     types.StringValue(s.Provider),
			Region:       types.StringValue(s.Region),
			Size:         types.StringValue(s.Size),
			PHPVersion:   types.StringValue(s.PHPVersion),
			DatabaseType: types.StringValue(s.DatabaseType),
			SSHPort:      types.Int64Value(int64(s.SSHPort)),
			IsReady:      types.BoolValue(s.IsReady),
			Revoked:      types.BoolValue(s.Revoked != nil && *s.Revoked),
			CreatedAt:    types.StringValue(s.CreatedAt),
		}
		if s.IPAddress != nil {
			model.IPAddress = types.StringValue(*s.IPAddress)
		} else {
			model.IPAddress = types.StringNull()
		}
		if s.PrivateIPAddress != nil {
			model.PrivateIPAddress = types.StringValue(*s.PrivateIPAddress)
		} else {
			model.PrivateIPAddress = types.StringNull()
		}
		serverModels = append(serverModels, model)
	}
	state.Servers = serverModels

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func filterForgeServers(servers []forge_client.Server, filters []Filter) []forge_client.Server {
	if len(filters) == 0 {
		return servers
	}

	var filtered []forge_client.Server
	for _, s := range servers {
		match := true
		for _, f := range filters {
			switch f.Name.ValueString() {
			case "name":
				if !matchesFilter(s.Name, f.Values) {
					match = false
				}
			case "provider":
				if !matchesFilter(s.Provider, f.Values) {
					match = false
				}
			case "region":
				if !matchesFilter(s.Region, f.Values) {
					match = false
				}
			case "ip_address":
				if s.IPAddress == nil || !matchesFilter(*s.IPAddress, f.Values) {
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
